package mikan_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/wetor/AnimeGo/internal/animego/anidata"
	"github.com/wetor/AnimeGo/internal/animego/anidata/mikan"
	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/test"
)

const testdata = "mikan"

func TestMain(m *testing.M) {
	fmt.Println("begin")
	log.Init(&log.Options{
		File:  "data/test.log",
		Debug: true,
	})
	test.HookAll(testdata, nil)
	defer test.UnHook()

	db := cache.NewBolt()
	db.Open("data/bolt.db")
	anidata.Init(&anidata.Options{Cache: db})

	m.Run()

	db.Close()
	_ = log.Close()
	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func TestMikan_Parse_ParseCache(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name       string
		args       args
		wantEntity *mikan.Entity
		wantErr    error
		wantErrStr string
	}{
		// TODO: Add test cases.
		{
			name:       "海贼王",
			args:       args{url: "https://mikanani.me/Home/Episode/18b60d48a72c603b421468aade7fdd0868ff2f2f"},
			wantEntity: &mikan.Entity{MikanID: 228, BangumiID: 975},
			wantErr:    nil,
		},
		{
			name:       "欢迎来到实力至上主义的教室 第二季",
			args:       args{url: "https://mikanani.me/Home/Episode/8849c25e05d6e2623b5333bc78d3a489a9b1cc59"},
			wantEntity: &mikan.Entity{MikanID: 2775, BangumiID: 371546},
			wantErr:    nil,
		},
		{
			name:       "err_request",
			args:       args{url: "https://mikanani.me/Home/Episode/err_request_not_found"},
			wantEntity: nil,
			wantErr:    &exceptions.ErrRequest{},
			wantErrStr: "解析Mikan信息失败: 请求 Mikan 失败",
		},
		{
			name:       "err_parse_mikan_url_html",
			args:       args{url: "https://mikanani.me/Home/Episode/err_parse_mikan_url_html"},
			wantEntity: nil,
			wantErr:    &exceptions.ErrMikanParseHTML{},
			wantErrStr: "解析Mikan信息失败: 解析 MikanUrl 失败，解析网页错误",
		},
		{
			name:       "err_parse_not_mikan_id_html",
			args:       args{url: "https://mikanani.me/Home/Episode/err_parse_not_mikan_id_html"},
			wantEntity: nil,
			wantErr:    &exceptions.ErrMikanParseHTML{},
			wantErrStr: "解析Mikan信息失败: 解析 MikanID 失败，解析网页错误",
		},
		{
			name:       "err_parse_mikan_id_html",
			args:       args{url: "https://mikanani.me/Home/Episode/err_parse_mikan_id_html"},
			wantEntity: nil,
			wantErr:    &exceptions.ErrMikanParseHTML{},
			wantErrStr: "解析Mikan信息失败: 解析 MikanID 失败，解析网页错误",
		},
		{
			name:       "err_parse_pub_group_id_html",
			args:       args{url: "https://mikanani.me/Home/Episode/err_parse_pub_group_id_html"},
			wantEntity: nil,
			wantErr:    &exceptions.ErrMikanParseHTML{},
			wantErrStr: "解析Mikan信息失败: 解析 PubGroupID 失败，解析网页错误",
		},
		{
			name:       "err_parse_bangumi_id_html_114514",
			args:       args{url: "https://mikanani.me/Home/Episode/err_parse_bangumi_id_html_114514"},
			wantEntity: nil,
			wantErr:    &exceptions.ErrMikanParseHTML{},
			wantErrStr: "解析Mikan BangumiID失败: 解析 BangumiID 失败，解析网页错误",
		},
	}
	m := &mikan.Mikan{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEntity, err := m.ParseCache(tt.args.url)
			if tt.wantErr != nil {
				assert.IsType(t, tt.wantErr, errors.Cause(err))
				assert.EqualError(t, err, tt.wantErrStr)
			} else {
				assert.NoError(t, err)
				assert.Equalf(t, tt.wantEntity, gotEntity, "ParseCache(%v)", tt.args.url)
			}
		})

		t.Run(tt.name, func(t *testing.T) {
			gotEntity, err := m.Parse(tt.args.url)
			if tt.wantErr != nil {
				assert.IsType(t, tt.wantErr, errors.Cause(err))
				assert.EqualError(t, err, tt.wantErrStr)
			} else {
				assert.NoError(t, err)
				assert.Equalf(t, tt.wantEntity, gotEntity, "Parse(%v)", tt.args.url)
			}
		})
	}
}
