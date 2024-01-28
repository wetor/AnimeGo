package exceptions

import "fmt"

type ErrTimer struct {
	Message string
}

func (e ErrTimer) Error() string {
	return e.Message
}

type ErrTimerExistTask struct {
	Name    string
	Message string
}

func (e ErrTimerExistTask) Error() string {
	if len(e.Message) == 0 {
		return fmt.Sprintf("任务 %s 已存在", e.Name)
	}
	return fmt.Sprintf("任务 %s 已存在，%s", e.Name, e.Message)
}

func (e ErrTimerExistTask) Exist() bool {
	return true
}

type ErrTimerRun struct {
	Name    string
	Message string
}

func (e ErrTimerRun) Error() string {
	if len(e.Message) == 0 {
		return fmt.Sprintf("任务 %s 执行失败", e.Name)
	}
	return fmt.Sprintf("任务 %s 执行失败，%s", e.Name, e.Message)
}
