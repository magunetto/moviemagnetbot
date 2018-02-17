package movie

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/ryanbradynd05/go-tmdb"
)

const (
	tmdbURLFmt = "https://www.themoviedb.org/%s/%d"
)

var (
	tapi *tmdb.TMDb

	errTMDbSearchNoResult = errors.New("No movies found on TMDb, please check your input")
	errTMDbSearchError    = errors.New("An error occurred while finding movies, please try again")
)

// InitTMDb init TMDb API
func InitTMDb() {
	tapi = tmdb.Init(os.Getenv("TMDB_API_TOKEN"))
}

// Search searches movies on TMDb
func Search(keyword string, limit int) (*[]Movie, error) {

	options := make(map[string]string)
	options["page"] = "1"
	result, err := tapi.SearchMulti(keyword, options)
	if err != nil {
		log.Printf("error while querying tmdb: %s", err)
		return nil, errTMDbSearchError
	}
	if len(result.Results) == 0 {
		return nil, errTMDbSearchNoResult
	}

	return newMoviesBySearch(result, limit), nil
}

func newMoviesBySearch(result *tmdb.MultiSearchResults, limit int) *[]Movie {

	movies := []Movie{}

	for i, r := range result.GetMoviesResults() {
		if i == limit {
			break
		}
		m := New()
		m.TMDbID = r.ID
		m.Title = r.Title
		m.Date = r.ReleaseDate
		m.TMDbURL = fmt.Sprintf(tmdbURLFmt, r.MediaType, r.ID)
		movies = append(movies, m)
	}

	for i, r := range result.GetTvResults() {
		if i == limit {
			break
		}
		m := New()
		m.TMDbID = r.ID
		m.Title = r.Name
		m.Date = r.FirstAirDate
		m.TMDbURL = fmt.Sprintf(tmdbURLFmt, r.MediaType, r.ID)
		movies = append(movies, m)
	}

	return &movies
}
