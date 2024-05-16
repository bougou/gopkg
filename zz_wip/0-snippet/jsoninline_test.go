package snippet

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Project struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Author interface {
	Commit()
}

// Anonymous struct fields are usually marshaled as if their inner exported fields
// were fields in the outer struct, subject to the usual Go visibility rules
// amended as described in the next paragraph.

// JiraHttpReqField 跟 Project 结构体是平级关系，如果很多 struct 需要 Project 里面的字段，
// 可以直接 inline Project, 减少重复定义
type JiraHttpReqField struct {
	Project     `json:",inline"`
	Key         string `json:"key"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
}

type JiraHttpReqField1 struct {
	Project
	Key         string `json:"key"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
}

// JiraHttpReqField2 test An anonymous struct field with a name given in
// its JSON tag is treated as having that name, rather than being anonymous.
type JiraHttpReqField2 struct {
	Project     `json:"project"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
}

// JiraHttpReqField3 test An anonymous struct field of interface type is
// treated the same as having that type as its name, rather than being anonymous.
type JiraHttpReqField3 struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Author
}

func Test_JsonInline(t *testing.T) {
	p := Project{
		Key:   "key",
		Value: "value",
	}

	j := &JiraHttpReqField{
		Project:     p,
		Summary:     "summary",
		Description: "description",
	}
	d, _ := json.Marshal(j)
	fmt.Println(string(d))
	// {"key":"key","value":"value","summary":"summary","description":"description"}

	j1 := &JiraHttpReqField{
		Project:     p,
		Summary:     "summary",
		Description: "description",
	}
	d1, _ := json.Marshal(j1)
	fmt.Println(string(d1))
	// {"key":"key","value":"value","summary":"summary","description":"description"}

	j2 := &JiraHttpReqField2{
		Project:     p,
		Summary:     "summary",
		Description: "description",
	}
	d2, _ := json.Marshal(j2)
	fmt.Println(string(d2))
	// {"project":{"key":"key","value":"value"},"summary":"summary","description":"description"}

	j3 := &JiraHttpReqField3{
		Summary:     "summary",
		Description: "description",
	}
	d3, _ := json.Marshal(j3)
	fmt.Println(string(d3))
	// {"summary":"summary","description":"description","Author":null}
}
