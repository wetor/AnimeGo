package configs

import (
	"bytes"
	"os"
	"testing"
)

func TestUpdateConfig(t *testing.T) {
	t.Skip("Skipping update config")

	_ = os.Setenv("ANIMEGO_CONFIG_VERSION", "1.1.0")
	file, _ := os.ReadFile("data/animego_100.yaml")
	_ = os.WriteFile("data/animego.yaml", file, 0666)
	UpdateConfig("data/animego.yaml", false)

	want, _ := os.ReadFile("data/animego_110.yaml")
	got, _ := os.ReadFile("data/animego.yaml")
	if !bytes.Equal(got, want) {
		t.Errorf("UpdateConfig() = %s, want %s", got, want)
	}
	//_ = os.Remove("data/animego.yaml")
}
