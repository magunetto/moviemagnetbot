package uri

import (
	"net/url"
	"path"
	"strings"
)

const (
	eD2kPrefix = "ed2k://"
)

// NewURI creates a URI from given string.
func NewURI(s string) *URI {
	scheme := SchemeUnknown
	u, err := url.Parse(s)
	if err != nil {
		if strings.HasPrefix(strings.ToLower(s), eD2kPrefix) {
			scheme = SchemeED2k
		}
	} else {
		scheme = u.Scheme
	}
	return &URI{url: u, Scheme: scheme}
}

// URI contains higher level information from a url.URL.
type URI struct {
	url    *url.URL
	Scheme string
}

// IsValid returns true if the URI has a known scheme.
func (u *URI) IsValid() bool {
	return u.Scheme != ""
}

// DisplayName returns a human-friendly name of the object that the URI is pointing to.
// NOTE: If there's none, the full URI will be returned.
func (u *URI) DisplayName() string {
	name := ""
	switch u.Scheme {
	case SchemeMagnet:
		name = u.url.Query().Get("dn")
	case SchemeFTP:
		name = path.Base(u.url.Path)
		if name == "/" || name == "." {
			name = ""
		}
	}

	if name == "" {
		name = u.url.String()
	}
	return name
}
