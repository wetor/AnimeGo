package task

import (
	"archive/zip"
	"io"
	"os"
	"sync"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

const (
	CDN1               = "https://ghproxy.com/"
	ArchiveReleaseBase = "https://github.com/wetor/AnimeGoData/releases/download/archive/"
	Subject            = "bolt_sub.zip"
	SubjectDB          = "bolt_sub.db"

	Cron              = "0 0 12 * * 3" // 每周三12点
	MaxModifyTimeHour = 12             // 首次启动时，是否执行任务的最长修改时间
	MinFileSizeKB     = 512            // 首次启动时，是否执行任务的最小文件大小
)

var firstRun = true

type BangumiTask struct {
	parser     *cron.Parser
	cron       string
	cache      api.CacheOpener
	cacheMutex *sync.Mutex
}

type BangumiOptions struct {
	Cache      api.CacheOpener
	CacheMutex *sync.Mutex
}

func NewBangumiTask(opts *BangumiOptions) *BangumiTask {
	return &BangumiTask{
		parser:     &SecondParser,
		cron:       Cron,
		cache:      opts.Cache,
		cacheMutex: opts.CacheMutex,
	}
}

func (t *BangumiTask) Name() string {
	return "BangumiCache"
}

func (t *BangumiTask) Cron() string {
	return t.cron
}

func (t *BangumiTask) SetVars(vars models.Object) {

}

func (t *BangumiTask) NextTime() time.Time {
	next, err := t.parser.Parse(t.Cron())
	if err != nil {
		log.DebugErr(err)
	}
	return next.Next(time.Now())
}

func (t *BangumiTask) download(url, name string) (string, error) {
	var err error
	req := gorequest.New()
	_, data, errs := req.Get(url).EndBytes()
	if errs != nil {
		log.DebugErr(errs[0])
		err = errors.WithStack(&exceptions.ErrSchedule{Message: "使用ghproxy下载失败: " + name})
		log.Warnf("%s", err)
		return "", err
	}
	file := xpath.Join(constant.CachePath, name)
	err = os.WriteFile(file, data, 0644)
	if err != nil {
		log.DebugErr(err)
		err = errors.WithStack(&exceptions.ErrSchedule{Message: "保存到文件失败: " + name})
		log.Warnf("%s", err)
		return "", err
	}
	return file, nil
}

func (t *BangumiTask) unzip(filename string) (err error) {
	zr, err := zip.OpenReader(filename)
	if err != nil {
		log.DebugErr(err)
		err = errors.WithStack(&exceptions.ErrSchedule{Message: "载入zip文件失败"})
		return err
	}

	// 遍历 zr ，将文件写入到磁盘
	for _, file := range zr.File {
		path_ := xpath.Join(constant.CachePath, file.Name)

		// 如果是目录，就创建目录
		if file.FileInfo().IsDir() {
			err = os.MkdirAll(path_, file.Mode())
			if err != nil {
				log.DebugErr(err)
				err = errors.WithStack(&exceptions.ErrSchedule{Message: "创建文件夹失败: " + path_})
				return err
			}
		}

		// 获取到 Reader
		fr, err := file.Open()
		if err != nil {
			log.DebugErr(err)
			err = errors.WithStack(&exceptions.ErrSchedule{Message: "读取zip内文件失败"})
			return err
		}

		// 创建要写出的文件对应的 Write
		fw, err := os.OpenFile(path_, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
		if err != nil {
			log.DebugErr(err)
			err = errors.WithStack(&exceptions.ErrSchedule{Message: "打开文件失败: " + path_})
			return err
		}

		_, err = io.Copy(fw, fr)
		if err != nil {
			log.DebugErr(err)
			err = errors.WithStack(&exceptions.ErrSchedule{Message: "写入文件失败: " + path_})
			return err
		}

		_ = fw.Close()
		_ = fr.Close()
	}
	_ = zr.Close()
	err = os.Remove(filename)
	if err != nil {
		log.DebugErr(err)
		err = errors.WithStack(&exceptions.ErrSchedule{Message: "删除文件失败: " + filename})
		return err
	}
	return nil
}

// Run
//
//	@Description:
//	@receiver *BangumiTask
//	@param opts ...interface{}
//	  opts[0] bool 是否启动时执行
func (t *BangumiTask) Run(args models.Object) (err error) {
	db := xpath.Join(constant.CachePath, SubjectDB)
	stat, err := os.Stat(db)
	// 首次启动时，若
	// 上次修改时间小于 MinModifyTimeHour 小时，且文件大小大于 MinFileSizeKB kb
	// 则不执行
	if firstRun && err == nil &&
		time.Now().Unix()-stat.ModTime().Unix() <= MaxModifyTimeHour*60*60 && stat.Size() > MinFileSizeKB*1024 {
		firstRun = false
		return
	}
	subUrl := CDN1 + ArchiveReleaseBase + Subject
	file, err := t.download(subUrl, Subject)
	if err != nil {
		return err
	}
	t.cacheMutex.Lock()
	t.cache.Close()
	err = t.unzip(file)
	if err != nil {
		return err
	}
	// 重新加载bolt
	t.cache.Open(db)
	t.cacheMutex.Unlock()
	if utils.FileSize(db) <= MinFileSizeKB*1024 {
		return errors.WithStack(&exceptions.ErrSchedule{Message: "缓存文件小于512KB"})
	}
	return nil
}
