package config

import (
	"fmt"
	"testing"
)

func TestNewConfig(t *testing.T) {
	c := NewConfig("/Users/wetor/GoProjects/GoBangumi/data/config/conf.yaml")
	fmt.Println(c)
}
