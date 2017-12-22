package main

import (
	"bytes"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	rarbg "github.com/umayr/go-torrentapi"
	bot "gopkg.in/tucnak/telebot.v2"
)

const (
	helpText       = "What movies do you like? Just send IMDb links to me"
	noIMDbText     = "No IMDb info found, please send correct IMDb links"
	noTorrentsText = "We have no torrents for this movie now, please come back later"
	errorText      = "An error occurred, please try again"
	dlPrefix       = "/dl"
)

func main() {

	// init telebot
	b, err := bot.NewBot(bot.Settings{
		Token:  os.Getenv("MOVIE_MAGNET_BOT_TOKEN"),
		Poller: &bot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatalf("error while creating telebot: %s", err)
	}
	log.Printf("telebot inited")

	// init rarbg
	api, err := rarbg.New()
	if err != nil {
		log.Fatalf("error while querying rarbg: %s", err)
	}
	log.Printf("rarbg inited")

	// magnet links cache
	magnets := make(map[int64]string)

	// handlers
	b.Handle("/start", func(m *bot.Message) {
		b.Send(m.Sender, helpText)
	})
	b.Handle("/help", func(m *bot.Message) {
		b.Send(m.Sender, helpText)
	})

	b.Handle(bot.OnText, func(m *bot.Message) {
		log.Printf("someone queried: %s", m.Text)

		// download requst
		if strings.HasPrefix(m.Text, dlPrefix) {
			handleDownload(b, m, magnets)
			return
		}

		// search request
		handleSearch(b, m, magnets, api)
	})

	b.Start()
}

func handleDownload(b *bot.Bot, m *bot.Message, magnets map[int64]string) {
	commands := strings.Split(m.Text, dlPrefix)
	if len(commands) < 2 {
		b.Send(m.Sender, noIMDbText)
		return
	}
	time, err := strconv.Atoi(commands[1])
	if err != nil {
		log.Printf("error while parsing timestamp: %s", err)
		b.Send(m.Sender, noIMDbText)
		return
	}
	magnet := magnets[int64(time)]
	if magnet == "" {
		b.Send(m.Sender, noIMDbText)
		return
	}

	b.Send(m.Sender, "`"+magnet+"`", &bot.SendOptions{ParseMode: bot.ModeMarkdown})
}

func handleSearch(b *bot.Bot, m *bot.Message, magnets map[int64]string, api *rarbg.API) {
	re, _ := regexp.Compile("tt[0-9]{7}") // e.g. tt0137523
	imdbIDs := re.FindAllString(m.Text, -1)
	if len(imdbIDs) == 0 {
		b.Send(m.Sender, helpText)
		return
	}
	for _, id := range imdbIDs {
		result := new(bytes.Buffer)
		searchIMDb(result, magnets, id, api)
		b.Send(m.Sender, result.String(), &bot.SendOptions{ParseMode: bot.ModeMarkdown})
	}
}
