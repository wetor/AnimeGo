package parser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBangumiEp_Parse(t *testing.T) {
	ep, err := ParseEp("【幻樱字幕组】【4月新番】【古见同学有交流障碍症 第二季 Komi-san wa, Komyushou Desu. S02】【22】【GB_MP4】【1920X1080】")
	fmt.Println(ep, err)
}

func TestBangumiEp_Parse2(t *testing.T) {
	tests := []struct {
		title string
		want  int
	}{
		{title: "【幻樱字幕组】【4月新番】【古见同学有交流障碍症 第二季 Komi-san wa, Komyushou Desu. S02】【22】【GB_MP4】【1920X1080】", want: 22},
		{title: "【极影字幕社】LoveLive! 虹咲学园学园偶像同好会 第2期 第12集 GB_CN HEVC_opus 1080p [复制磁连]", want: 12},
		{title: "[虹咲学园烤肉同好会][Love Live! 虹咲学园学园偶像同好会 第二季][03][简日内嵌][特效歌词][WebRip][1080p][AVC AAC MP4]", want: 3},
		{title: "[LoliHouse] Love Live! 虹咲学园学园偶像同好会 第二季 / Love Live! Nijigasaki S2 - 09 [WebRip 1080p HEVC-10bit AAC][简繁内封字幕]", want: 9},

		{title: "[NaN-Raws]Love_Live！虹咲学园_学园偶像同好会_第二季[10][Bahamut][WEB-DL][1080P][AVC_AAC][CHT][MP4][bangumi.online]", want: 10},
		{title: "[Lilith-Raws x WitEx.io] Love Live！虹咲学园 学园偶像同好会 S02 - 08 [Baha][WEB-DL][1080p][AVC AAC][CHT][MP4]", want: 8},
		{title: "[ANi] Love Live！虹咲学园 学园偶像同好会 第二季 - 12 [1080P][Baha][WEB-DL][AAC AVC][CHT][MP4]", want: 12},
		{title: "[NC-Raws] Love Live！虹咲学园 学园偶像同好会 第二季 / Nijigasaki S2 - 12 (Baha 1920x1080 AVC AAC MP4)", want: 12},

		{title: "Love Live！虹咲学园偶像同好会 第2季 第8集/Love Live! Nijigasaki - 08 (1080P)(AVC AAC)(CHS)", want: 8},
		{title: "[酷漫404][辉夜大小姐想让我告白 一终极浪漫一][10][1080P][WebRip][简日双语][AVC AAC][MP4][字幕组招人内详]", want: 10},
		{title: "[云光字幕组]辉夜大小姐想让我告白 -超级浪漫- Kaguya-sama wa Kokurasetai S3 [08][简体双语][1080p]招募后期", want: 8},
		{title: "[猎户不鸽发布组] 辉夜大小姐想让我告白？第3季 ~超级浪漫 ~ Kaguya-sama wa Kokurasetai S3 [11] [1080p+] [简中] [2022年4月番] [复制磁连]", want: 11},
		{title: "【极影字幕社+辉夜汉化组】辉夜大小姐想让我告白 究极浪漫 第10集 GB_CN HEVC opus 1080p [复制磁连]", want: 10},
		{title: "[Skymoon-Raws] 辉夜姬想让人告白 一超级浪漫一 / Kaguya-sama wa Kokurasetai S03 - 11 [ViuTV][WEB-RIP][1080p][AVC AAC][CHT][SRT][MKV](先行版本) ", want: 11},

		{title: "[澄空学园&雪飘工作室][辉夜大小姐想让我告白 第三季 / かぐや様は告らせたい 三期 / Kaguya-sama wa Kokurasetai Season 3][05][720p][繁体内嵌]", want: 5},
		{title: "[MingY] 辉夜大小姐想让我告白-Ultra Romantic-​ / Kaguya-sama wa Kokurasetai​ S3 [01][1080p][CHS] [复制磁连]", want: 1},
	}
	for _, tt := range tests {
		g, err := ParseEp(tt.title)
		fmt.Println(g, tt.title)
		assert.Equal(t, err, nil)
		assert.Equal(t, g, tt.want)

	}
}
