package request

type Options struct {
	UserAgent string
	Proxy     string // 使用代理
	Retry     int    // 额外重试次数，默认为不做重试
	RetryWait int    // 重试等待时间，最小3秒
	Timeout   int    // 超时时间，最小3秒
	Debug     bool
	Host      map[string]*HostOptions
}

type HostOptions struct {
	Redirect string
	Header   map[string]string
	Params   map[string]string
	Cookie   map[string]string
}
