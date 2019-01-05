// visit https://beta.musicbrainz.org/doc/Development/XML_Web_Service/ for details

package musicBrainz

import (
	"fmt"
	"io"
	"net/http"
)

const (
	APIRoot = "http://musicbrainz.org/ws/2"
)

func isStringInStrings(a string, bs ...string) bool {
	for _, v := range bs {
		if a == v {
			return true
		}
	}
	return false
}

func DoReq(client *http.Client, req *http.Request, handler func(io.Reader) error) error {
	if req == nil {
		return fmt.Errorf("invalid request")
	}

	if client == nil {
		client = &http.Client{}
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return handler(res.Body)
}

func AsyncDoReq(client *http.Client, req *http.Request, handler func(io.Reader, error)) error {
	if req == nil {
		return fmt.Errorf("invalid request")
	}

	if client == nil {
		client = &http.Client{}
	}

	go func() {
		res, err := client.Do(req)
		if res != nil && res.Body != nil {
			defer res.Body.Close()
			handler(res.Body, err)
		} else {
			handler(nil, err)
		}
	}()

	return nil
}
