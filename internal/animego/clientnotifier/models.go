package clientnotifier

type Callback struct {
	Renamed func(data any) error
}

type Options struct {
	DownloadPath string
	SavePath     string
	Rename       string
	Callback     *Callback
}
