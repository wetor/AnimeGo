package configs_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/jinzhu/copier"

	"github.com/wetor/AnimeGo/configs"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/utils"
)

func TestDefaultConfig(t *testing.T) {
	_ = utils.CreateMutiDir("data")
	_ = configs.DefaultFile("data/animego_default.yaml")
}

func TestCopy(t *testing.T) {
	type s1 struct {
		A int
		B struct {
			C string
			D bool
		}
	}

	type s2 struct {
		C string
		B struct {
			A int
			D bool
			E string
		}
	}

	ss1 := s1{
		A: 10,
		B: struct {
			C string
			D bool
		}{C: "测试", D: true},
	}

	ss2 := s2{}
	err := copier.Copy(&ss2, &ss1)
	if err != nil {
		panic(err)
	}
	fmt.Println(ss1)
	fmt.Println(ss2)
}

func TestUpdateConfig_120(t *testing.T) {
	_ = utils.CreateMutiDir("data")
	configs.ConfigVersion = "1.2.0"
	file, _ := os.ReadFile("testdata/animego_110.yaml")
	_ = os.WriteFile("data/animego.yaml", file, 0666)
	configs.UpdateConfig("data/animego.yaml", false)

	want, _ := os.ReadFile("testdata/animego_120.yaml")
	got, _ := os.ReadFile("data/animego.yaml")
	if !bytes.Equal(got, want) {
		t.Errorf("UpdateConfig() = %s, want %s", got, want)
	}
	_ = os.Remove("data/animego.yaml")
}

func TestUpdateConfig_130(t *testing.T) {
	_ = utils.CreateMutiDir("data")
	configs.ConfigVersion = "1.3.0"
	file, _ := os.ReadFile("testdata/animego_120_2.yaml")
	_ = os.WriteFile("data/animego.yaml", file, 0666)
	configs.UpdateConfig("data/animego.yaml", false)

	want, _ := os.ReadFile("testdata/animego_130.yaml")
	got, _ := os.ReadFile("data/animego.yaml")
	if !bytes.Equal(got, want) {
		t.Errorf("UpdateConfig() = %s, want %s", got, want)
	}
	_ = os.Remove("data/animego.yaml")
}

func TestUpdateConfig_140(t *testing.T) {
	_ = utils.CreateMutiDir("data")
	configs.ConfigVersion = "1.4.0"
	file, _ := os.ReadFile("testdata/animego_130.yaml")
	_ = os.WriteFile("data/animego.yaml", file, 0666)
	configs.UpdateConfig("data/animego.yaml", false)

	want, _ := os.ReadFile("testdata/animego_140.yaml")
	got, _ := os.ReadFile("data/animego.yaml")
	if !bytes.Equal(got, want) {
		t.Errorf("UpdateConfig() = %s, want %s", got, want)
	}
	//_ = os.Remove("data/animego.yaml")
}

func TestConvertPluginInfo(t *testing.T) {

	plugins := configs.ConvertPluginInfo([]configs.PluginInfo{
		{
			Enable: true,
			Type:   "python",
			File:   "/aaa",
			Args: map[string]any{
				"Cron":  "1111",
				"Retry": 77,
			},
		},
		{
			Enable: true,
			Type:   "python",
			File:   "/aaa",
			Args: map[string]any{
				"Cron":  "777888",
				"Retry": 13,
			},
		},
	})
	fmt.Println(plugins)
}

func TestCopy1(t *testing.T) {
	plugins := []*configs.PluginInfo{
		{
			Enable: true,
			Type:   "python",
			File:   "/aaa",
			Args: map[string]any{
				"Cron":  "1111",
				"Retry": 77,
			},
		},
		{
			Enable: true,
			Type:   "python",
			File:   "/aaa",
			Args: map[string]any{
				"Cron":  "777888",
				"Retry": 13,
			},
		},
	}
	dst := make([]*models.Plugin, len(plugins))
	copier.Copy(&dst, &plugins)
	fmt.Println(dst)
}
