package main

import (
	"fmt"
	"github.com/chfanghr/librespot/Spotify"
	"github.com/go-audio/aiff"
	"github.com/go-audio/audio"
	"github.com/jfreymuth/oggvorbis"
	"io"
	"io/ioutil"
	"os"
	"strings"
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
			var baseName string
			if *simpleFileName {
				baseName = *saveFileTo + "/" + track.GetName()
			} else {
				baseName = *saveFileTo + "/" + track.GetAlbum().GetName() + "-" + track.GetName()
			}
			baseName = strings.ReplaceAll(baseName, ":", "_")
			makeError := func(err error) error { return fmt.Errorf("error occur while fetching %s : %s", track.GetName(), err) }
			if *saveOgg {
				buf, err := ioutil.ReadAll(audioFile)
				if err != nil {
					return makeError(err)
				}
				err = ioutil.WriteFile(baseName+".ogg", buf, 0666)
				if err != nil {
					return makeError(err)
				}
				_, err = audioFile.Seek(0, io.SeekStart)
				if err != nil {
					return makeError(err)
				}
			}
			decodedPCM, format, err := oggvorbis.ReadAll(audioFile)
			if err != nil {
				return makeError(err)
			}
			logger.Println(track.GetName(), "fetched successfully")
			intPCMBuf := PCMI16TOInt(PCMF32ToI16(decodedPCM))
			// bitrate:=format.Bitrate
			channels := format.Channels
			sampleRate := format.SampleRate
			intBuffer := &audio.IntBuffer{
				Data: intPCMBuf,
				Format: &audio.Format{
					NumChannels: channels,
					SampleRate:  sampleRate,
				},
				SourceBitDepth: 16,
			}
		try:
			outAiff, err := os.Create(baseName + ".aiff")
			if err != nil {
				if os.IsExist(err) {
					baseName += RandStringRunes(3) + ".aiff"
					goto try
				} else {
					return makeError(err)
				}
			}
			aiffEncoder := aiff.NewEncoder(outAiff, sampleRate, 16, channels)
			err = aiffEncoder.Write(intBuffer)
			if err != nil {
				return makeError(err)
			}
			err = aiffEncoder.Close()
			if err != nil {
				return makeError(err)
			}
			//TODO write metadata to saved file
			logger.Println(track.GetName(), "downloaded successfully")
			return nil
		}
	}
}
