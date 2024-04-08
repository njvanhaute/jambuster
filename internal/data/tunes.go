package data

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
	"jambuster.njvanhaute.com/internal/validator"
)

type Tune struct {
	ID            int64         `json:"id"`             // Unique integer ID for the tune
	CreatedAt     time.Time     `json:"-"`              // Timestamp for when the tune is added to our database
	Title         string        `json:"title"`          // Tune title
	Styles        []string      `json:"styles"`         // Slice of styles for the tune (Bluegrass, old time, Irish, etc.)
	Keys          []Key         `json:"keys"`           // Slice of keys for the tune (ex: A major, G minor)
	TimeSignature TimeSignature `json:"time_signature"` // Tune time signature
	Structure     string        `json:"structure"`      // Tune structure (ex: AABA)
	HasLyrics     bool          `json:"has_lyrics"`     // Whether or not the tune has lyrics
	Version       int32         `json:"version"`        // The version number starts at 1 and will be incremented each time the tune info is updated
}

type TuneModel struct {
	DB *sql.DB
}

func (t TuneModel) Insert(tune *Tune) error {
	query := `
		INSERT INTO tunes (title, styles, keys, time_signature, structure, has_lyrics)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, version`

	args := []any{tune.Title, pq.Array(tune.Styles), pq.Array(tune.Keys), tune.TimeSignature, tune.Structure, tune.HasLyrics}

	return t.DB.QueryRow(query, args...).Scan(&tune.ID, &tune.CreatedAt, &tune.Version)
}

func (t TuneModel) Get(id int64) (*Tune, error) {
	return nil, nil
}

func (t TuneModel) Update(tune *Tune) error {
	return nil
}

func (t TuneModel) Delete(id int64) error {
	return nil
}

func ValidateTune(v *validator.Validator, tune *Tune) {
	v.Check(tune.Title != "", "title", "must be provided")
	v.Check(len(tune.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(tune.Styles != nil, "styles", "must be provided")
	v.Check(len(tune.Styles) >= 1, "styles", "must contain at least 1 style")
	v.Check(len(tune.Styles) <= 5, "styles", "must not contain more than 5 styles")
	v.Check(validator.Unique(tune.Styles), "styles", "must not contain duplicate values")

	v.Check(tune.Keys != nil, "keys", "must be provided")
	v.Check(len(tune.Keys) >= 1, "keys", "must contain at least 1 key")
	v.Check(len(tune.Keys) <= 10, "keys", "must not contain more than 10 keys")
	v.Check(validator.Unique(tune.Keys), "keys", "must not contain duplicate values")

	v.Check(tune.TimeSignature != "", "time_signature", "must be provided")
	v.Check(len(tune.TimeSignature) >= 3, "time_signature", "must be at least 3 characters long")
	v.Check(len(tune.TimeSignature) <= 5, "time_signature", "must not be more than 5 characters long")

	v.Check(tune.Structure != "", "structure", "must be provided")
	v.Check(len(tune.Structure) >= 1, "structure", "must be at least 1 character long")

}
