package uri

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMagnet(t *testing.T) {
	t.Parallel()
	assert.True(t, IsMagnet("magnet:?xt=urn:btih:c12fe1c06bba254a9dc9f519b335aa7c1367a88a"))
	assert.True(t, IsMagnet("MAGNET:?xt=urn:btih:c12fe1c06bba254a9dc9f519b335aa7c1367a88a"))
	assert.False(t, IsMagnet("magnet?????"))
	assert.False(t, IsMagnet(""))
}

func TestED2k(t *testing.T) {
	t.Parallel()
	assert.True(t, IsED2k("ed2k://|file|Shareaza_2.5.3.0_Win32.exe|6653348|7fb2bc10e0422a0e4f7e8613bd522c89|/|sources,252.191.193.62:6443|/"))
	assert.True(t, IsED2k("ED2K://|file|Shareaza_2.5.3.0_Win32.exe|6653348|7fb2bc10e0422a0e4f7e8613bd522c89|/|sources,252.191.193.62:6443|/"))
	assert.False(t, IsED2k("ed2k?????"))
	assert.False(t, IsED2k(""))
}

func TestFTP(t *testing.T) {
	t.Parallel()
	assert.True(t, IsFTP("ftp://user:pass@a.b.c/%2Ffoo/foobar%202k.mov"))
	assert.True(t, IsFTP("FTP://user:pass@a.b.c/%2Ffoo/foobar%202k.mov"))
	assert.False(t, IsFTP("ftp?????"))
	assert.False(t, IsFTP(""))
}

func TestIsValid(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		shouldBeValid bool
		url           string
	}{
		{true, "ftp://user:pass@a.b.c/%2Ffoo/foobar%202k.mov"},
		{true, "ed2k://|file|Shareaza_2.5.3.0_Win32.exe|6653348|7fb2bc10e0422a0e4f7e8613bd522c89|/|sources,252.191.193.62:6443|/"},
		{true, "magnet:?xt=urn:btih:c12fe1c06bba254a9dc9f519b335aa7c1367a88a"},
		// NOTE: Currently things like "edk2:?foo=bar" and "magnet://?foo=example" will be regarded as valid
		// I'm leaving it as it is for now, because:
		//   0. Following the standards should be a good practice
		//   1. I'm not sure if the same URI could be expressed in multiple forms
		//   2. URIs like these are unlikely to appear on the Internet, if they did, they should most likely be supported
		{false, "ed2k"},
		{false, "ftp"},
		{false, "magnet"},
		{false, "unknown"},
	} {
		assert.Equal(t, c.shouldBeValid, NewURI(c.url).IsValid())
	}
}
