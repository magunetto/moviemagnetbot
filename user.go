package main

import (
	"os"
	"time"

	"github.com/go-pg/pg/orm"
	"github.com/speps/go-hashids"
)

const (
	hashAlphabet = "0123456789abcdef"
)

// User (i.e. Downloader)
type User struct {
	ID            int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	TelegramID    int
	TelegramName  string
	FeedID        string
	FeedCheckedAt time.Time
	Torrents      []Torrent `pg:",many2many:user_torrents"`
}

// UserTorrent is about which user download what torrents
type UserTorrent struct {
	UserID               int
	TorrentID            int
	Torrent_DownloadedAt time.Time
}

func (u *User) create() (*User, error) {
	_, err := db.Model(u).
		Where("telegram_id= ?telegram_id").
		OnConflict("DO NOTHING").
		SelectOrInsert()
	if err != nil {
		return nil, err
	}
	return u.generateFeedID()
}

func (u *User) generateFeedID() (*User, error) {
	u, err := u.newFeedID()
	if err != nil {
		return nil, err
	}
	err = u.update()
	return u, err
}

func (u *User) newFeedID() (*User, error) {
	hd := hashids.NewData()
	hd.Salt = os.Getenv("MOVIE_MAGNET_BOT_SALT")
	hd.Alphabet = hashAlphabet
	h, err := hashids.NewWithData(hd)
	if err != nil {
		return nil, err
	}
	feed, err := h.Encode([]int{u.TelegramID})
	if err != nil {
		return nil, err
	}
	u.FeedID = feed
	return u, nil
}

func (u *User) getByTelegramID() (*User, error) {
	u, err := u.create()
	if err != nil {
		return nil, err
	}
	err = db.Model(u).Where("telegram_id = ?", u.TelegramID).Select()
	return u, err
}

func (u *User) getByFeedID() (*User, error) {
	err := db.Model(u).Where("feed_id = ?", u.FeedID).Select()
	return u, err
}

func (u *User) appendTorrent(t *Torrent) error {
	u, err := u.getByTelegramID()
	if err != nil {
		return err
	}
	return db.Insert(&UserTorrent{u.ID, t.ID, time.Now()})
}

func (u *User) getTorrents() ([]Torrent, error) {
	err := db.Model(u).Column("user.*", "Torrents").Where("id = ?", u.ID).Select()
	return u.Torrents, err
}

func (u *User) update() error {
	return db.Update(u)
}

func (u *User) renewFeedChecked() error {
	u.FeedCheckedAt = time.Now()
	return u.update()
}

func (u *User) isFeedActive() bool {
	return time.Now().Sub(u.FeedCheckedAt) < feedCheckThreshold
}

// BeforeInsert hook
func (u *User) BeforeInsert(db orm.DB) error {
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	return nil
}

// BeforeUpdate hook
func (u *User) BeforeUpdate(db orm.DB) error {
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = time.Now()
	}
	return nil
}
