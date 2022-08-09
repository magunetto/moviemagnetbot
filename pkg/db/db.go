package db

import (
	"log"
	"os"

	"github.com/go-pg/pg/v10"
)

// DB is *pg.DB
var DB *pg.DB

// Init DB
func Init() {

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// localhost
		DB = pg.Connect(&pg.Options{User: "postgres"})
	} else {
		// heroku
		opt, err := pg.ParseURL(dbURL)
		if err != nil {
			log.Fatalf("error while parsing url: %s", err)
		}
		DB = pg.Connect(opt)
	}

}
