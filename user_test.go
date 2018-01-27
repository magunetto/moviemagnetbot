package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserFeedActive(t *testing.T) {
	u := &User{
		FeedChecked: time.Now().Add(-feedCheckThreshold + 1*time.Second),
	}
	isActive := u.isFeedActive()
	assert.Equal(t, true, isActive)
}

func TestUserFeedInactive(t *testing.T) {
	u := &User{
		FeedChecked: time.Now().Add(-feedCheckThreshold),
	}
	isActive := u.isFeedActive()
	assert.Equal(t, false, isActive)

}
