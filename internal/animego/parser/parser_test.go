package parser_test

import (
	"fmt"
	"github.com/brahma-adshonor/gohook"
	"github.com/wetor/AnimeGo/internal/animego/anidata"
	"github.com/wetor/AnimeGo/internal/animego/anisource"
	"github.com/wetor/AnimeGo/internal/animego/anisource/mikan"
	"github.com/wetor/AnimeGo/internal/animego/parser"
	parserPlugin "github.com/wetor/AnimeGo/internal/animego/parser/plugin"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/json"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/request"
	"io"
	"net/url"
	"os"
	"path"
	"sync"
	"testing"
)

func HookGetWriter(uri string, w io.Writer) error {
	log.Infof("Mock HTTP GET %s", uri)
	id := path.Base(uri)
	jsonData, err := os.ReadFile(path.Join("testdata", id))
	if err != nil {
		return err
	}
	_, err = w.Write(jsonData)
	if err != nil {
		return err
	}
	return nil
}

func HookGet(uri string, body interface{}) error {
	log.Infof("Mock HTTP GET %s", uri)
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}
	id := u.Query().Get("with_text_query")
	if len(id) == 0 {
		id = path.Base(u.Path)
	}
	p := path.Join("testdata", id)

	jsonData, err := os.ReadFile(p)
	if err != nil {
		return err
	}
	_ = json.Unmarshal(jsonData, body)
	return nil
}

func TestMain(m *testing.M) {
	fmt.Println("begin")
	_ = os.RemoveAll("data")
	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
	})
	_ = gohook.Hook(request.GetWriter, HookGetWriter, nil)
	_ = gohook.Hook(request.Get, HookGet, nil)

	plugin.Init(&plugin.Options{
		Path:  "../../../assets/plugin",
		Debug: true,
	})

	b := cache.NewBolt()
	b.Open("data/bolt.db")
	anisource.Init(&anisource.Options{
		Options: &anidata.Options{
			Cache: b,
			CacheTime: map[string]int64{
				"mikan":      int64(7 * 24 * 60 * 60),
				"bangumi":    int64(3 * 24 * 60 * 60),
				"themoviedb": int64(14 * 24 * 60 * 60),
			},
		},
	})
	parser.Init(&parser.Options{
		TMDBFailSkip:           false,
		TMDBFailUseTitleSeason: true,
		TMDBFailUseFirstSeason: true,
	})
	bangumiCache := cache.NewBolt(true)
	bangumiCache.Open("../../../test/testdata/bolt_sub.bolt")
	bangumiCache.Add("bangumi_sub")
	mutex := sync.Mutex{}
	anidata.Init(&anidata.Options{
		Cache:            b,
		BangumiCache:     bangumiCache,
		BangumiCacheLock: &mutex,
	})
	request.Init(&request.Options{
		Proxy: "http://127.0.0.1:7890",
	})
	m.Run()

	bangumiCache.Close()
	//_ = os.RemoveAll("data")
	fmt.Println("end")
}

func TestParse(t *testing.T) {
	p := parserPlugin.NewParserPlugin(&models.Plugin{
		Enable: true,
		Type:   "builtin",
		File:   "builtin_parser.py",
	}, true)
	mgr := parser.NewManager(p, mikan.Mikan{})
	e := mgr.Parse(&models.ParseOptions{
		Title:      "[猎户不鸽压制] 万事屋斋藤先生转生异世界 / 斋藤先生无所不能 Benriya Saitou-san, Isekai ni Iku [01-12] [合集] [WebRip 1080p] [繁中内嵌] [H265 AAC] [2023年1月番] [4.8 GB]",
		TorrentUrl: "https://mikanani.me/Download/20230328/061af0fb9d93214b33179b040517cf9d858c2ffd.torrent",
		MikanUrl:   "https://mikanani.me/Home/Episode/061af0fb9d93214b33179b040517cf9d858c2ffd",
	})
	d, _ := json.Marshal(e)
	fmt.Println(string(d))

}
