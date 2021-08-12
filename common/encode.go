package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

func Encode(format string, obj interface{}) ([]byte, error) {
	var b bytes.Buffer

	switch format {
	case "yaml":
		e := yaml.NewEncoder(&b)
		e.SetIndent(2)
		if err := e.Encode(obj); err != nil {
			msg := fmt.Sprintf("yaml encode failed, err: %s", err)
			return nil, errors.New(msg)
		}
	case "json", "":
		e := json.NewEncoder(&b)
		e.SetIndent("", "  ")
		if err := e.Encode(obj); err != nil {
			msg := fmt.Sprintf("json encode failed, err: %s", err)
			return nil, errors.New(msg)
		}
	case "json-compact":
		e := json.NewEncoder(&b)
		e.SetIndent("", "")
		if err := e.Encode(obj); err != nil {
			msg := fmt.Sprintf("compact json encode failed, err: %s", err)
			return nil, errors.New(msg)
		}
	default:
		msg := fmt.Sprintf("unsupported output format (%s)", format)
		return nil, errors.New(msg)
	}

	return b.Bytes(), nil
}
