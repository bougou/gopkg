package deepmerge

import (
	"fmt"
	"testing"

	"github.com/bougou/gopkg/deepcopy"
	"github.com/imdario/mergo"
	"github.com/niemeyer/pretty"
)

func Test_DeepMerge(t *testing.T) {

	m := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": 100,
				"d": 200,
				"e": 300,
			},
			"c": 10,
		},
	}

	n := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": 1000,
				"d": 2000,
				"f": 4000,
			},
			"c": []int{10, 20},
		},
	}

	dst, _ := deepcopy.Map(m)
	if err := mergo.Merge(&dst, n, mergo.WithOverride); err != nil {
		fmt.Println(err)
	}

	pretty.Printf("dst: %v\n", dst)
	pretty.Printf("m: %v\n", m)
	pretty.Printf("n: %v\n", n)
}
