package task

import (
	"archive/zip"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/robfig/cron/v3"

	"github.com/wetor/AnimeGo/internal/utils"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/log"
)

const (
	CDN1               = "https://ghproxy.com/"
	ArchiveReleaseBase = "https://github.com/wetor/AnimeGoData/releases/download/archive/"
	Subject            = "bolt_sub.zip"
	SubjectDB          = "bolt_sub.db"

	Cron              = "0 0 12 * * 3" // 每周三12点
	MaxModifyTimeHour = 24             // 首次启动时，是否执行任务的最长修改时间
	MinFileSizeKB     = 512            // 首次启动时，是否执行任务的最小文件大小

	RetryNum  = 3  // 失败重试次数
	RetryWait = 60 // 失败重试等待时间，秒
)

type BangumiTask struct {
	parser   *cron.Parser
	cron     string
	savePath string
}

func NewBangumiTask() *BangumiTask {
	return &BangumiTask{
		savePath: DBDir,
		cron:     Cron,
		parser:   &SecondParser,
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
		log.Debugf("", errors.NewAniErrorD(errs))
		log.Errorf("使用ghproxy下载%s失败", name)
		return ""
	}
	file := path.Join(t.savePath, name)
	err := os.WriteFile(file, data, 0644)
	if err != nil {
		log.Debugf("", errors.NewAniErrorD(err))
		log.Errorf("保存文件到%s失败", name)
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
	db := path.Join(t.savePath, SubjectDB)
	stat, err := os.Stat(db)
	// 上次修改时间小于 MinModifyTimeHour 小时，且文件大小大于 MinFileSizeKB kb，跳过
	if force && err == nil &&
		time.Now().Unix()-stat.ModTime().Unix() <= MaxModifyTimeHour*60*60 && stat.Size() > MinFileSizeKB*1024 {
		return
	}
	subUrl := CDN1 + ArchiveReleaseBase + Subject
	file := t.download(subUrl, Subject)
	BangumiCacheLock.Lock()
	BangumiCache.Close()
	t.unzip(file)
	// 重新加载bolt
	BangumiCache.Open(db)
	BangumiCacheLock.Unlock()
	if utils.FileSize(db) <= MinFileSizeKB*1024 {
		errors.NewAniError("缓存文件小于512KB").TryPanic()
	}
}
