package animego

import (
	"context"
	"fmt"
	"github.com/wetor/AnimeGo/internal/store"
	"github.com/wetor/AnimeGo/test"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	test.TestInit()

	m.Run()
	fmt.Println("end")
}

func TestMikanProcess(t *testing.T) {

	m := NewAnimeGo()
	ctx, cancel := context.WithCancel(context.Background())
	m.Run(ctx)
	store.WG.Add(2)
	go func() {
		time.Sleep(5 * time.Second)
		cancel()
	}()

	store.WG.Wait()
}
