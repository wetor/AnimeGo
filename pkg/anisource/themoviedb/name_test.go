package themoviedb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestRegxStepOne(t *testing.T) {
	var str string

	str = "测试番剧VII测试"
	// 匹配I II III IV V VI VII VIII IX X
	step3 := regexp.MustCompile(`\s?(V|X|IX|IV|V?I{1,3}).*$`)
	str = step3.ReplaceAllString(str, "")
	fmt.Println(str)

}

func TestRegxStep(t *testing.T) {

	type args struct {
		step int
		str  string
	}
	tests := []struct {
		name string
		args args
		has  bool
		want string
	}{
		{name: "step 0", args: args{step: 0, str: "测试番剧 10期"}, has: true, want: "测试番剧"},
		{name: "step 0", args: args{step: 0, str: "测试番剧第2季"}, has: true, want: "测试番剧"},
		{name: "step 0", args: args{step: 0, str: "测试番剧八篇"}, has: true, want: "测试番剧"},
		{name: "step 0", args: args{step: 0, str: "测试番剧 第二部"}, has: true, want: "测试番剧"},
		{name: "step 0", args: args{step: 0, str: "测试番剧 第2"}, has: false, want: ""},
		{name: "step 0", args: args{step: 0, str: "测试番剧 篇"}, has: false, want: ""},

		{name: "step 1", args: args{step: 1, str: "测试番剧 2nd Season"}, has: true, want: "测试番剧"},
		{name: "step 1", args: args{step: 1, str: "测试番剧10thSeason"}, has: true, want: "测试番剧"},
		{name: "step 1", args: args{step: 1, str: "测试番剧Season 3"}, has: true, want: "测试番剧"},
		{name: "step 1", args: args{step: 1, str: "测试番剧 Season 10"}, has: true, want: "测试番剧"},
		{name: "step 1", args: args{step: 1, str: "测试番剧 2dn Season"}, has: false, want: ""},
		{name: "step 1", args: args{step: 1, str: "测试番剧 Season"}, has: false, want: ""},

		{name: "step 2", args: args{step: 2, str: "水浒传之聚义篇"}, has: false, want: ""},
		{name: "step 2", args: args{step: 2, str: "EUREKA/交響詩篇エウレカセブン ハイエボリューション"}, has: false, want: ""},
		{name: "step 2", args: args{step: 2, str: "魔法使いの嫁 詩篇.75 稲妻ジャックと妖精事件"}, has: true, want: "魔法使いの嫁"},
		{name: "step 2", args: args{step: 2, str: "蟲師 特別篇 日蝕む翳"}, has: true, want: "蟲師"},
		{name: "step 2", args: args{step: 2, str: "めぞん一刻 完結篇"}, has: true, want: "めぞん一刻"},
		{name: "step 2", args: args{step: 2, str: "宇宙戦艦ヤマト2199 第二章「太陽圏の死闘」"}, has: true, want: "宇宙戦艦ヤマト2199"},
		{name: "step 2", args: args{step: 2, str: "明星志願3：甜蜜樂章"}, has: false, want: ""},
		{name: "step 2", args: args{step: 2, str: "Re:ゼロから始める異世界生活 第四章 聖域と強欲の魔女"}, has: true, want: "Re:ゼロから始める異世界生活"},
		{name: "step 2", args: args{step: 2, str: "幻魔大戦 -神話前夜の章-"}, has: true, want: "幻魔大戦"},

		{name: "step 3", args: args{step: 3, str: "天外魔境II 卍MARU"}, has: true, want: "天外魔境"},
		{name: "step 3", args: args{step: 3, str: "Baldur's Gate II: Shadow of Amn"}, has: true, want: "Baldur's Gate"},
		{name: "step 3", args: args{step: 3, str: "グローランサーIV ~Wayfarer of the time~"}, has: true, want: "グローランサー"},
		{name: "step 3", args: args{step: 3, str: "提督の決断IV"}, has: true, want: "提督の決断"},
		{name: "step 3", args: args{step: 3, str: "Baldur's Gate 3: Shadow of Amn"}, has: true, want: "Baldur's Gate"},
		{name: "step 3", args: args{step: 3, str: "グローランサー4 ~Wayfarer of the time~"}, has: true, want: "グローランサー"},
		{name: "step 3", args: args{step: 3, str: "提督の決断4"}, has: true, want: "提督の決断"},
		{name: "step 3", args: args{step: 3, str: "明星志願3：甜蜜樂章"}, has: true, want: "明星志願"},
		{name: "step 3", args: args{step: 3, str: "カードファイト!! ヴァンガード will+Dress"}, has: true, want: "カードファイト!! ヴァンガード"},
	}
	for _, tt := range tests {
		t.Log(tt)
		has := nameRegxStep[tt.args.step].MatchString(tt.args.str)
		assert.Equal(t, tt.has, has)
		if has {
			got := nameRegxStep[tt.args.step].ReplaceAllString(tt.args.str, "")
			assert.Equal(t, tt.want, got)
		}
	}
}
