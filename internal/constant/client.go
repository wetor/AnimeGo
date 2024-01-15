package constant

type TorrentState string

const (
	StateUnknown     TorrentState = "unknown"      //未知
	StateWaiting     TorrentState = "waiting"      // 等待
	StateDownloading TorrentState = "downloading"  // 下载中
	StatePausing     TorrentState = "pausing"      // 暂停中
	StateMoving      TorrentState = "moving"       // 移动中
	StateSeeding     TorrentState = "seeding"      // 做种中
	StateComplete    TorrentState = "complete"     // 完成下载
	StateError       TorrentState = "error"        // 错误
	StateNotFound    TorrentState = "notfound"     // 不存在
	StateAdding      TorrentState = "adding"       // 添加下载项
	StateInit        TorrentState = "initializing" // 初始化
)

const ChanRetryConnect = 1 // 重连消息
