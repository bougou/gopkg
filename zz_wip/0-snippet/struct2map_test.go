package snippet

import (
	"fmt"
	"testing"

	"github.com/fatih/structs"
)

func Test_Struct2Map(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int64  `json:"age"`
		Dead bool   `json:"dead"`
	}

	p := &Person{
		Name: "Tom",
		Age:  20,
		Dead: false,
	}
	fmt.Println("object:", p)

	pMap := structs.Map(p)
	fmt.Println("default:", pMap) // use struct field name as key of map

	// structs package default use "structs" tag name to tweak the struct data
	// but we can change it to "json" tag name
	// if the tag not found on the field, use Field name as key of map
	pStruct := structs.New(p)
	pStruct.TagName = "json"
	pMap = pStruct.Map()
	fmt.Println("tag json:", pMap) // use tag name as key of map

	for k, v := range pMap {
		fmt.Printf("key: %v, value: %v, value type: %T\n", k, v, v)
	}

}
