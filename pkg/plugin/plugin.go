package plugin

var (
	Path string
)

type Options struct {
	Path string
}

func Init(opts *Options) {
	Path = opts.Path
}
