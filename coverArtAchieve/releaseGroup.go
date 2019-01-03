package coverArtAchieve

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type ReleaseGroup struct {
	Release string `json:"release"`
	Images  []struct {
		Edit       int    `json:"edit"`
		ID         string `json:"id"`
		Image      string `json:"image"`
		Thumbnails struct {
			Num250  string `json:"250"`
			Num500  string `json:"500"`
			Num1200 string `json:"1200"`
			Small   string `json:"small"`
			Large   string `json:"large"`
		} `json:"thumbnails"`
		Comment  string   `json:"comment"`
		Approved bool     `json:"approved"`
		Front    bool     `json:"front"`
		Types    []string `json:"types"`
		Back     bool     `json:"back"`
	} `json:"images"`
}

func GetReleaseGroupCoverArt(mbid string) (rel ReleaseGroup, err error) {
	cli := &http.Client{}
	url := RootURL + "/release-group/" + mbid

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

func GetReleaseGroupFrontCoverArt(mbid string) (buf []byte, err error) {
	cli := &http.Client{}
	url := RootURL + "/release-group/" + mbid + "/front"

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

func GetReleaseGroupFrontCoverArtWithSize(mbid, size string) (buf []byte, err error) {
	if size != CoverSize250 && size != CoverSize500 && size != CoverSize1200 {
		return nil, NewErrBadSize(size)
	}

	cli := &http.Client{}
	url := RootURL + "/release-group/" + mbid + "/front-" + size

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
