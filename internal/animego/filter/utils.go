package filter

import (
	"bytes"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/request"
	"github.com/wetor/AnimeGo/pkg/utils"
)

const (
	MagnetType  = "magnet"
	TorrentType = "torrent"
)

func ParseFeedItem(item *models.FeedItem) error {
	if len(item.Download) == 0 {
		return errors.NewAniError("download 链接为空")
	}
	hash, isTorrent := utils.TorrentHash(item.Download)
	if len(hash) == 40 {
		item.DownloadType = TorrentType
		item.Hash = hash
	} else if isTorrent {
		item.DownloadType = TorrentType
		b := bytes.NewBuffer(nil)
		err := request.GetWriter(item.Download, b)
		if err != nil {
			return errors.NewAniErrorD(err)
		}
		hash, err = utils.TorrentHashReader(b)
		if err != nil {
			return errors.NewAniErrorD(err)
		}
		item.Hash = hash
	} else if hash = utils.MagnetHash(item.Download); len(hash) == 40 {
		item.DownloadType = MagnetType
		item.Hash = hash
	} else {
		return errors.NewAniErrorf("不支持的链接 %s", item.Download)
	}
	return nil
}
