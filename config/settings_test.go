package config

import (
	"GoBangumi/models"
	"fmt"
	"testing"
)

func TestSettings_Tag(t *testing.T) {
	setting := Settings{
		TagSrc: "{year}年{quarter}月新番,{quarter_index},{quarter_name}季,第{ep}集,周{week},{week_name}",
	}
	str := setting.Tag(&models.Bangumi{
		AirDate: "2022-04-11",
		BangumiEp: &models.BangumiEp{
			Ep: 10,
		},
	})
	fmt.Println(str)
}
