package main

import (
	_ "github.com/go-audio/aiff"
	_ "github.com/jfreymuth/oggvorbis"
)

type URIType int

//TODO
type URI struct{}

//TODO
func URIFromString(string) URI { panic(nil) }

//TODO
func (u URI) Type() URIType { panic(nil) }

//TODO
func (u URI) Base62ID() string { panic(nil) }
