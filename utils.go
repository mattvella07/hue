package hue

import (
	"fmt"
	"reflect"
	"strings"
)

// formatSlice formats an int slice as a JSON string
func (h *Connection) formatSlice(sli []int) string {
	str := "["
	for _, s := range sli {
		str += fmt.Sprintf("\"%d\",", s)
	}
	str = str[:len(str)-1]
	str += "]"

	return str
}

// formatStruct formats a struct as a JSON string
func (h *Connection) formatStruct(data interface{}) string {
	str := "{"
	d := reflect.ValueOf(data)
	t := d.Type()

	for i := 0; i < d.NumField(); i++ {
		str += fmt.Sprintf("\"%s\": ", strings.ToLower(t.Field(i).Name))

		switch d.Field(i).Kind() {
		case reflect.String:
			str += fmt.Sprintf("\"%s\",", d.Field(i).Interface())
		// Check for slice and array
		case reflect.Struct:
			str += fmt.Sprintf("%s,", h.formatStruct(d.Field(i).Interface()))
		default:
			str += fmt.Sprintf("%v,", d.Field(i).Interface())
		}
	}
	str = str[:len(str)-1]
	str += "}"

	return str
}
