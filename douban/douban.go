package douban

import (
	"errors"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
)

// ErrIMDbURLMissing is error of can not find IMDb URL
var ErrIMDbURLMissing = errors.New("Can not find IMDb links on the page")

// Movie object
type Movie struct {
	imdbID string
}

// NewMovie returns a new Movie
func NewMovie() Movie {
	return Movie{}
}

// FetchFromURL fetches HTML from a URL
func (m *Movie) FetchFromURL(url string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	pageHTML, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return err
	}
	return m.ParseHTML(pageHTML)
}

var reIMDbURL = regexp.MustCompile(`http(s)?:\/\/www\.imdb\.com\/title\/tt\d{7}`)

// ParseHTML searches IMDb URL in a HTML
func (m *Movie) ParseHTML(html []byte) error {
	url := reIMDbURL.Find(html)
	if url == nil {
		return ErrIMDbURLMissing
	}
	m.imdbID = path.Base(string(url))
	return nil
}

// IMDbID returns IMDb ID of a Movie
func (m Movie) IMDbID() string {
	return m.imdbID
}
