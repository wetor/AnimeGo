package qbittorrent

const (
	QbtError              = "error"              // Some error occurred, applies to paused torrents
	QbtMissingFiles       = "missingFiles"       // Torrent data files is missing
	QbtUploading          = "uploading"          // Torrent is being seeded and data is being transferred
	QbtPausedUP           = "pausedUP"           // Torrent is paused and has finished downloading
	QbtQueuedUP           = "queuedUP"           // Queuing is enabled and torrent is queued for upload
	QbtStalledUP          = "stalledUP"          // Torrent is being seeded, but no connection were made
	QbtCheckingUP         = "checkingUP"         // Torrent has finished downloading and is being checked
	QbtForcedUP           = "forcedUP"           // Torrent is forced to uploading and ignore queue limit
	QbtAllocating         = "allocating"         // Torrent is allocating disk space for download
	QbtDownloading        = "downloading"        // Torrent is being downloaded and data is being transferred
	QbtMetaDL             = "metaDL"             // Torrent has just started downloading and is fetching metadata
	QbtPausedDL           = "pausedDL"           // Torrent is paused and has NOT finished downloading
	QbtQueuedDL           = "queuedDL"           // Queuing is enabled and torrent is queued for download
	QbtStalledDL          = "stalledDL"          // Torrent is being downloaded, but no connection were made
	QbtCheckingDL         = "checkingDL"         // Same as checkingUP, but torrent has NOT finished downloading
	QbtForcedDL           = "forcedDL"           // Torrent is forced to downloading to ignore queue limit
	QbtCheckingResumeData = "checkingResumeData" // Checking resume data on qBt startup
	QbtMoving             = "moving"             // Torrent is moving to another location
	QbtUnknown            = "unknown"            // Unknown status
)
