package constant

type NotifyState int

const (
	NotifyOnInit     NotifyState = iota // animeGo
	NotifyOnStart                       // qBittorrent, aria2
	NotifyOnDownload                    // animeGo
	NotifyOnPause                       // qBittorrent, aria2
	NotifyOnStop                        // qBittorrent, aria2
	NotifyOnSeeding                     // qBittorrent, aria2
	NotifyOnComplete                    // qBittorrent, aria2
	NotifyOnError                       // qBittorrent, aria2
)
