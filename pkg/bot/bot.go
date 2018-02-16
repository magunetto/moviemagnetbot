package bot

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/tucnak/telebot.v2"

	"github.com/magunetto/moviemagnetbot/pkg/movie"
	"github.com/magunetto/moviemagnetbot/pkg/torrent"
	"github.com/magunetto/moviemagnetbot/pkg/user"
)

const (
	replyHelp       = "What movies do you like? Try me with the title, or just send the IMDb / Douban links"
	replyNoIMDbIDs  = "We encountered an error while finding IMDb IDs for you: "
	replyNoPubStamp = "We could not find this magnet link, please check your input"
	replyNoTorrent  = "We encountered an error while finding this magnet link"
	replyFeedTips   = "Auto-download every link you requested by subscribing your RSS feed: `%s`"
	replyTaskAdded  = "Task added to your feed, it will start soon"

	cmdPrefixDown = "/dl"
	cmdPrefixTMDb = "/tmdb"

	itemsPerMovieSearch = 5
)

// Run init bot, register handlers, and start the bot
func Run() {

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
			tmdbTorrentSearchHandler(b, m)
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
	t := &torrent.Torrent{PubStamp: int64(pubStamp)}
	t, err = t.GetByPubStamp()
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
	u := &user.User{
		TelegramID:   m.Sender.ID,
		TelegramName: getUserName(m.Sender),
	}
	err = u.AppendTorrent(t)
	if err != nil {
		log.Printf("error while adding torrent for user: %s", err)
		return
	}

	if u.IsFeedActive() {
		_, err = b.Send(m.Sender, replyTaskAdded)
		if err != nil {
			log.Printf("error while sending message: %s", err)
		}
		return
	}
	_, err = b.Send(m.Sender, fmt.Sprintf(replyFeedTips, u.FeedURL()),
		&telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
	if err != nil {
		log.Printf("error while sending message: %s", err)
	}
}

func searchHandler(b *telebot.Bot, m *telebot.Message) {
	imdbIDs, err := movie.SearchIMDbID(m.Text)
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
		imdbTorrentSearchHandler(b, m, id)
	}
}

func getUserName(user *telebot.User) string {
	userName := user.Username
	if userName == "" {
		userName = strings.TrimSpace(fmt.Sprintf("%s %s", user.FirstName, user.LastName))
	}
	return userName
}
