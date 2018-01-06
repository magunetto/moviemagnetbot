package main

import (
	"bytes"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	rarbg "github.com/idealhack/go-torrentapi"
	bot "gopkg.in/tucnak/telebot.v2"
)

const (
	host           = "https://moviemagnetbot.herokuapp.com"
	helpText       = "What movies do you like? Just send IMDb links to me"
	noIMDbText     = "No IMDb info found, please send correct IMDb links"
	noTorrentsText = "We have no torrents for this movie now, please come back later"
	errorText      = "An error occurred, please try again"
	taskAddedText  = "You can use this magnet link as you want, but we also have a RSS feed for you: /feed"
	feedText       = "Your RSS feed URL is " + host + "/tasks/%d.xml, you can subscribe it in your download tools"
	dlPrefix       = "/dl"
)

func handleDownload(b *bot.Bot, m *bot.Message) {

	// get `PubDate` from command, e.g. /dl1514983115
	commands := strings.Split(m.Text, dlPrefix)
	if len(commands) < 2 {
		b.Send(m.Sender, noIMDbText)
		return
	}
	pubDate, err := strconv.Atoi(commands[1])
	if err != nil {
		log.Printf("error while parsing timestamp: %s", err)
		b.Send(m.Sender, noIMDbText)
		return
	}

	// get task by `PubDate`
	task := &Task{PubDate: int64(pubDate)}
	task, err = task.getByPubDate()
	if err != nil {
		log.Printf("error while getting task: %s", err)
		b.Send(m.Sender, noIMDbText)
		return
	}
	magnet := &task.Magnet
	b.Send(m.Sender, "`"+*magnet+"`", &bot.SendOptions{ParseMode: bot.ModeMarkdown})

	// save the task for user
	user := &User{TelegramID: m.Sender.ID}
	task.Created = time.Now()
	user.appendTask(task)
	b.Send(m.Sender, taskAddedText)
}

func handleSearch(b *bot.Bot, m *bot.Message, api *rarbg.API) {
	re, _ := regexp.Compile("tt[0-9]{7}") // e.g. tt0137523
	imdbIDs := re.FindAllString(m.Text, -1)
	if len(imdbIDs) == 0 {
		b.Send(m.Sender, helpText)
		return
	}
	for _, id := range imdbIDs {
		result := new(bytes.Buffer)
		searchIMDb(result, id, api)
		b.Send(m.Sender, result.String(), &bot.SendOptions{ParseMode: bot.ModeMarkdown})
	}
}
