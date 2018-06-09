package torrent

func SaveTorrentFromLink(link string) (*Torrent, error) {
	t := &Torrent{Magnet: link}
	return t.create()
}
