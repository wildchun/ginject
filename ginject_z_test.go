// Package ginject
//
//	@Description
//	@Author  hd_0411_qxc  2024/2/23 13:40
//	@Update  hd_0411_qxc  2024/2/23 13:40
package ginject

import (
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gmeta"
)

type Object struct {
	gmeta.Meta `prefix:"app"`
	appName    string `inject:"name" def:"app_default"`
	version    string `inject:"version" def:"1.0.0"`
	mqtt       struct {
		gmeta.Meta `prefix:"mqtt"`
		broker     string `inject:"broker" def:"tcp://10.147.198.110:1883"`
		clientId   string `inject:"clientId" def:""`
		username   string `inject:"username" def:""`
		password   string `inject:"password" def:""`
	}
	number struct {
		number1 int   `inject:"number" def:"1"`
		number2 int8  `inject:"number" def:"2"`
		number3 int16 `inject:"number" def:"3"`
		number4 int32 `inject:"number" def:"4"`
		number5 int64 `inject:"number" def:"5"`
	}
}

func TestInject(t *testing.T) {
	obj := &Object{}
	inj := New(g.Cfg())
	if err := inj.Apply(obj); err != nil {
		t.Error("injector failed", err)
	}
	g.Dump(obj)
}

func TestInjectSlice(t *testing.T) {
	var withSlice struct {
		struts []struct {
			name string `inject:"name"`
			age  int    `inject:"age"`
			list []bool `inject:"list"`
		} `inject:"structs"`
	}
	inj := New(g.Cfg())
	if err := inj.Apply(&withSlice); err != nil {
		t.Error("injector failed", err)
	}
	g.Dump(&withSlice)
}

func TestInjectWithPointer(t *testing.T) {
	var obj struct {
		Str *string `inject:"app.name"`
	}
	inj := New(g.Cfg())
	if err := inj.Apply(&obj); err != nil {
		t.Error("injector failed", err)
	}
	g.Dump(obj)
}
