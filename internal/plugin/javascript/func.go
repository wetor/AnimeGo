package javascript

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/wetor/AnimeGo/internal/animego/anidata/mikan"
	"github.com/wetor/AnimeGo/internal/animego/anisource"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin/public"
	"github.com/wetor/AnimeGo/pkg/errors"
)

// initFunc
//  @Description: 获取js需要注册的函数列表
//  @return models.Object
//
func (js JavaScript) initFunc() models.Object {
	return models.Object{
		"print": js.Print,
		"sleep": js.Sleep,
		"os": models.Object{
			"readFile": js.ReadFile,
		},
		"goLog": models.Object{
			"debug": zap.S().Debug,
			"info":  zap.S().Info,
			"error": zap.S().Error,
		},
		"animeGo": models.Object{
			"parseName":    js.ParseName,
			"getMikanInfo": js.GetMikanInfo,
		},
	}
}

func (js JavaScript) initVar() models.Object {
	return models.Object{
		"variable": models.Object{
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

func (js JavaScript) ParseName(name string) (episode *public.Episode) {
	episode = public.ParserName(name)
	if episode.Ep == 0 {
		panic(js.ToValue(errors.NewAniError("解析ep信息失败")))
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
