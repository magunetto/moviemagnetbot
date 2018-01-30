package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
)

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
	torrents, err := u.getTorrents()
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
