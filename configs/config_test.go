package configs

import (
	"fmt"
	"github.com/wetor/AnimeGo/internal/models"
	"testing"
)

func TestNewConfig(t *testing.T) {
	c := Init("/Users/wetor/GoProjects/AnimeGo/data/config/conf.yaml")
	fmt.Println(c)
}

func TestSettings_Tag(t *testing.T) {
	setting := Setting{}
	setting.TagSrc = "{year}年{quarter}月新番,{quarter_index},{quarter_name}季,第{ep}集,周{week},{week_name}"
	str := setting.Tag(&models.AnimeEntity{
		AirDate: "2022-04-11",
		AnimeEp: &models.AnimeEp{
			Ep: 10,
		},
	})
	fmt.Println(str)
}
