package coverArtAchieve

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Release struct {
	Images []struct {
		Types      []string `json:"types"`
		Front      bool     `json:"front"`
		Back       bool     `json:"back"`
		Edit       int      `json:"edit"`
		Image      string   `json:"image"`
		Comment    string   `json:"comment"`
		Approved   bool     `json:"approved"`
		ID         string   `json:"id"`
		Thumbnails struct {
			Num250  string `json:"250"`
			Num500  string `json:"500"`
			Num1200 string `json:"1200"`
			Small   string `json:"small"`
			Large   string `json:"large"`
		} `json:"thumbnails"`
	} `json:"images"`
	Release string `json:"release"`
}

func GetReleaseCoverArt(mbid string) (rel Release, err error) {
	cli := &http.Client{}
	url := RootURL + "/release/" + mbid

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	res, err := cli.Do(req)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		return rel, NewErrStatusCode(res.StatusCode)
	}
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&rel)
	return
}

func GetReleaseFrontCoverArt(mbid string) (buf []byte, err error) {
	cli := &http.Client{}
	url := RootURL + "/release/" + mbid + "/front"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	res, err := cli.Do(req)
	if err != nil {
		return
	}
	if res.StatusCode != http.StatusOK {
		return nil, NewErrStatusCode(res.StatusCode)
	}

	buf, err = ioutil.ReadAll(res.Body)
	return
}

func GetReleaseBackCoverArt(mbid string) (buf []byte, err error) {
	cli := &http.Client{}
	url := RootURL + "/release/" + mbid + "/back"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	res, err := cli.Do(req)
	if err != nil {
		return
	}
	if res.StatusCode != http.StatusOK {
		return nil, NewErrStatusCode(res.StatusCode)
	}

	buf, err = ioutil.ReadAll(res.Body)
	return
}
