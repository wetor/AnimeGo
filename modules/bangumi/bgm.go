package bangumi

import (
	"GoBangumi/config"
	"GoBangumi/models"
	"GoBangumi/models/bgm/res"
	"GoBangumi/utils"
	"fmt"
	"github.com/golang/glog"
)

const (
	BangumiBaseApi = "https://api.bgm.tv" // Bangumi 域名
)

var BangumiInfoApi = func(id int) string {
	return fmt.Sprintf("%s/v0/subjects/%d", BangumiBaseApi, id)
}

type Bgm struct {
	Info *models.Bangumi
}

func NewBgm() Bangumi {
	return &Bgm{}
}
func (b *Bgm) Parse(opt *models.BangumiParseOptions) *models.Bangumi {
	url_ := BangumiInfoApi(opt.ID)
	resp := &res.SubjectV0{}
	status, err := utils.ApiGet(url_, resp, config.Proxy())
	if err != nil {
		glog.Errorln(err)
		return nil
	}
	if status != 200 {
		glog.Errorln("Status:", status)
		return nil
	}
	b.Info = &models.Bangumi{
		ID:      int(resp.ID),
		Name:    resp.NameCN,
		NameJp:  resp.Name,
		AirDate: *resp.Date,
		Date:    *resp.Date, // TODO: ep播放日期
		Eps:     int(resp.Eps),
	}
	return b.Info
}
