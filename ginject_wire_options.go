// Package ginject
//
//	@Description
//	@Author  hd_0411_qxc  2024/3/8 13:56
//	@Update  hd_0411_qxc  2024/3/8 13:56
package ginject

type AutoWireOptions struct {
	// SkipUnExported specifies whether to skip the un-exported fields when injecting.
	SkipUnExported bool
	// ErrorOnUnmatched specifies whether to return error if there's unmatched field.
	ErrorOnUnmatched bool
}

var defaultOptions = &AutoWireOptions{
	SkipUnExported:   false,
	ErrorOnUnmatched: true,
}
