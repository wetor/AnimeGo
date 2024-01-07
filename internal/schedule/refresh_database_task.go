package schedule

import (
	"time"

	"github.com/robfig/cron/v3"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/log"
)

type RefreshTask struct {
	parser   *cron.Parser
	cron     string
	database api.Database
}

type RefreshOptions struct {
	Database api.Database
	Cron     string
}

func NewRefreshTask(opts *RefreshOptions) *RefreshTask {
	return &RefreshTask{
		parser:   &SecondParser,
		cron:     opts.Cron,
		database: opts.Database,
	}
}

func (t *RefreshTask) Name() string {
	return "RefreshDatabase"
}

func (t *RefreshTask) Cron() string {
	return t.cron
}

func (t *RefreshTask) SetVars(vars models.Object) {

}

func (t *RefreshTask) NextTime() time.Time {
	next, err := t.parser.Parse(t.Cron())
	if err != nil {
		log.DebugErr(err)
	}
	return next.Next(time.Now())
}

func (t *RefreshTask) Run(args models.Object) (err error) {
	return t.database.Scan()
}
