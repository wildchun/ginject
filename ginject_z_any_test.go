// Package ginject
//
//	@Description
//	@Author  hd_0411_qxc  2024/3/8 14:02
//	@Update  hd_0411_qxc  2024/3/8 14:02
package ginject

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gmeta"
)

type Object struct {
	gmeta.Meta `prefix:"app"`
	appName    string `autowire:"app.name" default:"app_default"`
	version    string `autowire:"app.version" default:"1.0.0"`
	mqtt       struct {
		broker string `autowire:"app.mqtt.broker" default:"tcp://10.147.198.110:1883"`
	}
	number struct {
		number1 int   `autowire:"app.number" default:"1"`
		number2 int8  `autowire:"app.number" default:"2"`
		number3 int16 `autowire:"app.number" default:"3"`
		number4 int32 `autowire:"app.number" default:"4"`
		number5 int64 `autowire:"app.number" default:"5"`
	}
	list []int `autowire:"list"`
}

func TestAny(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c := &AutoWire{
			adapter: g.Cfg(),
		}
		obj := &Object{}
		t.Assert(c.AutoWire(obj), nil)
		g.Dump(obj)
	})
}
func TestAny2(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var s = struct{ foo int }{100}
		var i int

		structValue := reflect.ValueOf(&s).Elem()
		fieldValue := structValue.Field(0)
		setValue := reflect.ValueOf(&i).Elem()

		// use of NewAt() method
		fieldValue = reflect.NewAt(fieldValue.Type(), unsafe.Pointer(fieldValue.UnsafeAddr())).Elem()
		setValue.Set(fieldValue)
		fieldValue.Set(setValue)

		fmt.Println(fieldValue)
	})
}
