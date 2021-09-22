package common

import (
	"strconv"

	"github.com/alecthomas/units"
)

// Size is an int64
type Size int64

func (s *Size) UnmarshalTOML(b []byte) error {
	var err error
	if len(b) == 0 {
		return nil
	}

	str := string(b)
	if b[0] == '"' || b[0] == '\'' {
		str, err = strconv.Unquote(str)
		if err != nil {
			return err
		}
	}

	val, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		*s = Size(val)
		return nil
	}

	val, err = units.ParseStrictBytes(str)
	if err != nil {
		return err
	}
	*s = Size(val)
	return nil
}

func (s *Size) UnmarshalText(text []byte) error {
	return s.UnmarshalTOML(text)
}
