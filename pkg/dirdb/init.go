package dirdb

var (
	DefaultExt    = ".json"
	DefaultDB  DB = &JsonDB{}
)

type Options struct {
	DefaultExt string
	DefaultDB  DB
}

func Init(opts *Options) {
	if len(opts.DefaultExt) != 0 {
		DefaultExt = opts.DefaultExt
	}
	if opts.DefaultDB != nil {
		DefaultDB = opts.DefaultDB
	}
}
