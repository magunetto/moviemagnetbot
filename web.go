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
	user := &User{FeedID: mux.Vars(r)["user"]}
	user, err := user.getByFeedID()
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
	for _, t := range user.Tasks {
		feed.Items = append(feed.Items, &feeds.Item{
			Title:   t.Title,
			Link:    &feeds.Link{Href: t.Magnet},
			Created: t.Created,
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
}
