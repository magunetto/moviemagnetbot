package bot

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"regexp"

	telebot "gopkg.in/tucnak/telebot.v2"

	"github.com/magunetto/moviemagnetbot/pkg/movie"
)

func searchIMDbIDsFromMessage(text string) ([]string, error) {
	imdbIDs := []string{}
	// Douban
	movieLinks := findDoubanMovieURLs(text)
	for _, url := range movieLinks {
		movie := movie.New()
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
	reDoubanMovieURL = regexp.MustCompile(`http(s)?\:\/\/movie\.douban\.com\/subject\/[0-9]+`)
	reDoubanAppURL   = regexp.MustCompile(`http(s)?\:\/\/www\.douban\.com\/doubanapp\/dispatch\?uri\=\/movie\/[0-9]+`)
	reIMDbID         = regexp.MustCompile(`tt[0-9]{7}`) // e.g. tt0137523
)

func findDoubanMovieURLs(s string) (urls []string) {
	urls = append(urls, reDoubanMovieURL.FindAllString(s, -1)...)
	urls = append(urls, reDoubanAppURL.FindAllString(s, -1)...)
	return
}

func findIMDbIDs(s string) []string {
	return reIMDbID.FindAllString(s, -1)
}

func movieSearchHandler(b *telebot.Bot, m *telebot.Message) {
	buf := new(bytes.Buffer)
	isSingleResult := renderMovieResult(buf, m.Text)
	_, err := b.Send(m.Sender, buf.String(),
		&telebot.SendOptions{ParseMode: telebot.ModeMarkdown, DisableWebPagePreview: !isSingleResult})
	if err != nil {
		log.Printf("error while sending message: %s", err)
	}
}

func renderMovieResult(w io.Writer, keyword string) bool {
	fmt.Fprintf(w, "§ %s\n", keyword)
	movies, err := movie.Search(keyword, itemsPerMovieSearch)
	if err != nil {
		fmt.Fprintln(w, err)
		return false
	}
	renderMovies(w, movies)
	return len(*movies) == 1
}

func renderMovies(w io.Writer, movies *[]movie.Movie) {
	for _, m := range *movies {
		command := fmt.Sprintf("%s%d", cmdPrefixTMDb, m.TMDbID)
		if m.Date != "" {
			fmt.Fprintf(w, "%s (%s)\n", m.Title, m.Date[0:4])
		} else {
			fmt.Fprintf(w, "%s\n", m.Title)
		}
		fmt.Fprintf(w, "▸ %s [¶](%s)\n", command, m.TMDbURL)
	}
}
