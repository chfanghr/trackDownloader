package musicBrainz

import (
	"encoding/json"
	"net/http"
	"time"
)

type SearchRecordingResult struct {
	Created    time.Time `json:"created"`
	Count      int       `json:"count"`
	Offset     int       `json:"offset"`
	Recordings []struct {
		ID           string      `json:"id"`
		Score        int         `json:"score"`
		Title        string      `json:"title"`
		Length       int         `json:"length"`
		Video        interface{} `json:"video"`
		ArtistCredit []struct {
			Joinphrase string `json:"joinphrase,omitempty"`
			Artist     struct {
				ID       string `json:"id"`
				Name     string `json:"name"`
				SortName string `json:"sort-name"`
				Aliases  []struct {
					SortName  string      `json:"sort-name"`
					TypeID    string      `json:"type-id"`
					Name      string      `json:"name"`
					Locale    interface{} `json:"locale"`
					Type      string      `json:"type"`
					Primary   interface{} `json:"primary"`
					BeginDate interface{} `json:"begin-date"`
					EndDate   interface{} `json:"end-date"`
				} `json:"aliases"`
			} `json:"artist"`
		} `json:"artist-credit"`
		Releases []struct {
			ID           string `json:"id"`
			Count        int    `json:"count"`
			Title        string `json:"title"`
			Status       string `json:"status"`
			ArtistCredit []struct {
				Artist struct {
					ID             string `json:"id"`
					Name           string `json:"name"`
					SortName       string `json:"sort-name"`
					Disambiguation string `json:"disambiguation"`
				} `json:"artist"`
			} `json:"artist-credit"`
			ReleaseGroup struct {
				ID          string `json:"id"`
				TypeID      string `json:"type-id"`
				Title       string `json:"title"`
				PrimaryType string `json:"primary-type"`
			} `json:"release-group"`
			Date          string `json:"date,omitempty"`
			Country       string `json:"country,omitempty"`
			ReleaseEvents []struct {
				Date string `json:"date"`
				Area struct {
					ID            string   `json:"id"`
					Name          string   `json:"name"`
					SortName      string   `json:"sort-name"`
					Iso31661Codes []string `json:"iso-3166-1-codes"`
				} `json:"area"`
			} `json:"release-events,omitempty"`
			TrackCount int `json:"track-count"`
			Media      []struct {
				Position int    `json:"position"`
				Format   string `json:"format"`
				Track    []struct {
					ID     string `json:"id"`
					Number string `json:"number"`
					Title  string `json:"title"`
					Length int    `json:"length"`
				} `json:"track"`
				TrackCount  int `json:"track-count"`
				TrackOffset int `json:"track-offset"`
			} `json:"media"`
			Disambiguation string `json:"disambiguation,omitempty"`
		} `json:"releases"`
		Isrcs []string `json:"isrcs"`
		Tags  []struct {
			Count int    `json:"count"`
			Name  string `json:"name"`
		} `json:"tags"`
	} `json:"recordings"`
}

func SearchRecordingWithISRC(ISRC string) (ress *SearchRecordingResult, err error) {
	ress = &SearchRecordingResult{}
	url := RootURL + "/recording?query=isrc:" + ISRC
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(ress)
	if err != nil {
		return nil, err
	}
	return
}
