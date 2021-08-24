package tree

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/bougou/gopkg/common"
)

func Test_ObjectList(t *testing.T) {
	d := []map[string]interface{}{}
	if err := json.Unmarshal([]byte(mockData), &d); err != nil {
		t.Errorf("json Unmarshal failed, err: %s", err)
	}

	objects := ObjectList{}

	for _, v := range d {
		object := NewMapObject("id", v)
		objects = append(objects, object)
	}

	node, err := objects.Tree("root", "propA", "propB")
	if err != nil {
		t.Error(err)
	}

	b, err := common.Encode("yaml", node)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(b))
}
