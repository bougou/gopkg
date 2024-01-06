package unit

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/bougou/gopkg/common"
)

var (
	unitMultiplier = map[string]float64{
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
		"E":  math.Pow(1000, 6),
		"Ei": math.Pow(1024, 6),
	}

	unitRank = []string{"", "Ki", "Mi", "Gi", "Ti", "Pi", "Ei"}
)

const valuePattern = `^(?P<Number>[0-9.,]+)\s*(?P<Unit>m||[KMGTP]i?)$`

func ConvertStrToFloat64(s string) float64 {
	r := regexp.MustCompile(valuePattern)

	n := common.GetPatternCaptured(r, s, "Number")
	u := common.GetPatternCaptured(r, s, "Unit")

	n = strings.ReplaceAll(n, ",", "")
	f, _ := strconv.ParseFloat(n, 64)

	multiplier, ok := unitMultiplier[u]
	if !ok {
		return f
	}
	return f * multiplier
}

// ConvertFloat64ToStr stringiy float64 value to human friendly string
func ConvertFloat64ToStr(f float64) string {
	if f < 1 {
		return fmt.Sprintf("%dm", int64(f*1000))
	}

	var i int
	for i = range unitRank {
		if f < 1024 {
			break
		}
		f = f / 1024
	}
	var preferedUnit string = unitRank[i]

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
