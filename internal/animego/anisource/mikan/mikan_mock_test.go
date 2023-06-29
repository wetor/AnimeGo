package mikan_test

import (
	"github.com/brahma-adshonor/gohook"
	"github.com/pkg/errors"

	"github.com/wetor/AnimeGo/internal/animego/anisource"
	"github.com/wetor/AnimeGo/internal/api"
)

type MikanMock struct{}

func (a *MikanMock) Parse(options any) (any, error) {
	return nil, errors.New("ErrNotImplemented")
}

func (a *MikanMock) ParseCache(options any) (any, error) {
	return nil, errors.New("ErrNotImplemented")
}

type BangumiMock struct{}

func (a *BangumiMock) Get(id int, filters any) (any, error) {
	return nil, errors.New("ErrNotImplemented")
}

func (a *BangumiMock) GetCache(id int, filters any) (any, error) {
	return nil, errors.New("ErrNotImplemented")
}

type ThemoviedbMock struct{}

func (a *ThemoviedbMock) Search(name string) (int, error) {
	return 0, errors.New("ErrNotImplemented")
}

func (a *ThemoviedbMock) SearchCache(name string) (int, error) {
	return 0, errors.New("ErrNotImplemented")
}

func (a *ThemoviedbMock) Get(id int, filters any) (any, error) {
	return nil, errors.New("ErrNotImplemented")
}

func (a *ThemoviedbMock) GetCache(id int, filters any) (any, error) {
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
	var err error
	err = gohook.Hook(anisource.Mikan, Mikan, nil)
	if err != nil {
		panic(err)
	}
	err = gohook.Hook(anisource.Bangumi, Bangumi, nil)
	if err != nil {
		panic(err)
	}
	err = gohook.Hook(anisource.Themoviedb, Themoviedb, nil)
	if err != nil {
		panic(err)
	}
}

func UnHook() {
	_ = gohook.UnHook(anisource.Mikan)
	_ = gohook.UnHook(anisource.Bangumi)
	_ = gohook.UnHook(anisource.Themoviedb)
}
