package ginject

type Options struct {
	ErrorOnUnmatched bool
	ErrorOnTypeErr   bool
	CreatIfPointer   bool
}

var defaultOptions = &Options{
	ErrorOnUnmatched: true,
	ErrorOnTypeErr:   true,
}
