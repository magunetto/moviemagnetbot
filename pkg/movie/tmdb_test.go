package movie

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitTMDbOK(t *testing.T) {
	os.Setenv("TMDB_API_TOKEN", "523587afe262c34af9ee7794c5f8de81")
	InitTMDb()
}

func TestSearchMovieAndTVOK(t *testing.T) {
	InitTMDb()
	movies, err := Search("La La Land", 1)
	assert.NoError(t, err)
	assert.Len(t, *movies, 2)
	assert.Equal(t, (*movies)[0].TMDbID, 313369)
}

func TestSearchMovieAndPersonOK(t *testing.T) {
	InitTMDb()
	movies, err := Search("Lady Bird", 2)
	assert.NoError(t, err)
	assert.Len(t, *movies, 1)
	assert.Equal(t, (*movies)[0].TMDbID, 391713)
}

func TestSearchNoResult(t *testing.T) {
	InitTMDb()
	movies, err := Search("abcdefghijklmn", 1)
	assert.Nil(t, movies)
	assert.Error(t, err, errTMDbSearchNoResult)
}
