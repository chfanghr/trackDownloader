package musicBrainz

import (
	"encoding/json"
	"net/http"
)

type LookupISRCResult struct {
	ISRC       string `json:"isrc"`
	Recordings []struct {
		Length         int64  `json:"length"`
		Title          string `json:"title"`
		Disambiguation string `json:"disambiguation"`
		ID             string `json:"id"`
		Video          bool   `json:"video"`
	} `json:"recordings"`
}

func LookupWithISRC(ISRC string) (ress *LookupISRCResult, err error) {
	ress = &LookupISRCResult{}
	url := RootURL + "/isrc/" + ISRC
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
