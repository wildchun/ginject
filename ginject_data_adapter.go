// Package ginject
//
//	@Description
//	@Author  hd_0411_qxc  2024/3/8 17:12
//	@Update  hd_0411_qxc  2024/3/8 17:12
package ginject

import (
	"context"

	"github.com/gogf/gf/v2/container/gvar"
)

type DataAdapter interface {
	// MustGet acts as function Get, but it panics if error occurs.
	MustGet(ctx context.Context, pattern string, def ...interface{}) *gvar.Var
}

type DataAdapterWrapper struct {
	Second DataAdapter
	First  map[string]*gvar.Var
}

func (d *DataAdapterWrapper) MustGet(ctx context.Context, pattern string, def ...interface{}) *gvar.Var {
	if v, ok := d.First[pattern]; ok {
		return v
	}
	return d.Second.MustGet(ctx, pattern, def...)
}
