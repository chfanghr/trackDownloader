package main

import (
	"fmt"
	"github.com/chfanghr/trackDownloader/Spotify"
	"io/ioutil"
	"os"
	"os/exec"
)

type Track struct {
	ISRC string

	ID      string
	Title   string
	Artists []string
	Album   string
	Date    string
	Country string

	SpotifyURI string
	CoverURL   string

	AudioFile    string
	CoverFile    string
	MetaDataFile string
}

func downloadTrackInternal(track *Spotify.Track) error {
	downloadWaitGroup.Add(1)
	defer downloadWaitGroup.Done()
	ss := session

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