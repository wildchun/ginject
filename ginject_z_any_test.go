// Package ginject
//
//	@Description
//	@Author  hd_0411_qxc  2024/3/8 14:02
//	@Update  hd_0411_qxc  2024/3/8 14:02
package ginject

import (
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

type Object struct {
	appName string `autowire:"app.name" default:"app_default"`
	version string `autowire:"app.version" default:"1.0.0"`
	Mqtt    struct {
		Broker string `autowire:"app.mqtt.broker" default:"tcp://10.147.198.110:1883"`
	}
	Number struct {
		number1 int   `autowire:"app.number" default:"1"`
		number2 int8  `autowire:"app.number" default:"2"`
		number3 int16 `autowire:"app.number" default:"3"`
		number4 int32 `autowire:"app.number" default:"4"`
		number5 int64 `autowire:"app.number" default:"5"`
	}
	list     []int      `autowire:"list"`
	listlist [][]string `autowire:"listlist"`
}

func TestAutoWireSkipUnexport(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c := NewAutoWire()
		obj := &Object{}
		t.Assert(c.AutoWire(obj, &AutoWireOptions{
			SkipUnExported:   true,
			ErrorOnUnmatched: true,
		}), nil)
		g.Dump(obj)
	})
}
func TestAutoWireWithUnexport(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c := NewAutoWire()
		obj := &Object{}
		t.Assert(c.AutoWire(obj, &AutoWireOptions{
			SkipUnExported:   false,
			ErrorOnUnmatched: true,
		}), nil)
		g.Dump(obj)
	})
}
