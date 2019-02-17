package movie

import (
	"regexp"
)

// SearchIMDbID finds IMDb IDs in given text
func SearchIMDbID(text string) ([]string, error) {
	imdbIDs := []string{}
	// Douban
	movieLinks := findDoubanMovieURLs(text)
	for _, url := range movieLinks {
		movie := New()
		if err := movie.FetchFromURL(url); err != nil {
			return nil, err
		}
		imdbIDs = append(imdbIDs, movie.IMDbID())
	}
	// IMDb
	imdbIDs = append(imdbIDs, findIMDbIDs(text)...)
	return imdbIDs, nil
}

var (
	reDoubanMovieURL  = regexp.MustCompile(`http(s)?\:\/\/movie\.douban\.com\/subject\/[0-9]+`)
	reDoubanAppURL    = regexp.MustCompile(`http(s)?\:\/\/www\.douban\.com\/doubanapp\/dispatch\?uri\=\/movie\/[0-9]+`)
	reDoubanAppNewURL = regexp.MustCompile(`http(s)?\:\/\/www\.douban\.com\/doubanapp\/dispatch\/movie\/[0-9]+`)
	reIMDbID          = regexp.MustCompile(`tt[0-9]{7}`) // e.g. tt0137523
)

func findDoubanMovieURLs(s string) (urls []string) {
	urls = append(urls, reDoubanMovieURL.FindAllString(s, -1)...)
	urls = append(urls, reDoubanAppURL.FindAllString(s, -1)...)
	urls = append(urls, reDoubanAppNewURL.FindAllString(s, -1)...)
	return
}

func findIMDbIDs(s string) []string {
	return reIMDbID.FindAllString(s, -1)
}
