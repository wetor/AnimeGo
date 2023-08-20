package exceptions

import "fmt"

type ErrDatabaseDirNotFound struct {
	Dir string
}

func (e ErrDatabaseDirNotFound) Error() string {
	return fmt.Sprintf("数据文件夹未找到: %s", e.Dir)
}

type ErrDatabaseDBNotFound struct {
	Name   string
	Season int
	Ep     int
}

func (e ErrDatabaseDBNotFound) Error() string {
	if e.Season > 0 && e.Ep > 0 {
		return fmt.Sprintf("数据未找到: %s[S%d][E%d]", e.Name, e.Season, e.Ep)
	} else if e.Season > 0 {
		return fmt.Sprintf("数据未找到: %s[S%d]", e.Name, e.Season)
	} else {
		return fmt.Sprintf("数据未找到: %s", e.Name)
	}
}

type ErrDatabaseEpisodeExist struct {
	Name   string
	Season int
	Ep     int
}

func (e ErrDatabaseEpisodeExist) Error() string {
	return fmt.Sprintf("数据已存在: %s[S%d][E%d]", e.Name, e.Season, e.Ep)
}
