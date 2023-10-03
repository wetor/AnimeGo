package models

type RenameResult struct {
	Index     int    `json:"index"`
	Scrape    bool   `json:"scrape"` // 是否进行刮削
	Filename  string `json:"filename"`
	AnimeDir  string `json:"anime_dir"`
	SeasonDir string `json:"season_dir"`
}

type RenameAllResult struct {
	Results   []*RenameResult
	Name      string `json:"name"`
	AnimeDir  string `json:"anime_dir"`
	SeasonDir string `json:"season_dir"`
}

func (r RenameAllResult) Scrape() bool {
	for _, res := range r.Results {
		if res.Scrape {
			return true
		}
	}
	return false
}

func (r RenameAllResult) Filenames() []string {
	result := make([]string, len(r.Results))
	for i, res := range r.Results {
		result[i] = res.Filename
	}
	return result
}

type RenameCallback func(*RenameResult)
type CompleteCallback func(*RenameAllResult)

type RenameOptions struct {
	Name             string           // 动画名
	Entity           *AnimeEntity     // 动画详情
	SrcDir           string           // 源文件夹
	DstDir           string           // 目标文件夹
	Mode             string           // 重命名模式
	RenameCallback   RenameCallback   // 重命名完成后回调
	CompleteCallback CompleteCallback // 完成重命名所有流程后回调
}
