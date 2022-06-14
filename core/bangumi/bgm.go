package bangumi

import (
	"GoBangumi/bgm/res"
	"GoBangumi/model"
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
	Info *model.Bangumi
}

func NewBgm() Bangumi {
	return &Bgm{}
}
func (b *Bgm) Parse(opt *model.BangumiParseOptions) *model.Bangumi {
	url_ := BangumiInfoApi(opt.ID)
	resp := &res.SubjectV0{}
	status, err := utils.ApiGet(url_, resp)
	if err != nil {
		glog.Errorln(err)
		return nil
	}
	if status != 200 {
		glog.Errorln("Status:", status)
		return nil
	}
	b.Info = &model.Bangumi{
		ID:     int(resp.ID),
		Name:   resp.NameCN,
		NameJp: resp.Name,
		Date:   *resp.Date,
		Eps:    int(resp.Eps),
	}
	return b.Info
}
