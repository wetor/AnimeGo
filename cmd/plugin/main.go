package main

import (
	"context"
	"flag"
	"os"
	"sync"
	"time"

	"github.com/wetor/AnimeGo/cmd/common"
	filterPlugin "github.com/wetor/AnimeGo/internal/animego/filter/plugin"
	renamerPlugin "github.com/wetor/AnimeGo/internal/animego/renamer/plugin"
	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/internal/schedule"
	"github.com/wetor/AnimeGo/internal/schedule/task"
	"github.com/wetor/AnimeGo/pkg/json"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

var (
	ctx, cancel = context.WithCancel(context.Background())
	wg          sync.WaitGroup

	pDebug           bool
	pFile            string
	pPlugin          string
	pType            string
	pPythonEntryFunc string
	pArgsJson        string
	pVarsJson        string
	pFilterInputFile string
	pRenameAnime     string
	pRenameFilename  string
)

func main() {
	common.PrintInfo()

	flag.BoolVar(&pDebug, "debug", true, "Debug模式，将会显示更多的日志")
	flag.StringVar(&pFile, "file", "", "插件脚本文件")
	flag.StringVar(&pPlugin, "plugin", "python", "插件类型，支持['python', 'filter', 'rename', 'schedule', 'feed']")
	flag.StringVar(&pType, "type", "python", "插件脚本类型，支持['python']")
	flag.StringVar(&pArgsJson, "args", "", "插件入口函数默认参数，json格式字符串")
	flag.StringVar(&pVarsJson, "vars", "", "插件全局变量，json格式字符串")
	flag.StringVar(&pPythonEntryFunc, "python_entry", "main", "python插件入口函数名")
	flag.StringVar(&pFilterInputFile, "filter_input", "", "filter插件输入json文件")
	flag.StringVar(&pRenameAnime, "rename_anime", "", "rename插件用于重命名的动画信息，json格式字符串")
	flag.StringVar(&pRenameFilename, "rename_filename", "", "rename插件的待重命名文件名")
	flag.Parse()

	common.RegisterExit(doExit)
	Main()
}

func doExit() {
	log.Infof("正在退出...")
	cancel()
	go func() {
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}()
}

func pluginPython(info *models.Plugin) map[string]any {
	p := plugin.LoadPlugin(&plugin.LoadPluginOptions{
		Plugin:    info,
		EntryFunc: pPythonEntryFunc,
	})
	return p.Run(pPythonEntryFunc, map[string]any{})
}

func pluginFilter(info *models.Plugin) []*models.FeedItem {
	data, err := os.ReadFile(pFilterInputFile)
	if err != nil {
		panic(err)
	}
	items := make([]*models.FeedItem, 0)
	err = json.Unmarshal(data, &items)
	if err != nil {
		panic(err)
	}

	f := filterPlugin.NewFilterPlugin(info)
	return f.Filter(items)
}

func pluginRename(info *models.Plugin) *models.RenameResult {
	anime := &models.AnimeEntity{}
	err := json.Unmarshal([]byte(pRenameAnime), anime)
	if err != nil {
		panic(err)
	}
	r := renamerPlugin.NewRenamePlugin(info)
	return r.Rename(anime, pRenameFilename)
}

func pluginSchedule(s *schedule.Schedule, info *models.Plugin) {
	s.Add(&schedule.AddTaskOptions{
		Name:     xpath.Base(pFile),
		StartRun: false,
		Task: task.NewScheduleTask(&task.ScheduleOptions{
			Plugin: info,
		}),
	})
}

func pluginFeed(s *schedule.Schedule, info *models.Plugin, callback func(items []*models.FeedItem)) {
	s.Add(&schedule.AddTaskOptions{
		Name:     xpath.Base(pFile),
		StartRun: false,
		Task: task.NewFeedTask(&task.FeedOptions{
			Plugin:   info,
			Callback: callback,
		}),
	})
}

func Main() {
	var err error
	if len(pFile) == 0 {
		panic("需要参数 file")
	}

	log.Init(&log.Options{
		File:  "plugin.log",
		Debug: pDebug,
	})
	dir, _ := os.Getwd()
	plugin.Init(&plugin.Options{
		Path:  dir,
		Debug: pDebug,
	})

	pluginInfo := &models.Plugin{
		Enable: true,
		Type:   pType,
		File:   pFile,
	}
	args := make(map[string]any)
	if len(pArgsJson) > 0 {
		err = json.Unmarshal([]byte(pArgsJson), &args)
		if err != nil {
			log.Warnf("插件入口函数默认参数解析错误: %v", err)
		}
		pluginInfo.Args = args
	}
	vars := make(map[string]any)
	if len(pVarsJson) > 0 {
		err = json.Unmarshal([]byte(pVarsJson), &vars)
		if err != nil {
			log.Warnf("插件全局变量解析错误: %v", err)
		}
		pluginInfo.Vars = vars
	}
	switch pPlugin {
	case constant.PluginTemplatePython:
		result := pluginPython(pluginInfo)
		log.Info(result)
	case constant.PluginTemplateRename:
		result := pluginRename(pluginInfo)
		log.Infof("rename结果: %v -> %v, tvshow.nfo位置: %v", pRenameFilename, result.Filepath, result.TVShowDir)
	case constant.PluginTemplateFilter:
		result := pluginFilter(pluginInfo)
		for i, item := range result {
			jsonData, _ := json.Marshal(item)
			log.Infof("[%d] filter结果: %v", i, string(jsonData))
		}
	case constant.PluginTemplateSchedule, constant.PluginTemplateFeed:
		s := schedule.NewSchedule(&schedule.Options{
			WG: &wg,
		})
		if pPlugin == constant.PluginTemplateSchedule {
			pluginSchedule(s, pluginInfo)
		} else if pPlugin == constant.PluginTemplateFeed {
			pluginFeed(s, pluginInfo, func(items []*models.FeedItem) {
				for i, item := range items {
					jsonData, _ := json.Marshal(item)
					log.Infof("[%d] feed结果: %v", i, string(jsonData))
				}
			})
		}
		s.Start(ctx)
		wg.Wait()
	default:
		panic("不支持的插件类型 " + pPlugin)
	}

}
