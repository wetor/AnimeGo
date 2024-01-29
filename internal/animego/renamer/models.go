package renamer

import (
	"sync"

	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/models"
)

type Options struct {
	WG            *sync.WaitGroup
	RefreshSecond int
}

type RenameTask struct {
	// 只读
	Src            string // 原名
	Dst            string
	Mode           string
	StateChan      chan constant.TorrentState
	RenameCallback models.RenameCallback // 重命名完成后回调
	Result         *models.RenameResult

	// 读写
	Enable      bool
	RenameState int
	State       constant.TorrentState
	ErrCount    int
}

type RenameTaskGroup struct {
	Keys             []string
	RenameResult     *models.RenameAllResult
	CompleteCallback models.CompleteCallback // 完成重命名所有流程后回调
}
