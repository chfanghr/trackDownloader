package musicBrainz

import (
	"fmt"
	"net/http"
)

type LookupResult struct {
	EndArea  interface{} `json:"end_area"`
	Name     string      `json:"name"`
	GenderID interface{} `json:"gender-id"`
	SortName string      `json:"sort-name"`
	LifeSpan struct {
		End   string `json:"end"`
		Ended bool   `json:"ended"`
		Begin string `json:"begin"`
	} `json:"life-span"`
	ID   string `json:"id"`
	Area struct {
		SortName       string   `json:"sort-name"`
		ID             string   `json:"id"`
		Disambiguation string   `json:"disambiguation"`
		Iso31661Codes  []string `json:"iso-3166-1-codes"`
		Name           string   `json:"name"`
	} `json:"area"`
	TypeID         string        `json:"type-id"`
	Country        string        `json:"country"`
	Disambiguation string        `json:"disambiguation"`
	Gender         interface{}   `json:"gender"`
	Type           string        `json:"type"`
	Isnis          []string      `json:"isnis"`
	Ipis           []interface{} `json:"ipis"`
	BeginArea      struct {
		ID             string `json:"id"`
		SortName       string `json:"sort-name"`
		Name           string `json:"name"`
		Disambiguation string `json:"disambiguation"`
	} `json:"begin_area"`
}

func BuildLookupRequest(entity, mbid string, inc ...string) (req *http.Request, err error) {
	var url string
	if mbid == "" {
		err = fmt.Errorf("mbid shouldn't be empty")
		return
	}
	if !isStringInStrings(entity, "area", "artist", "event", "instrument", "label", "place", "recording", "release", "release-group", "series", "work", "url") {
		err = fmt.Errorf("invalid entity %s", entity)
		return
	}
	url = fmt.Sprintf("%s/%s", APIRoot, entity)
	if len(inc) > 0 {
		var incs string
		for _, v := range inc {
			switch entity {
			case "artist":
				if isStringInStrings(v, "recordings", "releases", "release-groups", "works") {
					break
				}
			case "collection":
				if isStringInStrings(v, "user-collections") {
					break
				}
			case "label":
				if isStringInStrings(v, "releases") {
					break
				}
			case "recording":
				if isStringInStrings(v, "artists", "releases") {
					break
				}
			case "release":
				if isStringInStrings(v, "artists", "collections", "labels", "recordings", "release-groups") {
					break
				}
			case "release-group":
				if isStringInStrings(v, "artists", "releases") {
					break
				}
			default:
				err = fmt.Errorf("unspported inc %s for entity %s", v, entity)
				return
			}
			if incs != "" {
				incs += "+"
			}
			incs += v
		}
		incs = "?inc=" + incs
		url += incs
	}
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Add("Accept", "application/json")
	return
}

func BuildNonMBIDLookups(entity, value string, inc ...string) (req *http.Request, err error) {
	var url string
	if value == "" {
		err = fmt.Errorf("value shouldn't be empty")
		return
	}
	if !isStringInStrings(entity, "isrc", "iswc") {
		err = fmt.Errorf("unsupported entity")
		return
	}
	url = APIRoot + "/" + value
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Add("Accept", "application/json")
	return
}
