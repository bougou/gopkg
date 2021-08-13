package common

import (
	"bytes"
	"io"
)

type Encoder interface {
	Encode(obj interface{}, w io.Writer) error
	Identifier() string
}

type Decoder interface {
	Decode(data []byte) (interface{}, error)
}

type Serializer interface {
	Encoder
	Decoder
}

// codec binds an encoder and decoder.
type codec struct {
	Encoder
	Decoder
}

// Encode is a convenience wrapper for encoding obj to []byte from an Encoder
func Encode(e Encoder, obj interface{}) ([]byte, error) {
	// TODO: reuse buffer
	buf := &bytes.Buffer{}
	if err := e.Encode(obj, buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode is a convenience wrapper for decoding data into an obj.
func Decode(d Decoder, data []byte) (obj interface{}, err error) {
	obj, err = d.Decode(data)
	return
}
