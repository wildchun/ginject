// Package ginject
//
//	@Description
//	@Author  hd_0411_qxc  2024/2/26 9:54
//	@Update  hd_0411_qxc  2024/2/26 9:54
package ginject

import (
	"github.com/gogf/gf/v2/os/gstructs"
)

var pathTag = [...]string{"inject", "inj", "value", "bind"}
var defTag = [...]string{"default", "def"}

func getFieldPath(field gstructs.Field) string {
	for _, tag := range pathTag {
		if val, ok := field.TagLookup(tag); ok {
			return val
		}
	}
	return ""
}

func getFieldDefault(field gstructs.Field) string {
	for _, tag := range defTag {
		if val, ok := field.TagLookup(tag); ok {
			return val
		}
	}
	return ""
}

func joinPath(prefix, sub string) string {
	if prefix == "" {
		return sub
	}
	if sub == "" {
		return prefix
	}
	return prefix + "." + sub
}
