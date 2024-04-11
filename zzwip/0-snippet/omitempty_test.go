package snippet

import (
	"encoding/json"
	"fmt"
	"testing"
)

type S struct {
	A string  `json:"a"`
	B string  `json:"b"`
	C string  `json:"c"`
	D string  `json:"d,omitempty"`
	E *string `json:"e,omitempty"`
}

func Test_f1(t *testing.T) {
	d := `{"a": "x", "b": "", "d": ""}`

	s := S{}
	if err := json.Unmarshal([]byte(d), &s); err != nil {
		fmt.Println(err)
	}
	fmt.Println("a: ", s.A)
	fmt.Println("b: ", s.B)
	fmt.Println("C: ", s.C)
	fmt.Println("D: ", s.D)
	if s.E != nil {
		fmt.Println("E: ", *s.E)
	} else {
		fmt.Println("E:", s.E)
	}

	// a:  0xc0000102f0
	// b:  0xc000010300
	// C:  <empty>
	// D:  <empty>
	// E:  <nil>

}

func Test_f2(t *testing.T) {

	d := `{"a": "x", "e": ""}`

	s := S{}
	if err := json.Unmarshal([]byte(d), &s); err != nil {
		fmt.Println(err)
	}
	fmt.Println("a: ", s.A)
	fmt.Println("b: ", s.B)
	fmt.Println("C: ", s.C)
	fmt.Println("D: ", s.D)
	fmt.Println("E: ", *s.E)

	// a:  0xc000010330
	// b:  <empty>
	// C:  <empty>
	// D:  <empty>
	// E:  <nil>
}

func Test_f3(t *testing.T) {
	s := S{
		A: "hello",
		D: "",
	}

	b, _ := json.Marshal(s)

	fmt.Println(string(b))
}
