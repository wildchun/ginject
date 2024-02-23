// Package ginject
//
//	@Description
//	@Author  hd_0411_qxc  2024/2/23 13:17
//	@Update  hd_0411_qxc  2024/2/23 13:17
package ginject

import (
	"reflect"
	"unsafe"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gstructs"
)

// GetPtrUnexportFiled
//
// Get the pointer of the unexported field of the struct
func GetPtrUnexportFiled(s interface{}, filed string) reflect.Value {
	v := reflect.ValueOf(s).Elem().FieldByName(filed)
	return reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem()
}

// SetPtrUnexportFiled
//
// Set the pointer of the unexported field of the struct
// if the val‘s kind is different from the field's kind, return an error
func SetPtrUnexportFiled(s interface{}, filed string, val interface{}) error {
	v := GetPtrUnexportFiled(s, filed)
	rv := reflect.ValueOf(val)
	if v.Kind() != rv.Kind() {
		return gerror.Newf("invalid kind, expected kind: %v, got kind:%v", v.Kind(), rv.Kind())
	}
	v.Set(rv)
	return nil
}

func GetWriteAbleReflectValue(d interface{}, f gstructs.Field) reflect.Value {
	if f.IsExported() {
		return f.Value
	} else {
		return GetPtrUnexportFiled(d, f.Name())
	}

}
