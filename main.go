//                           _                                       __  __          __
//     ____ ___  ____ _   __(_)__  ____ ___  ____ _____ _____  ___  / /_/ /_  ____  / /_
//    / __ `__ \/ __ \ | / / / _ \/ __ `__ \/ __ `/ __ `/ __ \/ _ \/ __/ __ \/ __ \/ __/
//   / / / / / / /_/ / |/ / /  __/ / / / / / /_/ / /_/ / / / /  __/ /_/ /_/ / /_/ / /_
//  /_/ /_/ /_/\____/|___/_/\___/_/ /_/ /_/\__,_/\__, /_/ /_/\___/\__/_.___/\____/\__/
// https://github.com/moviemagnet/moviemagnetbot /___/
//

package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	rarbg "github.com/moviemagnet/go-torrentapi"
	"github.com/moviemagnet/tmdb"
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

	// init tmdb
	tapi = tmdb.New()
	tapi.APIKey = os.Getenv("TMDB_API_TOKEN")
	log.Printf("tmdb inited")

	// init rarbg
	rapi, err = rarbg.Init()
	if err != nil {
		log.Fatalf("error while querying rarbg: %s", err)
	}
	log.Printf("rarbg inited")

	// bot handlers
	b.Handle("/start", func(m *bot.Message) {
		b.Send(m.Sender, replyHelp)
	})
	b.Handle("/help", func(m *bot.Message) {
		b.Send(m.Sender, replyHelp)
	})
	b.Handle(bot.OnText, func(m *bot.Message) {
		log.Printf("@%s: %s", m.Sender.Username, m.Text)

		// download requst
		if strings.HasPrefix(m.Text, cmdPrefixDown) {
			downloadHandler(b, m)
			return
		}
		// tmdb search
		if strings.HasPrefix(m.Text, cmdPrefixTMDB) {
			tmdbHandler(b, m)
			return
		}

		// search request
		searchHandler(b, m)
	})
	// bot loop
	go b.Start()

	// http handlers
	r := mux.NewRouter()
	r.HandleFunc("/tasks/{user}.xml", feedHandler)
	// http loop
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}
