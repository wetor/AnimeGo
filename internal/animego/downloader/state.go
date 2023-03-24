package downloader

import (
	"github.com/wetor/AnimeGo/internal/models"
)

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

const (
	StateUnknown     models.TorrentState = "unknown"     //未知
	StateWaiting     models.TorrentState = "waiting"     // 等待
	StateDownloading models.TorrentState = "downloading" // 下载中
	StatePausing     models.TorrentState = "pausing"     // 暂停中
	StateMoving      models.TorrentState = "moving"      // 移动中
	StateSeeding     models.TorrentState = "seeding"     // 做种中
	StateComplete    models.TorrentState = "complete"    // 完成下载
	StateError       models.TorrentState = "error"       // 错误
	StateNotFound    models.TorrentState = "notfound"    // 不存在
	StateAdding      models.TorrentState = "adding"      // 添加下载项
)

// StateMap
//
//	@Description: 下载器状态转换
//	@param clientState string
//	@return models.TorrentState
func StateMap(clientState string) models.TorrentState {
	switch clientState {
	case QbtAllocating, QbtMetaDL, QbtStalledDL,
		QbtCheckingDL, QbtCheckingResumeData, QbtQueuedDL,
		QbtForcedUP, QbtQueuedUP:
		// 若进度为100，则下载完成
		return StateWaiting
	case QbtDownloading, QbtForcedDL:
		return StateDownloading
	case QbtMoving:
		return StateMoving
	case QbtUploading, QbtStalledUP:
		// 已下载完成
		return StateSeeding
	case QbtPausedDL:
		return StatePausing
	case QbtPausedUP, QbtCheckingUP:
		// 已下载完成
		return StateComplete
	case QbtError, QbtMissingFiles:
		return StateError
	case QbtUnknown:
		return StateUnknown
	default:
		return StateUnknown
	}
}
