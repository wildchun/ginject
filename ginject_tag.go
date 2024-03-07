// Package ginject
//
//	@Description
//	@Author  hd_0411_qxc  2024/2/26 9:54
//	@Update  hd_0411_qxc  2024/2/26 9:54
package ginject

import (
	"github.com/gogf/gf/v2/os/gstructs"
)

var pathTag = []string{"inject", "inj", "value", "bind", "json"}
var defTag = []string{"default", "def"}

type tagDesc struct {
	valid        bool
	path         *string
	defaultValue *string
}

func findTagValue(fi gstructs.Field, tagNames ...string) *string {
	for _, tag := range tagNames {
		if val, exist := fi.TagLookup(tag); exist {
			return &val
		}
	}
	return nil
}

func findTagDesc(fi gstructs.Field) *tagDesc {
	pathName := findTagValue(fi, pathTag...)
	if pathName == nil {
		return nil
	}
	defContent := findTagValue(fi, defTag...)
	return &tagDesc{
		valid:        true,
		path:         pathName,
		defaultValue: defContent,
	}
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

func (t *tagDesc) JoinPrefix(prefix string) {
	if t.path != nil {
		*t.path = joinPath(prefix, *t.path)
	}
}

func (t *tagDesc) WithSubPath(path string) string {
	if t.path != nil {
		return joinPath(*t.path, path)
	}
	return ""
}

func getInjectTagValue(field gstructs.Field) string {
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
