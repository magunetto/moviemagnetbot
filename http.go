package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
)

// RunHTTPServer register handlers and start HTTP server
func RunHTTPServer() {

	// http handlers
	r := mux.NewRouter()
	r.HandleFunc("/tasks/{user}.xml", feedHandler)

	// http loop
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}

func feedHandler(w http.ResponseWriter, r *http.Request) {

	// query user by feed id
	u := &User{FeedID: mux.Vars(r)["user"]}
	u, err := u.getByFeedID()
	if err != nil {
		log.Printf("error while getting user: %s", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// generate feed for user
	feed := &feeds.Feed{
		Title:   fmt.Sprintf("Movie Magnet Bot feed"),
		Link:    &feeds.Link{Href: host + r.URL.String()},
		Created: time.Now(),
	}
	torrents, err := u.getTorrents(itemsPerFeed)
	if err != nil {
		log.Printf("error while getting torrents: %s", err)
	}
	for _, t := range torrents {
		feed.Items = append(feed.Items, &feeds.Item{
			Title:   t.Title,
			Link:    &feeds.Link{Href: t.Magnet},
			Created: t.DownloadedAt,
		})
	}
	rss, err := feed.ToRss()
	if err != nil {
		log.Printf("error while generating feed: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// send the feed
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(rss))

	err = u.renewFeedChecked()
	if err != nil {
		log.Printf("error while updating user: %s", err)
	}
}
