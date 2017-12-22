package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	rarbg "github.com/umayr/go-torrentapi"
	bot "gopkg.in/tucnak/telebot.v2"
)

func main() {

	// init telebot
	b, err := bot.NewBot(bot.Settings{
		Token:  os.Getenv("MOVIE_MAGNET_BOT_TOKEN"),
		Poller: &bot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatalf("error while creating telebot: %s", err)
		return
	}

	// init rarbg
	api, err := rarbg.New()
	if err != nil {
		log.Fatalf("error while querying rarbg: %s", err)
		return
	}

	// handlers
	b.Handle("/start", func(m *bot.Message) {
		b.Send(m.Sender, "I find movies for you.\nTalk to me like this:\n/imdb tt0137523")
	})

	b.Handle("/imdb", func(m *bot.Message) {
		buf := new(bytes.Buffer)
		imdbHandler(buf, m, api)

		b.Send(m.Sender, buf.String(),
			&bot.SendOptions{ParseMode: bot.ModeMarkdown})
	})

	b.Start()
}

func imdbHandler(w io.Writer, m *bot.Message, api *rarbg.API) {

	// get keyword from message
	keyword := strings.Split(m.Text, " ")[1]
	log.Printf("someone searched: %s", keyword)

	// search it
	results, err := search(api, "imdb", keyword)
	if err != nil {
		fmt.Fprintln(w, "error while querying torrentapi: ", err)
		return
	}
	fmt.Fprintln(w, "Seeders / Leechers / Size / File Name")
	for _, r := range results {
		magnet := strings.Split(r.Download, "&")[0]
		fmt.Fprintf(w, "`%d` / `%d` / `%s` / `%s` / `%s`\n", r.Seeders, r.Leechers, humanizeSize(r.Size), r.Title, magnet)
	}
	return
}
