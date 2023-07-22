package qrcode

import (
	"testing"
)

func Test_1(t *testing.T) {
	if err := qr(); err != nil {
		panic(err)
	}
}
