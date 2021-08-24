package tree

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/bougou/gopkg/common"
)

func Test_GroupObjectsByPropValue(t *testing.T) {
	d := []map[string]interface{}{}
	if err := json.Unmarshal([]byte(mockData), &d); err != nil {
		t.Errorf("json Unmarshal failed, err: %s", err)
	}

	m, err := GroupObjectsByPropValue(d, "id", "propA")
	if err != nil {
		t.Errorf("GroupObjectsByPropValue failed, err: %s", err)
	}

	b, err := json.Marshal(m)
	if err != nil {
		t.Errorf("json marshal failed, err: %s", err)
	}

	fmt.Println("grouped by propA:", string(b))
}

func Test_Treefiy(t *testing.T) {
	d := []map[string]interface{}{}
	if err := json.Unmarshal([]byte(mockData), &d); err != nil {
		t.Errorf("json Unmarshal failed, err: %s", err)
	}

	tree, err := Treeify("root", d, "id", "propA", "propB")
	if err != nil {
		t.Errorf("Treeify failed, err: %s", err)
	}

	y, err := common.Encode("yaml", tree)
	if err != nil {
		t.Errorf("json marshal failed, err: %s", err)
	}

	fmt.Println("grouped by propA/propB:", string(y))

}
