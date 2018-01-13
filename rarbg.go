package main

import (
	"fmt"
	"io"
	"log"
	"time"

	rarbg "github.com/idealhack/go-torrentapi"
)

const (
	ranked = true            // Should results be ranked
	sort   = "seeders"       // Sort order (seeders, leechers, last)
	format = "json_extended" // Format (json, json_extended)
	limit  = 25              // Limit of results (25, 50, 100)
)

func searchIMDb(w io.Writer, id string, api *rarbg.API) {

	// search torrents by imdb id
	results, err := search(api, "imdb", id)
	if err != nil {
		log.Printf("error while querying rarbg: %s", err)
		fmt.Fprintln(w, replyRarbgErr)
		return
	}
	if len(results) == 0 {
		log.Printf("no torrents found for this movie: %s", id)
		fmt.Fprintln(w, replyNoTorrents)
		return
	}
	fmt.Fprintf(w, "§ `%s`\n", id)

	// output torrents
	for _, r := range results {

		// use `PubDate` as an unique command for each torrent
		t, err := time.Parse("2006-01-02 15:04:05 +0000", r.PubDate)
		if err != nil {
			log.Printf("error while parsing date: %s", err)
		}
		command := fmt.Sprintf("%s%d", cmdPrefixDown, t.Unix())
		fmt.Fprintf(w, "*%d*↑ *%d*↓ ∑`%s` %s\n", r.Seeders, r.Leechers, humanizeSize(r.Size), command)
		fmt.Fprintf(w, "[¶](%s) %s\n", r.InfoPage, r.Title)

		// create task for torrent
		task := &Task{
			Title:   r.Title,
			Magnet:  r.Download,
			PubDate: t.Unix(),
		}
		_, err = task.create()
		if err != nil {
			log.Printf("error while creating task: %s", err)
		}
	}
}

func search(api *rarbg.API, clue string, keyword string) (results rarbg.TorrentResults, err error) {
	switch clue {
	case "tvdb":
		api.SearchTVDB(keyword)
	case "imdb":
		api.SearchIMDb(keyword)
	case "search":
		api.SearchString(keyword)
	}

	api.Ranked(ranked).Sort(sort).Format(format).Limit(limit)
	results, err = api.Search()
	return
}

func humanizeSize(s uint64) string {
	size := float64(s)
	switch {
	case size < 1024:
		return fmt.Sprintf("%d", uint64(size))
	case size < 1024*1014:
		return fmt.Sprintf("%.2fk", size/1024)
	case size < 1024*1024*1024:
		return fmt.Sprintf("%.2fM", size/1024/1024)
	default:
		return fmt.Sprintf("%.2fG", size/1024/1024/1024)
	}
}
