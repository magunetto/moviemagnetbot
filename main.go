package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	rarbg "github.com/idealhack/go-torrentapi"
	bot "gopkg.in/tucnak/telebot.v2"
)

func main() {

	// init db
	initModel()
	log.Printf("postgres inited")

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
	api, err := rarbg.Init()
	if err != nil {
		log.Fatalf("error while querying rarbg: %s", err)
	}
	log.Printf("rarbg inited")

	// bot handlers
	b.Handle("/start", func(m *bot.Message) {
		b.Send(m.Sender, helpText)
	})
	b.Handle("/help", func(m *bot.Message) {
		b.Send(m.Sender, helpText)
	})
	b.Handle(bot.OnText, func(m *bot.Message) {
		log.Printf("@%s: %s", m.Sender.Username, m.Text)

		// download requst
		if strings.HasPrefix(m.Text, dlPrefix) {
			handleDownload(b, m)
			return
		}

		// search request
		handleSearch(b, m, api)
	})
	// bot loop
	go b.Start()

	// http handlers
	r := mux.NewRouter()
	r.HandleFunc("/tasks/{user}.xml", feedHandler)
	// http loop
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}
