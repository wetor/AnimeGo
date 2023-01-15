package task

import (
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/wetor/AnimeGo/pkg/errors"
)

type JSPluginTask struct {
	parser *cron.Parser
	cron   string
}

func NewJSPluginTask(parser *cron.Parser) *JSPluginTask {
	return &JSPluginTask{
		cron:   "*/5 * * * * ?", // 5s 执行一次
		parser: parser,
	}
}

func (t *JSPluginTask) Name() string {
	return "JavaScript Plugin"
}

func (t *JSPluginTask) Cron() string {
	return t.cron
}

func (t *JSPluginTask) NextTime() time.Time {
	next, err := t.parser.Parse(t.cron)
	errors.NewAniErrorD(err).TryPanic()
	return next.Next(time.Now())
}

func (t *JSPluginTask) Run(force bool) {
	defer errors.HandleError(func(err error) {
		zap.S().Error(err)
	})
	zap.S().Infof("[定时任务] %s 开始执行", t.Name())

	zap.S().Infof("[定时任务] %s 执行完毕，下次执行时间: %s", t.Name(), t.NextTime())
}
