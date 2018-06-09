package bot

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"

	telebot "gopkg.in/tucnak/telebot.v2"

	"github.com/magunetto/moviemagnetbot/pkg/torrent"
)

func imdbTorrentSearchHandler(b *telebot.Bot, m *telebot.Message, id string) {
	buf := new(bytes.Buffer)
	renderTorrentResult(buf, "imdb "+id)
	torrentResultHandler(buf, b, m)
}

func tmdbTorrentSearchHandler(b *telebot.Bot, m *telebot.Message) {
	id := strings.Split(m.Text, "@")[0]
	id = strings.TrimPrefix(id, cmdPrefixTMDb)
	buf := new(bytes.Buffer)
	renderTorrentResult(buf, "tmdb "+id)
	torrentResultHandler(buf, b, m)
}

func linkAddHandler(b *telebot.Bot, m *telebot.Message) {
	t, err := torrent.SaveTorrentFromLink(m.Text)
	if err != nil {
		log.Printf("error while adding magnet link: %s", err)
		_, err = b.Send(m.Chat, replyNotAdded)
		if err != nil {
			log.Printf("error while sending message: %s", err)
		}
		return
	}
	saveTorrentForUser(b, m, t)
}

func torrentResultHandler(s fmt.Stringer, b *telebot.Bot, m *telebot.Message) {
	_, err := b.Send(m.Chat, s.String(), &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
	if err != nil {
		log.Printf("error while sending message: %s", err)
	}
}

func renderTorrentResult(w io.Writer, keyword string) {
	fmt.Fprintf(w, "§ %s\n", keyword)
	torrents, err := torrent.Search(keyword, itemsInTorrentResult)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	renderTorrents(w, torrents)
}

func renderTorrents(w io.Writer, torrents *[]torrent.Torrent) {
	for _, t := range *torrents {
		command := fmt.Sprintf("%s%d", cmdPrefixDown, t.PubStamp)
		fmt.Fprintf(w, "%s\n", t.Title)
		fmt.Fprintf(w, "▸ *%d*↑ *%d*↓ `%s` %s [¶](%s)\n",
			t.Seeders, t.Leechers, t.HumanizeSize(), command, t.InfoPage)
	}
}
