package utils

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestSettings_Tag(t *testing.T) {
	tag := "{year}年{quarter}月新番,{quarter_index},{quarter_name}季,第{ep}集,周{week},{week_name}"
	got := Tag(tag, "2022-04-11", 10)
	want := "2022年4月新番,2,春季,第10集,周1,星期一"
	if got != want {
		t.Errorf("Tag() = %v, want %v", got, want)
	}
}
