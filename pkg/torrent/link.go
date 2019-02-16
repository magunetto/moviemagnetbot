package torrent

// SaveTorrentFromLink creates a Torrent from a link
func SaveTorrentFromLink(link string) (*Torrent, error) {
	t := &Torrent{Magnet: link}
	return t.create()
}
