// Package ginject
//
//	@Description
//	@Author  hd_0411_qxc  2024/3/8 17:43
//	@Update  hd_0411_qxc  2024/3/8 17:43
package ginject

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type AutoWire struct {
	adapter DataAdapter
}

// doWireStructInput is the input parameter for doAutoWireStruct.
type doWireStructInput struct {
	structValue reflect.Value
	prefix      string
	fieldPath   string
	opt         *AutoWireOptions
}

// doWireValueInput is the input parameter for doAutoWireValue.
type doWireValueInput struct {
	value         reflect.Value
	valuePath     string
	fieldPath     string
	defaultString *string
	opt           *AutoWireOptions
}

// doWireSliceInput is the input parameter for doAutoWireSlice.
type doWireSliceInput struct {
	value     reflect.Value
	valuePath string
	fieldPath string
	opt       *AutoWireOptions
}

// pathJoin joins the prefix and sub path with dot., like "prefix.sub"
func pathJoin(prefix, sub string) string {
	if prefix == "" {
		return sub
	}
	if sub == "" {
		return prefix
	}
	return prefix + "." + sub
}

// NewAutoWireWithAdapter creates and returns a new AutoWire object with given adapter.
func NewAutoWireWithAdapter(adapter DataAdapter) *AutoWire {
	return &AutoWire{
		adapter: adapter,
	}
}

// NewAutoWire creates and returns a new AutoWire object. the adapter is set to g.Cfg() by default.
func NewAutoWire() *AutoWire {
	return &AutoWire{
		adapter: g.Cfg(),
	}
}

// AutoWire automatically injects configuration data from adapter to the struct object.
// The parameter `data` should be a pointer of struct object. if `opt` is given, it overwrites the default options.
func (c *AutoWire) AutoWire(data interface{}, opt ...*AutoWireOptions) error {
	options := defaultOptions
	if len(opt) > 0 {
		options = opt[0]
	}
	structValue := reflect.ValueOf(data)
	if structValue.Kind() != reflect.Ptr || structValue.Elem().Kind() != reflect.Struct {
		return gerror.New("the data must be a valid pointer of struct")
	}
	structValue = structValue.Elem()
	if c.adapter == nil {
		return gerror.New("adapter is not set")
	}
	return c.doAutoWireStruct(doWireStructInput{
		structValue: structValue,
		prefix:      "",
		fieldPath:   "",
		opt:         options,
	})
}

// SetDataAdapter sets the data adapter for the AutoWire object.
func (c *AutoWire) SetDataAdapter(adapter DataAdapter) {
	c.adapter = adapter
}

// GetDataAdapter returns the data adapter of the AutoWire object.
func (c *AutoWire) GetDataAdapter() DataAdapter {
	return c.adapter
}

func (c *AutoWire) doAutoWireStruct(in doWireStructInput) error {
	structType := in.structValue.Type()
	for fiIdx := 0; fiIdx < structType.NumField(); fiIdx++ {
		field := structType.Field(fiIdx)
		if field.PkgPath != "" && in.opt.SkipUnExported {
			continue
		}
		fieldValue := in.structValue.Field(fiIdx)
		tag := field.Tag.Get("autowire")
		if (tag == "" || tag == "-") && fieldValue.Kind() != reflect.Struct {
			continue
		}
		newPrefix := pathJoin(in.prefix, tag)
		newFiPath := pathJoin(in.fieldPath, field.Name)
		if field.PkgPath != "" {
			fieldValue = reflect.NewAt(fieldValue.Type(), unsafe.Pointer(fieldValue.UnsafeAddr())).Elem()
		}
		var def *string = nil
		if defaultStr, ok := field.Tag.Lookup("default"); ok {
			def = &defaultStr
		}
		if err := c.doAutoWireValue(doWireValueInput{
			value:         fieldValue,
			valuePath:     newPrefix,
			fieldPath:     newFiPath,
			defaultString: def,
			opt:           in.opt,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (c *AutoWire) doAutoWireSlice(in doWireSliceInput) error {
	vArray := c.adapter.MustGet(nil, in.valuePath, nil)
	if vArray.IsNil() || !vArray.IsSlice() {
		if in.opt.ErrorOnUnmatched {
			return gerror.Newf("no path matched for field %v, path: %v,not slice", in.fieldPath, in.valuePath)
		}
		return nil
	}
	vLen := len(vArray.Slice())
	if vLen == 0 {
		return nil
	}
	newSlice := reflect.MakeSlice(in.value.Type(), vLen, vLen)
	for i := 0; i < vLen; i++ {
		newFiPath := in.valuePath + fmt.Sprintf("[%d]", i)
		newPrefix := pathJoin(in.valuePath, fmt.Sprintf("%d", i))
		if err := c.doAutoWireValue(doWireValueInput{
			value:         newSlice.Index(i),
			valuePath:     newPrefix,
			fieldPath:     newFiPath,
			defaultString: nil,
			opt:           in.opt,
		}); err != nil {
			return err
		}
	}
	in.value.Set(newSlice)
	return nil
}

func (c *AutoWire) doAutoWireValue(in doWireValueInput) error {
	if in.value.Kind() == reflect.Slice {
		return c.doAutoWireSlice(doWireSliceInput{
			value:     in.value,
			valuePath: in.valuePath,
			fieldPath: in.fieldPath,
			opt:       in.opt,
		})
	} else if in.value.Kind() == reflect.Struct {
		return c.doAutoWireStruct(doWireStructInput{
			structValue: in.value,
			prefix:      in.valuePath,
			fieldPath:   in.fieldPath,
			opt:         in.opt,
		})
	}
	// fmt.Println("[*] wire field:", in.fieldPath, "with value path:", in.valuePath)
	v := c.adapter.MustGet(nil, in.valuePath, nil)
	if v.IsNil() {
		if in.defaultString == nil {
			if in.opt.ErrorOnUnmatched {
				return gerror.Newf("no path matched for field %v, path: %v", in.valuePath, in.fieldPath)
			}
			return nil
		} else {
			v = gvar.New(*in.defaultString)
		}
	}
	switch in.value.Kind() {
	case reflect.String:
		in.value.SetString(v.String())
	case reflect.Bool:
		in.value.SetBool(v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		in.value.SetInt(v.Int64())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		in.value.SetUint(v.Uint64())
	default:
		if in.opt.ErrorOnUnmatched {
			return gerror.Newf("unsupported kind: %s", in.value.Kind().String())
		}
	}
	return nil
}
