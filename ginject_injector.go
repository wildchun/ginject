// Package ginject
//
//	@Description
//	@Author  hd_0411_qxc  2024/2/23 13:13
//	@Update  hd_0411_qxc  2024/2/23 13:13
package ginject

import (
	"context"
	"reflect"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/util/gmeta"
)

type DataSource interface {
	MustGet(ctx context.Context, pattern string, def ...interface{}) *gvar.Var
}

type injector struct {
	prefix string
	data   DataSource
}

type Injector interface {
	Apply(d interface{}) error
	Sub(prefix string) Injector
}

func New(data DataSource) Injector {
	return &injector{
		prefix: "",
		data:   data,
	}
}

func NewWithPrefix(data DataSource, prefix string) Injector {
	return &injector{
		prefix: prefix,
		data:   data,
	}
}

func writeableReflectValue(d interface{}, f gstructs.Field) reflect.Value {
	if f.IsExported() {
		return f.Value
	}
	return GetPtrUnexportFiled(d, f.Name())
}

func (inj *injector) SetDataSource(data DataSource) {
	inj.data = data
}

func (inj *injector) DataSource() DataSource {
	return inj.data
}

func (inj *injector) Sub(prefix string) Injector {
	return &injector{
		prefix: joinPath(inj.prefix, prefix),
		data:   inj.data,
	}
}

func (inj *injector) Apply(d interface{}) error {
	defer func() {
		if r := recover(); r != nil {
			println(r)
			return
		}
	}()
	if inj.data == nil {
		return gerror.New("data source is nil")
	}
	v := reflect.ValueOf(d)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return gerror.New("need struct")
	}
	return inj.doApplyStruct(context.TODO(), d, inj.prefix)
}

func (inj *injector) doApplyStruct(ctx context.Context, d interface{}, prefix string) error {
	// reset prefix if the struct has gmeta.Meta
	prefix = joinPath(prefix, gmeta.Get(d, "prefix").String())
	allFields, _ := gstructs.Fields(gstructs.FieldsInput{
		Pointer:         d,
		RecursiveOption: 0,
	})
	for _, field := range allFields {
		var (
			bind string
			def  string
		)
		if !field.Value.CanAddr() {
			continue
		}
		// if the field is a struct, do recursive injection and skip the next steps
		if field.Kind() == reflect.Struct {
			_ = inj.doApplyStruct(ctx, GetPtrUnexportFiled(d, field.Name()).Addr().Interface(), prefix)
			continue
		}
		if bind = getFieldPath(field); bind == "" {
			// if the field has no bind tag, it does not need to be injected
			continue
		}
		// get the value from the data source
		bindPath := joinPath(prefix, bind)
		def = getFieldDefault(field)
		// get the reflect.Value of the field which is writeable
		rv := GetWriteAbleReflectValue(d, field)
		switch rv.Kind() {
		case reflect.String:
			val := inj.data.MustGet(ctx, bindPath, def).String()
			rv.SetString(val)
		case reflect.Bool:
			val := inj.data.MustGet(ctx, bindPath, def).Bool()
			rv.SetBool(val)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val := inj.data.MustGet(ctx, bindPath, def).Int64()
			rv.SetInt(val)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val := inj.data.MustGet(ctx, bindPath, def).Uint64()
			rv.SetUint(val)
		case reflect.Slice:
			val := inj.data.MustGet(ctx, bindPath, def)
			if !val.IsSlice() {
				continue
			}
			// get elem kind
			elemKind := reflect.New(field.Type().Elem()).Elem().Kind()
			switch elemKind {
			case reflect.String:
				rv.Set(reflect.ValueOf(val.Strings()))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				rv.Set(reflect.ValueOf(val.Ints()))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				rv.Set(reflect.ValueOf(val.Uints()))
			default:
				continue
			}
		default:
			continue
		}
	}
	return nil
}
