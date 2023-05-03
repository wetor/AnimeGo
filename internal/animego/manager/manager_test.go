package manager_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/brahma-adshonor/gohook"
	"github.com/stretchr/testify/assert"

	"github.com/wetor/AnimeGo/assets"
	"github.com/wetor/AnimeGo/internal/animego/manager"
	"github.com/wetor/AnimeGo/internal/animego/renamer"
	renamerPlugin "github.com/wetor/AnimeGo/internal/animego/renamer/plugin"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/pkg/xpath"
	"github.com/wetor/AnimeGo/test"
)

const (
	DownloadPath = "data/download"
	SavePath     = "data/save"
	ContentFile  = "file.mp4"

	HookTimeUnix = 100
)

var (
	qbt          *ClientMock
	rename       *renamer.Manager
	mgr          *manager.Manager
	renamePlugin *renamerPlugin.Rename
	db           *cache.Bolt
	out          *bytes.Buffer
)

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
	plugin.Init(&plugin.Options{
		Path:  assets.TestPluginPath(),
		Debug: true,
	})
	qbt = &ClientMock{}
	wg := sync.WaitGroup{}
	manager.Init(&manager.Options{
		Downloader: manager.Downloader{
			UpdateDelaySecond:      1,
			DownloadPath:           DownloadPath,
			SavePath:               SavePath,
			Category:               "test",
			Tag:                    "test",
			AllowDuplicateDownload: false,
			SeedingTimeMinute:      0,
			Rename:                 "wait_move",
		},
		WG: &wg,
	})
	db = cache.NewBolt()
	db.Open("data/test.db")
	db.Add("name2status")
	db.Put("name2status", "test[第1季][第1集]", &models.DownloadStatus{
		Hash:       "0000a4042b0bac2406b71023fdfe5e9054ebb832",
		State:      "complete",
		Path:       []string{SavePath + "/test/test.mp4"},
		Init:       true,
		Renamed:    true,
		Downloaded: true,
		Scraped:    true,
		Seeded:     true,
		ExpireAt:   0,
	}, 0)
	renamer.Init(&renamer.Options{
		WG:                &wg,
		UpdateDelaySecond: 1,
	})
	renamePlugin = renamerPlugin.NewRenamePlugin(&models.Plugin{
		Enable: true,
		Type:   "python",
		File:   "rename/builtin_rename.py",
	})
	rename = renamer.NewManager(renamePlugin)

	mgr = manager.NewManager(qbt, db, rename)

	m.Run()
	db.Close()
	_ = log.Close()
	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func download(name string, season int, ep []int) (files []string, fullname, hash string) {
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
			Ep:  e,
			Src: fmt.Sprintf("%s/src_%d.mp4", name, e),
		})
	}
	fullname = anime.FullName()
	qbt.MockAddName(fullname, hash, anime.FilePathSrc())
	mgr.Download(anime)
	srcFiles := anime.FilePathSrc()
	files = make([]string, len(anime.Ep))
	log.ReInt(&log.Options{
		Debug: false,
	})
	for i := range anime.Ep {
		files[i] = xpath.Join(SavePath, renamePlugin.Rename(anime, i, srcFiles[i]).Filepath)
	}
	log.ReInt(&log.Options{
		Debug: true,
	})
	return
}

func initTest() (*sync.WaitGroup, func()) {
	wg := sync.WaitGroup{}
	renamer.WG = &wg
	manager.WG = &wg
	ctx, cancel := context.WithCancel(context.Background())
	qbt.MockInit(nil)
	qbt.Start(ctx)
	rename.Start(ctx)
	mgr.Start(ctx)
	return &wg, cancel
}

func TestManager_Success(t *testing.T) {
	out.Reset()
	wg, cancel := initTest()

	var file1 []string
	go func() {
		time.Sleep(8 * time.Second)
		cancel()
	}()
	manager.Conf.Rename = "wait_move"
	{
		log.Info("下载 1")
		file1, _, _ = download("动画1", 1, []int{1, 2, 3})
	}
	wg.Wait()
	for _, f := range file1 {
		assert.FileExists(t, f)
		_ = os.Remove(f)
	}
	fmt.Println(out.String())
	test.LogBatchCompare(out,
		"下载 1",
		"接收到下载项",
		"开始下载",
		3, // plugin
		3, // 下载进度
		3, // 移动
		3, // 写入元数据文件
		"移动完成",
		"正常退出",
	)
}

func TestManager_Exist(t *testing.T) {
	out.Reset()
	wg, cancel := initTest()

	var file1 []string
	go func() {
		time.Sleep(7 * time.Second)
		cancel()
	}()
	manager.Conf.Rename = "move"
	{
		log.Info("下载 1")
		file1, _, _ = download("动画1", 1, []int{1, 2, 4})

	}
	{
		log.Info("下载 1")
		file1, _, _ = download("动画1", 1, []int{1, 2, 4})
	}
	wg.Wait()
	for _, f := range file1 {
		assert.FileExists(t, f)
		_ = os.Remove(f)
	}

	fmt.Println(out.String())
	test.LogBatchCompare(out,
		"下载 1",
		"下载 1",
		"接收到下载项",
		"开始下载",
		"接收到下载项",
		"取消下载",
		3, // plugin
		2, // 下载进度
		3, // 移动
		3, // 写入元数据文件
		1, // 下载进度
		"移动完成",
		"正常退出",
	)
}

func TestManager_DeleteFile_ReDownload(t *testing.T) {
	out.Reset()
	wg, cancel := initTest()

	var file1 []string
	go func() {
		time.Sleep(13 * time.Second)
		cancel()
	}()
	manager.Conf.Rename = "move"
	go func() {
		time.Sleep(300 * time.Millisecond)
		{
			log.Info("下载 1")
			file1, _, _ = download("动画1", 1, []int{1})

		}
		time.Sleep(6*time.Second + 300*time.Millisecond)
		{
			log.Info("删除 1 文件")
			for _, f := range file1 {
				_ = os.Remove(f)
				assert.NoFileExists(t, f)
			}
		}
		time.Sleep(500 * time.Millisecond)
		{
			log.Info("下载 1")
			file1, _, _ = download("动画1", 1, []int{1})
		}
	}()
	wg.Wait()
	for _, f := range file1 {
		assert.FileExists(t, f)
		_ = os.Remove(f)
	}

	fmt.Println(out.String())
	test.LogBatchCompare(out,
		"下载 1",
		"接收到下载项",
		"开始下载",
		1, // plugin
		2, // 下载进度
		1, // 移动
		1, // 写入元数据文件
		1, // 下载进度
		"移动完成",
		"删除 1",
		"下载 1",
		"接收到下载项",
		"开始下载",
		1, // plugin
		2, // 下载进度
		1, // 移动
		1, // 写入元数据文件
		1, // 下载进度
		"移动完成",
		"正常退出",
	)
}

func TestManager_DeleteCache_ReDownload(t *testing.T) {
	out.Reset()
	wg, cancel := initTest()

	var file1 []string
	var name1 string
	go func() {
		time.Sleep(13 * time.Second)
		cancel()
	}()
	manager.Conf.Rename = "move"
	go func() {
		time.Sleep(300 * time.Millisecond)
		{
			log.Info("下载 1")
			file1, name1, _ = download("动画1", 1, []int{1})

		}
		time.Sleep(6*time.Second + 300*time.Millisecond)
		{
			log.Info("删除 1 缓存")
			mgr.DeleteCache(name1)
			for _, f := range file1 {
				assert.FileExists(t, f)
			}
		}
		time.Sleep(500 * time.Millisecond)
		{
			log.Info("下载 1")
			file1, _, _ = download("动画1", 1, []int{1})
		}
	}()
	wg.Wait()
	for _, f := range file1 {
		assert.FileExists(t, f)
		_ = os.Remove(f)
	}

	fmt.Println(out.String())
	test.LogBatchCompare(out,
		"下载 1",
		"接收到下载项",
		"开始下载",
		1, // plugin
		2, // 下载进度
		1, // 移动
		1, // 写入元数据文件
		1, // 下载进度
		"移动完成",
		"删除 1",
		"下载 1",
		"接收到下载项",
		"开始下载",
		1, // plugin
		2, // 下载进度
		1, // 移动
		1, // 写入元数据文件
		1, // 下载进度
		"移动完成",
		"正常退出",
	)
}

func TestManager_WaitClient(t *testing.T) {
	out.Reset()
	wg, cancel := initTest()

	var file1, file2 []string
	go func() {
		time.Sleep(10 * time.Second)
		cancel()
	}()
	go func() {
		{
			log.Info("Client离线")
			qbt.MockSetError(ErrorConnectedFailed, true)
		}
		time.Sleep(2*time.Second + 300*time.Millisecond)
		{
			log.Info("Client恢复")
			qbt.MockSetError(ErrorConnectedFailed, false)
		}
	}()
	manager.Conf.Rename = "link_delete"
	go func() {
		time.Sleep(300 * time.Millisecond)
		{
			log.Info("下载 1")
			file1, _, _ = download("动画1", 1, []int{1})

		}
		time.Sleep(1*time.Second + 300*time.Millisecond)
		{
			log.Info("下载 2")
			file2, _, _ = download("动画2", 1, []int{1})
		}
	}()

	wg.Wait()
	for _, f := range file1 {
		assert.FileExists(t, f)
		_ = os.Remove(f)
	}
	for _, f := range file2 {
		assert.FileExists(t, f)
		_ = os.Remove(f)
	}

	fmt.Println(out.String())
	test.LogBatchCompare(out,
		"Client离线",
		"下载 1",
		"等待连接到下载器。已接收到1个下载项",
		"下载 2",
		"等待连接到下载器。已接收到2个下载项",
		"Client恢复",
		"接收到下载项",
		"开始下载",
		"接收到下载项",
		"开始下载",
		1, // plugin
		1, // 下载进度
		1, // plugin
		3, // 下载进度
		2, // 移动
		2, // 下载进度
		4, // 移动
		2, // 写入元数据文件
		"移动完成",
		"移动完成",
		"正常退出",
	)
}

func TestManager_WaitClient_FullChan(t *testing.T) {
	out.Reset()
	wg, cancel := initTest()

	var file []string
	go func() {
		time.Sleep(12 * time.Second)
		cancel()
	}()
	go func() {
		{
			log.Info("Client离线")
			qbt.MockSetError(ErrorConnectedFailed, true)
		}
		time.Sleep(5 * time.Second)
		{
			log.Info("Client恢复")
			qbt.MockSetError(ErrorConnectedFailed, false)
		}
	}()
	manager.Conf.Rename = "move"
	go func() {
		time.Sleep(300 * time.Millisecond)
		{
			for i := 1; i <= manager.DownloadChanDefaultCap; i++ {
				log.Infof("下载 %d", i)
				f, _, _ := download("动画1", 1, []int{i})
				file = append(file, f...)
			}
		}
		{
			log.Info("下载 2")
			f, _, _ := download("动画2", 1, []int{1})
			file = append(file, f...)
			log.Info("下载 3")
			f, _, _ = download("动画3", 1, []int{1})
			file = append(file, f...)
		}
	}()

	wg.Wait()
	for _, f := range file {
		assert.FileExists(t, f)
		_ = os.Remove(f)
	}

	fmt.Println(out.String())
	test.LogBatchCompare(out,
		"Client离线",
		12, // 下载
		"等待连接到下载器。已接收到10个下载项",
		"等待连接到下载器。已接收到10个下载项",
		"Client恢复",
		26,
		24, // plugin
		12, // 下载进度
		12, // 移动
		12, // 写入元数据文件
		12, // 下载进度
		12, // 移动完成
		"正常退出",
	)
}

func TestManager_ReStart_NotDownloaded(t *testing.T) {
	out.Reset()
	var file1 []string
	{
		wg, cancel := initTest()
		go func() {
			time.Sleep(5 * time.Second)
			cancel()
		}()
		time.Sleep(300 * time.Millisecond)
		manager.Conf.Rename = "move"
		{
			log.Info("下载 1")
			file1, _, _ = download("动画1", 1, []int{1})

		}
		wg.Wait()
		for _, f := range file1 {
			assert.FileExists(t, f)
		}
	}
	time.Sleep(1 * time.Second)
	{
		log.Info("重启")
		wg, cancel := initTest()
		go func() {
			time.Sleep(3 * time.Second)
			cancel()
		}()
		time.Sleep(300 * time.Millisecond)
		{
			log.Info("下载 1")
			file1, _, _ = download("动画1", 1, []int{1})
		}
		wg.Wait()
	}

	for _, f := range file1 {
		assert.FileExists(t, f)
		_ = os.Remove(f)
	}

	fmt.Println(out.String())
	test.LogBatchCompare(out,
		"下载 1",
		"接收到下载项",
		"开始下载",
		1,
		2,
		1, // 移动
		1,
		1,
		"正常退出",
		"正常退出",
		"重启",
		2,
		"下载 1",
		"接收到下载项",
		"发现已下载",
		"取消下载，不允许重复",
		"正常退出",
	)
}

func MockUnix1() int64 {
	return HookTimeUnix
}

func MockUnix2() int64 {
	return HookTimeUnix + manager.AddingExpireSecond
}

func TestManager_AddFailed(t *testing.T) {
	out.Reset()
	var file1 []string

	wg, cancel := initTest()
	go func() {
		time.Sleep(10 * time.Second)
		cancel()
	}()
	manager.Conf.Rename = "move"

	log.Infof("Hook utils.Unix() = %v", HookTimeUnix)
	_ = gohook.Hook(utils.Unix, MockUnix1, nil)
	defer gohook.UnHook(utils.Unix)

	qbt.MockSetError(ErrorAddFailed, true)
	{
		log.Info("下载 1, 添加失败")
		file1, _, _ = download("动画1", 1, []int{1})
	}

	time.Sleep(1*time.Second + 300*time.Millisecond)

	qbt.MockSetError(ErrorAddFailed, false)
	{
		log.Info("下载 1, 重复下载")
		file1, _, _ = download("动画1", 1, []int{1})
	}
	time.Sleep(1*time.Second + 300*time.Millisecond)

	log.Infof("Hook utils.Unix() = %v", HookTimeUnix+manager.AddingExpireSecond)
	_ = gohook.Hook(utils.Unix, MockUnix2, nil)
	time.Sleep(1*time.Second + 300*time.Millisecond)

	{
		log.Info("下载 1")
		file1, _, _ = download("动画1", 1, []int{1})
	}
	wg.Wait()
	for _, f := range file1 {
		assert.FileExists(t, f)
		_ = os.Remove(f)
	}
	fmt.Println(out.String())

	test.LogBatchCompare(out,
		"Hook",
		"下载 1, 添加失败",
		"接收到下载项",
		"开始下载",
		"下载 1, 重复下载",
		"接收到下载项",
		"取消下载，不允许重复",
		"Hook",
		"下载 1",
		"接收到下载项",
		"开始下载",
		1,
		2,
		1, // 移动
		1,
		1,
		"移动完成",
		"正常退出",
	)
}
