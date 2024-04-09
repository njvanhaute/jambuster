package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return t.DB.QueryRowContext(ctx, query, args...).Scan(&tune.ID, &tune.CreatedAt, &tune.Version)
}

func (t TuneModel) Get(id int64) (*Tune, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, title, styles, keys, time_signature, structure, has_lyrics, version
		FROM tunes
		WHERE id = $1`

	var tune Tune
	var keyStrings []string

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := t.DB.QueryRowContext(ctx, query, id).Scan(
		&tune.ID,
		&tune.CreatedAt,
		&tune.Title,
		pq.Array(&tune.Styles),
		pq.Array(&keyStrings),
		&tune.TimeSignature,
		&tune.Structure,
		&tune.HasLyrics,
		&tune.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	for _, keyString := range keyStrings {
		tune.Keys = append(tune.Keys, Key(keyString))
	}

	return &tune, nil
}

func (t TuneModel) GetAll(title string, styles []string, keys []string, timeSignature string, structure string, hasLyrics *bool, filters Filters) ([]*Tune, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, created_at, title, styles, keys, time_signature, structure, has_lyrics, version
		FROM tunes
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (styles @> $2 OR $2 = '{}')
		AND (keys @> $3 OR $3 = '{}')
		AND (time_signature = $4 OR $4 = '')
		AND (structure = $5 OR $5 = '')
		AND (has_lyrics = $6 OR $6 IS NULL)
		ORDER BY %s %s, id ASC
		LIMIT $7 OFFSET $8`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{title, pq.Array(styles), pq.Array(keys), timeSignature, structure, hasLyrics,
		filters.limit(), filters.offset()}

	rows, err := t.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	tunes := []*Tune{}

	for rows.Next() {
		var tune Tune
		var keyStrings []string

		err := rows.Scan(
			&totalRecords,
			&tune.ID,
			&tune.CreatedAt,
			&tune.Title,
			pq.Array(&tune.Styles),
			pq.Array(&keyStrings),
			&tune.TimeSignature,
			&tune.Structure,
			&tune.HasLyrics,
			&tune.Version,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		for _, keyString := range keyStrings {
			tune.Keys = append(tune.Keys, Key(keyString))
		}

		tunes = append(tunes, &tune)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return tunes, metadata, nil
}

func (t TuneModel) Update(tune *Tune) error {
	query := `
		UPDATE tunes
		SET title = $1, styles = $2, keys = $3, time_signature = $4, structure = $5, has_lyrics = $6, version = version + 1
		WHERE id = $7 AND version = $8
		RETURNING version`

	args := []any{
		tune.Title,
		pq.Array(tune.Styles),
		pq.Array(tune.Keys),
		tune.TimeSignature,
		tune.Structure,
		tune.HasLyrics,
		tune.ID,
		tune.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := t.DB.QueryRowContext(ctx, query, args...).Scan(&tune.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (t TuneModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM tunes
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := t.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

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
