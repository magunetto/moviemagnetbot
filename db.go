package main

import (
	"log"
	"os"

	"github.com/go-pg/pg"
)

var db *pg.DB

// InitModel init model and create schemas
func InitModel() {

	// init model
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

	// create schemas
	err := createSchema(db)
	if err != nil {
		log.Printf("error while creating schema: %s", err)
	}
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{
		&User{},
		&Torrent{},
		&UserTorrent{},
	} {
		err := db.CreateTable(model, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
