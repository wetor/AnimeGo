package configs_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wetor/AnimeGo/configs"
	"github.com/wetor/AnimeGo/pkg/utils"
)

func TestMain(m *testing.M) {
	_ = utils.CreateMutiDir("data")
	m.Run()
	_ = os.RemoveAll("data")
}

func TestDefaultConfig(t *testing.T) {
	_ = configs.DefaultFile("data/animego_default.yaml")
}

func TestUpdateConfig_120(t *testing.T) {
	configs.ConfigVersion = "1.2.0"
	file, _ := os.ReadFile("testdata/animego_110.yaml")
	_ = os.WriteFile("data/animego.yaml", file, 0666)
	configs.UpdateConfig("data/animego.yaml", false)

	want, _ := os.ReadFile("testdata/animego_120.yaml")
	got, _ := os.ReadFile("data/animego.yaml")
	want = bytes.ReplaceAll(want, []byte("\r\n"), []byte("\n"))
	got = bytes.ReplaceAll(got, []byte("\r\n"), []byte("\n"))
	assert.Equal(t, string(got), string(want))
}

func TestUpdateConfig_130(t *testing.T) {
	configs.ConfigVersion = "1.3.0"
	file, _ := os.ReadFile("testdata/animego_120.yaml")
	_ = os.WriteFile("data/animego.yaml", file, 0666)
	configs.UpdateConfig("data/animego.yaml", false)

	want, _ := os.ReadFile("testdata/animego_130.yaml")
	got, _ := os.ReadFile("data/animego.yaml")
	want = bytes.ReplaceAll(want, []byte("\r\n"), []byte("\n"))
	got = bytes.ReplaceAll(got, []byte("\r\n"), []byte("\n"))
	assert.Equal(t, string(got), string(want))
}

func TestUpdateConfig_140(t *testing.T) {
	configs.ConfigVersion = "1.4.0"
	file, _ := os.ReadFile("testdata/animego_130.yaml")
	_ = os.WriteFile("data/animego.yaml", file, 0666)
	configs.UpdateConfig("data/animego.yaml", false)

	want, _ := os.ReadFile("testdata/animego_140.yaml")
	got, _ := os.ReadFile("data/animego.yaml")
	want = bytes.ReplaceAll(want, []byte("\r\n"), []byte("\n"))
	got = bytes.ReplaceAll(got, []byte("\r\n"), []byte("\n"))
	assert.Equal(t, string(got), string(want))
}

func TestUpdateConfig_141(t *testing.T) {
	configs.ConfigVersion = "1.4.1"
	file, _ := os.ReadFile("testdata/animego_140.yaml")
	_ = os.WriteFile("data/animego.yaml", file, 0666)
	configs.UpdateConfig("data/animego.yaml", false)

	want, _ := os.ReadFile("testdata/animego_141.yaml")
	got, _ := os.ReadFile("data/animego.yaml")
	want = bytes.ReplaceAll(want, []byte("\r\n"), []byte("\n"))
	got = bytes.ReplaceAll(got, []byte("\r\n"), []byte("\n"))
	assert.Equal(t, string(got), string(want))
}
