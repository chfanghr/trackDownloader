package main

import "flag"

var (
	authBufFile         = flag.String("authBufFile", "", "path to authBuffer file")
	username            = flag.String("username", "", "name of Spotify account")
	password            = flag.String("password", "", "Password of Spotify account")
	deviceName          = flag.String("deviceName", "trackDownloader", "name of device")
	logFile             = flag.String("logFile", "", "path to log file")
	saveFileTo          = flag.String("saveFileTo", "./", "path to save audio file")
	saveAuthBufTo       = flag.String("saveAuthBufTo", "", "path to save authBuffer")
	saveAuthBufPassword = flag.String("saveAuthBufPassword", RandStringRunes(5), "Password of saved authBuffer file")
	targetsToSearch     = flag.String("search", "", "targets to search for,split by \",\"")
	albumURIsToView     = flag.String("viewAlbum", "", "URIs of albums to view,split by \",\"")
	trackURIsToView     = flag.String("viewTrack", "", "URIs of tracks to view,split by \",\"")
	artistURIsToView    = flag.String("viewArtist", "", "URIs of artists to view,split by \",\"")
	playlistURIsToView  = flag.String("viewPlaylist", "", "URIs of playlists to view,split by \",\"")
	URLsToView          = flag.String("view", "", "URIs to view,with format spotify:xxx:xxx,split by \",\"")
	trackURIsToDownload = flag.String("downloadTrack", "", "URIs of tracks to download,split by \",\"")
	albumURIsToDownload = flag.String("downloadAlbum", "", "URIs of albums to download,split by \",\"")
	URLsToDownload      = flag.String("download", "", "URIs to download,split by \",\"")
	limitOfSearchResult = flag.Int("searchResultLimit", 12, "limit of search result to shaw")
	quality             = flag.Int("quality", 320, "quality of audio file")
	viewRootPlaylist    = flag.Bool("viewRootPlaylists", false, "view root playlist or not")
	quiet               = flag.Bool("quiet", false, "output log to stdout or not")
	saveOgg             = flag.Bool("saveOgg", false, "save orginal ogg data for higher quality")
	simpleFileName      = flag.Bool("simpleName", false, "") //TODO
)

var downloadOneByOne = func() *bool {
	tmp := true
	return &tmp
}()
