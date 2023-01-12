package task

import (
	"archive/zip"
	"github.com/parnurzeal/gorequest"
	"github.com/robfig/cron/v3"
	"github.com/wetor/AnimeGo/internal/store"
	"github.com/wetor/AnimeGo/internal/utils"
	"github.com/wetor/AnimeGo/pkg/errors"
	"go.uber.org/zap"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"
)

const (
	CDN1               = "https://ghproxy.com/"
	ArchiveReleaseBase = "https://github.com/wetor/AnimeGoData/releases/download/archive/"
	Subject            = "bolt_sub.zip"
	SubjectDB          = "bolt_sub.db"
)

type BangumiTask struct {
	parser   *cron.Parser
	cron     string
	savePath string
}

func NewBangumiTask(savePath string, parser *cron.Parser) *BangumiTask {
	return &BangumiTask{
		savePath: savePath,
		cron:     "0 0 12 * * 3", // 每周三12点
		parser:   parser,
	}
}

func (t *BangumiTask) Name() string {
	return "BangumiCache"
}

func (t *BangumiTask) Cron() string {
	return t.cron
}

func (t *BangumiTask) NextTime() time.Time {
	next, err := t.parser.Parse(t.cron)
	errors.NewAniErrorD(err).TryPanic()
	return next.Next(time.Now())
}

func (t *BangumiTask) download(url, name string) string {

	req := gorequest.New()
	_, data, errs := req.Get(url).EndBytes()
	if errs != nil {
		zap.S().Debug(errors.NewAniErrorD(errs))
		zap.S().Errorf("使用ghproxy下载%s失败", name)
		return ""
	}
	file := path.Join(t.savePath, name)
	err := os.WriteFile(file, data, 0644)
	if err != nil {
		zap.S().Debug(errors.NewAniErrorD(err))
		zap.S().Errorf("保存文件到%s失败", name)
		return ""
	}
	return file
}

func (t *BangumiTask) unzip(filename string) {
	zr, err := zip.OpenReader(filename)
	errors.NewAniErrorD(err).TryPanic()

	// 遍历 zr ，将文件写入到磁盘
	for _, file := range zr.File {
		path_ := filepath.Join(t.savePath, file.Name)

		// 如果是目录，就创建目录
		if file.FileInfo().IsDir() {
			err = os.MkdirAll(path_, file.Mode())
			errors.NewAniErrorD(err).TryPanic()
			continue
		}

		// 获取到 Reader
		fr, err := file.Open()
		errors.NewAniErrorD(err).TryPanic()

		// 创建要写出的文件对应的 Write
		fw, err := os.OpenFile(path_, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
		errors.NewAniErrorD(err).TryPanic()

		_, err = io.Copy(fw, fr)
		errors.NewAniErrorD(err).TryPanic()

		fw.Close()
		fr.Close()
	}
	zr.Close()
	err = os.Remove(filename)
	errors.NewAniErrorD(err).TryPanic()
}

func (t *BangumiTask) Run(force bool) {
	defer errors.HandleError(func(err error) {
		zap.S().Error(err)
	})
	db := path.Join(t.savePath, SubjectDB)
	stat, _ := os.Stat(db)
	// 上次修改时间小于24小时，且文件大小大于512kb，跳过
	if force && time.Now().Unix()-stat.ModTime().Unix() <= 24*60*60 && stat.Size() > 512*1024 {
		return
	}
	zap.S().Infof("[定时任务] %s 开始执行", t.Name())
	subUrl := CDN1 + ArchiveReleaseBase + Subject
	file := t.download(subUrl, Subject)
	store.BangumiCacheLock.Lock()
	store.BangumiCache.Close()
	t.unzip(file)
	// 重新加载bolt
	store.BangumiCache.Open(db)
	store.BangumiCacheLock.Unlock()
	if utils.FileSize(db) <= 512*1024 {
		zap.S().Infof("[定时任务] %s 执行失败", t.Name())
	} else {

		zap.S().Infof("[定时任务] %s 执行完毕，下次执行时间: %s", t.Name(), t.NextTime())
	}
}
