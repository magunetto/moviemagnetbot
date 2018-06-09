package torrent

func SaveTorrentFromMagnet(magnet string) (*Torrent, error) {
	t := &Torrent{Magnet: magnet}
	return t.create()
}

func (t *Torrent) IsFromMangnet() bool {
	return t.Title == ""
}
