package snippet

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

func Test_1(t *testing.T) {
	s1 := "product-otottABC"

	// TrimLeft will remove ALL leading Unicode code points contained in the cutset.
	// cutset, the set is [p, r, o, d, u, c, t , -]
	cutset := "product-"
	fmt.Println(strings.TrimLeft(s1, cutset)) // Output: ABC

	prefix := "product-"
	fmt.Println(strings.TrimPrefix(s1, prefix)) // Output: otottABC
}

// 按行处理
func Test_Line(t *testing.T) {
	str := "first line\nsecond line\nthird line"
	scanner := bufio.NewScanner(strings.NewReader(str))
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}
}

func Test_ReadLineByLine(t *testing.T) {
	file, err := os.Open("a.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "|")

		id := strings.TrimSpace(fields[1])
		name := strings.TrimSpace(fields[2])

		desc := strings.TrimSpace(fields[3])
		fmt.Printf("update cOsConfig set `desc` = \"%s\" where id = %s and name = \"%s\";\n", desc, id, name)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func Test_Multilines(t *testing.T) {
	a := `First
	Second
	Third`

	if a != "First\n\tSecond\n\tThird" {
		t.Error("a string not expected")
	}

	b := `First
Second
Third`

	if b != "First\nSecond\nThird" {
		t.Error("b string not expected")
	}
}

// strip all non-alphanumeric chars
//
// Alphanumeric characters by definition only comprise the letters A to Z and the digits 0 to 9. Spaces and underscores are usually considered punctuation characters, so no, they shouldn't be allowed.

// If a field specifically says "alphanumeric characters, space and underscore", then they're included. Otherwise in most cases you generally assume they're not.

func StripNonAlphaNumeric(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') {
			result.WriteByte(b)
		}
	}
	return result.String()
}
