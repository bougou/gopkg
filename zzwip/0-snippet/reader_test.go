package snippet

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

// golang.org/src/io/io.go
// Reader is the interface that wraps the basic Read method.
//
// Read reads up to len(p) bytes into p. It returns the number of bytes
// read (0 <= n <= len(p)) and any error encountered. Even if Read
// returns n < len(p), it may use all of p as scratch space during the call.
// If some data is available but not len(p) bytes, Read conventionally
// returns what is available instead of waiting for more.
//
// When Read encounters an error or end-of-file condition after
// successfully reading n > 0 bytes, it returns the number of
// bytes read. It may return the (non-nil) error from the same call
// or return the error (and n == 0) from a subsequent call.
// An instance of this general case is that a Reader returning
// a non-zero number of bytes at the end of the input stream may
// return either err == EOF or err == nil. The next Read should
// return 0, EOF.
//
// Callers should always process the n > 0 bytes returned before
// considering the error err. Doing so correctly handles I/O errors
// that happen after reading some bytes and also both of the
// allowed EOF behaviors.
//
// Implementations of Read are discouraged from returning a
// zero byte count with a nil error, except when len(p) == 0.
// Callers should treat a return of 0 and nil as indicating that
// nothing happened; in particular it does not indicate EOF.
//
// Implementations must not retain p.
// type Reader interface {
// 	Read(p []byte) (n int, err error)
// }

type X struct {
	Name string
}

// 读取 X
type XReader struct {
	X
}

// 实现 io.Reader 接口
func (r *XReader) Read(p []byte) (n int, err error) {

	sr := strings.NewReader(r.Name)
	sn, err := sr.Read(p)
	fmt.Println(sn)

	fmt.Println(p)
	return len(p), nil
}

func Test_Read(t *testing.T) {
	x := X{
		Name: "hello world",
	}

	r := XReader{
		X: x,
	}

	p := make([]byte, 15)
	n, err := r.Read(p)
	fmt.Println(p)

	if err != nil {
		fmt.Println("error: ", err)
	} else {
		fmt.Println("read: ", n)
	}

	fmt.Println(p)
}

func T() {
	var r io.Reader

	// create reader from string
	r = strings.NewReader("this-is-string")

	// create reader from file
	// os.Open return *os.File which implements io.Reader
	// read completely into memory
	r, err := os.Open("/path/to/file")
	if err != nil {
		fmt.Print(err)
	}

	file, err := os.Open("/path/to/file")
	if err != nil {
		fmt.Println(err)
	}
	r = bufio.NewReader(file)

	// read content to out
	var out []byte
	r.Read(out)
}
