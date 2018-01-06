package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
)

func feedHandler(w http.ResponseWriter, r *http.Request) {

	// query user by telegram id
	id, err := strconv.Atoi(mux.Vars(r)["user"])
	if err != nil {
		log.Printf("error while parsing telegram id: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := &User{TelegramID: id}
	user, err = user.getByTelegram()
	if err != nil && user == nil {
		// no user with this telegram id
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("error while getting user: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// generate feed for user
	feed := &feeds.Feed{
		Title:       fmt.Sprintf("Movie Magnet Bot Feed #%d", id),
		Link:        &feeds.Link{Href: host + r.URL.String()},
		Description: fmt.Sprintf("Download tasks for user #%d", id),
		Created:     time.Now(),
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
