package feed_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/wetor/AnimeGo/internal/animego/feed"
	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/json"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/test"
)

const testdata = "feed"

func BeforeRss() {
	fmt.Println("begin")
	_ = utils.CreateMutiDir("data")
	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
	})

	test.HookGetWriter(testdata, nil)
	defer test.UnHook()

	raw, _ := test.GetData("feed", "Mikan.xml")
	raw = bytes.Replace(raw, []byte("\r\n"), []byte("\n"), -1)
	_ = os.WriteFile(test.GetDataPath("feed", "Mikan.xml"), raw, os.ModePerm)
}

func loadItems() []*models.FeedItem {
	var items []*models.FeedItem
	data, err := test.GetData("feed", "Mikan.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &items)
	if err != nil {
		panic(err)
	}
	return items
}

func TestRss_Parse(t *testing.T) {
	BeforeRss()
	defer After()
	type fields struct {
		url  string
		file string
		raw  []byte
	}

	raw, err := test.GetData("feed", "Mikan.xml")
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name       string
		fields     fields
		wantItems  []*models.FeedItem
		wantErr    error
		wantErrStr string
	}{
		// TODO: Add test cases.
		{
			name: "mikan_file",
			fields: fields{
				file: test.GetDataPath("feed", "Mikan.xml"),
			},
			wantItems: loadItems(),
		},
		{
			name: "mikan_raw",
			fields: fields{
				raw: raw,
			},
			wantItems: loadItems(),
		},
		{
			name: "skip_and_err_length",
			fields: fields{
				file: test.GetDataPath("feed", "skip_and_err_length.xml"),
			},
			wantItems: []*models.FeedItem{
				{
					MikanUrl:   "https://mikanani.me/Home/Episode/2076477d6a119fae9ad882ecc5fd697c1afaee75",
					Name:       "万事屋斋藤先生转生异世界",
					Date:       "2023-01-23",
					Type:       "application/x-bittorrent",
					TorrentUrl: "https://mikanani.me/Download/20230123/2076477d6a119fae9ad882ecc5fd697c1afaee75.torrent",
					Length:     int64(0),
				},
			},
		},
		{
			name: "err_request",
			fields: fields{
				url: "err_request",
			},
			wantErr:    &exceptions.ErrRequest{},
			wantErrStr: "请求 Rss 失败",
		},
		{
			name:       "null",
			fields:     fields{},
			wantErr:    &exceptions.ErrFeed{},
			wantErrStr: "Rss为空",
		},
		{
			name: "err_not_found",
			fields: fields{
				file: test.GetDataPath("feed", "err_not_found"),
			},
			wantErr:    &exceptions.ErrFeed{},
			wantErrStr: "打开Rss文件失败",
		},
		{
			name: "err_parse_feed",
			fields: fields{
				file: test.GetDataPath("feed", "err_parse_feed.xml"),
			},
			wantErr:    &exceptions.ErrFeed{},
			wantErrStr: "解析Rss失败",
		},
	}
	r := feed.NewRss()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotItems []*models.FeedItem
			var err error
			if len(tt.fields.raw) > 0 {
				gotItems, err = r.Parse(tt.fields.raw)
			} else if len(tt.fields.file) > 0 {
				gotItems, err = r.ParseFile(tt.fields.file)
			} else if len(tt.fields.url) > 0 {
				gotItems, err = r.ParseUrl(tt.fields.url)
			} else {
				gotItems, err = r.ParseFile("")
			}

			if tt.wantErr != nil {
				assert.IsType(t, tt.wantErr, errors.Cause(err))
				assert.EqualError(t, err, tt.wantErrStr)
			} else {
				assert.NoError(t, err)
				assert.Equalf(t, tt.wantItems, gotItems, "Parse()")
			}
		})
	}
}
