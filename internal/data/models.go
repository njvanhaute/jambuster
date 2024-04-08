package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Tunes TuneModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Tunes: TuneModel{DB: db},
	}
}
