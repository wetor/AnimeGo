package schedule

import (
	"fmt"
	"github.com/wetor/AnimeGo/test"
	"testing"
)

func TestNewSchedule(t *testing.T) {
	test.TestInit()
	s := NewSchedule()

	for _, ts := range s.List() {
		fmt.Println(ts)
	}
	s.Start(nil)
	select {}
}
