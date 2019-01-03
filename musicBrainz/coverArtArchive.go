package musicBrainz

import (
	"encoding/json"
	"net/http"
)

const CoverArtArchiveRootURL = "http://coverartarchive.org"

type GetCoversWithMBIDResult struct {
	Images []struct {
		Approved   bool   `json:"approved"`
		Back       bool   `json:"back"`
		Comment    string `json:"comment"`
		Edit       int64  `json:"edit"`
		Front      bool   `json:"front"`
		Id         int64  `json:"id"`
		Image      string `json:"image"`
		Thumbnails struct {
			Q1200 string `json:"1200"`
			Q250  string `json:"250"`
			Q500  string `json:"500"`
			Large string `json:"large"`
			Small string `json:"small"`
		} `json:"thumbnails"`
		Types []string `json:"types"`
	}
	Release string `json:"release"`
}

func GetCoversReleaseWithMBID(mbid string) (ress *GetCoversWithMBIDResult, err error) {
	ress = &GetCoversWithMBIDResult{}
	url := CoverArtArchiveRootURL + "/release/" + mbid
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept", "application/json")
	res, err := http.DefaultClient.Do(request)
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
