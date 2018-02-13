package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/magunetto/moviemagnetbot/movie"

	"gopkg.in/tucnak/telebot.v2"
)

const (
	userFeedTitle = "Movie Magnet Bot feed"
	userFeedURL   = "https://moviemagnetbot.herokuapp.com/tasks/%s.xml"

	replyHelp       = "What movies do you like? Try me with the title, or just send the IMDb / Douban links"
	replyRarbgErr   = "We encountered an error while finding magnet links, please try again"
	replyNoIMDbIDs  = "We encountered an error while finding IMDb IDs for you: "
	replyNoTorrents = "We have no magnet links for this movie now, please come back later"
	replyNoPubStamp = "We could not find this magnet link, please check your input"
	replyNoTorrent  = "We encountered an error while finding this magnet link"
	replyFeedTips   = "Auto-download every link you requested by subscribing your RSS feed: `%s`"
	replyTaskAdded  = "Task added to your feed, it will start soon"

	cmdPrefixDown = "/dl"
	cmdPrefixTMDb = "/tmdb"

	itemsPerMovieSearch = 5
	itemsPerFeed        = 20
	feedCheckThreshold  = 24 * time.Hour
)

// RunBot init bot, register handlers, and start the bot
func RunBot() {

	// init bot
	b, err := telebot.NewBot(telebot.Settings{
		Token:  os.Getenv("MOVIE_MAGNET_BOT_TOKEN"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatalf("error while creating telebot: %s", err)
	}

	// bot handlers
	b.Handle("/start", func(m *telebot.Message) {
		_, err := b.Send(m.Sender, replyHelp)
		if err != nil {
			log.Printf("error while sending message: %s", err)
		}
	})
	b.Handle("/help", func(m *telebot.Message) {
		_, err := b.Send(m.Sender, replyHelp)
		if err != nil {
			log.Printf("error while sending message: %s", err)
		}
	})
	b.Handle(telebot.OnText, func(m *telebot.Message) {
		log.Printf("%s: %s", getUserName(m.Sender), m.Text)

		// download requst
		if strings.HasPrefix(m.Text, cmdPrefixDown) {
			downloadHandler(b, m)
			return
		}
		// tmdb search
		if strings.HasPrefix(m.Text, cmdPrefixTMDb) {
			tmdbHandler(b, m)
			return
		}

		// search request
		searchHandler(b, m)
	})

	// bot loop
	b.Start()
}

func downloadHandler(b *telebot.Bot, m *telebot.Message) {

	// get `PubStamp` from command, e.g. /dl1514983115
	pubStampString := m.Text[len(cmdPrefixDown):len(m.Text)]
	pubStamp, err := strconv.Atoi(pubStampString)
	if err != nil {
		log.Printf("error while parsing timestamp: %s", err)
		_, err = b.Send(m.Sender, replyNoPubStamp)
		if err != nil {
			log.Printf("error while sending message: %s", err)
		}
		return
	}

	// get torrent by `PubStamp`
	t := &Torrent{PubStamp: int64(pubStamp)}
	t, err = t.getByPubStamp()
	if err != nil {
		log.Printf("error while getting torrent: %s", err)
		_, err = b.Send(m.Sender, replyNoTorrent)
		if err != nil {
			log.Printf("error while sending message: %s", err)
		}
		return
	}
	_, err = b.Send(m.Sender, fmt.Sprintf("`%s`", t.Magnet), &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
	if err != nil {
		log.Printf("error while sending message: %s", err)
	}

	// save the torrent for user
	u := &User{
		TelegramID:   m.Sender.ID,
		TelegramName: getUserName(m.Sender),
	}
	err = u.appendTorrent(t)
	if err != nil {
		log.Printf("error while adding torrent for user: %s", err)
		return
	}

	if u.isFeedActive() {
		_, err = b.Send(m.Sender, replyTaskAdded)
		if err != nil {
			log.Printf("error while sending message: %s", err)
		}
		return
	}
	_, err = b.Send(m.Sender, fmt.Sprintf(fmt.Sprintf(replyFeedTips, userFeedURL), u.FeedID),
		&telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
	if err != nil {
		log.Printf("error while sending message: %s", err)
	}
}

func tmdbHandler(b *telebot.Bot, m *telebot.Message) {
	tmdbID := m.Text[len(cmdPrefixTMDb):len(m.Text)]
	buffer := new(bytes.Buffer)
	fmt.Fprintf(buffer, "§ %s\n", m.Text)
	searchTorrents(buffer, "tmdb", tmdbID)
	_, err := b.Send(m.Sender, buffer.String(), &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
	if err != nil {
		log.Printf("error while sending message: %s", err)
	}
}

func searchHandler(b *telebot.Bot, m *telebot.Message) {
	imdbIDs, err := searchIMDbIDsFromMessage(m.Text)
	if err != nil {
		_, err = b.Send(m.Sender, replyNoIMDbIDs+err.Error())
		if err != nil {
			log.Printf("error while sending message: %s", err)
		}
		return
	}

	// No IMDb found in message, take it as keyword to search movies or TVs
	if len(imdbIDs) == 0 {
		movieSearchHandler(b, m)
		return
	}

	// Found IMDbs, search torrents for each of them
	for _, id := range imdbIDs {
		torrentSearchHandler(b, m, id)
	}
}

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
	isSingleResult := renderMovieSearchResult(buf, m.Text)
	_, err := b.Send(m.Sender, buf.String(),
		&telebot.SendOptions{ParseMode: telebot.ModeMarkdown, DisableWebPagePreview: !isSingleResult})
	if err != nil {
		log.Printf("error while sending message: %s", err)
	}
}

func renderMovieSearchResult(w io.Writer, keyword string) bool {
	fmt.Fprintf(w, "§ %s\n", keyword)
	movies, err := movie.SearchMovies(keyword, itemsPerMovieSearch)
	if err != nil {
		fmt.Fprintln(w, err)
		return false
	}
	renderMovies(w, movies)
	return len(movies) == 1
}

func renderMovies(w io.Writer, movies []movie.Movie) {
	for _, m := range movies {
		command := fmt.Sprintf("%s%d", cmdPrefixTMDb, m.TMDbID)
		if m.Date != "" {
			fmt.Fprintf(w, "%s (%s)\n", m.Title, m.Date[0:4])
		} else {
			fmt.Fprintf(w, "%s\n", m.Title)
		}
		fmt.Fprintf(w, "▸ %s [¶](%s)\n", command, m.TMDbURL)
	}
}

func torrentSearchHandler(b *telebot.Bot, m *telebot.Message, id string) {
	result := new(bytes.Buffer)
	fmt.Fprintf(result, "§ /%s\n", id)
	isSingleResult := searchTorrents(result, "imdb", id)
	_, err := b.Send(m.Sender, result.String(),
		&telebot.SendOptions{ParseMode: telebot.ModeMarkdown, DisableWebPagePreview: !isSingleResult})
	if err != nil {
		log.Printf("error while sending message: %s", err)
	}
}

func getUserName(user *telebot.User) string {
	userName := user.Username
	if userName == "" {
		userName = strings.TrimSpace(fmt.Sprintf("%s %s", user.FirstName, user.LastName))
	}
	return userName
}
