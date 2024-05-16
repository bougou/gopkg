package snippet

import "fmt"

type formatValue struct {
	format string
	value  interface{}
}

func fv(format string, value interface{}) formatValue {
	return formatValue{
		format: format,
		value:  value,
	}
}

func formatValuesTable(formatValues []formatValue) string {
	var format string
	var values []interface{}
	for k, v := range formatValues {
		if k == 0 {
			format += v.format
		} else {
			format += " | " + v.format
		}
		values = append(values, v.value)
	}
	return fmt.Sprintf(format, values...)
}
