package request

type InitOptions struct {
	Proxy     string // 使用代理
	Retry     int    // 额外重试次数，默认为不做重试
	RetryWait int    // 重试等待时间，最小3秒
	Timeout   int    // 超时时间，最小3秒
	Debug     bool
}

func (o *InitOptions) Default() {
	if o.RetryWait < 3 {
		o.RetryWait = 3
	}
	if o.Timeout < 3 {
		o.Timeout = 3
	}
}
