// Package ginject
//
//	@Description
//	@Author  hd_0411_qxc  2024/2/23 13:13
//	@Update  hd_0411_qxc  2024/2/23 13:13
package ginject

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"reflect"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/util/gconv"
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
	Apply(d interface{}, opt ...Options) error
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

func (inj *injector) Apply(d interface{}, opt ...Options) error {
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
		return gerror.New("injecting object should be a struct or a pointer to struct")
	}
	option := defaultOptions
	if len(opt) > 0 {
		option = &opt[0]
	}
	return inj.doApplyStruct(context.TODO(), d, inj.prefix, option)
}

func (inj *injector) doApplyStruct(ctx context.Context, d interface{}, prefix string, opt *Options) error {
	// reset prefix if the struct has gmeta.Meta
	prefix = joinPath(prefix, gmeta.Get(d, "prefix").String())
	structFields, _ := gstructs.Fields(gstructs.FieldsInput{
		Pointer:         d,
		RecursiveOption: 0,
	})
	for _, field := range structFields {
		// if the field is a struct, do recursive injection and skip the next steps
		if field.Kind() == reflect.Struct {
			_ = inj.doApplyStruct(ctx, GetPtrUnexportFiled(d, field.Name()).Addr().Interface(), prefix, opt)
			continue
		}

		tag := findTagDesc(field)
		if tag == nil {
			// if the field has no bind tag, it does not need to be injected
			continue
		}

		// get the value from the data source
		tag.JoinPrefix(prefix)

		//defValueStr = getFieldDefault(field)

		// get the reflect.Value of the field which is writeable
		// if the field is a pointer, create a new value and set it to the field
		rv := GetWriteAbleReflectValue(d, field)
		if field.Type().Kind() == reflect.Ptr {
			if !opt.CreatIfPointer {
				continue
			}
			val := reflect.New(field.Type().Elem())
			rv.Set(val)
			rv = val.Elem()
		}

		switch rv.Kind() {
		case reflect.String,
			reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			if err := inj.doApplyBaseType(rv, tag, opt); err != nil {
				return err
			}
		case reflect.Slice:
			inj.doApplySlice(ctx, rv, tag, opt)
		default:
			continue
		}
	}
	return nil
}

func (inj *injector) doApplySlice(ctx context.Context, v reflect.Value, tag *tagDesc, opt *Options) {
	val := inj.data.MustGet(ctx, *tag.path, nil)
	if !val.IsSlice() {
		return
	}
	// get elem kind
	elemKind := reflect.New(v.Type().Elem()).Elem().Kind()
	switch elemKind {
	case reflect.String,
		reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		vLen := len(val.Vars())
		newV := reflect.MakeSlice(v.Type(), vLen, vLen)
		for index := 0; index < vLen; index++ {
			idxPath := tag.WithSubPath(gconv.String(index))
			err := inj.doApplyBaseType(newV.Index(index), &tagDesc{
				valid:        true,
				path:         &idxPath,
				defaultValue: nil,
			}, opt)
			if err != nil {
				return
			}
		}
		v.Set(newV)
	case reflect.Struct:
		vLen := len(val.Vars())
		newV := reflect.MakeSlice(v.Type(), vLen, vLen)
		for index := 0; index < vLen; index++ {
			idxPath := tag.WithSubPath(gconv.String(index))
			_ = inj.doApplyStruct(ctx, newV.Index(index).Addr().Interface(), idxPath, opt)
		}
		v.Set(newV)
	default:
		return
	}
}

func (inj *injector) doApplyBaseType(value reflect.Value, tag *tagDesc, opt *Options) error {
	v := inj.data.MustGet(context.Background(), *tag.path, nil)

	if v.IsNil() {
		if tag.defaultValue != nil {
			v = gvar.New(*tag.defaultValue)
		} else if opt.ErrorOnUnmatched {
			return gerror.New("value not found")
		} else {
			return nil
		}
	}
	switch value.Kind() {
	case reflect.String:
		value.SetString(v.String())
	case reflect.Bool:
		value.SetBool(v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value.SetInt(v.Int64())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value.SetUint(v.Uint64())
	case reflect.Float32, reflect.Float64:
		value.SetFloat(v.Float64())
	default:
		if opt.ErrorOnTypeErr {
			return gerror.Newf("unsupported type:%v", value.Kind())
		} else {
			return nil
		}
	}
	return nil
}
