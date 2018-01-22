package main

import (
	"log"
	"os"
	"time"

	"github.com/go-pg/pg"
	"github.com/speps/go-hashids"
)

// Task can be appended by User
type Task struct {
	ID      int
	Title   string
	Magnet  string
	PubDate int64
	Created time.Time
}

// User (i.e. Downloader)
type User struct {
	ID          int
	TelegramID  int
	FeedID      string
	FeedChecked time.Time
	Tasks       []*Task
}

const (
	hashAlphabet = "0123456789abcdef"
)

var db *pg.DB

func initModel() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// localhost
		db = pg.Connect(&pg.Options{User: "postgres"})
	} else {
		// heroku
		opt, err := pg.ParseURL(dbURL)
		if err != nil {
			log.Fatalf("error while parsing url: %s", err)
		}
		db = pg.Connect(opt)
	}

	err := createSchema(db)
	if err != nil {
		log.Printf("error while creating schema: %s", err)
	}
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{&Task{}, &User{}} {
		err := db.CreateTable(model, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Task) create() (*Task, error) {
	_, err := db.Model(t).
		Where("pub_date= ?pub_date").
		OnConflict("DO NOTHING").
		SelectOrInsert()
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Task) getByPubDate() (*Task, error) {
	err := db.Model(t).Where("pub_date = ?", t.PubDate).Select()
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (u *User) create() (*User, error) {

	// create user
	_, err := db.Model(u).
		Where("telegram_id= ?telegram_id").
		OnConflict("DO NOTHING").
		SelectOrInsert()
	if err != nil {
		return nil, err
	}

	// generate a feed id for user
	u, err = u.newFeedID()
	if err != nil {
		return nil, err
	}
	return u, nil
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

func (u *User) appendTask(t *Task) error {
	u, err := u.create()
	if err != nil {
		return err
	}
	u.Tasks = append(u.Tasks, t)
	err = db.Update(u)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) getByFeedID() (*User, error) {
	err := db.Model(u).Where("feed_id = ?", u.FeedID).Select()
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (u *User) update() error {
	return db.Update(u)
}
