package torrent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHumanizeSizeInByte(t *testing.T) {
	size := humanizeSize(0)
	assert.Equal(t, "0", size)
}

func TestHumanizeSizeInKB(t *testing.T) {
	size := humanizeSize(1024)
	assert.Equal(t, "1.00K", size)
}

func TestHumanizeSizeInMB(t *testing.T) {
	size := humanizeSize(1024 * 1024)
	assert.Equal(t, "1.00M", size)
}

func TestHumanizeSizeInGB(t *testing.T) {
	size := humanizeSize(1024 * 1024 * 1024)
	assert.Equal(t, "1.00G", size)
}
