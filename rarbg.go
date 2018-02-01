package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	rarbg "github.com/magunetto/go-torrentapi"
)

const (
	ranked = true            // Should results be ranked
	sort   = "seeders"       // Sort order (seeders, leechers, last)
	format = "json_extended" // Format (json, json_extended)
	limit  = 25              // Limit of results (25, 50, 100)
)

var (
	rapi *rarbg.API
)

// InitRARBG init RARBG API
func InitRARBG() {
	api, err := rarbg.Init()
	if err != nil {
		log.Fatalf("error while querying rarbg: %s", err)
	}
	rapi = api
}

func searchTorrents(w io.Writer, service string, id string) (isSingleResult bool) {

	torrentResults, err := searchByServiceID(service, id)
	if err != nil {
		log.Printf("error while querying rarbg: %s", err)
		fmt.Fprintln(w, replyRarbgErr)
		return false
	}

	if len(torrentResults) == 0 {
		log.Printf("no torrents found for this movie: %s", id)
		fmt.Fprintln(w, replyNoTorrents)
		return false
	}

	renderTorrents(w, torrentResults)

	return len(torrentResults) == 1
}

func searchByServiceID(service string, id string) (rarbg.TorrentResults, error) {
	switch service {
	case "imdb":
		rapi.SearchIMDb(id)
	case "tmdb":
		rapi.SearchTheMovieDb(id)
	}
	rapi.Ranked(ranked).Sort(sort).Format(format).Limit(limit)
	return rapi.Search()
}

func renderTorrents(w io.Writer, trs []rarbg.TorrentResult) {

	for _, tr := range trs {
		t, err := saveTorrent(tr)
		if err != nil {
			log.Printf("error while saving torrent: %s", err)
			continue
		}
		t.renderTorrent(w)
	}
}

func saveTorrent(tr rarbg.TorrentResult) (*Torrent, error) {

	t := &Torrent{TorrentResult: tr}
	if tr.Title == "" {
		return t, errors.New("Torrent title should not be empty")
	}

	// use `PubDate` as an unique command for each torrent
	pubDate, err := time.Parse("2006-01-02 15:04:05 +0000", t.PubDate)
	if err != nil {
		return t, err
	}
	t.Title = tr.Title
	t.Magnet = tr.Download
	t.PubStamp = pubDate.Unix()

	return t.create()
}
