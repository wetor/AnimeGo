package timer

import (
	"context"
	"github.com/google/uuid"
	"github.com/wetor/AnimeGo/internal/api"
	"sync"
)

const (
	DefaultUpdateSecond = 1
	DefaultRetryCount   = 1
)

type TaskFunc func() error

type Options struct {
	Cache        api.Cacher
	RetryCount   int
	UpdateSecond int

	WG  *sync.WaitGroup
	Ctx context.Context
}

func (o *Options) Default() {
	if o.RetryCount == 0 {
		o.RetryCount = DefaultRetryCount
	}
	if o.UpdateSecond == 0 {
		o.UpdateSecond = DefaultUpdateSecond
	}
}

type AddOptions struct {
	Name     string
	Duration int64
	Func     TaskFunc
	Loop     bool
}

func (o *AddOptions) Default() {
	if len(o.Name) == 0 {
		o.Name = uuid.NewString()
	}
}
