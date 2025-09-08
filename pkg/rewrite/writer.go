package rewrite

import (
	"net/url"
)

// Create a interface that can be reused to re-write urls in other pkgs
type AdressWritter interface {
	ReWriteURL() (*url.URL, error)
}

// Gotta bring a logger here using a singleton
// Check which design pattern is this in here -> https://www.youtube.com/watch?v=BJatgOiiht4 factory?
// btw don't think that scheme will be s3 or gs, from log details, instead, it should be bucket.s3.amazonaws.com/path so the schema will always be https and case should change into if else
// Gotta experiment upload using blob, but probably it will be https://container.microsoft.net/blob
func ReWriteURL() (*url.URL, error) {}
