package data

import "time"

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
