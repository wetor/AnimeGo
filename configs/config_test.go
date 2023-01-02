package configs

import (
	"bytes"
	"github.com/wetor/AnimeGo/internal/models"
	"os"
	"testing"
)

func TestSettings_Tag(t *testing.T) {
	setting := Setting{}
	setting.TagSrc = "{year}年{quarter}月新番,{quarter_index},{quarter_name}季,第{ep}集,周{week},{week_name}"
	got := setting.Tag(&models.AnimeEntity{
		AirDate: "2022-04-11",
		Ep:      10,
	})
	want := "2022年4月新番,2,春季,第10集,周1,星期一"
	if got != want {
		t.Errorf("Tag() = %v, want %v", got, want)
	}
}

func TestDefaultConfig(t *testing.T) {
	want, err := os.ReadFile("data/default.json")
	if err != nil {
		panic(err)
	}
	got := DefaultDoc()

	if bytes.Compare(got, want) != 0 {
		t.Errorf("DefaultDoc() = %s, want %s", got, want)
	}
}

func TestUpdateConfig(t *testing.T) {
	_ = os.Setenv("ANIMEGO_CONFIG_VERSION", "1.1.0")
	file, _ := os.ReadFile("data/animego_100.yaml")
	_ = os.WriteFile("data/animego.yaml", file, 0666)
	UpdateConfig("data/animego.yaml", false)

	want, _ := os.ReadFile("data/animego_110.yaml")
	got, _ := os.ReadFile("data/animego.yaml")
	if bytes.Compare(got, want) != 0 {
		t.Errorf("UpdateConfig() = %s, want %s", got, want)
	}
	_ = os.Remove("data/animego.yaml")
}
