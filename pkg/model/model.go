package model

import (
	"github.com/go-pg/pg/v10"

	"github.com/magunetto/moviemagnetbot/pkg/torrent"
	"github.com/magunetto/moviemagnetbot/pkg/user"
)

// CreateSchema create tables for models
func CreateSchema(db *pg.DB) error {
	for _, model := range []interface{}{
		&torrent.Torrent{},
		&user.User{},
		&user.UserTorrent{},
	} {
		err := db.Model(model).CreateTable(nil)
		if err != nil {
			return err
		}
	}
	return nil
}
