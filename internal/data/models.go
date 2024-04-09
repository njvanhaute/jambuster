package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Tokens TokenModel
	Tunes  TuneModel
	Users  UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Tokens: TokenModel{DB: db},
		Tunes:  TuneModel{DB: db},
		Users:  UserModel{DB: db},
	}
}
