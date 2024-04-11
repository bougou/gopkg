package snippet

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type Human struct {
	Head string `json:"head" xml:"head"`
	Body string `json:"body" xml:"body"`
	Leg  string `json:"leg_j,omitempty" xml:"leg_x"`
}

func getFieldNameByTag(s interface{}, tagKey string, tagValue string) string {
	t := reflect.TypeOf(s)
	// 如果传入的是指针类型，则获取指针的基类型
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		panic("Bad type")
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		v := strings.Split(f.Tag.Get(tagKey), ",")[0] // use split to ignore tag "options"
		if v == tagValue {
			return f.Name
		}
	}
	return ""
}

func getTagNameByField(s interface{}, tagKey string, field string) string {
	t := reflect.TypeOf(s)
	// 如果传入的是指针类型，则获取指针的基类型
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		panic("Bad type")
	}

	f, ok := t.FieldByName(field)
	if !ok {
		panic("Field not found")
	}
	return strings.Split((f.Tag.Get(tagKey)), ",")[0]
}

func Test_Human(t *testing.T) {
	fmt.Println(getTagNameByField(Human{}, "json", "Head")) // head
	fmt.Println(getTagNameByField(&Human{}, "json", "Leg")) // leg_j
	fmt.Println(getTagNameByField(Human{}, "xml", "Leg"))   // leg_x

	fmt.Println(getFieldNameByTag(Human{}, "json", "head"))  // Head
	fmt.Println(getFieldNameByTag(&Human{}, "json", "body")) // Body
	fmt.Println(getFieldNameByTag(Human{}, "json", "leg_j")) // Leg
	fmt.Println(getFieldNameByTag(&Human{}, "xml", "leg_x")) // Leg
}
