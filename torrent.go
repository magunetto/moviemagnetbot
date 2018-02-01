package main

import (
	"fmt"
	"io"
	"time"

	rarbg "github.com/magunetto/go-torrentapi"
)

// Torrent (magnet links)
type Torrent struct {
	ID                  int
	Title               string
	Magnet              string
	PubStamp            int64
	Downloaders         []User    `pg:",many2many:user_torrents"`
	DownloadedAt        time.Time `sql:"-"`
	rarbg.TorrentResult `sql:"-"`
}

func (t *Torrent) create() (*Torrent, error) {
	_, err := db.Model(t).
		Where("pub_stamp= ?pub_stamp").
		OnConflict("DO NOTHING").
		SelectOrInsert()
	return t, err
}

func (t *Torrent) getByPubStamp() (*Torrent, error) {
	err := db.Model(t).Where("pub_stamp = ?", t.PubStamp).Select()
	return t, err
}

func (t *Torrent) renderTorrent(w io.Writer) {
	command := fmt.Sprintf("%s%d", cmdPrefixDown, t.PubStamp)
	fmt.Fprintf(w, "%s\n", t.Title)
	fmt.Fprintf(w, "▸ *%d*↑ *%d*↓ `%s` %s [¶](%s)\n",
		t.Seeders, t.Leechers, humanizeSize(t.Size), command, t.InfoPage)
}

func humanizeSize(s uint64) string {
	size := float64(s)
	switch {
	case size < 1024:
		return fmt.Sprintf("%d", uint64(size))
	case size < 1024*1014:
		return fmt.Sprintf("%.2fK", size/1024)
	case size < 1024*1024*1024:
		return fmt.Sprintf("%.2fM", size/1024/1024)
	default:
		return fmt.Sprintf("%.2fG", size/1024/1024/1024)
	}
}
