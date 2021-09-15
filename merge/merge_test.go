package merge

import (
	"fmt"
	"testing"

	"github.com/bougou/gopkg/common"
)

func Test_merge(t *testing.T) {
	tests := []struct {
		src map[string]interface{}
		dst map[string]interface{}
	}{

		{
			src: map[string]interface{}{
				"a": map[string]interface{}{
					"k": "hello",
				},
				"b": "some",
				"c": 12,
				"d": map[string]interface{}{
					"e": "test",
				},
			},

			dst: map[string]interface{}{
				"a": 100,
				"c": "hello",
				"d": map[string]interface{}{
					"e": "jump",
					"f": "hi",
				},
			},
		},
	}

	for _, tt := range tests {
		m := NewMap(tt.dst)
		m.Merge(tt.src, WithOverride)
		b, _ := common.Encode("json", tt.dst)
		fmt.Println(string(b))

		b1, _ := common.Encode("json", m.Value())
		fmt.Println(string(b1))

	}

}
