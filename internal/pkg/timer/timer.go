package timer

import (
	"context"
	"sync"
	"time"

	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/pkg/log"
)

const (
	Name = "Timer"
)

const (
	StatusStop    = "stop"
	StatusInit    = "init"
	StatusWait    = "wait"
	StatusExpired = "expired"
)

type Task struct {
	Name        string `json:"name"`         // 任务名
	Duration    int64  `json:"duration"`     // 执行定时，秒
	Start       int64  `json:"start"`        // 开始时间
	RunDuration int64  `json:"run_duration"` // 此次定时已执行时间，秒。序列化使用
	Status      string `json:"status"`       // 状态
	RetryCount  int    `json:"retry_count"`  // 剩余重试次数
	Loop        bool   `json:"loop"`         // 是否循环执行
}

type Timer struct {
	tasks map[string]*Task
	funcs map[string]TaskFunc

	sync.Mutex

	*Options
}

func NewTimer(opts *Options) *Timer {
	opts.Default()
	t := &Timer{
		tasks:   make(map[string]*Task),
		funcs:   make(map[string]TaskFunc),
		Options: opts,
	}
	t.Cache.Add(Name)
	err := t.Unmarshal()
	if err != nil {
		log.Warnf("[Timer] 载入缓存记录失败，可能不存在")
	}
	return t
}

func (t *Timer) HasTask(name string) bool {
	_, ok := t.tasks[name]
	return ok
}

func (t *Timer) AddTask(opts *AddOptions) (*Task, error) {
	t.Lock()
	defer t.Unlock()
	opts.Default()
	if t.HasTask(opts.Name) {
		return nil, exceptions.ErrTimerExistTask{Name: opts.Name}
	}
	task := &Task{
		Name:     opts.Name,
		Duration: opts.Duration,
		Status:   StatusInit,
		Loop:     opts.Loop,
	}
	t.tasks[task.Name] = task
	t.funcs[task.Name] = opts.Func

	return task, nil
}

func (t *Timer) Start(ctx context.Context) {
	t.WG.Add(1)
	go func() {
		defer t.WG.Done()
		for {
			select {
			case <-ctx.Done():
				t.Marshal()
				log.Debugf("[Timer] 正常退出 %s", Name)
				return
			default:
				t.update()
				time.Sleep(time.Duration(t.UpdateSecond) * time.Second)
			}
		}
	}()
}

func (t *Timer) update() {
	t.Lock()
	defer t.Unlock()
	var err error
	deleteTasks := make([]string, 0)
	now := time.Now().Unix()
	for _, task := range t.tasks {
		if task.Status == StatusStop {
			continue
		}

		if task.Status == StatusWait && now >= task.Start+task.Duration {
			// 执行任务
			log.Debugf("[Timer] 任务 %s 开始执行", task.Name)
			if f, ok := t.funcs[task.Name]; ok {
				err = f()
			} else {
				err = nil
				log.Warnf("[Timer] 任务 %s 执行失败，未注册执行函数，忽略", task.Name)
			}
			finish := false
			if err != nil {
				task.RetryCount--
				log.Debugf("[Timer] 任务 %s 执行失败，第 %d 次重试", task.Name, t.RetryCount-task.RetryCount)
				log.DebugErr(err)
			} else {
				finish = true
				log.Infof("[Timer] 任务 %s 执行成功", task.Name)
			}

			if task.Status != StatusExpired && task.RetryCount <= 0 {
				finish = true
				log.Warnf("[Timer] 任务 %s 执行失败，重试 %d 次", task.Name, t.RetryCount-task.RetryCount)
			}

			if finish {
				if task.Loop {
					task.Status = StatusInit
				} else {
					task.Status = StatusExpired
				}
			}
		}

		if task.Status == StatusExpired {
			deleteTasks = append(deleteTasks, task.Name)
		}

		if task.Status == StatusInit {
			task.Start = now
			task.Status = StatusWait
			task.RetryCount = t.RetryCount
			log.Debugf("[Timer] 任务 %s 已添加，下次执行： %d 秒后", task.Name, task.Duration)
		}
	}

	for _, id := range deleteTasks {
		delete(t.tasks, id)
		delete(t.funcs, id)
	}
}

func (t *Timer) Marshal() {
	t.Lock()
	defer t.Unlock()
	now := time.Now().Unix()
	for _, task := range t.tasks {
		task.RunDuration = (now - task.Start) % task.Duration
	}
	t.Cache.Put(Name, t.Name, t.tasks, 0)
}

func (t *Timer) Unmarshal() error {
	t.Lock()
	defer t.Unlock()
	t.tasks = make(map[string]*Task)
	err := t.Cache.Get(Name, t.Name, &t.tasks)
	if err != nil {
		return err
	}
	now := time.Now().Unix()
	for _, task := range t.tasks {
		task.Start = now - task.RunDuration
		task.RunDuration = 0
	}
	return nil
}

func (t *Timer) RegisterTaskFuncs(funcs map[string]TaskFunc) {
	for name, f := range funcs {
		t.funcs[name] = f
	}
}
