package movie

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/magunetto/go-tmdb"
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

	movies = append(movies, *(newMoviesByTMDbResults(result.GetMoviesResults(), limit))...)
	movies = append(movies, *(newTVsByTMDbResults(result.GetTvResults(), limit))...)

	for i, r := range result.GetPersonResults() {
		if i == limit {
			break
		}
		movies = append(movies, *(newMoviesByTMDbResults(r.GetMoviesKnownFor(), limit))...)
		movies = append(movies, *(newTVsByTMDbResults(r.GetTvKnownFor(), limit))...)
	}

	return &movies
}

func newMoviesByTMDbResults(results []tmdb.MultiSearchMovieInfo, limit int) *[]Movie {
	movies := []Movie{}
	for i, r := range results {
		if i == limit {
			break
		}
		m := newMovieByTMDbResult(&r)
		movies = append(movies, *m)
	}
	return &movies
}

func newTVsByTMDbResults(results []tmdb.MultiSearchTvInfo, limit int) *[]Movie {
	movies := []Movie{}
	for i, r := range results {
		if i == limit {
			break
		}
		m := newTVByTMDbResult(&r)
		movies = append(movies, *m)
	}
	return &movies
}

func newMovieByTMDbResult(r *tmdb.MultiSearchMovieInfo) *Movie {
	m := New()
	m.TMDbID = r.ID
	m.Title = r.Title
	m.Date = r.ReleaseDate
	m.TMDbURL = fmt.Sprintf(tmdbURLFmt, r.MediaType, r.ID)
	return &m
}

func newTVByTMDbResult(r *tmdb.MultiSearchTvInfo) *Movie {
	m := New()
	m.TMDbID = r.ID
	m.Title = r.Name
	m.Date = r.FirstAirDate
	m.TMDbURL = fmt.Sprintf(tmdbURLFmt, r.MediaType, r.ID)
	return &m
}
