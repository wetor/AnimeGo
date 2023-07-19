package renamer_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/wetor/AnimeGo/pkg/utils"

	"github.com/wetor/AnimeGo/internal/animego/downloader"
	"github.com/wetor/AnimeGo/internal/animego/renamer"
	renamerPlugin "github.com/wetor/AnimeGo/internal/animego/renamer/plugin"
	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/xpath"
	"github.com/wetor/AnimeGo/test"
)

const (
	DownloadPath = "data/download"
	SavePath     = "data/save"
)

var (
	out *bytes.Buffer
	r   *renamer.Manager
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	_ = os.RemoveAll("data")
	_ = os.MkdirAll(DownloadPath, os.ModePerm)
	_ = os.MkdirAll(SavePath, os.ModePerm)
	out = bytes.NewBuffer(nil)
	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
		Out:   out,
	})
	plugin.Init(&plugin.Options{
		Path:  "plugin/testdata",
		Debug: true,
	})
	wg := sync.WaitGroup{}
	renamer.Init(&renamer.Options{
		WG:                &wg,
		UpdateDelaySecond: 1,
	})
	p := renamerPlugin.NewRenamePlugin(&models.Plugin{
		Enable: true,
		Type:   "builtin",
		File:   "builtin_rename.py",
	})
	r = renamer.NewManager(p)

	m.Run()
	_ = log.Close()
	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func liseFiles(anime *models.AnimeEntity) []string {
	dst := anime.FilePath()
	result := make([]string, len(dst))
	for i := range result {
		result[i] = xpath.Join(SavePath, dst[i])
	}
	return result
}

func downloadFile(anime *models.AnimeEntity) {
	srcs := anime.FilePathSrc()
	for _, s := range srcs {
		_ = os.WriteFile(xpath.Join(DownloadPath, s), []byte{}, os.ModePerm)
	}
}

func rename(r *renamer.Manager, mode string, anime *models.AnimeEntity) error {
	downloadFile(anime)
	err := r.AddRenameTask(&models.RenameOptions{
		Entity: anime,
		SrcDir: DownloadPath,
		DstDir: SavePath,
		Mode:   mode,
		RenameCallback: func(result *models.RenameResult) {
			log.Infof("下载第%d集完成 %s", anime.Ep[result.Index].Ep, result.Filepath)
		},
		CompleteCallback: func(result *models.RenameResult) {
			log.Infof("下载完成 %s", anime.DirName())
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func initTest() (*sync.WaitGroup, func()) {
	wg := sync.WaitGroup{}
	renamer.WG = &wg
	ctx, cancel := context.WithCancel(context.Background())
	r.Start(ctx)
	return &wg, cancel
}

func TestManager_AddRenameTask(t *testing.T) {
	out.Reset()
	wg, cancel := initTest()

	anime1 := &models.AnimeEntity{
		NameCN: "动画1",
		Season: 1,
		Ep: []*models.AnimeEpEntity{
			{Type: models.AnimeEpNormal, Ep: 712, Src: "src_712.mp4"},
			{Type: models.AnimeEpNormal, Ep: 713, Src: "src_713.mp4"},
			{Type: models.AnimeEpUnknown, Ep: 0, Src: "src_714.mp4"},
		},
	}
	anime2 := &models.AnimeEntity{
		NameCN: "动画2",
		Season: 1,
		Ep: []*models.AnimeEpEntity{
			{Type: models.AnimeEpNormal, Ep: 1026, Src: "src_1026.mp4"},
		},
	}
	anime3 := &models.AnimeEntity{
		NameCN: "动画3",
		Season: 1,
		Ep: []*models.AnimeEpEntity{
			{Type: models.AnimeEpUnknown, Ep: 0, Src: "src_996.mp4"},
		},
	}
	anime4 := &models.AnimeEntity{
		NameCN: "动画4",
		Season: 2,
		Ep: []*models.AnimeEpEntity{
			{Type: models.AnimeEpNormal, Ep: 114, Src: "src_114.mp4"},
		},
	}

	files1 := liseFiles(anime1)
	files2 := liseFiles(anime2)
	files3 := liseFiles(anime3)
	files4 := liseFiles(anime4)

	var err error
	err = rename(r, "link_delete", anime1)
	assert.NoError(t, err)

	err = rename(r, "wait_move", anime2)
	assert.NoError(t, err)

	err = rename(r, "link", anime3)
	assert.NoError(t, err)

	err = rename(r, "move", anime4)
	assert.NoError(t, err)

	go func() {
		for i := range files1 {
			go func(i int) {
				_ = r.SetDownloadState(anime1.FullName(), i, downloader.StateSeeding)
				time.Sleep(1*time.Second + 100*time.Duration(i)*time.Millisecond)
				_ = r.SetDownloadState(anime1.FullName(), i, downloader.StateComplete)
			}(i)
		}
		time.Sleep(3 * time.Second)
		for i := range files2 {
			go func(i int) {
				_ = r.SetDownloadState(anime2.FullName(), i, downloader.StateSeeding)
				time.Sleep(1*time.Second + 100*time.Duration(i)*time.Millisecond)
				_ = r.SetDownloadState(anime2.FullName(), i, downloader.StateComplete)
			}(i)
		}
		time.Sleep(3 * time.Second)
		for i := range files3 {
			go func(i int) {
				_ = r.SetDownloadState(anime3.FullName(), i, downloader.StateSeeding)
				time.Sleep(1*time.Second + 100*time.Duration(i)*time.Millisecond)
				_ = r.SetDownloadState(anime3.FullName(), i, downloader.StateComplete)
			}(i)
		}
		time.Sleep(3 * time.Second)
		for i := range files4 {
			go func(i int) {
				_ = r.SetDownloadState(anime4.FullName(), i, downloader.StateSeeding)
				time.Sleep(1*time.Second + 100*time.Duration(i)*time.Millisecond)
				_ = r.SetDownloadState(anime4.FullName(), i, downloader.StateComplete)
			}(i)
		}
	}()

	go func() {
		time.Sleep(12 * time.Second)
		cancel()
	}()
	wg.Wait()

	for _, f := range files1 {
		assert.FileExists(t, f)
		_ = os.Remove(f)
	}
	for _, f := range files2 {
		assert.FileExists(t, f)
		_ = os.Remove(f)
	}
	for _, f := range files3 {
		assert.FileExists(t, f)
		_ = os.Remove(f)
	}
	for _, f := range files4 {
		assert.FileExists(t, f)
		_ = os.Remove(f)
	}

	fmt.Println(out.String())
	test.LogBatchCompare(out, nil,
		map[string]any{"Rename插件": 6, "[重命名] 链接": 3},
		[]string{"[重命名] 删除", "下载第712集完成 动画1"},
		[]string{"[重命名] 删除", "下载第713集完成 动画1"},
		[]string{"[重命名] 删除", "下载第0集完成 动画1"},
		"下载完成 动画1",
		[]string{"[重命名] 移动", "下载第1026集完成 动画2"},
		"下载完成 动画2",
		[]string{"[重命名] 链接", "下载第0集完成 动画3"},
		"下载完成 动画3",
		[]string{"[重命名] 移动", "下载第114集完成 动画4"},
		"下载完成 动画4",
		"正常退出",
	)
}

func TestManager_Method(t *testing.T) {
	out.Reset()
	wg, cancel := initTest()

	anime1 := &models.AnimeEntity{
		NameCN: "动画1",
		Season: 1,
		Ep: []*models.AnimeEpEntity{
			{Type: models.AnimeEpNormal, Ep: 1027, Src: "src_1027.mp4"},
		},
	}

	files1 := liseFiles(anime1)

	var err error
	err = rename(r, "link_delete", anime1)
	assert.NoError(t, err)

	name := "test"
	_, err = r.GetRenameTaskState(name)
	assert.ErrorAsf(t, err, &exceptions.ErrRename{}, "GetRenameTaskState(): %s", err)
	_, err = r.GetEpTaskState(name, 10)
	assert.ErrorAsf(t, err, &exceptions.ErrRename{}, "GetEpTaskState(): %s", err)
	err = r.SetDownloadState(name, 10, downloader.StateSeeding)
	assert.ErrorAsf(t, err, &exceptions.ErrRename{}, "SetDownloadState(): %s", err)

	name = anime1.FullName()
	if r.HasRenameTask(name) {
		_, err = r.GetRenameTaskState(name)
		assert.NoError(t, err)
		_, err = r.GetEpTaskState(name, 0)
		assert.NoError(t, err)
		_, err = r.GetEpTaskState(name, 10)
		assert.ErrorAsf(t, err, &exceptions.ErrRename{}, "GetEpTaskState(): %s", err)
		err = r.SetDownloadState(name, 10, downloader.StateSeeding)
		assert.ErrorAsf(t, err, &exceptions.ErrRename{}, "SetDownloadState(): %s", err)
	}

	for i := range files1 {
		go func(i int) {
			_ = r.SetDownloadState(name, i, downloader.StateSeeding)
			_ = r.SetDownloadState(name, i, downloader.StateComplete)
		}(i)
	}

	go func() {
		time.Sleep(5 * time.Second)
		cancel()
	}()
	wg.Wait()

	for _, f := range files1 {
		assert.FileExists(t, f)
		_ = os.Remove(f)
	}
	fmt.Println(out.String())
	test.LogBatchCompare(out, nil,
		"Rename插件",
		[]string{"[重命名] 链接", "[重命名] 删除", "下载第1027集完成 动画1"},
		"下载完成 动画1",
		"正常退出",
	)
}

func TestManager_LinkErr(t *testing.T) {
	out.Reset()
	wg, cancel := initTest()

	anime1 := &models.AnimeEntity{
		NameCN: "动画1",
		Season: 1,
		Ep: []*models.AnimeEpEntity{
			{Type: models.AnimeEpNormal, Ep: 1028, Src: "src_1028.mp4"},
		},
	}
	anime2 := &models.AnimeEntity{
		NameCN: "动画2",
		Season: 2,
		Ep: []*models.AnimeEpEntity{
			{Type: models.AnimeEpNormal, Ep: 1029, Src: "src_1029.mp4"},
		},
	}

	patches := test.HookSingle(utils.CreateLink, func(src, dst string) error {
		log.Infof("Hook CreateLink")
		return &os.LinkError{Op: "link", Old: src, New: dst, Err: errors.New("Hook CreateLink")}
	})
	defer patches.Reset()

	files1 := liseFiles(anime1)
	files2 := liseFiles(anime2)
	var err error
	err = rename(r, "link_delete", anime1)
	assert.NoError(t, err)
	err = rename(r, "link", anime2)
	assert.NoError(t, err)

	go func() {
		for i := range files1 {
			go func(i int) {
				_ = r.SetDownloadState(anime1.FullName(), i, downloader.StateSeeding)
				_ = r.SetDownloadState(anime1.FullName(), i, downloader.StateComplete)
			}(i)
		}
		time.Sleep(7 * time.Second)
		for i := range files2 {
			go func(i int) {
				_ = r.SetDownloadState(anime2.FullName(), i, downloader.StateSeeding)
				_ = r.SetDownloadState(anime2.FullName(), i, downloader.StateComplete)
			}(i)
		}
	}()

	go func() {
		time.Sleep(13 * time.Second)
		cancel()
	}()
	wg.Wait()

	fmt.Println(out.String())
	test.LogBatchCompare(out, nil,
		[]string{"Rename插件", "Rename插件"},
		[]string{"[重命名] 链接", "Hook CreateLink", "Hook CreateLink", "失败: 创建文件链接失败"},
		[]string{"[重命名] 链接", "Hook CreateLink", "Hook CreateLink", "失败: 创建文件链接失败"},
		[]string{"[重命名] 链接", "Hook CreateLink", "Hook CreateLink", "失败: 创建文件链接失败"},
		[]string{"[重命名] 失败，跳过流程", "下载第1028集完成", "下载完成 动画1"},
		[]string{"[重命名] 链接", "Hook CreateLink", "Hook CreateLink", "失败: 创建文件链接失败"},
		[]string{"[重命名] 链接", "Hook CreateLink", "Hook CreateLink", "失败: 创建文件链接失败"},
		[]string{"[重命名] 链接", "Hook CreateLink", "Hook CreateLink", "失败: 创建文件链接失败"},
		[]string{"[重命名] 失败，跳过流程", "下载第1029集完成", "下载完成 动画2"},
		"正常退出",
	)
}

func TestManager_MoveErr(t *testing.T) {
	out.Reset()
	wg, cancel := initTest()

	anime1 := &models.AnimeEntity{
		NameCN: "动画1",
		Season: 1,
		Ep: []*models.AnimeEpEntity{
			{Type: models.AnimeEpNormal, Ep: 1030, Src: "src_1030.mp4"},
		},
	}
	anime2 := &models.AnimeEntity{
		NameCN: "动画2",
		Season: 2,
		Ep: []*models.AnimeEpEntity{
			{Type: models.AnimeEpNormal, Ep: 1031, Src: "src_1031.mp4"},
		},
	}

	patches := test.HookSingle(utils.Rename, func(src, dst string) error {
		log.Infof("Hook Rename")
		return &os.LinkError{Op: "rename", Old: src, New: dst, Err: errors.New("Hook Rename")}
	})
	defer patches.Reset()

	files1 := liseFiles(anime1)
	files2 := liseFiles(anime2)
	var err error
	err = rename(r, "move", anime1)
	assert.NoError(t, err)
	err = rename(r, "wait_move", anime2)
	assert.NoError(t, err)

	go func() {
		for i := range files1 {
			go func(i int) {
				_ = r.SetDownloadState(anime1.FullName(), i, downloader.StateSeeding)
				_ = r.SetDownloadState(anime1.FullName(), i, downloader.StateComplete)
			}(i)
		}
		time.Sleep(7 * time.Second)
		for i := range files2 {
			go func(i int) {
				_ = r.SetDownloadState(anime2.FullName(), i, downloader.StateSeeding)
				_ = r.SetDownloadState(anime2.FullName(), i, downloader.StateComplete)
			}(i)
		}
	}()

	go func() {
		time.Sleep(14 * time.Second)
		cancel()
	}()
	wg.Wait()

	fmt.Println(out.String())
	test.LogBatchCompare(out, nil,
		[]string{"Rename插件", "Rename插件"},
		[]string{"[重命名] 移动", "Hook Rename", "Hook Rename", "失败: 重命名文件失败"},
		[]string{"[重命名] 移动", "Hook Rename", "Hook Rename", "失败: 重命名文件失败"},
		[]string{"[重命名] 移动", "Hook Rename", "Hook Rename", "失败: 重命名文件失败"},
		[]string{"[重命名] 失败，跳过流程", "下载第1030集完成", "下载完成 动画1"},
		[]string{"[重命名] 移动", "Hook Rename", "Hook Rename", "失败: 重命名文件失败"},
		[]string{"[重命名] 移动", "Hook Rename", "Hook Rename", "失败: 重命名文件失败"},
		[]string{"[重命名] 移动", "Hook Rename", "Hook Rename", "失败: 重命名文件失败"},
		[]string{"[重命名] 失败，跳过流程", "下载第1031集完成", "下载完成 动画2"},
		"正常退出",
	)
}

func TestManager_LinkDeleteErr(t *testing.T) {
	out.Reset()
	wg, cancel := initTest()

	anime1 := &models.AnimeEntity{
		NameCN: "动画1",
		Season: 1,
		Ep: []*models.AnimeEpEntity{
			{Type: models.AnimeEpNormal, Ep: 1032, Src: "src_1032.mp4"},
		},
	}

	patches := test.HookSingle(os.Remove, func(dst string) error {
		log.Infof("Hook Remove")
		return &os.PathError{Op: "remove", Path: dst, Err: errors.New("Hook Remove")}
	})

	files1 := liseFiles(anime1)
	var err error
	err = rename(r, "link_delete", anime1)
	assert.NoError(t, err)

	go func() {
		for i := range files1 {
			go func(i int) {
				_ = r.SetDownloadState(anime1.FullName(), i, downloader.StateSeeding)
				_ = r.SetDownloadState(anime1.FullName(), i, downloader.StateComplete)
			}(i)
		}
	}()

	go func() {
		time.Sleep(7 * time.Second)
		cancel()
	}()
	wg.Wait()
	// 可能已经移动完成，覆盖
	time.Sleep(1 * time.Second)
	patches.Reset()
	for _, f := range files1 {
		assert.FileExists(t, f)
	}
	wg, cancel = initTest()
	err = rename(r, "link_delete", anime1)
	assert.NoError(t, err)

	go func() {
		for i := range files1 {
			go func(i int) {
				_ = r.SetDownloadState(anime1.FullName(), i, downloader.StateSeeding)
				_ = r.SetDownloadState(anime1.FullName(), i, downloader.StateComplete)
			}(i)
		}
	}()

	go func() {
		time.Sleep(5 * time.Second)
		cancel()
	}()
	wg.Wait()

	// 可能已经移动完成，跳过
	time.Sleep(1 * time.Second)
	wg, cancel = initTest()
	patches = test.HookSingle(os.WriteFile, func(name string, data []byte, perm os.FileMode) error {
		return nil
	})
	defer patches.Reset()
	err = rename(r, "link_delete", anime1)
	assert.NoError(t, err)
	go func() {
		for i := range files1 {
			go func(i int) {
				_ = r.SetDownloadState(anime1.FullName(), i, downloader.StateSeeding)
				_ = r.SetDownloadState(anime1.FullName(), i, downloader.StateComplete)
			}(i)
		}
	}()

	go func() {
		time.Sleep(4 * time.Second)
		cancel()
	}()
	wg.Wait()

	fmt.Println(out.String())
	test.LogBatchCompare(out, nil,
		[]string{"Rename插件"},
		[]string{"[重命名] 链接", "[重命名] 删除", "Hook Remove", "Hook Remove", "失败: 删除文件失败"},
		[]string{"[重命名] 删除", "Hook Remove", "Hook Remove", "失败: 删除文件失败"},
		[]string{"[重命名] 删除", "Hook Remove", "Hook Remove", "失败: 删除文件失败"},
		[]string{"[重命名] 失败，跳过流程", "下载第1032集完成", "下载完成 动画1"},
		"正常退出",
		[]string{"Rename插件"},
		[]string{"[重命名] 可能已经移动完成，覆盖", "[重命名] 链接", "[重命名] 删除"},
		[]string{"下载第1032集完成", "下载完成 动画1"},
		"正常退出",
		[]string{"Rename插件"},
		[]string{"[重命名] 可能已经移动完成，跳过"},
		[]string{"下载第1032集完成", "下载完成 动画1"},
		"正常退出",
	)
}

func TestManager_NotFileErr(t *testing.T) {
	out.Reset()
	wg, cancel := initTest()

	anime1 := &models.AnimeEntity{
		NameCN: "动画1",
		Season: 1,
		Ep: []*models.AnimeEpEntity{
			{Type: models.AnimeEpNormal, Ep: 1033, Src: "src_1033.mp4"},
		},
	}

	files1 := liseFiles(anime1)

	var err error
	patches := test.HookSingle(os.WriteFile, func(name string, data []byte, perm os.FileMode) error {
		return nil
	})
	err = rename(r, "move", anime1)
	patches.Reset()
	assert.NoError(t, err)

	go func() {
		for i := range files1 {
			go func(i int) {
				_ = r.SetDownloadState(anime1.FullName(), i, downloader.StateSeeding)
				_ = r.SetDownloadState(anime1.FullName(), i, downloader.StateComplete)
			}(i)
		}
	}()
	go func() {
		time.Sleep(2*time.Second + 500*time.Millisecond)
		downloadFile(anime1)
	}()
	go func() {
		time.Sleep(6 * time.Second)
		cancel()
	}()
	wg.Wait()
	fmt.Println(out.String())
	test.LogBatchCompare(out, nil,
		"Rename插件",
		[]string{"失败: 未找到文件", "失败: 未找到文件", "[重命名] 移动"},
		"下载第1033集完成 动画1",
		"下载完成 动画1",
		"正常退出",
	)
}

func TestManager_OtherErr(t *testing.T) {
	out.Reset()
	wg, cancel := initTest()

	anime1 := &models.AnimeEntity{
		NameCN: "动画1",
		Season: 1,
		Ep: []*models.AnimeEpEntity{
			{Type: models.AnimeEpNormal, Ep: 1034, Src: "src_1034.mp4"},
		},
	}

	files1 := liseFiles(anime1)

	var err error
	err = rename(r, "test", anime1)
	assert.NoError(t, err)

	go func() {
		for i := range files1 {
			go func(i int) {
				_ = r.SetDownloadState(anime1.FullName(), i, downloader.StateSeeding)
				_ = r.SetDownloadState(anime1.FullName(), i, downloader.StateComplete)
			}(i)
		}
	}()
	go func() {
		time.Sleep(3 * time.Second)
		cancel()
	}()
	wg.Wait()

	fmt.Println(out.String())
	test.LogBatchCompare(out, nil,
		"Rename插件",
		"失败: 不支持的重命名模式",
		"下载完成 动画1",
		"正常退出",
	)
}
