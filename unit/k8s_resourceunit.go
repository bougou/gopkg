package unit

import (
	"fmt"
	"math"
	"regexp"
	"strconv"

	"github.com/bougou/gopkg/common"
)

var (
	k8sMultiplier = map[string]float64{
		"m":  0.001,
		"":   math.Pow(1024, 0),
		"K":  math.Pow(1000, 1),
		"Ki": math.Pow(1024, 1),
		"M":  math.Pow(1000, 2),
		"Mi": math.Pow(1024, 2),
		"G":  math.Pow(1000, 3),
		"Gi": math.Pow(1024, 3),
		"T":  math.Pow(1000, 4),
		"Ti": math.Pow(1024, 4),
		"P":  math.Pow(1000, 5),
		"Pi": math.Pow(1024, 5),
	}

	k8sRank = []string{"", "Ki", "Mi", "Gi", "Ti", "Pi"}
)

const K8SResourcePattern = `^(?P<Number>[0-9.]+)(?P<Unit>m||[KMGTP]i?)$`

func ParseK8SResourceStrToFloat64(s string) float64 {
	r := regexp.MustCompile(K8SResourcePattern)

	n := common.GetPatternCaptured(r, s, "Number")
	u := common.GetPatternCaptured(r, s, "Unit")

	f, _ := strconv.ParseFloat(n, 64)

	multiplier, ok := k8sMultiplier[u]
	if !ok {
		return f
	}
	return f * multiplier
}

// ParseK8SResourceFloat64ToStr stringiy float64 value to human friendly string
func ParseK8SResourceFloat64ToStr(f float64) string {
	if f < 1 {
		return fmt.Sprintf("%dm", int64(f*1000))
	}

	var i int
	for i = range k8sRank {
		if f < 1024 {
			break
		}
		f = f / 1024
	}
	var preferedUnit string = k8sRank[i]

	var res int64

	// 向下取整
	f1 := float64(int64(f))

	if math.Abs(f-f1) <= 0.05 {
		res = int64(f1)
	} else {
		res = int64(math.Ceil(f * 10 / 10))
	}

	return fmt.Sprintf("%d%s", res, preferedUnit)
}
