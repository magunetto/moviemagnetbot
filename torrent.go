package main

import (
	"time"
)

// Torrent (magnet links)
type Torrent struct {
	ID          int
	Title       string
	Magnet      string
	PubDate     int64
	Downloaders []User    `pg:",many2many:user_torrents"`
	Downloaded  time.Time `sql:"-"`
}

func (t *Torrent) create() (*Torrent, error) {
	_, err := db.Model(t).
		Where("pub_date= ?pub_date").
		OnConflict("DO NOTHING").
		SelectOrInsert()
	return t, err
}

func (t *Torrent) getByPubDate() (*Torrent, error) {
	err := db.Model(t).Where("pub_date = ?", t.PubDate).Select()
	return t, err
}
