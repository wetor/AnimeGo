package configs_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wetor/AnimeGo/configs"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/test"
)

const testdata = "config"

func TestMain(m *testing.M) {
	_ = utils.CreateMutiDir("data")
	log.Init(&log.Options{
		File: "data/log.log",
	})
	m.Run()
	log.Close()
	//_ = os.RemoveAll("data")
}

func TestDefaultConfig(t *testing.T) {
	_ = configs.DefaultFile("data/animego_default.yaml")
}

func EqualFile(t *testing.T, file1, file2 string) {
	want, _ := os.ReadFile(file1)
	got, _ := os.ReadFile(file2)
	want = bytes.ReplaceAll(want, []byte("\r\n"), []byte("\n"))
	got = bytes.ReplaceAll(got, []byte("\r\n"), []byte("\n"))
	assert.Equal(t, string(got), string(want))
}

func TestUpdateConfig_120(t *testing.T) {
	configs.ConfigVersion = "1.2.0"

	file, _ := test.GetData(testdata, "animego_110.yaml")
	_ = os.WriteFile("data/animego.yaml", file, 0666)
	configs.UpdateConfig("data/animego.yaml", false)

	EqualFile(t, "data/animego.yaml", test.GetDataPath(testdata, "animego_120.yaml"))
}

func TestUpdateConfig_130(t *testing.T) {
	configs.ConfigVersion = "1.3.0"
	file, _ := test.GetData(testdata, "animego_120.yaml")
	_ = os.WriteFile("data/animego.yaml", file, 0666)
	configs.UpdateConfig("data/animego.yaml", false)

	EqualFile(t, "data/animego.yaml", test.GetDataPath(testdata, "animego_130.yaml"))
}

func TestUpdateConfig_140(t *testing.T) {
	configs.ConfigVersion = "1.4.0"
	file, _ := test.GetData(testdata, "animego_130.yaml")
	_ = os.WriteFile("data/animego.yaml", file, 0666)
	configs.UpdateConfig("data/animego.yaml", false)

	EqualFile(t, "data/animego.yaml", test.GetDataPath(testdata, "animego_140.yaml"))
}

func TestUpdateConfig_141(t *testing.T) {
	configs.ConfigVersion = "1.4.1"
	file, _ := test.GetData(testdata, "animego_140.yaml")
	_ = os.WriteFile("data/animego.yaml", file, 0666)
	configs.UpdateConfig("data/animego.yaml", false)

	EqualFile(t, "data/animego.yaml", test.GetDataPath(testdata, "animego_141.yaml"))
}

func TestUpdateConfig_150(t *testing.T) {
	configs.ConfigVersion = "1.5.0"
	file, _ := test.GetData(testdata, "animego_141.yaml")
	_ = os.WriteFile("data/animego.yaml", file, 0666)
	configs.UpdateConfig("data/animego.yaml", false)

	EqualFile(t, "data/animego.yaml", test.GetDataPath(testdata, "animego_150.yaml"))
}

func TestUpdateConfig_151(t *testing.T) {
	configs.ConfigVersion = "1.5.1"
	file, _ := test.GetData(testdata, "animego_150.yaml")
	_ = os.WriteFile("data/animego.yaml", file, 0666)
	configs.UpdateConfig("data/animego.yaml", false)

	EqualFile(t, "data/animego.yaml", test.GetDataPath(testdata, "animego_151.yaml"))
}

func TestUpdateConfig_152(t *testing.T) {
	configs.ConfigVersion = "1.5.2"
	file, _ := test.GetData(testdata, "animego_151.yaml")
	_ = os.WriteFile("data/animego.yaml", file, 0666)
	configs.UpdateConfig("data/animego.yaml", false)

	EqualFile(t, "data/animego.yaml", test.GetDataPath(testdata, "animego_152.yaml"))
}

func TestUpdateConfig_160(t *testing.T) {
	configs.ConfigVersion = "1.6.0"
	file, _ := test.GetData(testdata, "animego_152.yaml")
	_ = os.WriteFile("data/animego.yaml", file, 0666)
	configs.UpdateConfig("data/animego.yaml", false)

	EqualFile(t, "data/animego.yaml", test.GetDataPath(testdata, "animego_160.yaml"))
}

func TestUpdateConfig_170(t *testing.T) {
	configs.ConfigVersion = "1.7.0"
	file, _ := test.GetData(testdata, "animego_160.yaml")
	_ = os.WriteFile("data/animego.yaml", file, 0666)
	configs.UpdateConfig("data/animego.yaml", false)

	EqualFile(t, "data/animego.yaml", test.GetDataPath(testdata, "animego_170.yaml"))
}

func TestInitEnvConfig(t *testing.T) {
	_ = os.Setenv("ANIMEGO_QBT_URL", "http://127.0.0.1:18080")
	_ = os.Setenv("ANIMEGO_QBT_DOWNLOAD_PATH", "7766/download")
	_ = os.Setenv("ANIMEGO_DOWNLOAD_PATH", "aw8da/test/download")
	_ = os.Setenv("ANIMEGO_WEB_PORT", "10086")
	f := test.GetDataPath(testdata, "animego_152.yaml")
	configs.InitEnvConfig(f, "data/animego.yaml")

	conf := configs.Load("data/animego.yaml")

	assert.Equal(t, conf.Setting.Client.QBittorrent.Url, "http://127.0.0.1:18080")
	assert.Equal(t, conf.Setting.Client.QBittorrent.DownloadPath, "7766/download")
	assert.Equal(t, conf.Setting.DownloadPath, "aw8da/test/download")
	assert.Equal(t, conf.Setting.WebApi.Port, 10086)
}
