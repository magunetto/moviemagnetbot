package main

import (
	"log"

	"github.com/magunetto/moviemagnetbot/pkg/bot"
	"github.com/magunetto/moviemagnetbot/pkg/db"
	"github.com/magunetto/moviemagnetbot/pkg/http"
	"github.com/magunetto/moviemagnetbot/pkg/model"
	"github.com/magunetto/moviemagnetbot/pkg/movie"
	"github.com/magunetto/moviemagnetbot/pkg/torrent"
)

func main() {

	db.Init()
	log.Printf("db inited")

	err := model.CreateSchema(db.DB)
	if err != nil {
		log.Printf("error while creating schema: %s", err)
	}

	movie.InitTMDb()
	log.Printf("tmdb inited")

	torrent.InitRARBG()
	log.Printf("rarbg inited")

	go bot.Run()
	log.Printf("bot started")

	go http.RunServer()
	log.Printf("http server started")

	select {}
}
