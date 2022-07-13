package movie

import (
	"errors"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
)

// userAgent is the User-Agent header for sending HTTP request
const userAgent = "moviemagnetbot/0.1"

// ErrIMDbURLMissing is error of can not find IMDb URL
var ErrIMDbURLMissing = errors.New("Can not find IMDb links on the page")

// FetchFromURL fetches movie IMDb from a URL
func (m *Movie) FetchFromURL(url string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", userAgent)
	res, err := client.Do(req)
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

var reIMDbString = regexp.MustCompile(`tt\d{8}|tt\d{7}`)

// parseHTML searches IMDb URL in a HTML
func (m *Movie) parseHTML(html []byte) error {
	url := reIMDbString.Find(html)
	if url == nil {
		return ErrIMDbURLMissing
	}
	m.imdbID = path.Base(string(url))
	return nil
}
