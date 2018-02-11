package movie

import (
	"errors"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
)

// ErrIMDbURLMissing is error of can not find IMDb URL
var ErrIMDbURLMissing = errors.New("Can not find IMDb links on the page")

// FetchFromURL fetches movie IMDb from a URL
func (m *Movie) FetchFromURL(url string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	pageHTML, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = res.Body.Close()
	if err != nil {
		return err
	}
	return m.parseHTML(pageHTML)
}

var reIMDbURL = regexp.MustCompile(`http(s)?:\/\/www\.imdb\.com\/title\/tt\d{7}`)

// parseHTML searches IMDb URL in a HTML
func (m *Movie) parseHTML(html []byte) error {
	url := reIMDbURL.Find(html)
	if url == nil {
		return ErrIMDbURLMissing
	}
	m.imdbID = path.Base(string(url))
	return nil
}
