package torrent

import (
	"errors"
	"log"
	"strings"
	"time"

	rarbg "github.com/magunetto/go-torrentapi"
)

const (
	app    = "moviemagnetbot" // app name is required
	ranked = true             // Should results be ranked
	sort   = "seeders"        // Sort order (seeders, leechers, last)
	format = "json_extended"  // Format (json, json_extended)
	limit  = 25               // Limit of results (25, 50, 100)
)

var (
	rapi *rarbg.API

	errRARBG      = errors.New("We encountered an error while finding magnet links, please try again")
	errNoTorrents = errors.New("We have no magnet links for this movie now, please come back later")
)

// InitRARBG init RARBG API
func InitRARBG() {
	api, err := rarbg.Init(app)
	if err != nil {
		log.Fatalf("error while querying rarbg: %s", err)
	}
	rapi = api
}

// Search take keyword to search torrents
func Search(keyword string, limit int) (*[]Torrent, error) {

	keywords := strings.Split(keyword, " ")

	torrentResults, err := searchByServiceID(keywords[0], keywords[1])
	if err != nil {
		log.Printf("error while querying rarbg: %s", err)
		return nil, errRARBG
	}

	if len(torrentResults) == 0 {
		log.Printf("no torrents found for this movie: %s", keywords[1])
		return nil, errNoTorrents
	}

	torrents, err := newTorrentsBySearch(&torrentResults, limit)
	if err != nil {
		return nil, err
	}

	return torrents, nil
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

func newTorrentsBySearch(trs *rarbg.TorrentResults, limit int) (*[]Torrent, error) {

	torrents := []Torrent{}

	for i, tr := range *trs {
		if i == limit {
			break
		}
		t, err := saveTorrent(tr)
		if err != nil {
			continue
		}
		torrents = append(torrents, *t)
	}

	return &torrents, nil
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
