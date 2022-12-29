package javascript

import (
	"fmt"
	"github.com/wetor/AnimeGo/internal/animego/anidata/mikan"
	"github.com/wetor/AnimeGo/internal/animego/anisource"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/third_party/poketto"
	"go.uber.org/zap"
	"os"
	"path"
	"strings"
	"time"
)

// initFunc
//  @Description: 获取js需要注册的函数列表
//  @return plugin.Object
//
func (js JavaScript) initFunc() plugin.Object {
	return plugin.Object{
		"print": js.Print,
		"sleep": js.Sleep,
		"os": plugin.Object{
			"readFile": js.ReadFile,
		},
		"goLog": plugin.Object{
			"debug": zap.S().Debug,
			"info":  zap.S().Info,
			"error": zap.S().Error,
		},
		"animeGo": plugin.Object{
			"parseName":    js.ParseName,
			"getMikanInfo": js.GetMikanInfo,
		},
	}
}

func (js JavaScript) initVar() plugin.Object {
	return plugin.Object{
		"variable": plugin.Object{
			"version": os.Getenv("ANIMEGO_VERSION"),
			"name":    currName,
		},
	}
}

func (js JavaScript) Print(params ...any) {
	fmt.Println(params...)
}

func (js JavaScript) Sleep(ms int64) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func (js *JavaScript) ReadFile(filename string) string {
	if strings.Index(filename, "../") >= 0 {
		panic("禁止使用'../'访问路径")
	}
	file, err := os.ReadFile(path.Join(currRootPath, filename))
	if err != nil {
		panic(js.ToValue(err))
	}
	return string(file)
}

func (js JavaScript) ParseName(name string) (episode *poketto.Episode) {
	episode = poketto.NewEpisode(name)
	episode.TryParse()
	if episode.ParseErr != nil {
		panic(js.ToValue(episode.ParseErr))
	}
	return
}

func (js JavaScript) GetMikanInfo(url string) *mikan.MikanInfo {
	defer errors.HandleError(func(err error) {
		panic(js.ToValue(err))
	})
	info := anisource.Mikan().CacheParseMikanInfo(url)
	return info
}
