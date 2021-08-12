package common

import (
	"bytes"
	"strconv"

	"github.com/alecthomas/units"
)

// Size just wraps an int64
type Size struct {
	Size int64
}

func (s *Size) UnmarshalTOML(b []byte) error {
	var err error
	b = bytes.Trim(b, `'`)

	val, err := strconv.ParseInt(string(b), 10, 64)
	if err == nil {
		s.Size = val
		return nil
	}
	uq, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	val, err = units.ParseStrictBytes(uq)
	if err != nil {
		return err
	}
	s.Size = val
	return nil
}
