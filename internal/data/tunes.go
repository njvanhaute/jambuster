package data

import (
	"time"

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
	Version       int32         `json:"version"`        // The version number starts at 1 and will be incremented each time the tune info is updated
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
