package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestfindDoubanMovieURLsOK(t *testing.T) {
	urls := findDoubanMovieURLs("view-source:https://movie.douban.com/subject/22265634/?from=showing")
	assert.Equal(t, []string{"https://movie.douban.com/subject/22265634"}, urls)
}
func TestfindDoubanMovieURLsMultiple(t *testing.T) {
	urls := findDoubanMovieURLs(`
		https://movie.douban.com/subject/22265634/?from=showing
		https://movie.douban.com/subject/22265635/?from=showing
		https://movie.douban.com/subject/22265636/?from=showing
	`)
	assert.Len(t, urls, 3)
	assert.Equal(t, []string{
		"https://movie.douban.com/subject/22265634",
		"https://movie.douban.com/subject/22265635",
		"https://movie.douban.com/subject/22265636",
	}, urls)
}

func TestfindDoubanMovieURLsMissing(t *testing.T) {
	urls := findDoubanMovieURLs("view-source:https://movie.douban2.com/subject/22265634/?from=showing")
	assert.Len(t, urls, 0)
}

func TestfindIMDbIDsOK(t *testing.T) {
	ids := findIMDbIDs(`<span class="pl">IMDb链接:</span> <a href="http://www.imdb.com/title/tt2527336" target="_blank" rel="nofollow">`)
	assert.Equal(t, []string{"tt2527336"}, ids)
}

func TestfindIMDbIDsMultiple(t *testing.T) {
	ids := findIMDbIDs(`
		<a href="http://www.imdb.com/title/tt2527336" target="_blank" rel="nofollow">
		<a href="http://www.imdb.com/title/tt2527337" target="_blank" rel="nofollow">
		<a href="http://www.imdb.com/title/tt2527338" target="_blank" rel="nofollow">`)
	assert.Len(t, ids, 3)
	assert.Equal(t, []string{
		"tt2527336",
		"tt2527337",
		"tt2527338",
	}, ids)
}

func TestfindIMDbIDsMissing(t *testing.T) {
	ids := findIMDbIDs(`<span class="pl">IMDb链接:</span> <a href="http://www.imdb.com/title/ttt2527336" target="_blank" rel="nofollow">`)
	assert.Len(t, ids, 0)
}

func TestSearchIMDbIDsFromMessageIMDbURLOK(t *testing.T) {
	ids, err := searchIMDbIDsFromMessage(`http://www.imdb.com/title/tt2527336`)
	assert.NoError(t, err)
	assert.Equal(t, []string{"tt2527336"}, ids)
}

func TestSearchIMDbIDsFromMessageDoubanURLOK(t *testing.T) {
	ids, err := searchIMDbIDsFromMessage(`https://movie.douban.com/subject/22265634/`)
	assert.NoError(t, err)
	assert.Equal(t, []string{"tt2527336"}, ids)
}
