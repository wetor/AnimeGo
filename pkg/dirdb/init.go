package dirdb

var (
	DefaultExt    = map[string]struct{}{".json": {}}
	DefaultDB  DB = &JsonDB{}
)

type Options struct {
	DefaultExt []string
	DefaultDB  DB
}

func Init(opts *Options) {
	if len(opts.DefaultExt) != 0 {
		DefaultExt = make(map[string]struct{}, len(opts.DefaultExt))
		for _, ext := range opts.DefaultExt {
			DefaultExt[ext] = struct{}{}
		}
	}
	if opts.DefaultDB != nil {
		DefaultDB = opts.DefaultDB
	}
}

func InExt(ext string) bool {
	_, has := DefaultExt[ext]
	return has
}
