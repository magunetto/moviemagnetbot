package main

import (
	"fmt"
	"io"
	"log"

	"github.com/moviemagnet/tmdb"
)

const (
	tmdbURL        = "https://www.themoviedb.org/%s/%d"
	itemsPerResult = 5
)

var (
	tapi *tmdb.TMDB
)

func searchMoviesAndTVs(w io.Writer, keyword string) (isSingleResult bool) {

	result, err := tapi.SearchMulti(keyword)
	if err != nil {
		log.Printf("error while querying tmdb: %s", err)
		fmt.Fprintln(w, replyTMDbErr)
		return false
	}
	if len(result.Results) == 0 {
		log.Printf("no movie found using this keyword: %s", keyword)
		fmt.Fprintln(w, replyNoTMDb)
		return false
	}

	for i, r := range result.Results {
		if i == itemsPerResult {
			break
		}
		title := r.Title
		date := r.ReleaseDate
		if r.MediaType == "tv" {
			title = r.Name
			date = r.FirstAirDate
		}
		command := fmt.Sprintf("%s%d", cmdPrefixTMDB, r.ID)
		url := fmt.Sprintf(tmdbURL, r.MediaType, r.ID)
		fmt.Fprintf(w, "%s (%s)\n", title, date[0:4])
		fmt.Fprintf(w, "~ %s [¶](%s)\n", command, url)
	}
	return len(result.Results) == 1
}
