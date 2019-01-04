package musicBrainz

import "fmt"

const CoverArtAPIRoot = "https://coverartarchive.org"

type ErrStatusCode struct {
	code int
}

func (e *ErrStatusCode) Error() string {
	return fmt.Sprint(e.code)
}

func NewErrStatusCode(code int) *ErrStatusCode {
	return &ErrStatusCode{code: code}
}

type ErrBadSize struct {
	badSize string
}

func (e *ErrBadSize) Error() string {
	return e.badSize
}

func NewErrBadSize(bad string) *ErrBadSize {
	return &ErrBadSize{badSize: bad}
}

const (
	CoverSize250  = "250"
	CoverSize500  = "500"
	CoverSize1200 = "1200"
)
