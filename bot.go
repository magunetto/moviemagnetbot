package main

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/moviemagnet/moviemagnetbot/douban"

	bot "gopkg.in/tucnak/telebot.v2"
)

const (
	host = "https://moviemagnetbot.herokuapp.com"

	replyHelp       = "What movies do you like? Try me with the title, or just send the IMDb / Douban links"
	replyRarbgErr   = "We encountered an error while finding magnet links, please try again"
	replyTMDbErr    = "We encountered an error while finding movies, please try again"
	replyNoIMDbIDs  = "We encountered an error while finding IMDb IDs for you: "
	replyNoTorrents = "We have no magnet links for this movie now, please come back later"
	replyNoPubDate  = "We could not find this magnet link, please check your input"
	replyNoTMDb     = "We could not find this movie on TMDb, please check your input"
	replyNoTorrent  = "We encountered an error while finding this magnet link"
	replyFeedTips   = "Auto-download every link you requested by subscribing " + host + "/tasks/%s.xml"
	replyTaskAdded  = "Task added to your feed, it will start soon"

	cmdPrefixDown = "/dl"
	cmdPrefixTMDB = "/tmdb"

	itemsPerFeed       = 20
	feedCheckThreshold = time.Duration(24 * time.Hour)
)

func downloadHandler(b *bot.Bot, m *bot.Message) {

	// get `PubDate` from command, e.g. /dl1514983115
	pubDateString := m.Text[len(cmdPrefixDown):len(m.Text)]
	pubDate, err := strconv.Atoi(pubDateString)
	if err != nil {
		log.Printf("error while parsing timestamp: %s", err)
		b.Send(m.Sender, replyNoPubDate)
		return
	}

	// get torrent by `PubDate`
	t := &Torrent{PubDate: int64(pubDate)}
	t, err = t.getByPubDate()
	if err != nil {
		log.Printf("error while getting torrent: %s", err)
		b.Send(m.Sender, replyNoTorrent)
		return
	}
	magnet := &t.Magnet
	b.Send(m.Sender, "`"+*magnet+"`", &bot.SendOptions{ParseMode: bot.ModeMarkdown})

	// save the torrent for user
	u := &User{
		TelegramID:   m.Sender.ID,
		TelegramName: m.Sender.Username,
	}
	err = u.appendTorrent(t)
	if err != nil {
		log.Printf("error while adding torrent for user: %s", err)
		return
	}

	if u.isFeedActive() {
		b.Send(m.Sender, replyTaskAdded)
		return
	}
	b.Send(m.Sender, fmt.Sprintf(replyFeedTips, u.FeedID))
}

func tmdbHandler(b *bot.Bot, m *bot.Message) {
	tmdbID := m.Text[len(cmdPrefixTMDB):len(m.Text)]
	buffer := new(bytes.Buffer)
	fmt.Fprintf(buffer, "§ %s\n", m.Text)
	searchTorrents(buffer, "tmdb", tmdbID)
	b.Send(m.Sender, buffer.String(), &bot.SendOptions{ParseMode: bot.ModeMarkdown})
}

func searchHandler(b *bot.Bot, m *bot.Message) {
	imdbIDs, err := searchIMDbIDsFromMessage(m.Text)
	if err != nil {
		b.Send(m.Sender, replyNoIMDbIDs+err.Error())
		return
	}

	// Search movies or tvs
	if len(imdbIDs) == 0 {
		result := new(bytes.Buffer)
		fmt.Fprintf(result, "§ %s\n", m.Text)
		isSingleResult := searchMoviesAndTVs(result, m.Text)
		b.Send(m.Sender, result.String(),
			&bot.SendOptions{ParseMode: bot.ModeMarkdown, DisableWebPagePreview: !isSingleResult})
		return
	}

	// Search torrents
	for _, id := range imdbIDs {
		result := new(bytes.Buffer)
		fmt.Fprintf(result, "§ /%s\n", id)
		isSingleResult := searchTorrents(result, "imdb", id)
		b.Send(m.Sender, result.String(),
			&bot.SendOptions{ParseMode: bot.ModeMarkdown, DisableWebPagePreview: !isSingleResult})
	}
}

func searchIMDbIDsFromMessage(text string) ([]string, error) {
	imdbIDs := []string{}
	// Douban
	movieLinks := findDoubanMovieURLs(text)
	for _, url := range movieLinks {
		movie := douban.NewMovie()
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
	reIMDbID         = regexp.MustCompile(`tt[0-9]{7}`) // e.g. tt0137523
)

func findDoubanMovieURLs(s string) []string {
	return reDoubanMovieURL.FindAllString(s, -1)
}

func findIMDbIDs(s string) []string {
	return reIMDbID.FindAllString(s, -1)
}
