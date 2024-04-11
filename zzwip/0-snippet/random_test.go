package snippet

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"
)

func Test_Random(t *testing.T) {
	Random()
}

func Random() {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%x\n", buf)
	fmt.Println(hex.EncodeToString(buf))
	fmt.Println(base64.StdEncoding.EncodeToString(buf))

	return

}
