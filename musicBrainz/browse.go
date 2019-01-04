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
			res = append(res, possibleBrowseIncForAllEntityExceptRecordingAndReleaseAndReleaseGroup...)
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
			return
		}()...),
	},
	"release-group": {
		strings.Split("artist,collection,release", ","),
		append([]string{}, func() (res []string) {
			res = append(res, possibleBrowseIncForAllEntity...)
			res = append(res, possibleBrowseIncForAllEntityExceptPlaceAndReleaseAndSeries...)
			res = append(res, possibleBrowseIncForAllEntityExceptRelease...)
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

func BuildBrowseRequest(entity string, param struct {
	linkedEntityAndValue []struct {
		linkedEntity string
		value        string
	}
	offset string
	limit  string
	inc    []string

	filterByType   string
	filterByStatus string
}) (req *http.Request, err error) {
	if _, ok := possibleBrowseEntity[entity]; !ok {
		err = fmt.Errorf("unsuportted entity %s", entity)
		return
	}
	var url string = APIRoot + "/" + entity
	var query string
	if len(param.linkedEntityAndValue) > 0 {
		var tmp string
		for _, v := range param.linkedEntityAndValue {
			if v.linkedEntity == "" {
				continue
			}
			if !isStringInStrings(v.linkedEntity, possibleBrowseEntity[entity].linkedEntity...) {
				err = fmt.Errorf("unsupported linked entity %s for entity %s", v.linkedEntity, entity)
				return
			}
			if tmp != "" {
				tmp += "&"
			}
			tmp = tmp + v.linkedEntity + "=" + v.value
		}
		if tmp != "" {
			query += tmp
		}
	}
	if len(param.inc) > 0 {
		var incs string
		for _, v := range param.inc {
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
	if param.offset != "" {
		if query != "" {
			query += "&"
		}
		query = query + "offset=" + param.offset
	}
	if param.limit != "" {
		if query != "" {
			query += "&"
		}
		query = query + "limit=" + param.limit
	}
	if query != "" {
		url = url + "?" + query
	}
	return
}
