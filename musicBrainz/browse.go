package musicBrainz

import (
	"fmt"
	"net/http"
	"strings"
)

var possibleReleaseStatus = []string{"official", "promotion", "bootleg", "pseudo-release"}

var possibleReleaseType = []string{"nat", "album", "single", "ep", "compilation", "soundtrack", "spokenword", "interview", "audiobook", "live", "remix", "other"}

var possibleBrowseIncForAllEntity = []string{"annotation"}

var possibleBrowseIncForAllEntityExceptRelease = []string{"tags"}

var possibleBrowseIncForAllEntityExceptPlaceAndReleaseAndSeries = []string{"ratings"}

var possibleBrowseIncForAllEntityExceptRecordingAndReleaseAndReleaseGroup = []string{"aliases"}

var possibleBrowseEntity = map[string]struct {
	linkedEntity []string
	inc          []string
}{
	"area": {
		[]string{"collection"},
		append([]string{}, func() (res []string) {
			res = append(res, possibleBrowseIncForAllEntity...)
			res = append(res, possibleBrowseIncForAllEntityExceptPlaceAndReleaseAndSeries...)
			res = append(res, possibleBrowseIncForAllEntityExceptRelease...)
			res = append(res, possibleBrowseIncForAllEntityExceptRecordingAndReleaseAndReleaseGroup...)
			return
		}()...),
	},
	"artist": {
		strings.Split("area,collection,recording,release,release-group,work", ","),
		append([]string{}, func() (res []string) {
			res = append(res, possibleBrowseIncForAllEntity...)
			res = append(res, possibleBrowseIncForAllEntityExceptPlaceAndReleaseAndSeries...)
			res = append(res, possibleBrowseIncForAllEntityExceptRelease...)
			res = append(res, possibleBrowseIncForAllEntityExceptRecordingAndReleaseAndReleaseGroup...)
			return
		}()...),
	},
	"collection": {
		strings.Split("area,artist,editor,event,label,place,recording,release,release-group,work", ","),
		append([]string{}, func() (res []string) {
			res = append(res, possibleBrowseIncForAllEntity...)
			res = append(res, possibleBrowseIncForAllEntityExceptPlaceAndReleaseAndSeries...)
			res = append(res, possibleBrowseIncForAllEntityExceptRelease...)
			res = append(res, possibleBrowseIncForAllEntityExceptRecordingAndReleaseAndReleaseGroup...)
			return
		}()...),
	},
	"event": {
		strings.Split("area,artist,collection,place", ","),
		append([]string{}, func() (res []string) {
			res = append(res, possibleBrowseIncForAllEntity...)
			res = append(res, possibleBrowseIncForAllEntityExceptPlaceAndReleaseAndSeries...)
			res = append(res, possibleBrowseIncForAllEntityExceptRelease...)
			res = append(res, possibleBrowseIncForAllEntityExceptRecordingAndReleaseAndReleaseGroup...)
			return
		}()...),
	},
	"instrument": {
		[]string{"collection"},
		append([]string{}, func() (res []string) {
			res = append(res, possibleBrowseIncForAllEntity...)
			res = append(res, possibleBrowseIncForAllEntityExceptPlaceAndReleaseAndSeries...)
			res = append(res, possibleBrowseIncForAllEntityExceptRelease...)
			res = append(res, possibleBrowseIncForAllEntityExceptRecordingAndReleaseAndReleaseGroup...)
			return
		}()...),
	},
	"label": {
		strings.Split("area,collection,release", ","),
		append([]string{}, func() (res []string) {
			res = append(res, possibleBrowseIncForAllEntity...)
			res = append(res, possibleBrowseIncForAllEntityExceptPlaceAndReleaseAndSeries...)
			res = append(res, possibleBrowseIncForAllEntityExceptRelease...)
			res = append(res, possibleBrowseIncForAllEntityExceptRecordingAndReleaseAndReleaseGroup...)
			return
		}()...),
	},
	"recording": {
		strings.Split("artist,collection,release", ","),
		append([]string{}, func() (res []string) {
			res = append(res, possibleBrowseIncForAllEntity...)
			res = append(res, possibleBrowseIncForAllEntityExceptPlaceAndReleaseAndSeries...)
			res = append(res, possibleBrowseIncForAllEntityExceptRelease...)
			res = append(res, "artist-credit", "isrcs")
			return
		}()...),
	},
	"work": {
		strings.Split("artist,collection", ","),
		append([]string{}, func() (res []string) {
			res = append(res, possibleBrowseIncForAllEntity...)
			res = append(res, possibleBrowseIncForAllEntityExceptPlaceAndReleaseAndSeries...)
			res = append(res, possibleBrowseIncForAllEntityExceptRelease...)
			res = append(res, possibleBrowseIncForAllEntityExceptRecordingAndReleaseAndReleaseGroup...)
			return
		}()...),
	},
	"url": {
		[]string{"resource"},
		append([]string{}, func() (res []string) {
			res = append(res, possibleBrowseIncForAllEntity...)
			res = append(res, possibleBrowseIncForAllEntityExceptPlaceAndReleaseAndSeries...)
			res = append(res, possibleBrowseIncForAllEntityExceptRelease...)
			res = append(res, possibleBrowseIncForAllEntityExceptRecordingAndReleaseAndReleaseGroup...)
			return
		}()...),
	},
	"place": {
		strings.Split("area,collection", ","),
		append([]string{}, func() (res []string) {
			res = append(res, possibleBrowseIncForAllEntity...)
			res = append(res, possibleBrowseIncForAllEntityExceptRelease...)
			res = append(res, possibleBrowseIncForAllEntityExceptRecordingAndReleaseAndReleaseGroup...)
			return
		}()...),
	},
	"release": {
		strings.Split("area,artist,collection,label,track,track_artist,recording,release-group", ","),
		append([]string{}, func() (res []string) {
			res = append(res, possibleBrowseIncForAllEntity...)
			res = append(res, "artist-credit", "labels", "recordings", "release-groups", "media", "discids", "isrcs")
			return
		}()...),
	},
	"release-group": {
		strings.Split("artist,collection,release", ","),
		append([]string{}, func() (res []string) {
			res = append(res, possibleBrowseIncForAllEntity...)
			res = append(res, possibleBrowseIncForAllEntityExceptPlaceAndReleaseAndSeries...)
			res = append(res, possibleBrowseIncForAllEntityExceptRelease...)
			res = append(res, "artist-credit")
			return
		}()...),
	},
	"series": {
		[]string{"collection"},
		append([]string{}, func() (res []string) {
			res = append(res, possibleBrowseIncForAllEntity...)
			res = append(res, possibleBrowseIncForAllEntityExceptRelease...)
			return
		}()...),
	},
}

type BuildBrowseRequestParamLinkedEntityAndValue struct {
	LinkedEntity string
	Value        string
}

type BuildBrowseRequestParam struct {
	LinkedEntityAndValue []BuildBrowseRequestParamLinkedEntityAndValue
	Offset               string
	Limit                string
	Inc                  []string

	FilterByType   []string
	FilterByStatus []string
}

func BuildBrowseRequest(entity string, param BuildBrowseRequestParam) (req *http.Request, err error) {
	if _, ok := possibleBrowseEntity[entity]; !ok {
		err = fmt.Errorf("unsuportted entity %s", entity)
		return
	}
	var url string = APIRoot + "/" + entity
	var query string

	if len(param.LinkedEntityAndValue) > 0 {
		var tmp string
		for _, v := range param.LinkedEntityAndValue {
			if v.LinkedEntity == "" {
				continue
			}
			if !isStringInStrings(v.LinkedEntity, possibleBrowseEntity[entity].linkedEntity...) {
				err = fmt.Errorf("unsupported linked entity %s for entity %s", v.LinkedEntity, entity)
				return
			}
			if tmp != "" {
				tmp += "&"
			}
			tmp = tmp + v.LinkedEntity + "=" + v.Value
		}
		if tmp != "" {
			query += tmp
		}
	}
	if len(param.Inc) > 0 {
		var incs string
		for _, v := range param.Inc {
			if v == "" {
				continue
			}
			if !isStringInStrings(v, possibleBrowseEntity[entity].inc...) {
				err = fmt.Errorf("unsupported inc %s for entity %s", v, entity)
				return
			}
			if incs != "" {
				incs += "+"
			}
			incs += v
		}
		if incs != "" {
			if query != "" {
				query += "&"
			}
			query += incs
		}
	}
	if isStringInStrings(entity, "release", "release-group") {
		if len(param.FilterByType) > 0 {
			var filterTypes string
			for _, v := range param.FilterByType {
				if v == "" {
					continue
				}
				if !isStringInStrings(v, possibleReleaseType...) {
					err = fmt.Errorf("unsupported filter type %s", v)
					return
				}
				if filterTypes != "" {
					filterTypes += "|"
				}
				filterTypes += v
			}
			if filterTypes != "" {
				if query != "" {
					query += "&"
				}
				query = query + "type=" + filterTypes
			}
		}
	}
	if isStringInStrings(entity, "release") {
		if len(param.FilterByStatus) > 0 {
			var filterStatus string
			for _, v := range param.FilterByStatus {
				if v == "" {
					continue
				}
				if !isStringInStrings(v, possibleReleaseStatus...) {
					err = fmt.Errorf("unsupported filter status %s", v)
					return
				}
				if filterStatus != "" {
					filterStatus += "|"
				}
				filterStatus += v
			}
			if filterStatus != "" {
				if query != "" {
					query += "&"
				}
				query = query + "status=" + filterStatus
			}
		}
	}
	if param.Offset != "" {
		if query != "" {
			query += "&"
		}
		query = query + "offset=" + param.Offset
	}
	if param.Limit != "" {
		if query != "" {
			query += "&"
		}
		query = query + "limit=" + param.Limit
	}
	if query != "" {
		url = url + "?" + query
	}
	return http.NewRequest("GET", url, nil)
}
