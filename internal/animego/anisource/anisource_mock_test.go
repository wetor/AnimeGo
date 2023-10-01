package anisource_test

import (
	"github.com/pkg/errors"
	"github.com/wetor/AnimeGo/internal/animego/anisource"

	"github.com/wetor/AnimeGo/internal/animego/anidata/bangumi"
	"github.com/wetor/AnimeGo/internal/animego/anidata/mikan"
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/exceptions"
)

type MikanMock struct{}

func (a *MikanMock) Parse(options any) (any, error) {
	return nil, errors.New("ErrNotImplemented")
}

func (a *MikanMock) ParseCache(options any) (any, error) {
	url := options.(string)
	switch url {
	case "err_mikan":
		return nil, &exceptions.ErrRequest{Name: url}
	case "err_bangumi":
		return &mikan.Entity{
			MikanID:   114,
			BangumiID: 114,
		}, nil
	case "err_themoviedb_search":
		return &mikan.Entity{
			MikanID:   514,
			BangumiID: 514,
		}, nil
	case "err_themoviedb_get":
		return &mikan.Entity{
			MikanID:   1919,
			BangumiID: 1919,
		}, nil
	}
	return nil, errors.New("ErrNotImplemented")
}

type BangumiMock struct{}

func (a *BangumiMock) Get(id int, filters any) (any, error) {
	return nil, errors.New("ErrNotImplemented")
}

func (a *BangumiMock) GetCache(id int, filters any) (any, error) {
	switch id {
	case 114:
		return nil, &exceptions.ErrRequest{Name: "err_bangumi"}
	case 514:
		return &bangumi.Entity{
			ID:      514,
			Name:    "err_themoviedb_search",
			NameCN:  "err_themoviedb_search",
			AirDate: "1919-05-14",
			Eps:     2,
		}, nil
	case 1919:
		return &bangumi.Entity{
			ID:      1919,
			Name:    "err_themoviedb_get",
			NameCN:  "err_themoviedb_get",
			AirDate: "1919-05-14",
			Eps:     2,
		}, nil
	}
	return nil, errors.New("ErrNotImplemented")
}

type ThemoviedbMock struct{}

func (a *ThemoviedbMock) Search(name string) (int, error) {
	return 0, errors.New("ErrNotImplemented")
}

func (a *ThemoviedbMock) SearchCache(name string) (int, error) {
	switch name {
	case "err_themoviedb_search":
		return 0, &exceptions.ErrRequest{Name: "err_themoviedb_search"}
	case "err_themoviedb_get":
		return 666, nil
	}
	return 0, errors.New("ErrNotImplemented")
}

func (a *ThemoviedbMock) Get(id int, filters any) (any, error) {
	return nil, errors.New("ErrNotImplemented")
}

func (a *ThemoviedbMock) GetCache(id int, filters any) (any, error) {
	switch id {
	case 666:
		return 0, &exceptions.ErrRequest{Name: "err_themoviedb_get"}
	}
	return nil, errors.New("ErrNotImplemented")
}

func Mikan() api.AniDataParse {
	return &MikanMock{}
}

func Bangumi() api.AniDataGet {
	return &BangumiMock{}
}

func Themoviedb(key string) api.AniDataSearchGet {
	return &ThemoviedbMock{}
}

func Hook() {
	anisource.MikanInstance = Mikan()
	anisource.BangumiInstance = Bangumi()
	anisource.ThemoviedbInstance = Themoviedb("")
}

func UnHook() {
	anisource.MikanInstance = nil
	anisource.BangumiInstance = nil
	anisource.ThemoviedbInstance = nil
}
