package musicBrainz

import "strings"

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
