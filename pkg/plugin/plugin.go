package plugin

var (
	Path  string
	Debug bool
)

type Options struct {
	Path  string
	Debug bool
}

func Init(opts *Options) {
	Path = opts.Path
	Debug = opts.Debug
}
