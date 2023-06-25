package exceptions

import "fmt"

type ErrSchedule struct {
	Message string
}

func (e ErrSchedule) Error() string {
	return e.Message
}

type ErrScheduleAdd struct {
	Name string
}

func (e ErrScheduleAdd) Error() string {
	return fmt.Sprintf("添加定时任务 %s 失败", e.Name)
}

type ErrScheduleRun struct {
	Name    string
	Message any
}

func (e ErrScheduleRun) Error() string {
	return fmt.Sprintf("定时任务 %s 执行失败: %v", e.Name, e.Message)
}

type ErrSchedulePluginGetVar struct {
	Name    string
	Message string
}

func (e ErrSchedulePluginGetVar) Error() string {
	return fmt.Sprintf("定时任务插件 %s 变量: %s", e.Name, e.Message)
}
