package movie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindDoubanMovieURLsOK(t *testing.T) {
	urls := findDoubanMovieURLs("view-source:https://movie.douban.com/subject/22265634/?from=showing")
	assert.Equal(t, []string{"https://movie.douban.com/subject/22265634"}, urls)
}

func TestFindDoubanMovieURLsMultiple(t *testing.T) {
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

func TestFindDoubanMovieURLsMissing(t *testing.T) {
	urls := findDoubanMovieURLs("view-source:https://movie.douban2.com/subject/22265634/?from=showing")
	assert.Len(t, urls, 0)
}

func TestFindIMDbIDsOK(t *testing.T) {
	ids := findIMDbIDs(`<span class="pl">IMDb链接:</span> <a href="http://www.imdb.com/title/tt2527336" target="_blank" rel="nofollow">`)
	assert.Equal(t, []string{"tt2527336"}, ids)
}

func TestFindIMDbIDsMultiple(t *testing.T) {
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

func TestFindIMDbIDsMissing(t *testing.T) {
	ids := findIMDbIDs(`<span class="pl">IMDb链接:</span> <a href="http://www.imdb.com/title/tt252733" target="_blank" rel="nofollow">`)
	assert.Len(t, ids, 0)
}

func TestSearchIMDbIDIMDbURLOK(t *testing.T) {
	ids, err := SearchIMDbID(`http://www.imdb.com/title/tt2527336`)
	assert.NoError(t, err)
	assert.Equal(t, []string{"tt2527336"}, ids)
}

func TestSearchIMDbIDMultipleURLOK(t *testing.T) {
	ids, err := SearchIMDbID(`
		http://www.imdb.com/title/tt2527336
		https://movie.douban.com/subject/22265634
	`)
	assert.NoError(t, err)
	assert.Equal(t, []string{"tt2527336", "tt2527336"}, ids)
}

func TestSearchIMDbIDDoubanURLOK(t *testing.T) {
	ids, err := SearchIMDbID(`https://movie.douban.com/subject/22265634/`)
	assert.NoError(t, err)
	assert.Equal(t, []string{"tt2527336"}, ids)
}
