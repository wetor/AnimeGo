package main

import (
	"context"
	"flag"
	"os"
	"sync"

	filterPlugin "github.com/wetor/AnimeGo/internal/animego/filter/plugin"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/internal/schedule"
	"github.com/wetor/AnimeGo/internal/schedule/task"
	"github.com/wetor/AnimeGo/pkg/json"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

var (
	wg  sync.WaitGroup
	ctx context.Context

	pFile            string
	pPlugin          string
	pType            string
	pPythonEntryFunc string
	pArgsJson        string
	pVarsJson        string
	pFilterInputFile string
)

func main() {
	flag.StringVar(&pFile, "file", "", "插件脚本文件")
	flag.StringVar(&pPlugin, "plugin", "python", "插件类型，支持['python', 'filter', 'schedule', 'feed']")
	flag.StringVar(&pType, "type", "python", "插件脚本类型，支持['python']")
	flag.StringVar(&pArgsJson, "args", "", "插件入口函数默认参数，json格式字符串")
	flag.StringVar(&pVarsJson, "vars", "", "插件全局变量，json格式字符串")
	flag.StringVar(&pPythonEntryFunc, "python_entry", "main", "python插件入口函数名")
	flag.StringVar(&pFilterInputFile, "filter_input", "", "filter插件输入json文件")
	flag.Parse()

	if len(pFile) == 0 {
		panic("需要参数 file")
	}

	log.Init(&log.Options{
		File:  "plugin.log",
		Debug: true,
	})
	ctx = context.Background()

	var err error
	args := make(map[string]any)
	vars := make(map[string]any)
	if len(pArgsJson) > 0 {
		err = json.Unmarshal([]byte(pArgsJson), &args)
		if err != nil {
			log.Warnf("插件入口函数默认参数解析错误: %v", err)
		}
	}
	if len(pVarsJson) > 0 {
		err = json.Unmarshal([]byte(pVarsJson), &vars)
		if err != nil {
			log.Warnf("插件全局变量解析错误: %v", err)
		}
	}
	dir, _ := os.Getwd()
	plugin.Init(&plugin.Options{
		Path: dir,
	})
	switch pPlugin {
	case "python":
		p := plugin.LoadPlugin(&plugin.LoadPluginOptions{
			Plugin: &models.Plugin{
				Enable: true,
				Type:   pType,
				File:   pFile,
				Args:   args,
				Vars:   vars,
			},
			EntryFunc: pPythonEntryFunc,
		})
		result := p.Run(pPythonEntryFunc, map[string]any{})
		log.Info(result)
	case "filter":
		data, err := os.ReadFile(pFilterInputFile)
		if err != nil {
			panic(err)
		}
		items := make([]*models.FeedItem, 0)
		err = json.Unmarshal(data, &items)
		if err != nil {
			panic(err)
		}

		f := filterPlugin.NewFilterPlugin([]models.Plugin{
			{
				Enable: true,
				Type:   pType,
				File:   pFile,
				Args:   args,
				Vars:   vars,
			},
		})
		result := f.Filter(items)
		for i, item := range result {
			jsonData, _ := json.Marshal(item)
			log.Infof("[%d] filter结果: %v", i, string(jsonData))
		}
	case "schedule", "feed":
		s := schedule.NewSchedule(&schedule.Options{
			WG: &wg,
		})
		if pPlugin == "schedule" {
			s.Add(&schedule.AddTaskOptions{
				Name:     xpath.Base(pFile),
				StartRun: false,
				Task: task.NewScheduleTask(&task.ScheduleOptions{
					Plugin: &models.Plugin{
						Enable: true,
						Type:   pType,
						File:   pFile,
						Args:   args,
						Vars:   vars,
					},
				}),
			})
		} else if pPlugin == "feed" {
			s.Add(&schedule.AddTaskOptions{
				Name:     xpath.Base(pFile),
				StartRun: false,
				Task: task.NewFeedTask(&task.FeedOptions{
					Plugin: &models.Plugin{
						Enable: true,
						Type:   pType,
						File:   pFile,
						Args:   args,
						Vars:   vars,
					},
					Callback: func(items []*models.FeedItem) {
						for i, item := range items {
							jsonData, _ := json.Marshal(item)
							log.Infof("[%d] 接收到feed下载项: %v", i, string(jsonData))
						}
					},
				}),
			})
		}

		s.Start(ctx)
		wg.Wait()
	default:
		panic("不支持的插件类型 " + pPlugin)
	}

}
