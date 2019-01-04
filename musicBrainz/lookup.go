package musicBrainz

import (
	"fmt"
	"net/http"
)

var possibleLookupEntityAndInc = map[string][]string{
	"area":          []string{},
	"artist":        []string{"recordings", "releases", "release-groups", "works"},
	"event":         []string{},
	"instrument":    []string{},
	"label":         []string{"releases"},
	"place":         []string{},
	"recording":     []string{"artists", "releases", "recording-level-rels", "work-level-rels"},
	"release":       []string{"artists", "collections", "labels", "recordings", "release-groups"},
	"release-group": []string{"artists", "releases"},
	"series":        []string{},
	"work":          []string{},
	"url":           []string{},
}

var possibleLookupIncForAllEntity = []string{"discids", "media", "isrcs", "artist-credit", "various-artists", "aliases", "annotation", "tags", "ratings", "area-rels", "artist-rels", "event-rels", "instrument-rels", "label-rels", "place-rels", "recording-rels", "release-rels", "release-group-rels", "series-rels", "url-rels", "work-rels"}

var possibleNonMBIDLookupEntity = []string{"isrc", "iswc"}

func BuildLookupRequest(entity, mbid string, inc ...string) (req *http.Request, err error) {
	var url string
	if mbid == "" {
		err = fmt.Errorf("mbid shouldn't be empty")
		return
	}
	if _, ok := possibleLookupEntityAndInc[entity]; !ok {
		err = fmt.Errorf("invalid entity %s", entity)
		return
	}
	url = APIRoot + "/" + entity + "/" + mbid
	if len(inc) > 0 {
		var incs string
		for _, v := range inc {
			if v == "" {
				continue
			}
			if !isStringInStrings(v, possibleLookupIncForAllEntity...) {
				if isStringInStrings(v, possibleLookupEntityAndInc[entity]...) {
					err = fmt.Errorf("unspported inc %s for entity %s", v, entity)
					return
				}
			}
			if incs != "" {
				incs += "+"
			}
			incs += v
		}
		if incs != "" {
			incs = "?inc=" + incs
			url += incs
		}
	}
	return http.NewRequest("GET", url, nil)
}

func BuildNonMBIDLookupRequest(entity, id string, inc ...string) (req *http.Request, err error) {
	var url string
	if id == "" {
		err = fmt.Errorf("id shouldn't be empty")
		return
	}
	if !isStringInStrings(entity, possibleNonMBIDLookupEntity...) {
		err = fmt.Errorf("unsupported entity")
		return
	}
	url = APIRoot + "/" + id
	if len(inc) > 0 {
		var incs string
		for _, v := range inc {
			if v == "" {
				continue
			}
			if incs != "" {
				incs += "+"
			}
			incs += v
		}
		if incs != "" {
			incs = "?inc=" + incs
			url += incs
		}
	}
	return http.NewRequest("GET", url, nil)
}

type BuildLookupDiscidRequesParam struct {
	Inc []string
	Toc []string
}

func BuildLookupDiscidRequest(discid string, param BuildLookupDiscidRequesParam) (req *http.Request, err error) {
	var url = APIRoot + "/discid/"
	if discid == "" {
		if !(len(param.Inc) > 0 && len(param.Toc) > 0) {
			err = fmt.Errorf("one of discid and param should be non-empty")
			return
		}
		url += "-"
	} else {
		url += discid
	}
	var ext string
	if len(param.Inc) > 0 {
		var extInc string
		for _, v := range param.Inc {
			if v == "" {
				continue
			}
			if extInc != "" {
				extInc += "+"
			}
			extInc += v
		}
		if extInc != "" {
			ext += "inc="
			ext += extInc
		}
	}
	if len(param.Toc) > 0 {
		var extToc string
		for _, v := range param.Inc {
			if v == "" {
				continue
			}
			if extToc != "" {
				extToc += "+"
			}
			extToc += v
		}
		if extToc != "" {
			if ext != "" {
				ext += "&"
			}
			ext += "toc="
			ext += extToc
		}
	}
	if ext != "" {
		url = "?" + ext
	}
	return http.NewRequest("GET", url, nil)
}
