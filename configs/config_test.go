package configs

import (
	"fmt"
	"github.com/wetor/AnimeGo/internal/models"
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	c := Init("../assets/config/animego.yaml")
	fmt.Println(c)
}

func TestSettings_Tag(t *testing.T) {
	setting := Setting{}
	setting.TagSrc = "{year}年{quarter}月新番,{quarter_index},{quarter_name}季,第{ep}集,周{week},{week_name}"
	str := setting.Tag(&models.AnimeEntity{
		AirDate: "2022-04-11",
		Ep:      10,
	})
	fmt.Println(str)
}

func TestDefaultConfig(t *testing.T) {
	os.Setenv("ANIMEGO_CONFIG_VERSION", "1.0.0")
	os.WriteFile("../assets/default.json", DefaultDoc(), 0666)
}

func TestUpdateConfig(t *testing.T) {
	os.Setenv("ANIMEGO_CONFIG_VERSION", "1.1.0")

	UpdateConfig("../data/animego.yaml")
}
