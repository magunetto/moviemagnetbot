package movie

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitTMDbEnvKeyOK(t *testing.T) {
	os.Setenv("TMDB_API_TOKEN", "test")
	InitTMDb()
	assert.Len(t, tapi.APIKey, 4)
}

func TestInitTMDbDefaultKeyOK(t *testing.T) {
	os.Setenv("TMDB_API_TOKEN", "")
	InitTMDb()
	assert.Len(t, tapi.APIKey, 32)
}

func TestSearchMovieAndTVOK(t *testing.T) {
	InitTMDb()
	movies, err := Search("黑镜", 3)
	assert.NoError(t, err)
	assert.Len(t, *movies, 3)
	assert.Equal(t, (*movies)[0].TMDbID, 42009)
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
