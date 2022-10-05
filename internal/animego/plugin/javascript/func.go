package javascript

import (
	"AnimeGo/internal/animego/anisource"
	"AnimeGo/pkg/anisource/mikan"
	"AnimeGo/third_party/poketto"
	"fmt"
	"go.uber.org/zap"
	"os"
	"time"
)

// initFunc
//  @Description: 获取js需要注册的函数列表
//  @return Object
//
func (js JavaScript) initFunc() Object {
	return Object{
		"print": js.Print,
		"sleep": js.Sleep,
		"os": Object{
			"readFile": js.ReadFile,
			"getPwd":   js.GetPwd,
		},
		"goLog": Object{
			"debug": zap.S().Debug,
			"info":  zap.S().Info,
			"error": zap.S().Error,
		},
		"animeGo": Object{
			"parseName":    js.ParseName,
			"getMikanInfo": js.GetMikanInfo,
			"test":         js.Test,
		},
	}
}

func (js JavaScript) Print(params ...interface{}) {
	fmt.Println(params...)
}

func (js JavaScript) Sleep(ms int64) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func (js JavaScript) ReadFile(filename string) string {
	file, err := os.ReadFile(filename)
	if err != nil {
		panic(js.ToValue(err))
	}
	return string(file)
}

func (js JavaScript) GetPwd() string {
	pwd, _ := os.Getwd()
	return pwd
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
	info, err := anisource.Mikan().CacheParseMikanInfo(url)
	if err != nil {
		panic(js.ToValue(err))
	}
	return info
}

func (js JavaScript) Test() {
	panic(js.ToValue("异常测试"))
}
