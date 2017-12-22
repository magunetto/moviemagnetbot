package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	rarbg "github.com/umayr/go-torrentapi"
)

const (
	ranked = true            // Should results be ranked
	sort   = "seeders"       // Sort order (seeders, leechers, last)
	format = "json_extended" // Format (json, json_extended)
	limit  = 25              // Limit of results (25, 50, 100)
)

func searchIMDb(w io.Writer, id string, api *rarbg.API) {

	// search torrents for a movie
	results, err := search(api, "imdb", id)
	if err != nil {
		log.Printf("error while querying torrentapi: %s", err)
		fmt.Fprintln(w, errorText)
		return
	}
	fmt.Fprintln(w, "Results for *", id, "*")
	fmt.Fprintln(w, "Seeders / Leechers / Size / File Name")

	// output every torrent
	for _, r := range results {
		magnet := strings.Split(r.Download, "&")[0]
		fmt.Fprintf(w, "`%d` / `%d` / `%s` / `%s` / `%s`\n",
			r.Seeders, r.Leechers, humanizeSize(r.Size), r.Title, magnet)
	}
}

func search(api *rarbg.API, clue string, keyword string) (results rarbg.TorrentResults, err error) {
	switch clue {
	case "tvdb":
		api.SearchTVDB(keyword)
	case "imdb":
		api.SearchImDB(keyword)
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
