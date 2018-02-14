package user

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserFeedActive(t *testing.T) {
	u := &User{
		FeedCheckedAt: time.Now().Add(-feedCheckThreshold + 1*time.Second),
	}
	isActive := u.IsFeedActive()
	assert.Equal(t, true, isActive)
}

func TestUserFeedInactive(t *testing.T) {
	u := &User{
		FeedCheckedAt: time.Now().Add(-feedCheckThreshold),
	}
	isActive := u.IsFeedActive()
	assert.Equal(t, false, isActive)
}
