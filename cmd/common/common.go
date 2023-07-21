package common

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/wetor/AnimeGo/configs"
	"github.com/wetor/AnimeGo/internal/constant"
	pkgLog "github.com/wetor/AnimeGo/pkg/log"
)

var (
	version   = "dev"
	buildTime = "dev"
)

func init() {
	err := os.Setenv("ANIMEGO_VERSION", fmt.Sprintf("%s-%s", version, buildTime))
	if err != nil {
		panic(err)
	}
}

func PrintInfo() {
	fmt.Println(`--------------------------------------------------
    ___            _                   ______     
   /   |   ____   (_)____ ___   ___   / ____/____ 
  / /| |  / __ \ / // __ \__ \ / _ \ / / __ / __ \
 / ___ | / / / // // / / / / //  __// /_/ // /_/ /
/_/  |_|/_/ /_//_//_/ /_/ /_/ \___/ \____/ \____/
    `)
	fmt.Printf("AnimeGo %s\n", os.Getenv("ANIMEGO_VERSION"))
	fmt.Printf("AnimeGo config v%s\n", configs.ConfigVersion)
	fmt.Printf("%s\n", constant.AnimeGoGithub)
	fmt.Println("--------------------------------------------------")
}

func RegisterExit(doExit func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		for s := range sigs {
			switch s {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT:
				pkgLog.Infof("收到退出信号: %v", s)
				doExit()
			default:
				pkgLog.Infof("收到其他信号: %v", s)
			}
		}
	}()
}
