package torrent

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMagnetTitle(t *testing.T) {
	t.Parallel()
	const dn = "Once.Upon.a.Time.%5foo%20bar%5D"
	const link = "magnet:?xt=urn:btih:foo&dn=" + dn + "&tr=http%3A%2F%2Ftracker.example.com%3A80%2Fannounce&tr=udp%3A%2F%2F9.example.com%3A2950&tr=udp%3A%2F%2F9.example.com%3A2800"
	dnUnescaped, err := url.QueryUnescape(dn)
	assert.NoError(t, err)

	tor := newTorrentFromLink(link)

	assert.NotEqual(t, dn, tor.Title)
	assert.Equal(t, dnUnescaped, tor.Title)
}

func TestFTPTitle(t *testing.T) {
	t.Parallel()
	const filename = "Once.Upon.a.Time.%5foo%20bar%5D"
	const link = "ftp://1.1.1.1/" + filename
	filenameUnescaped, err := url.QueryUnescape(filename)
	assert.NoError(t, err)

	tor := newTorrentFromLink(link)

	assert.NotEqual(t, filename, tor.Title)
	assert.Equal(t, filenameUnescaped, tor.Title)
}

func TestDefaultTitleIsFullURL(t *testing.T) {
	t.Parallel()
	for _, link := range []string{
		// unknown scheme
		"unknown://user:pass@1.1.1.1/foo/bar%20foo?bar#foo.bar",
		// valid magnet without a valid dn field
		"magnet:?xt=urn:btih:foo&tr=http%3A%2F%2Ftracker.example.com%3A80%2Fannounce&tr=udp%3A%2F%2F9.example.com%3A2950&tr=udp%3A%2F%2F9.example.com%3A2800",
		"magnet:?xt=urn:btih:foo&dn=&tr=http%3A%2F%2Ftracker.example.com%3A80%2Fannounce&tr=udp%3A%2F%2F9.example.com%3A2950&tr=udp%3A%2F%2F9.example.com%3A2800",
		"magnet:?xt=urn:btih:foo&dn&tr=http%3A%2F%2Ftracker.example.com%3A80%2Fannounce&tr=udp%3A%2F%2F9.example.com%3A2950&tr=udp%3A%2F%2F9.example.com%3A2800",
		// FTP full site, are these really going to happen?
		"ftp://a.b.c",
		"ftp://1.1.1.1",
		"ftp://1.1.1.1/",
	} {
		tor := newTorrentFromLink(link)
		assert.Equal(t, link, tor.Title)
	}
}
