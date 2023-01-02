package public

type Episode struct {
	TitleRaw   string
	Name       string
	NameCN     string
	NameEN     string
	NameJP     string
	Season     int
	SeasonRaw  string
	Ep         int
	Sub        string
	Group      string
	Definition string
	Source     string
}

type Options struct {
	PluginPath string
}
