package main

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/wetor/AnimeGo/cmd/common"
	feedPlugin "github.com/wetor/AnimeGo/internal/animego/feed"
	filterPlugin "github.com/wetor/AnimeGo/internal/animego/filter"
	renamerPlugin "github.com/wetor/AnimeGo/internal/animego/renamer"
	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/internal/schedule"
	"github.com/wetor/AnimeGo/pkg/json"
	"github.com/wetor/AnimeGo/pkg/log"
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
	pRenameInput     string
)

func main() {
	common.PrintInfo()
	// -file=assets/plugin/rename/rename.py -plugin=rename -rename_input=assets/plugin/rename/testdata.json
	flag.BoolVar(&pDebug, "debug", true, "Debug模式，将会显示更多的日志")
	flag.StringVar(&pFile, "file", "", "插件脚本文件")
	flag.StringVar(&pPlugin, "plugin", "python", "插件类型，支持['python', 'filter', 'rename', 'schedule', 'feed']")
	flag.StringVar(&pType, "type", "python", "插件脚本类型，支持['python']")
	flag.StringVar(&pArgsJson, "args", "", "插件入口函数默认参数，json格式字符串")
	flag.StringVar(&pVarsJson, "vars", "", "插件全局变量，json格式字符串")
	flag.StringVar(&pPythonEntryFunc, "python_entry", "main", "python插件入口函数名")
	flag.StringVar(&pFilterInputFile, "filter_input", "", "filter插件要过滤的内容json文件")
	flag.StringVar(&pRenameInput, "rename_input", "", "rename插件用于重命名的动画信息json文件")
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

func pluginPython(entryFunc string, info *models.Plugin) map[string]any {
	p, err := plugin.LoadPlugin(&plugin.LoadPluginOptions{
		Plugin:    info,
		EntryFunc: entryFunc,
	})
	if err != nil {
		panic(err)
	}
	result, err := p.Run(entryFunc, map[string]any{})
	if err != nil {
		panic(err)
	}
	return result
}

func pluginFilter(items []*models.FeedItem, info *models.Plugin) []*models.FeedItem {
	f := filterPlugin.NewFilterPlugin(info)
	res, _ := f.FilterAll(items)
	return res
}

func pluginRename(anime *models.AnimeEntity, info *models.Plugin) []*models.RenameResult {
	result := make([]*models.RenameResult, len(anime.Ep))
	r := renamerPlugin.NewRenamePlugin(info)
	var err error
	for i, ep := range anime.Ep {
		result[i], err = r.Rename(anime, i, ep.Src)
		if err != nil {
			panic(err)
		}
	}
	return result
}

func pluginSchedule(file string, s *schedule.Schedule, info *models.Plugin) {
	t, _ := schedule.NewScheduleTask(&schedule.PluginOptions{
		Plugin: info,
	})
	s.Add(&schedule.AddTaskOptions{
		Name:     filepath.Base(file),
		StartRun: false,
		Task:     t,
	})
}

func pluginFeed(file string, s *schedule.Schedule, info *models.Plugin, callback func(items []*models.FeedItem) error) {
	t, _ := schedule.NewFeedTask(&schedule.FeedOptions{
		Plugin:   info,
		Callback: callback,
	})

	s.Add(&schedule.AddTaskOptions{
		Name:     filepath.Base(file),
		StartRun: false,
		Task:     t,
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
		Feed:  feedPlugin.NewRss(),
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
		result := pluginPython(pPythonEntryFunc, pluginInfo)
		log.Info(result)
	case constant.PluginTemplateRename:
		data, err := os.ReadFile(pRenameInput)
		if err != nil {
			panic(err)
		}
		anime := &models.AnimeEntity{}
		err = json.Unmarshal(data, anime)
		if err != nil {
			panic(err)
		}
		result := pluginRename(anime, pluginInfo)
		log.Infof("rename结果: ")
		for _, r := range result {
			log.Infof("    [%d] %v -> %v, tvshow.nfo位置: %v", r.Index, anime.Ep[r.Index].Src, r.Filename, r.AnimeDir)
		}
	case constant.PluginTemplateFilter:
		data, err := os.ReadFile(pFilterInputFile)
		if err != nil {
			panic(err)
		}
		items := make([]*models.FeedItem, 0)
		err = json.Unmarshal(data, &items)
		if err != nil {
			panic(err)
		}
		result := pluginFilter(items, pluginInfo)
		for i, item := range result {
			jsonData, _ := json.Marshal(item)
			log.Infof("[%d] filter结果: %v", i, string(jsonData))
		}
	case constant.PluginTemplateSchedule, constant.PluginTemplateFeed:
		s := schedule.NewSchedule(&schedule.Options{
			WG: &wg,
		})
		if pPlugin == constant.PluginTemplateSchedule {
			pluginSchedule(pFile, s, pluginInfo)
		} else if pPlugin == constant.PluginTemplateFeed {
			pluginFeed(pFile, s, pluginInfo, func(items []*models.FeedItem) error {
				for i, item := range items {
					jsonData, _ := json.Marshal(item)
					log.Infof("[%d] feed结果: %v", i, string(jsonData))
				}
				return nil
			})
		}
		s.Start(ctx)
		wg.Wait()
	default:
		panic("不支持的插件类型 " + pPlugin)
	}

}
