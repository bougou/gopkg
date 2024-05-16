package snippet

import (
	"encoding/json"
	"fmt"
	"reflect"
)

var j string = `{
	"kstring": "hello world",
	"kinterger": 100,
	"kfloat": 12.34,
	"kstringfloat": "100",
	"kstringfloat2": "100",
	"kmap": {
		"kstring": "hello world",
		"kinterger": 100,
		"kfloat": 12.34
	},
	"klist": [
		{
			"kstring": "hello world",
			"kinterger": 100,
			"kfloat": 12.34
		},
		{
			"kstring": "hello world",
			"kinterger": 100,
			"kfloat": 12.34
		}
	]
}`

func json2map() {

	m := map[string]interface{}{}

	json.Unmarshal([]byte(j), &m)
	fmt.Println(m)

	for k, v := range m {
		switch t := v.(type) {
		case string, int32, int64, float32, float64, map[string]interface{}, []interface{}:
			fmt.Printf("key: %s, type: %s, value: %v\n", k, reflect.TypeOf(v), t)

		default:
			fmt.Println("unknown")
		}
	}

	// To unmarshal JSON into an interface value, Unmarshal stores one of these in the interface value:
	// bool, for JSON booleans
	// float64, for JSON numbers
	// string, for JSON strings
	// []interface{}, for JSON arrays
	// map[string]interface{}, for JSON objects
	// nil for JSON null

	// 如果要访问 map 内层字段，需要逐层使用 type assert
	_klist := m["klist"].([]interface{})
	_klist0 := _klist[0].(map[string]interface{})
	_klist0kstring := _klist0["kstring"]
	fmt.Println(_klist0kstring)
}

func json2struct() {
	// 转换工具 https://mholt.github.io/json-to-go/
	// 定义 struct，明确每个字段、内层字段的类型
	type jStruct struct {
		Kstring       string  `json:"kstring"`
		Kinterger     int     `json:"kinterger"`
		Kfloat        float64 `json:"kfloat"`
		KStringFloat  float64 `json:"kstringfloat"`         // if json value is string, unmarshal will fill it 0.
		KStringFloat2 float64 `json:"kstringfloat2,string"` // recogize string value and auto convert to float64

		// The "string" option signals that a field is stored as JSON inside a JSON-encoded string. It applies only to fields of string, floating point, integer, or boolean types. This extra level of encoding is sometimes used when communicating with JavaScript programs:

		Kmap struct {
			Kstring   string  `json:"kstring"`
			Kinterger int     `json:"kinterger"`
			Kfloat    float64 `json:"kfloat"`
		} `json:"kmap"`
		Klist []struct {
			Kstring   string  `json:"kstring"`
			Kinterger int     `json:"kinterger"`
			Kfloat    float64 `json:"kfloat"`
		} `json:"klist"`
	}

	s := jStruct{}
	json.Unmarshal([]byte(j), &s)
	fmt.Println(s)
	// NOTE !!!!!!!
	fmt.Println(s.KStringFloat)  // 0
	fmt.Println(s.KStringFloat2) // 100

	// 访问时，不需要使用 type assert
	fmt.Println(s.Klist[0].Kstring)

	b, _ := json.MarshalIndent(s, "", "  ")
	fmt.Println(string(b))

}

type Cache map[string]string

func test3() {
	c := Cache{}
	c["hello"] = "world"
	c["test"] = "test1"
	c["a"] = "b"

	delete(c, "test")
	fmt.Println(c)
}
