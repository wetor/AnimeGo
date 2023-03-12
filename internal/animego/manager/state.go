package manager

import (
	"github.com/wetor/AnimeGo/internal/animego/downloader/qbittorrent"
	"github.com/wetor/AnimeGo/internal/models"
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

// stateMap
//
//	@Description: 下载器状态转换
//	@param clientState string
//	@return models.TorrentState
func stateMap(clientState string) models.TorrentState {
	switch clientState {
	case qbittorrent.QbtAllocating, qbittorrent.QbtMetaDL, qbittorrent.QbtStalledDL,
		qbittorrent.QbtCheckingDL, qbittorrent.QbtCheckingResumeData, qbittorrent.QbtQueuedDL,
		qbittorrent.QbtForcedUP, qbittorrent.QbtQueuedUP:
		// 若进度为100，则下载完成
		return StateWaiting
	case qbittorrent.QbtDownloading, qbittorrent.QbtForcedDL:
		return StateDownloading
	case qbittorrent.QbtMoving:
		return StateMoving
	case qbittorrent.QbtUploading, qbittorrent.QbtStalledUP:
		// 已下载完成
		return StateSeeding
	case qbittorrent.QbtPausedDL:
		return StatePausing
	case qbittorrent.QbtPausedUP, qbittorrent.QbtCheckingUP:
		// 已下载完成
		return StateComplete
	case qbittorrent.QbtError, qbittorrent.QbtMissingFiles:
		return StateError
	case qbittorrent.QbtUnknown:
		return StateUnknown
	default:
		return StateUnknown
	}
}
