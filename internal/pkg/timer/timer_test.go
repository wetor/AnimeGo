package timer_test

import (
	"context"
	"fmt"
	"github.com/wetor/AnimeGo/pkg/log"
	"sync"
	"testing"
	"time"

	"github.com/wetor/AnimeGo/internal/pkg/timer"
)

func TestNewTimer(t *testing.T) {
	log.Init(&log.Options{
		Debug: true,
	})
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(8 * time.Second)
		cancel()
	}()

	tm := timer.NewTimer(&timer.Options{
		UpdateSecond: 1,
		RetryCount:   2,
		WG:           &wg,
		Ctx:          ctx,
	})

	tm.Start()
	tm.AddTask(&timer.AddOptions{
		Name:     "Failed",
		Duration: 2,
		Loop:     true,
		Func: func() error {
			fmt.Println("[task] run fail 2")
			return fmt.Errorf("error2")
		},
	})
	tm.AddTask(&timer.AddOptions{
		Name:     "Success",
		Duration: 2,
		Loop:     true,
		Func: func() error {
			fmt.Println("[task] run 2")
			return nil
		},
	})
	tm.AddTask(&timer.AddOptions{
		Name:     "Once",
		Duration: 5,
		Func: func() error {
			fmt.Println("[task] run 5")
			return nil
		},
	})
	wg.Wait()

}
