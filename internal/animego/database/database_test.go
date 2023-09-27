package database_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wetor/AnimeGo/assets"
	"github.com/wetor/AnimeGo/internal/animego/database"
	"github.com/wetor/AnimeGo/internal/animego/renamer"
	renamerPlugin "github.com/wetor/AnimeGo/internal/animego/renamer/plugin"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/dirdb"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/test"
)

const (
	DownloadPath = "data/download"
	SavePath     = "data/save"
)

var (
	rename       *renamer.Manager
	renamePlugin *renamerPlugin.Rename
	db           *cache.Bolt
	out          *bytes.Buffer

	dbManager *database.Database
)

type DownloaderMock struct {
}

func (m *DownloaderMock) Delete(hash string) error {
	log.Infof("Delete %v", hash)
	return nil
}

func TestMain(m *testing.M) {
	fmt.Println("begin")
	_ = os.RemoveAll("data")
	_ = utils.CreateMutiDir(DownloadPath)
	_ = utils.CreateMutiDir(SavePath)
	out = bytes.NewBuffer(nil)
	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
		Out:   out,
	})
	wg := sync.WaitGroup{}
	plugin.Init(&plugin.Options{
		Path:  assets.TestPluginPath(),
		Debug: true,
	})
	dirdb.Init(&dirdb.Options{
		DefaultExt: []string{".a_json", ".s_json", ".e_json"}, // anime, season
	})
	database.Init(&database.Options{
		DownloaderConf: database.DownloaderConf{
			DownloadPath: DownloadPath,
			SavePath:     SavePath,
			Rename:       "link_delete",
		},
	})
	renamer.Init(&renamer.Options{
		WG:            &wg,
		RefreshSecond: 1,
	})
	renamePlugin = renamerPlugin.NewRenamePlugin(&models.Plugin{
		Enable: true,
		Type:   "python",
		File:   "rename/builtin_rename.py",
	})
	rename = renamer.NewManager(renamePlugin)
	db = cache.NewBolt()
	db.Open("data/test.db")
	var err error
	downloader := &DownloaderMock{}
	dbManager, err = database.NewDatabase(db, rename, &database.Callback{
		Renamed: func(data any) error {
			return downloader.Delete(data.(string))
		},
	})
	if err != nil {
		panic(err)
	}

	m.Run()
	db.Close()
	_ = log.Close()
	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func initTest(clean bool) (*sync.WaitGroup, func()) {
	if clean {
		_ = os.RemoveAll("data")
	}
	wg := sync.WaitGroup{}
	renamer.WG = &wg
	ctx, cancel := context.WithCancel(context.Background())
	rename.Init()
	rename.Start(ctx)
	dbManager.Init()
	_ = dbManager.Scan()
	return &wg, cancel
}

func AddItem(name string, season int, ep []int) (hash string) {
	hash = utils.MD5([]byte(fmt.Sprintf("%v%v%v", name, season, ep)))
	anime := &models.AnimeEntity{
		NameCN: name,
		Season: season,
		Torrent: &models.AnimeTorrent{
			Hash: hash,
		},
	}
	anime.Ep = make([]*models.AnimeEpEntity, 0, len(ep))
	for _, e := range ep {
		anime.Ep = append(anime.Ep, &models.AnimeEpEntity{
			Type: models.AnimeEpNormal,
			Ep:   e,
			Src:  fmt.Sprintf("%s/src_%d.mp4", name, e),
		})
	}
	err := dbManager.Add(anime)
	if err != nil {
		panic(err)
	}
	for _, e := range anime.Ep {
		f := path.Join(DownloadPath, e.Src)
		err := utils.CreateMutiDir(path.Dir(f))
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(f, []byte{}, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	return hash
}

func TestOnDownloadStart(t *testing.T) {
	out.Reset()
	wg, cancel := initTest(true)
	hash := AddItem("动画1", 2, []int{1, 2, 3})
	dbManager.OnDownloadStart([]models.ClientEvent{
		{Hash: hash},
	})
	fmt.Println(out.String())
	test.LogBatchCompare(out, test.MatchContainsRegexp,
		"OnDownloadStart",
		map[string]any{`\[Plugin\] Rename插件.*? src_\d.mp4`: 3},
	)
	out.Reset()

	dbManager.OnDownloadSeeding([]models.ClientEvent{
		{Hash: hash},
	})

	time.Sleep(1*time.Second + 500*time.Millisecond)
	fmt.Println(out.String())
	test.LogBatchCompare(out, test.MatchContainsRegexp,
		"OnDownloadSeeding",
		map[string]any{`write data/save/动画1/S02/E00\d.e_json`: 3},
		map[string]any{`\[重命名\] 链接「data/download/动画1/src_\d.mp4」`: 3},
	)
	out.Reset()
	dbManager.OnDownloadComplete([]models.ClientEvent{
		{Hash: hash},
	})
	go func() {
		time.Sleep(1*time.Second + 500*time.Millisecond)
		cancel()
	}()
	wg.Wait()
	assert.FileExists(t, path.Join(SavePath, "动画1", "anime.a_json"))
	fmt.Println(out.String())
	test.LogBatchCompare(out, test.MatchContainsRegexp,
		"OnDownloadComplete",
		map[string]any{`write data/save/动画1/S02/E00\d.e_json`: 3},
		map[string]any{`\[重命名\] 删除「data/download/动画1/src_\d.mp4」`: 3},
		"移动完成「动画1[第2季][1-3集]」",
		"write data/save/动画1/anime.a_json",
		"write data/save/动画1/S02/anime.s_json",
		map[string]any{`write data/save/动画1/S02/E00\d.e_json`: 3},
		"写入元数据文件「data/save/动画1/tvshow.nfo」",
		"刮削完成: 动画1[第2季][1-3集]",
		map[string]any{`write data/save/动画1/S02/E00\d.e_json`: 3},
		"Delete",
		"正常退出",
	)
}

// TestOnDownloadExistAnime
//
//	下载已存在的剧集
func TestOnDownloadExistAnime(t *testing.T) {
	out.Reset()
	wg, cancel := initTest(true)

	// 下载 1 2 3
	hash := AddItem("动画1", 2, []int{1, 2, 3})
	dbManager.OnDownloadStart([]models.ClientEvent{
		{Hash: hash},
	})
	dbManager.OnDownloadSeeding([]models.ClientEvent{
		{Hash: hash},
	})
	time.Sleep(1*time.Second + 500*time.Millisecond)
	dbManager.OnDownloadComplete([]models.ClientEvent{
		{Hash: hash},
	})
	time.Sleep(2 * time.Second)

	fmt.Println(out.String())
	test.LogBatchCompare(out, test.MatchContainsRegexp,
		"OnDownloadStart",
		map[string]any{`\[Plugin\] Rename插件.*? src_\d.mp4`: 3},
		"OnDownloadSeeding",
		map[string]any{`write data/save/动画1/S02/E00\d.e_json`: 3},
		map[string]any{`\[重命名\] 链接「data/download/动画1/src_\d.mp4」`: 3},
		"OnDownloadComplete",
		map[string]any{`write data/save/动画1/S02/E00\d.e_json`: 3},
		map[string]any{`\[重命名\] 删除「data/download/动画1/src_\d.mp4」`: 3},
		"移动完成「动画1[第2季][1-3集]」",
		"write data/save/动画1/anime.a_json",
		"write data/save/动画1/S02/anime.s_json",
		map[string]any{`write data/save/动画1/S02/E00\d.e_json`: 3},
		"写入元数据文件「data/save/动画1/tvshow.nfo」",
		"刮削完成: 动画1[第2季][1-3集]",
		map[string]any{`write data/save/动画1/S02/E00\d.e_json`: 3},
		"Delete",
	)
	out.Reset()

	// 下载3 4 5
	// 已下载3，跳过
	hash2 := AddItem("动画1", 2, []int{3, 4, 5})
	dbManager.OnDownloadStart([]models.ClientEvent{
		{Hash: hash2},
	})
	dbManager.OnDownloadSeeding([]models.ClientEvent{
		{Hash: hash2},
	})
	time.Sleep(1*time.Second + 500*time.Millisecond)
	dbManager.OnDownloadComplete([]models.ClientEvent{
		{Hash: hash2},
	})
	time.Sleep(2 * time.Second)

	fmt.Println(out.String())
	test.LogBatchCompare(out, test.MatchContainsRegexp,
		"OnDownloadStart",
		map[string]any{`\[Plugin\] Rename插件.*? src_\d.mp4`: 3},
		"发现部分已下载，跳过此部分重命名: data/download/动画1/src_3.mp4",
		"OnDownloadSeeding",
		map[string]any{`write data/save/动画1/S02/E00\d.e_json`: 2},
		map[string]any{`\[重命名\] 链接「data/download/动画1/src_\d.mp4」`: 2},
		"OnDownloadComplete",
		map[string]any{`write data/save/动画1/S02/E00\d.e_json`: 2},
		map[string]any{`\[重命名\] 删除「data/download/动画1/src_\d.mp4」`: 2},
		"移动完成「动画1[第2季][3-5集]」",
		"write data/save/动画1/anime.a_json",
		"write data/save/动画1/S02/anime.s_json",
		map[string]any{`write data/save/动画1/S02/E00\d.e_json`: 2},
		"写入元数据文件「data/save/动画1/tvshow.nfo」",
		"刮削完成: 动画1[第2季][3-5集]",
		map[string]any{`write data/save/动画1/S02/E00\d.e_json`: 2},
		"Delete",
	)
	out.Reset()

	// 下载 1 3 5
	// 全部已下载，跳过
	hash3 := AddItem("动画1", 2, []int{1, 3, 5})
	dbManager.OnDownloadStart([]models.ClientEvent{
		{Hash: hash3},
	})
	dbManager.OnDownloadSeeding([]models.ClientEvent{
		{Hash: hash3},
	})
	time.Sleep(1*time.Second + 500*time.Millisecond)
	dbManager.OnDownloadComplete([]models.ClientEvent{
		{Hash: hash3},
	})
	time.Sleep(2 * time.Second)

	fmt.Println(out.String())
	test.LogBatchCompare(out, test.MatchContainsRegexp,
		"OnDownloadStart",
		map[string]any{`\[Plugin\] Rename插件.*? src_\d.mp4`: 3},
		map[string]any{`发现部分已下载，跳过此部分重命名: .*?src_\d.mp4`: 3, "OnDownloadSeeding": 1},
		"OnDownloadComplete",
		map[string]any{`重命名任务不存在，可能已经完成`: 3},
	)
	out.Reset()
	// 结束
	go func() {
		time.Sleep(3 * time.Second)
		cancel()
	}()
	wg.Wait()

	fmt.Println(out.String())
	test.LogBatchCompare(out, nil, "正常退出")

	exist := dbManager.IsExist(&models.AnimeEntity{
		NameCN: "动画1",
		Season: 2,
		Ep: []*models.AnimeEpEntity{
			{Type: models.AnimeEpNormal, Ep: 1},
			{Type: models.AnimeEpNormal, Ep: 2},
			{Type: models.AnimeEpNormal, Ep: 3},
			{Type: models.AnimeEpNormal, Ep: 4},
			{Type: models.AnimeEpNormal, Ep: 5},
		},
	})
	assert.Equal(t, exist, true)

	exist = dbManager.IsExist(&models.AnimeEntity{
		NameCN: "动画1",
		Season: 2,
		Ep: []*models.AnimeEpEntity{
			{Type: models.AnimeEpNormal, Ep: 1},
			{Type: models.AnimeEpNormal, Ep: 2},
			{Type: models.AnimeEpNormal, Ep: 3},
			{Type: models.AnimeEpNormal, Ep: 6},
		},
	})
	assert.Equal(t, exist, false)

}

// TestOnDownloadRestartOnSeed
//
//	做种过程中重启
func TestOnDownloadRestartOnSeed(t *testing.T) {
	out.Reset()
	wg, cancel := initTest(true)

	// 下载 1 2 3
	hash := AddItem("动画1", 2, []int{1, 2, 3})
	dbManager.OnDownloadStart([]models.ClientEvent{
		{Hash: hash},
	})
	dbManager.OnDownloadSeeding([]models.ClientEvent{
		{Hash: hash},
	})

	// 结束
	go func() {
		time.Sleep(2 * time.Second)
		cancel()
	}()
	wg.Wait()

	fmt.Println(out.String())
	test.LogBatchCompare(out, test.MatchContainsRegexp,
		"OnDownloadStart",
		map[string]any{`\[Plugin\] Rename插件.*? src_\d.mp4`: 3},
		"OnDownloadSeeding",
		map[string]any{`write data/save/动画1/S02/E00\d.e_json`: 3},
		map[string]any{`\[重命名\] 链接「data/download/动画1/src_\d.mp4」`: 3},
		"正常退出",
	)
	out.Reset()

	wg, cancel = initTest(false)
	dbManager.OnDownloadStart([]models.ClientEvent{
		{Hash: hash},
	})
	dbManager.OnDownloadSeeding([]models.ClientEvent{
		{Hash: hash},
	})
	time.Sleep(1*time.Second + 500*time.Millisecond)
	dbManager.OnDownloadComplete([]models.ClientEvent{
		{Hash: hash},
	})

	// 结束
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()
	wg.Wait()

	fmt.Println(out.String())
	test.LogBatchCompare(out, test.MatchContainsRegexp,
		"OnDownloadStart",
		map[string]any{`\[Plugin\] Rename插件.*? src_\d.mp4`: 3},
		"OnDownloadSeeding",
		map[string]any{
			`\[重命名\] 可能已经移动完成，覆盖:「data/download/动画1/src_\d.mp4」`: 3,
			`\[重命名\] 链接「data/download/动画1/src_\d.mp4」`:           3,
		},
		"OnDownloadComplete",
		map[string]any{`write data/save/动画1/S02/E00\d.e_json`: 3},
		map[string]any{`\[重命名\] 删除「data/download/动画1/src_\d.mp4」`: 3},
		"移动完成「动画1[第2季][1-3集]」",
		"write data/save/动画1/anime.a_json",
		"write data/save/动画1/S02/anime.s_json",
		map[string]any{`write data/save/动画1/S02/E00\d.e_json`: 3},
		"写入元数据文件「data/save/动画1/tvshow.nfo」",
		"刮削完成: 动画1[第2季][1-3集]",
		map[string]any{`write data/save/动画1/S02/E00\d.e_json`: 3},
		"Delete",
		"正常退出",
	)

	exist := dbManager.IsExist(&models.AnimeEntity{
		NameCN: "动画1",
		Season: 2,
		Ep: []*models.AnimeEpEntity{
			{Type: models.AnimeEpNormal, Ep: 1},
			{Type: models.AnimeEpNormal, Ep: 2},
			{Type: models.AnimeEpNormal, Ep: 3},
		},
	})
	assert.Equal(t, exist, true)

}

// TestOnDownloadStep
//
//	分步下载
func TestOnDownloadStep(t *testing.T) {
	out.Reset()
	wg, cancel := initTest(true)

	// 下载 1
	hash := AddItem("动画1", 2, []int{1})
	dbManager.OnDownloadStart([]models.ClientEvent{
		{Hash: hash},
	})
	dbManager.OnDownloadSeeding([]models.ClientEvent{
		{Hash: hash},
	})
	time.Sleep(1*time.Second + 500*time.Millisecond)
	dbManager.OnDownloadComplete([]models.ClientEvent{
		{Hash: hash},
	})

	time.Sleep(1 * time.Second)
	cancel()
	wg.Wait()

	fmt.Println(out.String())
	out.Reset()

	wg, cancel = initTest(false)
	// 下载 2
	hash = AddItem("动画1", 2, []int{2})
	dbManager.OnDownloadStart([]models.ClientEvent{
		{Hash: hash},
	})
	dbManager.OnDownloadSeeding([]models.ClientEvent{
		{Hash: hash},
	})
	time.Sleep(1*time.Second + 500*time.Millisecond)
	dbManager.OnDownloadComplete([]models.ClientEvent{
		{Hash: hash},
	})

	time.Sleep(1 * time.Second)
	cancel()
	wg.Wait()
	fmt.Println(out.String())

	exist := dbManager.IsExist(&models.AnimeEntity{
		NameCN: "动画1",
		Season: 2,
		Ep: []*models.AnimeEpEntity{
			{Type: models.AnimeEpNormal, Ep: 1},
			{Type: models.AnimeEpNormal, Ep: 2},
		},
	})
	assert.Equal(t, exist, true)

}
