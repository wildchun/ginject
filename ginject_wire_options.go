// Package ginject
//
//	@Description
//	@Author  hd_0411_qxc  2024/3/8 13:56
//	@Update  hd_0411_qxc  2024/3/8 13:56
package ginject

type ApplyOptions struct {
	// SkipUnExported specifies whether to skip the un-exported fields when injecting.
	SkipUnExported bool
	// ExitOnError specifies whether to exit the process when error occurs.
	ExitOnError bool
	// ErrorOnUnmatched specifies whether to return error if there's unmatched field.
	ErrorOnUnmatched bool
	// ErrorOnTypeErr specifies whether to return error if there's type error.
	CreatIfPointer bool
}

var defaultOptions = &ApplyOptions{
	SkipUnExported:   false,
	ErrorOnUnmatched: true,
	CreatIfPointer:   true,
	ExitOnError:      true,
}
