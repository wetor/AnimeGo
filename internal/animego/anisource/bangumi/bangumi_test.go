package bangumi_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/wetor/AnimeGo/internal/animego/anisource/bangumi"
	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/internal/pkg/request"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/test"
)

var (
	bangumiInst *bangumi.Bangumi
	ctx, cancel = context.WithCancel(context.Background())
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	_ = utils.CreateMutiDir("data")
	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
	})
	host := test.MockBangumiStart(ctx)
	request.Init(&request.Options{
		Host: map[string]*request.HostOptions{
			constant.BangumiHost: {
				Redirect: host,
			},
		},
	})
	mutex := sync.Mutex{}

	db := cache.NewBolt()
	db.Open("data/bolt.db")

	bangumiCache := cache.NewBolt(true)
	bangumiCache.Open(test.GetDataPath("", "bolt_sub.bolt"))
	bangumiInst = bangumi.NewBangumi(&bangumi.Options{
		Cache:            db,
		BangumiCache:     bangumiCache,
		BangumiCacheLock: &mutex,
	})
	m.Run()

	cancel()
	db.Close()
	bangumiCache.Close()
	_ = log.Close()
	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func TestBangumi_Get_GetCache(t *testing.T) {
	type args struct {
		bangumiID int
		ep        int
	}
	tests := []struct {
		name       string
		onlyCache  bool
		args       args
		wantEntity *bangumi.Entity
		wantErr    error
		wantErrStr string
	}{
		// TODO: Add test cases.
		{
			name:       "联盟空军航空魔法音乐队 光辉魔女",
			args:       args{bangumiID: 253047, ep: 3},
			wantEntity: &bangumi.Entity{ID: 253047, NameCN: "联盟空军航空魔法音乐队 光辉魔女", Name: "連盟空軍航空魔法音楽隊 ルミナスウィッチーズ", Eps: 12, AirDate: "2022-07-03", Type: 0, Platform: 0},
			wantErr:    nil,
		},
		{
			name:       "欢迎来到实力至上主义的教室 第二季",
			args:       args{bangumiID: 371546, ep: 5},
			wantEntity: &bangumi.Entity{ID: 371546, NameCN: "欢迎来到实力至上主义教室 第二季", Name: "ようこそ実力至上主義の教室へ 2nd Season", Eps: 13, AirDate: "2022-07-04", Type: 0, Platform: 0},
			wantErr:    nil,
		},
		{
			name:       "CLANNAD",
			onlyCache:  true,
			args:       args{bangumiID: 51, ep: 5},
			wantEntity: &bangumi.Entity{ID: 51, NameCN: "CLANNAD", Name: "CLANNAD -クラナド-", Eps: 22, AirDate: "2007-10-04", Type: 2, Platform: 1},
			wantErr:    nil,
		},
		{
			name:       "err_request",
			args:       args{bangumiID: 114514, ep: 5},
			wantEntity: nil,
			wantErr:    &exceptions.ErrRequest{},
			wantErrStr: "获取Bangumi信息失败: 请求 Bangumi 失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEntity, err := bangumiInst.GetCache(tt.args.bangumiID, nil)
			if tt.wantErr != nil {
				assert.IsType(t, tt.wantErr, errors.Cause(err))
				assert.EqualError(t, err, tt.wantErrStr)
			} else {
				assert.NoError(t, err)
				assert.Equalf(t, tt.wantEntity, gotEntity, "GetCache(%v)", tt.args.bangumiID)
			}
		})
		if tt.onlyCache {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			gotEntity, err := bangumiInst.Get(tt.args.bangumiID, nil)
			if tt.wantErr != nil {
				assert.IsType(t, tt.wantErr, errors.Cause(err))
				assert.EqualError(t, err, tt.wantErrStr)
			} else {
				assert.NoError(t, err)
				assert.Equalf(t, tt.wantEntity, gotEntity, "Get(%v)", tt.args.bangumiID)
			}
		})
	}
}

func TestBangumi_Search_SearchCache(t *testing.T) {
	t.Skip("跳过Bangumi Search测试")
	type args struct {
		name string
	}
	tests := []struct {
		name       string
		onlyCache  bool
		args       args
		wantID     int
		wantErr    error
		wantErrStr string
	}{
		{
			name:    "联盟空军航空魔法音乐队 光辉魔女",
			args:    args{name: "联盟空军航空魔法音乐队 光辉魔女"},
			wantID:  253047,
			wantErr: nil,
		},
		{
			name:    "欢迎来到实力至上主义的教室 第二季",
			args:    args{name: "欢迎来到实力至上主义的教室 第二季"},
			wantID:  371546,
			wantErr: nil,
		},
		{
			name:      "CLANNAD",
			onlyCache: true,
			args:      args{name: "CLANNAD"},
			wantID:    51,
			wantErr:   nil,
		},
		{
			name:       "err_request",
			args:       args{name: "联盟空军航空魔法音乐队 光辉魔女"},
			wantID:     0,
			wantErr:    &exceptions.ErrRequest{},
			wantErrStr: "获取Bangumi信息失败: 请求 Bangumi 失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEntity, err := bangumiInst.SearchCache(tt.args.name, nil)
			if tt.wantErr != nil {
				assert.IsType(t, tt.wantErr, errors.Cause(err))
				assert.EqualError(t, err, tt.wantErrStr)
			} else {
				assert.NoError(t, err)
				assert.Equalf(t, tt.wantID, gotEntity, "SearchCache(%v)", tt.args.name)
			}
		})
		if tt.onlyCache {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			gotEntity, err := bangumiInst.Search(tt.args.name, nil)
			if tt.wantErr != nil {
				assert.IsType(t, tt.wantErr, errors.Cause(err))
				assert.EqualError(t, err, tt.wantErrStr)
			} else {
				assert.NoError(t, err)
				assert.Equalf(t, tt.wantID, gotEntity, "Search(%v)", tt.args.name)
			}
		})
	}
}
