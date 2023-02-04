package javascript

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/wetor/AnimeGo/internal/animego/anidata/mikan"
	"github.com/wetor/AnimeGo/internal/animego/anisource"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin/public"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/try"
)

// initFunc
//  @Description: 获取js需要注册的函数列表
//  @return models.Object
//
func (p JavaScript) initFunc() models.Object {
	return models.Object{
		"print": p.Print,
		"sleep": p.Sleep,
		"os": models.Object{
			"readFile": p.ReadFile,
		},
		"goLog": models.Object{
			"debug": log.Debug,
			"info":  log.Info,
			"warn":  log.Warn,
			"error": log.Error,
		},
		"animeGo": models.Object{
			"parseName":    p.ParseName,
			"getMikanInfo": p.GetMikanInfo,
		},
	}
}

func (p JavaScript) initVar() models.Object {
	return models.Object{
		"variable": models.Object{
			"version": os.Getenv("ANIMEGO_VERSION"),
			"name":    currName,
		},
	}
}

func (p JavaScript) Print(params ...any) {
	fmt.Println(params...)
}

func (p JavaScript) Sleep(ms int64) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func (p *JavaScript) ReadFile(filename string) string {
	if strings.Index(filename, "../") >= 0 {
		panic("禁止使用'../'访问路径")
	}
	file, err := os.ReadFile(path.Join(currRootPath, filename))
	if err != nil {
		panic(p.ToValue(err))
	}
	return string(file)
}

func (p JavaScript) ParseName(name string) (episode *models.TitleParsed) {
	episode = public.ParserName(name)
	if episode.Ep == 0 {
		panic(p.ToValue(errors.NewAniError("解析ep信息失败")))
	}
	return
}

func (p JavaScript) GetMikanInfo(url string) (info *mikan.MikanInfo) {
	try.This(func() {
		info = anisource.Mikan().CacheParseMikanInfo(url)
	}).Catch(func(err try.E) {
		panic(p.ToValue(err))
	})
	return
}
