package movie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	movieHTMLPart = `
<span class="pl">类型:</span> <span property="v:genre">动作</span> / <span property="v:genre">科幻</span> / <span property="v:genre">冒险</span><br/>
<span class="pl">官方网站:</span> <a href="http://starwars.com" rel="nofollow" target="_blank">starwars.com</a><br/>
<span class="pl">制片国家/地区:</span> 美国<br/>
<span class="pl">语言:</span> 英语<br/>
<span class="pl">上映日期:</span> <span property="v:initialReleaseDate" content="2018-01-05(中国大陆)">2018-01-05(中国大陆)</span> / <span property="v:initialReleaseDate" content="2017-12-15(美国)">2017-12-15(美国)</span><br/>
<span class="pl">片长:</span> <span property="v:runtime" content="152">152分钟</span><br/>
<span class="pl">又名:</span> 星球大战：最后绝地武士(港) / 星球大战8 / 星战8 / Star Wars: Episode VIII<br/>
<span class="pl">IMDb链接:</span> <a href="http://www.imdb.com/title/tt2527336" target="_blank" rel="nofollow">tt2527336</a><br>

IMDb link in Comments: "http://www.imdb.com/title/tt2527337" target="_blank" rel="nofollow">tt2527337</a><br>
`

	movieHTMLMissingIMDbURL = `
<span class="pl">类型:</span> <span property="v:genre">动作</span> / <span property="v:genre">科幻</span> / <span property="v:genre">冒险</span><br/>
<span class="pl">官方网站:</span> <a href="http://starwars.com" rel="nofollow" target="_blank">starwars.com</a><br/>
<span class="pl">制片国家/地区:</span> 美国<br/>
<span class="pl">语言:</span> 英语<br/>
<span class="pl">上映日期:</span> <span property="v:initialReleaseDate" content="2018-01-05(中国大陆)">2018-01-05(中国大陆)</span> / <span property="v:initialReleaseDate" content="2017-12-15(美国)">2017-12-15(美国)</span><br/>
<span class="pl">片长:</span> <span property="v:runtime" content="152">152分钟</span><br/>
<span class="pl">又名:</span> 星球大战：最后绝地武士(港) / 星球大战8 / 星战8 / Star Wars: Episode VIII<br/>
`

	movieURL = "https://movie.douban.com/subject/1293181/"
)

func TestMovieFetchFromURLOK(t *testing.T) {
	m := New()
	err := m.FetchFromURL(movieURL)
	assert.NoError(t, err)
	assert.Equal(t, "tt0054215", m.IMDbID())
}

func TestMovieParseHTMLOK(t *testing.T) {
	m := New()
	err := m.parseHTML([]byte(movieHTMLPart))
	assert.NoError(t, err)
	assert.Equal(t, "tt2527336", m.IMDbID())
}

func TestMovieParseHTMLMissingIMDbURL(t *testing.T) {
	m := New()
	err := m.parseHTML([]byte(movieHTMLMissingIMDbURL))
	assert.Equal(t, err, ErrIMDbURLMissing)
	assert.Equal(t, "", m.IMDbID())
}
