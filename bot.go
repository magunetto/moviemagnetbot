package main

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/idealhack/moviemagnetbot/douban"

	bot "gopkg.in/tucnak/telebot.v2"
)

const (
	host = "https://moviemagnetbot.herokuapp.com"

	replyHelp       = "What movies do you like? Send IMDb or Douban links to me"
	replyRarbgErr   = "We encountered an error while finding magnet links, please try again"
	replyNoIMDbIDs  = "We encountered an error while finding IMDb IDs for you: "
	replyNoTorrents = "We have no magnet links for this movie now, please come back later"
	replyNoPubDate  = "We could not find this magnet link, please check your input"
	replyNoTorrent  = "We encountered an error while finding this magnet link"
	replyFeedTips   = "Auto-download every link you requested by subscribing " + host + "/tasks/%s.xml"
	replyTaskAdded  = "Task added to your feed, it will start soon"

	cmdPrefixDown = "/dl"

	itemsPerFeed       = 20
	feedCheckThreshold = time.Duration(24 * time.Hour)
)

func downloadHandler(b *bot.Bot, m *bot.Message) {

	// get `PubDate` from command, e.g. /dl1514983115
	commands := strings.Split(m.Text, cmdPrefixDown)
	if len(commands) < 2 {
		b.Send(m.Sender, replyNoPubDate)
		return
	}
	pubDate, err := strconv.Atoi(commands[1])
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
	u := &User{TelegramID: m.Sender.ID}
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

func searchHandler(b *bot.Bot, m *bot.Message) {
	imdbIDs, err := searchIMDbIDsFromMessage(m.Text)
	if err != nil {
		b.Send(m.Sender, replyNoIMDbIDs+err.Error())
		return
	}
	if len(imdbIDs) == 0 {
		b.Send(m.Sender, replyHelp)
		return
	}

	for _, id := range imdbIDs {
		result := new(bytes.Buffer)
		searchIMDb(result, id)
		b.Send(m.Sender, result.String(), &bot.SendOptions{ParseMode: bot.ModeMarkdown})
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
