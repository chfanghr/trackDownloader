package musicBrainz

import (
	"net/http"
	"reflect"
	"testing"
)

func TestBuildSearchRequest(t *testing.T) {
	type args struct {
		searchType string
		queryMsg   string
		limit      string
		offset     string
	}
	tests := []struct {
		name    string
		args    args
		wantReq *http.Request
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotReq, err := BuildSearchRequest(tt.args.searchType, tt.args.queryMsg, tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildSearchRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotReq, tt.wantReq) {
				t.Errorf("BuildSearchRequest() = %v, want %v", gotReq, tt.wantReq)
			}
		})
	}
}

func TestSearchReleaseWithIsrc(t *testing.T) {
	type args struct {
		isrc   string
		limit  string
		offset string
	}
	tests := []struct {
		name    string
		args    args
		wantRi  []RecordingInformation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRi, err := SearchReleaseWithIsrc(tt.args.isrc, tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchReleaseWithIsrc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRi, tt.wantRi) {
				t.Errorf("SearchReleaseWithIsrc() = %v, want %v", gotRi, tt.wantRi)
			}
		})
	}
}

func TestSearchMostSimilarRelease(t *testing.T) {
	type args struct {
		infos            []RecordingInformation
		date             string
		title            string
		artist           string
		releaseGroupType string
		releaseGroup     string
		ReleaseStatus    string
		ReleaseType      string
	}
	tests := []struct {
		name    string
		args    args
		wantId  string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotId, err := SearchMostSimilarRelease(tt.args.infos, tt.args.date, tt.args.title, tt.args.artist, tt.args.releaseGroupType, tt.args.releaseGroup, tt.args.ReleaseStatus, tt.args.ReleaseType)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchMostSimilarRelease() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotId != tt.wantId {
				t.Errorf("SearchMostSimilarRelease() = %v, want %v", gotId, tt.wantId)
			}
		})
	}
}
