package manager

import (
	"GoBangumi/internal/downloader/qbittorent"
	"GoBangumi/internal/models"
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
)

// stateMap
//  @Description: 下载器状态转换
//  @param clientState string
//  @return models.TorrentState
//
func stateMap(clientState string) models.TorrentState {
	switch clientState {
	case qbittorent.QbtAllocating, qbittorent.QbtMetaDL, qbittorent.QbtStalledDL,
		qbittorent.QbtCheckingDL, qbittorent.QbtCheckingResumeData, qbittorent.QbtQueuedDL,
		qbittorent.QbtForcedUP, qbittorent.QbtQueuedUP:
		// 若进度为100，则下载完成
		return StateWaiting
	case qbittorent.QbtDownloading, qbittorent.QbtForcedDL:
		return StateDownloading
	case qbittorent.QbtMoving:
		return StateMoving
	case qbittorent.QbtUploading, qbittorent.QbtStalledUP:
		// 已下载完成
		return StateSeeding
	case qbittorent.QbtPausedDL:
		return StatePausing
	case qbittorent.QbtPausedUP, qbittorent.QbtCheckingUP:
		// 已下载完成
		return StateComplete
	case qbittorent.QbtError, qbittorent.QbtMissingFiles:
		return StateError
	case qbittorent.QbtUnknown:
		return StateUnknown
	default:
		return StateUnknown
	}
}
