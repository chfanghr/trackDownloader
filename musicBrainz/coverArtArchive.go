package musicBrainz

import (
	"encoding/json"
	"net/http"
)

const CoverArtArchiveRootURL = "http://coverartarchive.org"

type GetCoverReleaseWithMBIDResult struct {
	Images []struct {
		Approved   bool   `json:"approved"`
		Back       bool   `json:"back"`
		Comment    string `json:"comment"`
		Edit       int    `json:"edit"`
		Front      bool   `json:"front"`
		ID         int64  `json:"id"`
		Image      string `json:"image"`
		Thumbnails struct {
			Num250  string `json:"250"`
			Num500  string `json:"500"`
			Num1200 string `json:"1200"`
			Large   string `json:"large"`
			Small   string `json:"small"`
		} `json:"thumbnails"`
		Types []string `json:"types"`
	} `json:"images"`
	Release string `json:"release"`
}

func GetCoverReleaseWithMBID(mbid string) (ress *GetCoverReleaseWithMBIDResult, err error) {
	ress = &GetCoverReleaseWithMBIDResult{}
	url := CoverArtArchiveRootURL + "/release/" + mbid
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept", "application/json")
	client := &http.Client{}
	res, err := client.Do(request)
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
