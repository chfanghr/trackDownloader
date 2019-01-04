package musicBrainz

import (
	"fmt"
	"net/http"
	"strings"
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
		url = url + "?" + query
	}
	return http.NewRequest("GET", url, nil)
}
