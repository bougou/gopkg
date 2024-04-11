package snippet

import (
	"fmt"
	"strings"
	"testing"
)

func Test_BuildString(t *testing.T) {
	var buf strings.Builder

	buf.WriteString("abc")
	buf.Write([]byte{'p'})
	buf.WriteString("empty")

	fmt.Println(buf.String())
}
