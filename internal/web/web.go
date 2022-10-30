package web

type InitOptions struct {
	Debug bool
}

var (
	Debug bool
)

func Init(opt *InitOptions) {
	Debug = opt.Debug
}
