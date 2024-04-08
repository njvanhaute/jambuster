package data

import "time"

type Tune struct {
	ID            int64     `json:"id"`             // Unique integer ID for the tune
	CreatedAt     time.Time `json:"-"`              // Timestamp for when the tune is added to our database
	Title         string    `json:"title"`          // Tune title
	Styles        []string  `json:"styles"`         // Slice of styles for the movie (Bluegrass, old time, Irish, etc.)
	Key           string    `json:"key"`            // Tune key (ex: A major)
	TimeSignature string    `json:"time_signature"` // Tune time signature
	Structure     string    `json:"structure"`      // Tune structure (ex: AABA)
	Version       int32     `json:"version"`        // The version number starts at 1 and will be incremented each time the tune info is updated
}
