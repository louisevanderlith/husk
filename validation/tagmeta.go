package validation

import (
	"reflect"
	"strconv"
	"strings"
)

type tagMeta struct {
	Required bool
	Size     int
	Type     reflect.Kind
	PropName string
}

//GetMeta returns the `hsk` meta tags associated with a field.
func GetMeta(tag string, kind reflect.Kind) tagMeta {
	result := tagMeta{}
	parts := strings.Split(tag, ";")

	required := !strings.Contains(tag, "null")
	result.Required = required
	result.Type = kind

	hasSize := strings.Contains(tag, "size")

	if hasSize {
		rawSize := getFromTag(parts, "size")
		sSize := strings.Replace(strings.Replace(rawSize, "size(", "", -1), ")", "", -1)

		size, err := strconv.ParseInt(sSize, 10, 32)

		if err == nil {
			result.Size = int(size)
		}
	}

	return result
}

func getFromTag(list []string, name string) string {
	var result string
	for _, v := range list {
		if strings.Contains(v, name) {
			result = v
			break
		}
	}

	return result
}
