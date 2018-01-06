package douban

import (
	"errors"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
)

var ErrIMDbUrlMissing = errors.New("Can not find IMDb links on the page")

type Movie struct {
	imdbID string
}

func NewMovie() Movie {
	return Movie{}
}

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

func (m *Movie) ParseHTML(html []byte) error {
	url := reIMDbURL.Find(html)
	if url == nil {
		return ErrIMDbUrlMissing
	}
	m.imdbID = path.Base(string(url))
	return nil
}

func (m Movie) IMDbID() string {
	return m.imdbID
}
