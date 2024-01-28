package timer_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/log"

	"github.com/wetor/AnimeGo/internal/pkg/timer"
)

var (
	bolt api.Cacher
	tm   *timer.Timer
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	log.Init(&log.Options{
		File:  "data/test.log",
		Debug: true,
	})
	bolt = cache.NewBolt()
	bolt.Open("data/bolt.db")
	m.Run()
	bolt.Close()
	_ = log.Close()
	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func initTest(clean bool) (*sync.WaitGroup, func()) {
	if clean {
		_ = os.RemoveAll("data")
	}
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	tm = timer.NewTimer(&timer.Options{
		Name:         "Test",
		Cache:        bolt,
		UpdateSecond: 1,
		RetryCount:   2,
		WG:           &wg,
	})
	tm.Start(ctx)
	return &wg, cancel
}

func TestNewTimer(t *testing.T) {
	wg, cancel := initTest(false)

	go func() {
		time.Sleep(8 * time.Second)
		cancel()
	}()

	_, _ = tm.AddTask(&timer.AddOptions{
		Name:     "Failed",
		Duration: 2,
		Loop:     true,
		Func: func() error {
			fmt.Println("[task] run fail 2")
			return fmt.Errorf("error2")
		},
	})
	_, _ = tm.AddTask(&timer.AddOptions{
		Name:     "Success",
		Duration: 2,
		Loop:     true,
		Func: func() error {
			fmt.Println("[task] run 2")
			return nil
		},
	})
	_, _ = tm.AddTask(&timer.AddOptions{
		Name:     "Once",
		Duration: 5,
		Func: func() error {
			fmt.Println("[task] run 5")
			return nil
		},
	})
	wg.Wait()

}

func TestTimerMarshal(t *testing.T) {
	wg, cancel := initTest(false)
	_ = wg
	_, _ = tm.AddTask(&timer.AddOptions{
		Name:     "Success",
		Duration: 5,
		Loop:     true,
		Func: func() error {
			fmt.Println("[task] run 2")
			return nil
		},
	})
	log.Info("Sleep 7")
	time.Sleep(7 * time.Second)
	log.Info("Stop")
	cancel()
	log.Info("Sleep 2")
	time.Sleep(2 * time.Second)
	wg, cancel = initTest(false)
	log.Info("Start")
	tm.RegisterTaskFuncs(map[string]timer.TaskFunc{
		"Success": func() error {
			fmt.Println("[task] run 2")
			return nil
		},
	})
	time.Sleep(9 * time.Second)
	cancel()
}
