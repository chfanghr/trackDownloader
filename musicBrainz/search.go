package musicBrainz

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var possibleSearchType = strings.Split("artist,release,release-group,recording,work,label,track,annotation,area,cdstub,event,instrument,label,place,series,tag", ",")

func BuildSearchRequest(searchType, queryMsg, limit, offset string) (req *http.Request, err error) {
	if !isStringInStrings(searchType, possibleSearchType...) {
		err = fmt.Errorf("unsupported search type %s", searchType)
	}
	url := APIRoot + "/" + searchType + "/"
	var query string
	if queryMsg != "" {
		query += queryMsg
	}
	if limit != "" {
		if query != "" {
			query += "&"
		}
		query = query + "limit=" + limit
	}
	if offset != "" {
		if query != "" {
			query += "&"
		}
		query = query + "offset=" + offset
	}
	if query != "" {
		url = url + "?query=" + query
	}
	return http.NewRequest("GET", url, nil)
}

type RecordingInformation struct {
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
}

func SearchReleaseWithIsrc(isrc, limit, offset string) (ri []RecordingInformation, err error) {
	type resultType struct {
		Created    time.Time              `json:"created"`
		Count      int                    `json:"count"`
		Offset     int                    `json:"offset"`
		Recordings []RecordingInformation `json:"recordings"`
	}

	var res resultType
	if limit == "" {
		limit = "100"
	}
	req, err := BuildSearchRequest("recording", "isrc:"+isrc, limit, offset)
	if err != nil {
		return
	}
	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	ress, err := client.Do(req)
	if err != nil {
		return
	}
	fmt.Println(req.URL)
	decoder := json.NewDecoder(ress.Body)
	defer ress.Body.Close()
	err = decoder.Decode(&res)
	if err != nil {
		return
	}
	ri = res.Recordings
	return
}

func SearchMostSimilarRelease(infos []RecordingInformation, date, title, artist, releaseGroupType, releaseGroup, ReleaseStatus, ReleaseType string) (id string, err error) {
	type releaseInfo struct {
		title                          string
		status                         string
		date                           string
		releaseGroup, releaseGroupType string
		releaseStatus, releaseType     string
		artist                         string
	}
	//TODO:
	panic("Implement me")
}
