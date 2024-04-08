package data

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidTimeSignatureFormat = errors.New("invalid time signature format")

type TimeSignature string

func (ts *TimeSignature) UnmarshalJSON(jsonValue []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidTimeSignatureFormat
	}

	parts := strings.Split(unquotedJSONValue, "/")

	if len(parts) != 2 {
		return ErrInvalidTimeSignatureFormat
	}

	_, err = strconv.Atoi(parts[0])
	if err != nil {
		return ErrInvalidTimeSignatureFormat
	}

	_, err = strconv.Atoi(parts[1])
	if err != nil {
		return ErrInvalidTimeSignatureFormat
	}

	*ts = TimeSignature(unquotedJSONValue)
	return nil
}
