package main

import (
	"fmt"
	"testing"

	"github.com/bougou/gopkg/common"
	"github.com/davecgh/go-spew/spew"
	"github.com/kr/pretty"
	"gopkg.in/yaml.v3"
)

type Person2 struct {
	Name    yaml.Node
	Address yaml.Node `yaml:"address"`
}

func (p *Person2) UnmarshalYAML(value *yaml.Node) error {
	printNode(value)
	return nil
}

func printNode(value *yaml.Node) {
	pretty.Printf(">>> node: %# v\n", value)
}

func Test_Person2(t *testing.T) {
	data := `
# name line 1
# name line 2
Name: John Doe  # name comment
# line 3
# line 4

# line 5
# line 6
address:   # line 7
# line 8
# line 9

  # line 10
  # line 11

  # line 12
  # line 13
  street: 123 E 3rd St  # line 14
  city: Denver
  state: CO
  zip: 81526
`

	var person Person2
	if err := yaml.Unmarshal([]byte(data), &person); err != nil {
		t.Error(err)
	}
	spew.Dump(person)

	fmt.Println("address:", person.Address)

	b, err := common.Encode("yaml", person)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("yaml:")
	fmt.Println(string(b))

}
