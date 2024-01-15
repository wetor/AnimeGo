package transmission

const (
	// TorrentStatusStopped represents a stopped torrent
	TorrentStatusStopped = "stopped"
	// TorrentStatusCheckWait represents a torrent queued for files checking
	TorrentStatusCheckWait = "waiting to check files"
	// TorrentStatusCheck represents a torrent which files are currently checked
	TorrentStatusCheck = "checking files"
	// TorrentStatusDownloadWait represents a torrent queue to download
	TorrentStatusDownloadWait = "waiting to download"
	// TorrentStatusDownload represents a torrent currently downloading
	TorrentStatusDownload = "downloading"
	// TorrentStatusSeedWait represents a torrent queued to seed
	TorrentStatusSeedWait = "waiting to seed"
	// TorrentStatusSeed represents a torrent currently seeding
	TorrentStatusSeed = "seeding"
	// TorrentStatusIsolated represents a torrent which can't find peers
	TorrentStatusIsolated = "can't find peers"

	TorrentStatusComplete = "complete stopped"
)

const (
	IdleModeGlobal    int64 = iota // follow the global settings
	IdleModeSingle                 // override the global settings, seeding until a certain idle time
	IdleModeUnlimited              // override the global settings, seeding regardless of activity
)
