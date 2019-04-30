package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/chfanghr/librespot"
	"github.com/chfanghr/librespot/Spotify"
	"github.com/chfanghr/librespot/core"
	"github.com/chfanghr/librespot/utils"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	logger            *log.Logger     = nil
	session           *core.Session   = nil
	downloadWaitGroup *sync.WaitGroup = nil
	//vorbisPath                        = ""
	//globalContext     context.Context    = nil
	//cancelFunc  context.CancelFunc       = nil
	version                              = "DEBUG"
	realQuality Spotify.AudioFile_Format = Spotify.AudioFile_OGG_VORBIS_320
	nullDevPath                          = os.DevNull
)

func downloadErrorHandler(err error) {
	if err != nil {
		logger.Println("error occur while downloading :", err.Error())
	}
}

func processURI(uri string) (uritype, id string, err error) {
	errInvalidURI := errors.New("invalid uri")
	ress := strings.Split(uri, ":")
	if len(ress) != 3 || ress[0] != "spotify" {
		err = errInvalidURI
		return
	}

	switch ress[1] {
	case "track", "album", "playlist":
		uritype = ress[1]
	default:
		err = errInvalidURI
		return
	}

	id = ress[2]
	return
}

func downloadWithURLs() {
	tmp := getURLsToDownload()
	for _, v := range tmp {
		uriType, id, err := processURI(v)
		if err != nil {
			logger.Println("error occur while parsing URI : invalid URI", v)
			continue
		}
		switch uriType {
		case "album":
			if *albumURIsToDownload != "" {
				*albumURIsToDownload += ","
			}
			*albumURIsToDownload += id
			break
		case "track":
			if *trackURIsToDownload != "" {
				*trackURIsToView += ","
			}
			*trackURIsToDownload += id
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
	return downloadTrackInternal(track)
}

func setupLogger() {
	nullDev, err := os.OpenFile(nullDevPath, os.O_WRONLY|os.O_APPEND, 0666)
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
		*username, *password, err = UnpackAuthBuf(buf)
		logger.Println(*username, "login with authBuffer")
		if err != nil {
			logger.Fatalln("error occur while reading authBuffer file :", err)
		}
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
		//ab := &authBuf{
		//	Username: *username,
		//	Password: *password,
		//}
		//buf, err := ab.Encrypt(*saveAuthBufPassword)
		buf, err := PackAuthBuf(*username, *password, *saveAuthBufPassword)
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

//func getTracksURIsInsidePlaylistToDownload() []string {
//	return strings.Split(*playlistURIsToDownload, ",")
//}

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

func viewWithURIs() {
	tmp := getURLsToView()
	for _, v := range tmp {
		uritype, id, err := processURI(v)
		if err != nil {
			logger.Println("error occur while parsing URI : invalid URI", v)
			continue
		}
		switch uritype {
		case "album":
			if *albumURIsToView != "" {
				*albumURIsToView += ","
			}
			*albumURIsToView += id
			break
		case "track":
			if *trackURIsToView != "" {
				*trackURIsToView += ","
			}
			*trackURIsToView += id
			break
		case "artist":
			if *artistURIsToView != "" {
				*albumURIsToView += ","
			}
			*artistURIsToView += id
			break
		case "playlist":
			if *playlistURIsToView != "" {
				*playlistURIsToView += ","
			}
			*playlistURIsToView += id
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
	setRealQuality()
	loadDownloadJobs()
	time.Sleep(time.Second)
	downloadWaitGroup.Wait()
}

func main() {
	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "TrackDownloader version : %s\n", version)
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage of trackDownloader :\n")
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
		viewWithURIs()
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
