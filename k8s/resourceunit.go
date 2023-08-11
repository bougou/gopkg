package k8s

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
		// "B", "KB", "MB", ... CAN NOT be used as k8s resource requests and limits value.
	}

	k8sRank = []string{"", "Ki", "Mi", "Gi", "Ti", "Pi"}
)

// K8SResourcePattern is the regex pattern to match the requests and limits value string
// for k8s resource. Like: 128974848, 129e6, 129M,  128974848000m, 123Mi
// The valid value DOES NOT allow space between the number and the unit.
const K8SResourcePattern = `^(?P<Number>[0-9.]+)(?P<Unit>m||[KMGTP]i?)$`

var (
	k8sResourceValueMatcher *regexp.Regexp
)

func init() {
	k8sResourceValueMatcher = regexp.MustCompile(K8SResourcePattern)
}

// ConvertResourceValueStringToFloat64 convert the k8s resource value from
// human friendly string to float64 number of the base unit.
// The conversion is lostless.
func ConvertResourceValueStringToFloat64(s string) float64 {
	number := common.GetPatternCaptured(k8sResourceValueMatcher, s, "Number")
	unit := common.GetPatternCaptured(k8sResourceValueMatcher, s, "Unit")

	f, _ := strconv.ParseFloat(number, 64)

	multiplier, ok := k8sMultiplier[unit]
	if !ok {
		return f
	}
	return f * multiplier
}

// ToString stringiy float64 value to human friendly string.
//
// The conversion may be not exact, eg:
//
//	1024*1024*1024 + 1024*1024*1024*0.049 -> "1GB"
//	1024*1024*1024 + 1024*1024*1024*0.050 -> "2GB"
//	1024*1024*1024 + 1024*1024*1024*0.051 -> "2GB"
func ConvertResourceValueFloat64ToString(f float64) string {
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
