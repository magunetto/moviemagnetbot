package main

import (
	"bytes"
	"log"
	"os"
	"regexp"
	"time"

	rarbg "github.com/umayr/go-torrentapi"
	bot "gopkg.in/tucnak/telebot.v2"
)

const (
	helpText  = "What movies do you like? Just send IMDb links to me"
	errorText = "An error occurred, please try again"
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

	// handlers
	b.Handle("/start", func(m *bot.Message) {
		b.Send(m.Sender, helpText)
	})
	b.Handle("/help", func(m *bot.Message) {
		b.Send(m.Sender, helpText)
	})

	b.Handle(bot.OnText, func(m *bot.Message) {
		log.Printf("someone queried: %s", m.Text)

		// find all IMDb id in the message
		re, _ := regexp.Compile("tt[0-9]{7}") // e.g. tt0137523
		imdbIDs := re.FindAllString(m.Text, -1)
		if len(imdbIDs) == 0 {
			b.Send(m.Sender, helpText)
			return
		}

		// search movies by them
		for _, id := range imdbIDs {
			result := new(bytes.Buffer)
			searchIMDb(result, id, api)
			b.Send(m.Sender, result.String(), &bot.SendOptions{ParseMode: bot.ModeMarkdown})
		}
	})

	b.Start()
}
