package configs_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wetor/AnimeGo/configs"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/test"
)

const testdata = "config"

func TestMain(m *testing.M) {
	_ = utils.CreateMutiDir("data")
	m.Run()
	_ = os.RemoveAll("data")
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

	file := test.GetData(testdata, "animego_110.yaml")
	_ = os.WriteFile("data/animego.yaml", file, 0666)
	configs.UpdateConfig("data/animego.yaml", false)

	EqualFile(t, "data/animego.yaml", test.GetDataPath(testdata, "animego_120.yaml"))
}

func TestUpdateConfig_130(t *testing.T) {
	configs.ConfigVersion = "1.3.0"
	file := test.GetData(testdata, "animego_120.yaml")
	_ = os.WriteFile("data/animego.yaml", file, 0666)
	configs.UpdateConfig("data/animego.yaml", false)

	EqualFile(t, "data/animego.yaml", test.GetDataPath(testdata, "animego_130.yaml"))
}

func TestUpdateConfig_140(t *testing.T) {
	configs.ConfigVersion = "1.4.0"
	file := test.GetData(testdata, "animego_130.yaml")
	_ = os.WriteFile("data/animego.yaml", file, 0666)
	configs.UpdateConfig("data/animego.yaml", false)

	EqualFile(t, "data/animego.yaml", test.GetDataPath(testdata, "animego_140.yaml"))
}

func TestUpdateConfig_141(t *testing.T) {
	configs.ConfigVersion = "1.4.1"
	file := test.GetData(testdata, "animego_140.yaml")
	_ = os.WriteFile("data/animego.yaml", file, 0666)
	configs.UpdateConfig("data/animego.yaml", false)

	EqualFile(t, "data/animego.yaml", test.GetDataPath(testdata, "animego_141.yaml"))
}

func TestUpdateConfig_151(t *testing.T) {
	configs.ConfigVersion = "1.5.0"
	file := test.GetData(testdata, "animego_141.yaml")
	_ = os.WriteFile("data/animego.yaml", file, 0666)
	configs.UpdateConfig("data/animego.yaml", false)

	EqualFile(t, "data/animego.yaml", test.GetDataPath(testdata, "animego_150.yaml"))
}
