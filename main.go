package main

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	rand2 "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/chfanghr/trackDownloader/Spotify"
	"github.com/chfanghr/trackDownloader/librespot"
	"github.com/chfanghr/trackDownloader/librespot/core"
	"github.com/chfanghr/trackDownloader/librespot/utils"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	authBufFile            = flag.String("authBufFile", "", "path to authBuffer file")
	authBufPassword        = flag.String("authBufPassword", "", "Password of authBuffer file")
	username               = flag.String("username", "", "name of Spotify account")
	password               = flag.String("password", "", "Password of Spotify account")
	deviceName             = flag.String("deviceName", "trackdl", "name of device")
	logFile                = flag.String("logFile", "", "path to log file")
	saveFileTo             = flag.String("saveFileTo", "./", "path to save audio file")
	vorbisComment          = flag.String("vorbisComment", "./vorbiscomment", "path to vorbisComment executable file")
	saveAuthBufTo          = flag.String("saveAuthBufTo", "", "path to save authBuffer")
	saveAuthBufPassword    = flag.String("saveAuthBufPassword", RandStringRunes(5), "Password of saved authBuffer file")
	targetsToSearch        = flag.String("search", "", "targets to search for,split by \",\"")
	albumURIsToView        = flag.String("viewAlbum", "", "URIs of albums to view,split by \",\"")
	trackURIsToView        = flag.String("viewTrack", "", "URIs of tracks to view,split by \",\"")
	artistURIsToView       = flag.String("viewArtist", "", "URIs of artists to view,split by \",\"")
	playlistURIsToView     = flag.String("viewPlaylist", "", "URIs of playlists to view,split by \",\"")
	URLsToView             = flag.String("view", "", "URIs to view,begin with https://open.spotify.com/,split by \",\"")
	trackURIsToDownload    = flag.String("downloadTrack", "", "URIs of tracks to download,split by \",\"")
	albumURIsToDownload    = flag.String("downloadAlbum", "", "URIs of albums to download,split by \",\"")
	playlistURIsToDownload = flag.String("downloadPlaylist", "", "URIs of playlists,split by \",\"")
	URLsToDownload         = flag.String("download", "", "URIs to download,split by \",\"")
	limitOfSearchResult    = flag.Int("searchResultLimit", 12, "limit of search result to shaw")
	quality                = flag.Int("quality", 320, "quality of audio file")
	viewRootPlaylist       = flag.Bool("viewRootPlaylists", false, "view root playlist or not")
	quiet                  = flag.Bool("quiet", false, "output log to stdout or not")
)

var (
	logger            *log.Logger     = nil
	session           *core.Session   = nil
	downloadWaitGroup *sync.WaitGroup = nil
	vorbisPath                        = ""
	//globalContext     context.Context    = nil
	cancelFunc  context.CancelFunc       = nil
	version                              = "DEBUG"
	realQuality Spotify.AudioFile_Format = Spotify.AudioFile_OGG_VORBIS_160
)

func RandStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type authBuf struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

func (a *authBuf) createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (a *authBuf) encrypt(password string) ([]byte, error) {
	buf, err := json.Marshal(a)
	if err != nil {
		return nil, errors.New("error occur while encrypt authBuf : " + err.Error())
	}
	block, _ := aes.NewCipher([]byte(a.createHash(password)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.New("error occur while encrypt authBuf : " + err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand2.Reader, nonce); err != nil {
		return nil, errors.New("error occur while encrypt authBuf : " + err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, buf, nil)
	return ciphertext, nil
}

func (a *authBuf) decrypt(password string, encrypted []byte) error {
	key := []byte(a.createHash(password))
	block, err := aes.NewCipher(key)
	if err != nil {
		return errors.New("error occur while decrypt authbuf : " + err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return errors.New("error occur while decrypt authbuf : " + err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := encrypted[:nonceSize], encrypted[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	err = json.Unmarshal(plaintext, a)
	if err != nil {
		return errors.New("error occur while decrypt authbuf : " + err.Error())
	}
	return nil
}

func (a *authBuf) GetPassword() string {
	return a.Password
}

func (a *authBuf) GetUsername() string {
	return a.Username
}

func (a *authBuf) SetPassword(p string) {
	a.Password = p
}

func (a *authBuf) SetUsername(u string) {
	a.Username = u
}

func downloadErrorHandler(err error) {
	if err != nil {
		logger.Println("error occur while downloading :", err.Error())
	}
}

func downloadWithURLs() {
	tmp := getURLsToDownload()
	for _, v := range tmp {
		str := strings.TrimPrefix(v, "https://open.spotify.com/")
		methodAndURI := strings.Split(str, "/")
		if len(methodAndURI) != 2 {
			logger.Println("error occur while parsing URI : invalid URI", v)
			continue
		}
		switch methodAndURI[0] {
		case "album":
			if *albumURIsToDownload != "" {
				*albumURIsToDownload += ","
			}
			*albumURIsToDownload += methodAndURI[1]
			break
		case "track":
			if *trackURIsToDownload != "" {
				*trackURIsToView += ","
			}
			*trackURIsToDownload += methodAndURI[1]
			break
		//case "playlist":
		//	if *playlistURIsToView != "" {
		//		*playlistURIsToView += ","
		//	}
		//	*playlistURIsToView += methodAndURI[1]
		//	break
		default:
			logger.Println("error occur while parsing URI : invalid URI", v)
		}
	}
}

func downloadAlbums() {
	logger.SetPrefix("download albums ")
	logger.SetFlags(log.LstdFlags)
	defer logger.SetPrefix(*deviceName + " ")
	defer logger.SetFlags(log.LstdFlags)
	targets := getAlbumsURIsToDownload()
	for _, v := range targets {
		logger.Printf("loading album %s for downloading\n", v)
		album, err := session.Mercury().GetAlbum(utils.Base62ToHex(v))
		if err != nil || album == nil {
			logger.Printf("error occur while gathering album data of %s : %s\n", v, err)
			continue
		}
		logger.Println("Title :", album.GetName())
		logger.Println("Artists :", func() string {
			var res string
			for i, artist := range album.GetArtist() {
				if i != 0 {
					res += ","
				}
				res += artist.GetName()
			}
			return res
		}())
		for _, d := range album.GetDisc() {
			for _, t := range d.GetTrack() {
				if *trackURIsToDownload != "" {
					*trackURIsToDownload += ","
				}
				*trackURIsToDownload += utils.ConvertTo62(t.GetGid())
			}
		}
	}
}

func downloadTrack(id string) error {
	ss := session
	logger.Println("loading track for downloading :", id)
	track, err := ss.Mercury().GetTrack(utils.Base62ToHex(id))
	if err != nil {
		return fmt.Errorf("error occur while loading track %s : %s", id, err)
	}
	if track.GetName() == "" {
		return fmt.Errorf("error occur while loading track %s : fail to get track", id)
	}
	go downloadTrackInternal(track)
	return nil
}

func downloadTrackInternal(track *Spotify.Track) error {
	downloadWaitGroup.Add(1)
	defer downloadWaitGroup.Done()
	ss := session

	// show information of tack

	logger.Println("Title :", track.GetName())
	logger.Println("Album :", track.GetAlbum().GetName())
	logger.Println("Artists :", func() string {
		var res string
		for i, artist := range track.GetAlbum().GetArtist() {
			if i != 0 {
				res += ","
			}
			res += artist.GetName()
		}
		return res
	}())
	logger.Println("External id :", func() string {
		var res string
		for i, eid := range track.GetExternalId() {
			if i != 0 {
				res += ","
			}
			res += "{ type : " + eid.GetTyp() + " , id : " + eid.GetId() + " }"
		}
		return res
	}())

	// select a file to download
	var selectedFile *Spotify.AudioFile
	for _, file := range track.GetFile() {
		if file.GetFormat() == realQuality {
			selectedFile = file
		}
	}

	if selectedFile == nil {
		logger.Println("failed to fetch", track.GetName(), ": cannot find audio file")
		return nil
	} else {
		// fetch audio file
		logger.Println(track.GetName(), "fetching audio file")
		audioFile, err := ss.Player().LoadTrack(selectedFile, track.GetGid())
		if err != nil {
			return fmt.Errorf("error occur while fetching %s : %s", track.GetName(), err)
		} else {
			buf, err := ioutil.ReadAll(audioFile)
			if err != nil {
				return fmt.Errorf("error occur while fetching %s : %s", track.GetName(), err)
			}
			logger.Println(track.GetName(), "fetched successfully")

			outputFile := *saveFileTo + "/" + track.GetAlbum().GetName() + "-" + track.GetName() + "-" + RandStringRunes(10) + ".ogg"

			err = ioutil.WriteFile(outputFile, buf, 0666)
			if err != nil {
				return fmt.Errorf("error occur while saving %s : %s", track.GetName(), err)
			}

			var metadata string
			metadata = metadata + "ALBUM=" + track.Album.GetName() + "\n"
			metadata = metadata + "ARTIST=" + func() string {
				var res string
				for i, artist := range track.GetAlbum().GetArtist() {
					if i != 0 {
						res += ","
					}
					res += artist.GetName()
				}
				return res
			}() + "\n"
			metadata = metadata + "TITLE=" + track.GetName() + "\n"
			metadata = metadata + "GENRE=" + func() string {
				var res string
				for i, v := range track.Album.GetGenre() {
					if i > 0 {
						res += ","
					}
					res += v
				}
				return res

			}() + "\n"
			metadata = metadata + "DATE=" + track.GetAlbum().GetDate().String() + "\n"
			metadata = metadata + "COPYRIGHT=" + func() string {
				var res string
				for i, v := range track.GetAlbum().GetCopyright() {
					if i > 0 {
						res += ","
					}
					res += v.GetText()
				}
				return res
			}()

			tmpMetadataFile := *saveFileTo + "/." + RandStringRunes(10) + ".metadata"
			if version != "DEBUG" {
				defer os.Remove(tmpMetadataFile)
			}
			err = ioutil.WriteFile(tmpMetadataFile, []byte(metadata), 0666)
			if err != nil {
				return fmt.Errorf("error occur while writing metadata to track %s : %s", track.GetName(), err)
			}

			vorbisCommentCommand := exec.Command(vorbisPath, "-a", outputFile, "-c", tmpMetadataFile)
			err = vorbisCommentCommand.Run()

			if err != nil {
				return fmt.Errorf("error occur while writing metadata to track %s : %s", track.GetName(), err)
			}

			logger.Println(track.GetName(), "downloaded successfully")
			return nil
		}
	}
}

//func downloadTrackInternal(track *Spotify.Track) error {
//	downloadChan := func() chan error {
//		ch := make(chan error)
//		go func() {
//			ch <- func() error {
//				downloadWaitGroup.Add(1)
//				defer downloadWaitGroup.Done()
//				ss := session
//
//				// show information of tack
//
//				logger.Println("Title :", track.GetName())
//				logger.Println("Album :", track.GetAlbum().GetName())
//				logger.Println("Artists :", func() string {
//					var res string
//					for i, artist := range track.GetAlbum().GetArtist() {
//						if i != 0 {
//							res += ","
//						}
//						res += artist.GetName()
//					}
//					return res
//				}())
//				logger.Println("External id :", func() string {
//					var res string
//					for i, eid := range track.GetExternalId() {
//						if i != 0 {
//							res += ","
//						}
//						res += "{ type : " + eid.GetTyp() + " , id : " + eid.GetId() + " }"
//					}
//					return res
//				}())
//
//				// select a file to download
//				var selectedFile *Spotify.AudioFile
//				for _, file := range track.GetFile() {
//					if file.GetFormat() == Spotify.AudioFile_Format(*quality) {
//						selectedFile = file
//					}
//				}
//
//				if selectedFile == nil {
//					logger.Println("failed to fetch", track.GetName(), ": cannot find audio file")
//					return nil
//				} else {
//					// fetch audio file
//					logger.Println(track.GetName(), "fetching audio file")
//					audioFile, err := ss.Player().LoadTrack(selectedFile, track.GetGid())
//					if err != nil {
//						return fmt.Errorf("error occur while fetching %s : %s", track.GetName(), err)
//					} else {
//						buf, err := ioutil.ReadAll(audioFile)
//						if err != nil {
//							return fmt.Errorf("error occur while fetching %s : %s", track.GetName(), err)
//						}
//						logger.Println(track.GetName(), "fetched successfully , save to tmp file")
//
//						// save to tmp file
//						tmpFile := *saveFileTo + "/." + RandStringRunes(10) + ".ogg"
//						err = ioutil.WriteFile(tmpFile, buf, 0666)
//						if err != nil {
//							return fmt.Errorf("error occur while saving tmp file of %s : %s", track.GetName(), err)
//						}
//						defer os.Remove(tmpFile)
//
//						// convert file
//						// Why do I use syscall.ForkExec?
//						// exec.run() seems to have some bugs
//						// it never execute vorbisComment
//						logger.Println(track.GetName(), "call vorbisComment to convert tmp file")
//						ffmpegArgs := []string{vorbisPath, "-i", tmpFile, *saveFileTo + "/" + track.GetAlbum().GetName() + "_" + track.GetName() + "_" + RandStringRunes(10) + ".mp3"}
//						workingDirectory, err := os.Getwd()
//						if err != nil {
//							return fmt.Errorf("error occur while converting %s : %s", track.GetName(), err)
//						}
//						pid, err := syscall.ForkExec(vorbisPath, ffmpegArgs, &syscall.ProcAttr{
//							Dir: workingDirectory,
//						})
//						if err != nil {
//							return fmt.Errorf("error occur while converting %s : %s", track.GetName(), err)
//						}
//						logger.Println(track.GetName(), "waiting for vorbisComment")
//						process, err := os.FindProcess(pid)
//						if err != nil {
//							return fmt.Errorf("error occur while converting %s : %s", track.GetName(), err)
//						}
//						_, err = process.Wait()
//						if err != nil {
//							return fmt.Errorf("error occur while converting %s : %s", track.GetName(), err)
//						}
//
//						logger.Println(track.GetName() + " converted successfully")
//						return nil
//					}
//				}
//			}()
//		}()
//		return ch
//	}()
//	select {
//	case <-globalContext.Done():
//		return errors.New("download job canceled")
//	case err := <-downloadChan:
//		return err
//	}
//}

func setupLogger() {
	nullDev, err := os.OpenFile("/dev/null", os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("setup logger failed :", err)
	}
	logDst := append([]io.Writer{}, nullDev)
	if !*quiet {
		logDst = append(logDst, os.Stdout)
	}
	if *logFile != "" {
		lf, err := os.OpenFile(*logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("setup logger failed :", err)
		}
		logDst = append(logDst, lf)
	}
	multiWriter := io.MultiWriter(logDst...)
	logger = log.New(multiWriter, *deviceName+" ", log.LstdFlags)
	core.SetLogger(logger)
}

func setupDownloadWaitGroup() {
	downloadWaitGroup = new(sync.WaitGroup)
}

func login() {
	if *authBufFile != "" {
		buf, err := ioutil.ReadFile(*authBufFile)
		if err != nil {
			logger.Fatalln("error occur while reading authBuffer file :", err)
		}
		ab := &authBuf{}
		err = ab.decrypt(*authBufPassword, buf)
		if err != nil {
			logger.Fatalln("error occur while reading authBuffer file :", err)
		}
		*username = ab.GetUsername()
		*password = ab.GetPassword()
		logger.Println(*username, "login with authBuffer")
	} else {
		if *username == "" {
			logger.Fatalln("please provide a nonempty username")
		}
		logger.Println(*username, "login with Password")
	}
	if *username == "" {
		logger.Fatalln("please provide a nonempty username")
	}
	if *password != "" {
		ss, err := librespot.Login(*username, *password, *deviceName)
		if err != nil {
			logger.Fatalln("error occur while logging in :", err.Error())
		}
		session = ss
	} else {
		logger.Fatalln("please provide Password or blob file")
	}
	if session == nil {
		logger.Fatalln("error occur while logging in : unknown error")
	}
	logger.Println(*username, "logged in successfully")
	if *saveAuthBufTo != "" {
		logger.Println("save authBuffer with Password :", *saveAuthBufPassword)
		ab := &authBuf{
			Username: *username,
			Password: *password,
		}
		buf, err := ab.encrypt(*saveAuthBufPassword)
		if err != nil {
			logger.Println("error occur while saving authBuffer :", err)
			return
		}
		err = ioutil.WriteFile(*saveAuthBufTo, buf, 0600)
		if err != nil {
			logger.Println("error occur while saving authBuffer :", err)
			return
		}
		logger.Println("authBuffer successfully saved")
	}
}

//func setupContext() {
//	globalContext, cancelFunc = context.WithCancel(context.TODO())
//}

func checkVorbisCommment() {
	logger.Println("check vorbisComment")
	if *vorbisComment == "" {
		logger.Fatalln("please provide a valid path to vorbisComment executable file")
	} else {
		path, err := exec.LookPath(*vorbisComment)
		if err != nil {
			logger.Fatalln("please provide a valid path to vorbisComment executable")
		}
		cmd := exec.Command(path, "-V")
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		err = cmd.Run()
		if err != nil {
			logger.Fatalln(" error occur while checking vorbisComment : ", err)
		}
		output := out.String()
		outputStrs := strings.Fields(string(output))
		if outputStrs[0] == "vorbiscomment" && outputStrs[1] == "from" && outputStrs[2] == "vorbis-tools" {
			logger.Println("using vobisComment version", outputStrs[3])
		} else {
			logger.Fatalln("error occur while checking vorbisComment :", errors.New("unknown vorbisComment version"))
		}
		vorbisPath = path
	}
}

func setRealQuality() {
	switch *quality {
	case 96:
		realQuality = Spotify.AudioFile_OGG_VORBIS_96
		break
	case 160:
		realQuality = Spotify.AudioFile_OGG_VORBIS_160
		break
	case 320:
		realQuality = Spotify.AudioFile_OGG_VORBIS_320
		break
	default:
		logger.Fatalln("please provide a valid quality value (96/160/320)")
	}
}

func setNumberOfRealSearchResult() {
	if *limitOfSearchResult < 0 {
		logger.Fatalln("limit of search result should be positive number")
	}
	if *limitOfSearchResult > 30 {
		*limitOfSearchResult = 30
	}
}

func getTracksURIsToDownload() []string {
	return strings.Split(*trackURIsToDownload, ",")
}

func getURLsToDownload() []string {
	return strings.Split(*URLsToDownload, ",")
}

func getAlbumsURIsToDownload() []string {
	return strings.Split(*albumURIsToDownload, ",")
}

func getTracksURIsInsidePlaylistToDownload() []string {
	return strings.Split(*playlistURIsToDownload, ",")
}

func getSearchTargets() []string {
	return strings.Split(*targetsToSearch, ",")
}

func getArtistsURIsToView() []string {
	return strings.Split(*artistURIsToView, ",")
}

func getAlbumsURIsToView() []string {
	return strings.Split(*albumURIsToView, ",")
}

func getTracksURIsToView() []string {
	return strings.Split(*trackURIsToView, ",")
}

func getPlaylistURIsToView() []string {
	return strings.Split(*playlistURIsToView, ",")
}

func getURLsToView() []string {
	return strings.Split(*URLsToView, ",")
}

func loadDownloadJobs() {
	for _, track := range getTracksURIsToDownload() {
		downloadErrorHandler(downloadTrack(track))
	}
}

func search() {
	setNumberOfRealSearchResult()
	logger.SetPrefix("search 		")
	logger.SetFlags(0)
	defer logger.SetPrefix(*deviceName + " ")
	defer logger.SetFlags(log.LstdFlags)
	targets := getSearchTargets()
	logger.Println()
	logger.Println("begin to do search jobs")
	logger.Println()
	for _, v := range targets {
		resp, err := session.Mercury().Search(v, *limitOfSearchResult, session.Country(), session.Username())
		if err != nil {
			logger.Println("error occur while searching for target", v, " :", err)
			continue
		}
		if resp == nil {
			logger.Println("error occur while searching for target", v, " :", errors.New("unknown error"))
			continue
		}
		res := resp.Results
		if res.Error != nil {
			logger.Println("error occur while searching for target", v, " :", err)
			continue
		}
		logger.Println("search result of", v, " :")

		logger.Printf("Tracks : %d (total %d)\n", len(res.Tracks.Hits), res.Tracks.Total)
		for _, track := range res.Tracks.Hits {
			logger.Printf("%s by %s =>  (%s)\n", track.Name, func() string {
				var res string
				for i, v := range track.Artists {
					if i > 0 {
						res += ","
					}
					res += v.Name
				}
				return res
			}(), track.Uri)
		}
		logger.Println()
		logger.Printf("Albums : %d (total %d)\n", len(res.Albums.Hits), res.Albums.Total)
		for _, album := range res.Albums.Hits {
			logger.Printf("%s by %s =>  (%s)\n", album.Name, func() string {
				var res string
				for i, v := range album.Artists {
					if i > 0 {
						res += ","
					}
					res += v.Name
				}
				return res
			}(), album.Uri)
		}
		logger.Println()
		logger.Printf("Artists : %d (total %d)\n", len(res.Artists.Hits), res.Artists.Total)
		for _, artist := range res.Artists.Hits {
			logger.Printf("%s =>  (%s)\n", artist.Name, artist.Uri)
		}
		logger.Println()
	}
	logger.Println("all search jobs done ")
}

func viewArtists() {
	logger.SetPrefix("view artists 		")
	logger.SetFlags(0)
	defer logger.SetPrefix(*deviceName + " ")
	defer logger.SetFlags(log.LstdFlags)
	targets := getArtistsURIsToView()
	logger.Println()
	logger.Println("begin to view artists")
	for _, v := range targets {
		logger.Println()
		logger.Println("begin to view artist", v)
		artist, err := session.Mercury().GetArtist(utils.Base62ToHex(v))
		if err != nil {
			logger.Println("error occur while getting artist information of id", v, " :", err)
			continue
		}
		if artist == nil {
			logger.Println("error occur while getting artist information of id", v, " :", errors.New("unknown error"))
			continue
		}
		logger.Printf("Artist : %s\n", artist.GetName())
		logger.Printf("Popularity : %d\n", artist.GetPopularity())
		logger.Printf("Genre : %s\n", artist.GetGenre())

		if artist.GetTopTrack() != nil && len(artist.GetTopTrack()) > 0 {
			tt := artist.GetTopTrack()[0]
			logger.Printf("Top tracks (country %s) : \n", tt.GetCountry())
			for _, t := range tt.GetTrack() {
				trackid := utils.ConvertTo62(t.GetGid())
				track, err := session.Mercury().GetTrack(utils.Base62ToHex(trackid))
				if err != nil {
					logger.Println("error occur while getting track information of id", v, " :", err)
					continue
				}
				if track == nil {
					logger.Println("error occur while getting track information of id", v, " :", errors.New("unknown error"))
					continue
				}
				var res string
				res += track.GetName()
				res += " => "
				res += trackid
				logger.Println(res)
			}
		}
		logger.Println()
	}
	logger.Println("end of viewing artists")
}

func viewAlbums() {
	logger.SetPrefix("view albums 		")
	logger.SetFlags(0)
	defer logger.SetPrefix(*deviceName + " ")
	defer logger.SetFlags(log.LstdFlags)
	targets := getAlbumsURIsToView()
	logger.Println()
	logger.Println("begin to view albums")
	for _, v := range targets {
		logger.Println()
		logger.Println("begin to view album", v)
		album, err := session.Mercury().GetAlbum(utils.Base62ToHex(v))
		if err != nil {
			logger.Println("error occur while getting album information of id", v, " :", err)
			continue
		}
		if album == nil {
			logger.Println("error occur while getting album information of id", v, " :", errors.New("unknown error"))
			continue
		}

		logger.Printf("Album : %s\n", album.GetName())
		logger.Printf("Popularity : %d\n", album.GetPopularity())
		logger.Printf("Genre : %s\n", album.GetGenre())
		logger.Printf("Date : %d-%d-%d\n", album.GetDate().GetYear(), album.GetDate().GetMonth(), album.GetDate().GetDay())
		logger.Printf("Label : %s\n", album.GetLabel())
		logger.Printf("Type : %s\n", album.GetTyp())
		logger.Printf("Artists : %s\n", func() string {
			var res string
			for i, v := range album.GetArtist() {
				if i > 0 {
					res += ","
				}
				res += v.GetName()
			}
			return res
		}())
		for _, disc := range album.GetDisc() {
			logger.Printf("\nDisc %d (%s): \n", disc.GetNumber(), disc.GetName())
			for _, track := range disc.GetTrack() {
				logger.Printf("%s => %s\n", track.GetName(), utils.ConvertTo62(track.GetGid()))
			}
		}
		logger.Println()
	}
	logger.Println("end of viewing albums")
}

func viewTracks() {
	logger.SetPrefix("view tracks 		")
	logger.SetFlags(0)
	defer logger.SetPrefix(*deviceName + " ")
	defer logger.SetFlags(log.LstdFlags)
	targets := getTracksURIsToView()
	logger.Println()
	logger.Println("begin to view tracks")

	for _, v := range targets {
		logger.Println()
		logger.Println("begin to view track", v)
		track, err := session.Mercury().GetTrack(utils.Base62ToHex(v))
		if err != nil {
			logger.Println("error occur while getting track information of id", v, " :", err)
			continue
		}
		if track == nil {
			logger.Println("error occur while getting track information of id", v, " :", errors.New("unknown error"))
			continue
		}
		logger.Println("Title :", track.GetName())
		logger.Println("Album :", track.GetAlbum().GetName())
		logger.Println("Artists :", func() string {
			var res string
			for i, artist := range track.GetAlbum().GetArtist() {
				if i != 0 {
					res += ","
				}
				res += artist.GetName()
			}
			return res
		}())
		logger.Println("External id :", func() string {
			var res string
			for i, eid := range track.GetExternalId() {
				if i != 0 {
					res += ","
				}
				res += "{ type : " + eid.GetTyp() + " , id : " + eid.GetId() + " }"
			}
			return res
		}())
		logger.Println()
	}
	logger.Println("end of viewing tracks")
}

func viewPlaylists() {
	logger.SetPrefix("view playlists 		")
	logger.SetFlags(0)
	defer logger.SetPrefix(*deviceName + " ")
	defer logger.SetFlags(log.LstdFlags)
	logger.Println()
	logger.Println("begin to view playlists")
	targets := getPlaylistURIsToView()
	for _, v := range targets {
		logger.Println()
		logger.Println("begin to view playlist", v)
		playlistURI := "user/" + session.Username() + "/playlist/" + v
		playlist, err := session.Mercury().GetPlaylist(playlistURI)
		if err != nil {
			logger.Println("error occur while getting playlist", playlistURI, " :", err)
			continue
		}
		if playlist == nil {
			logger.Println("error occur while getting playlist", playlistURI, " :", errors.New("unknown error"))
			continue
		}

		for _, v := range playlist.GetContents().GetItems() {
			tmpuri := v.GetUri()
			tmpuri = strings.TrimPrefix(tmpuri, "spotify:track:")

			track, err := session.Mercury().GetTrack(utils.Base62ToHex(tmpuri))
			if err != nil {
				logger.Println("error occur while getting track :", err)
				continue
			}
			if track == nil {
				logger.Println("error occur while getting track :", errors.New("unknown error"))
				continue
			}
			logger.Println(track.GetName(), "by", func() string {
				var res string
				for i, artist := range track.GetAlbum().GetArtist() {
					if i != 0 {
						res += ","
					}
					res += artist.GetName()
				}
				return res
			}(), "=>", tmpuri)
		}
		logger.Println()
	}
	logger.Println("end of viewing playlists")
}

func doViewRootPlayLists() {
	logger.SetPrefix("view root playlists 		")
	logger.SetFlags(0)
	defer logger.SetPrefix(*deviceName + " ")
	defer logger.SetFlags(log.LstdFlags)

	logger.Println()
	logger.Println("begin to view root playlist")
	playlist, err := session.Mercury().GetRootPlaylist(session.Username())

	if err != nil || playlist.Contents == nil {
		logger.Println("error occur while getting root list: ", err)
		return
	}

	items := playlist.Contents.Items
	for _, v := range items {
		id := strings.TrimPrefix(v.GetUri(), "spotify:")
		id = strings.Replace(id, ":", "/", -1)
		list, _ := session.Mercury().GetPlaylist(id)
		logger.Println(list.Attributes.GetName(), "=>", id)
	}
	logger.Println()
	logger.Println("end of viewing root playlist")
}

func viewWithURLs() {
	tmp := getURLsToView()
	for _, v := range tmp {
		str := strings.TrimPrefix(v, "https://open.spotify.com/")
		methodAndURI := strings.Split(str, "/")
		if len(methodAndURI) != 2 {
			logger.Println("error occur while parsing URI : invalid URI", v)
			continue
		}
		switch methodAndURI[0] {
		case "album":
			if *albumURIsToView != "" {
				*albumURIsToView += ","
			}
			*albumURIsToView += methodAndURI[1]
			break
		case "track":
			if *trackURIsToView != "" {
				*trackURIsToView += ","
			}
			*trackURIsToView += methodAndURI[1]
			break
		case "artist":
			if *artistURIsToView != "" {
				*albumURIsToView += ","
			}
			*artistURIsToView += methodAndURI[1]
			break
		case "playlist":
			if *playlistURIsToView != "" {
				*playlistURIsToView += ","
			}
			*playlistURIsToView += methodAndURI[1]
			break
		default:
			logger.Println("error occur while parsing URI : invalid URI", v)
		}
	}
}

func downloadTracks() {
	logger.SetPrefix("download tracks ")
	logger.SetFlags(log.LstdFlags)
	defer logger.SetPrefix(*deviceName + " ")
	defer logger.SetFlags(log.LstdFlags)
	setupDownloadWaitGroup()
	checkVorbisCommment()
	setRealQuality()
	//setupContext()
	loadDownloadJobs()
	//waitForDownloadJobDone()
	time.Sleep(time.Second)
	downloadWaitGroup.Wait()
}

func waitForDownloadJobDone() {
	signalChan := make(chan os.Signal)
	signal.Ignore(syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGQUIT)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGQUIT)

	jobChan := func() (tmp chan int) {
		go func() {
			downloadWaitGroup.Wait()
			tmp <- 0
		}()
		return
	}()
	<-time.After(time.Second)
	select {
	case sig := <-signalChan:
		logger.Println("receive signal :", sig)
		cancelFunc()
		logger.Println("all download job canceled")
		return
	case <-jobChan:
		logger.Println("all download job done")
		return
	}
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "TrackDownloader version : %s", version)
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of trackDownloader :\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	setupLogger()
	logger.Println("program version :", version)
	logger.Println("warning : at this time the program only support current user's playlist")
	logger.Println("source code hosted on https://github.com/chfanghr/trackDownloader")
	login()

	if *viewRootPlaylist {
		doViewRootPlayLists()
	}

	if *URLsToView != "" {
		viewWithURLs()
	}

	if *artistURIsToView != "" {
		viewArtists()
	}

	if *albumURIsToView != "" {
		viewAlbums()
	}

	if *trackURIsToView != "" {
		viewTracks()
	}

	if *playlistURIsToView != "" {
		viewPlaylists()
	}

	if *targetsToSearch != "" {
		search()
	}

	if *URLsToDownload != "" {
		downloadWithURLs()
	}

	if *albumURIsToDownload != "" {
		downloadAlbums()
	}

	if *trackURIsToDownload != "" {
		downloadTracks()
	}

	logger.Println("all jobs done,exit")
}
