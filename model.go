package main

import (
	"log"
	"os"

	"github.com/go-pg/pg"
)

// Task can be appended by User
type Task struct {
	ID      int
	Title   string
	Magnet  string
	PubDate int64
}

// User (i.e. Downloader)
type User struct {
	ID         int
	TelegramID int
	Tasks      []*Task
}

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
			log.Fatalf("error while connecting database: %s", err)
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

func (task *Task) create() {
	_, err := db.Model(task).
		Where("pub_date= ?pub_date").
		OnConflict("DO NOTHING").
		SelectOrInsert()
	if err != nil {
		log.Printf("error while creating task: %s", err)
	}
}

func (task *Task) getByPubDate() (*Task, error) {
	err := db.Model(task).Where("pub_date = ?", task.PubDate).Select()
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (user *User) create() {
	_, err := db.Model(user).
		Where("telegram_id= ?telegram_id").
		OnConflict("DO NOTHING").
		SelectOrInsert()
	if err != nil {
		log.Printf("error while getting user: %s", err)
	}
}

func (user *User) appendTask(task *Task) {
	user.create()
	user.Tasks = append(user.Tasks, task)
	err := db.Update(user)
	if err != nil {
		log.Printf("error while appending task: %s", err)
	}
}
