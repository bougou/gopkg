package common

import "strconv"

type Number struct {
	Value float64
}

func (n *Number) UnmarshalTOML(b []byte) error {
	value, err := strconv.ParseFloat(string(b), 64)
	if err != nil {
		return err
	}

	n.Value = value
	return nil
}
