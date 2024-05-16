package snippet

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"text/scanner"
)

func Test_Rune(t *testing.T) {
	const a = `你好，世界`
	fmt.Printf("%s\n", a) // 你好，世界

	for _, r := range a {
		fmt.Println(r, string(r), strconv.QuoteRune(r))
	}
	// 20320  你   '你'
	// 22909  好   '好'
	// 65292  ，   '，'
	// 19990  世   '世'
	// 30028  界   '界'

	fmt.Println(reflect.TypeOf(a)) // string

	// string to []rune
	r := []rune(a)
	fmt.Println(reflect.TypeOf(r)) // []int32
	fmt.Println(r)                 // [20320 22909 65292 19990 30028]

	s := string(r)
	fmt.Println(s)            // 你好，世界
	fmt.Println(string(r[0])) // 你

	for i, r := range s {
		fmt.Printf("%d - %c (%v)\n", i, r, r) // Note, the i value
		// 0 - 你 (20320)
		// 3 - 好 (22909)
		// 6 - ， (65292)
		// 9 - 世 (19990)
		// 12 - 界 (30028)
	}

	var sc scanner.Scanner
	reader := strings.NewReader(s) // 将 string 创建为 io.Reader
	sc.Init(reader)                // 用 io.Reader 初始化 scanner.Scanner

	fmt.Println("===")
	i := 0
	for r := sc.Next(); i < 4; {
		fmt.Println(sc.Position)
		fmt.Println(r)
		fmt.Println(string(r))
		fmt.Println(r == scanner.Ident)
		fmt.Println(sc.TokenText())
		i++

	}

	fmt.Println("===")

	r1 := sc.Next()
	fmt.Println(r1) // 20320

	r2 := sc.Next()
	fmt.Println(r2) // 22909

	r3 := sc.Scan()
	fmt.Println(r3) // 65292

	r4 := sc.Scan()
	fmt.Println(r4) // -2
	fmt.Println(r4 == scanner.Ident)
	fmt.Println(sc.TokenText()) // 世界

	r5 := sc.Scan()
	fmt.Println(r5) // -1
	fmt.Println(r5 == scanner.Ident)
	fmt.Println(r5 == scanner.Char)
	fmt.Println(sc.TokenText()) // 空值

}

// the following are digest from scanner package

// The result of Scan is one of these tokens or a Unicode character.
const (
	EOF = -(iota + 1)
	Ident
	Int
	Float
	Char
	String
	RawString
	Comment

	// internal use only
	skipComment
)

var tokenString = map[rune]string{
	EOF:       "EOF",
	Ident:     "Ident",
	Int:       "Int",
	Float:     "Float",
	Char:      "Char",
	String:    "String",
	RawString: "RawString",
	Comment:   "Comment",
}

// TokenString returns a printable string for a token or Unicode character.
func TokenString(tok rune) string {
	if s, found := tokenString[tok]; found {
		return s
	}
	return fmt.Sprintf("%q", string(tok))
}
