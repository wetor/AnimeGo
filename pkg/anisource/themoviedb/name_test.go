package themoviedb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegxStepOne(t *testing.T) {
	res := SimilarText("ダンジョンに出会いを求めるのは間違っているだろうか", "ダンジョンに出会いを求めるのは間違っているだろうか IV")
	fmt.Println(res)
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
		{name: "step 3", args: args{step: 3, str: "オーバーロードIV"}, has: true, want: "オーバーロード"},
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

//func TestRemoveNameSuffix(t *testing.T) {
//	type args struct {
//		name string
//		step int
//	}
//	tests := []struct {
//		name         string
//		args         args
//		wantNextName string
//		wantNextStep int
//		wantErr      assert.ErrorAssertionFunc
//	}{
//		// TODO: Add test cases.
//		{name: "step 0", args: args{step: 0, name: "测试番剧 10期"}, wantNextStep: 1, wantNextName: "测试番剧"},
//		{name: "step 0", args: args{step: 0, name: "测试番剧第2季"}, wantNextStep: 1, wantNextName: "测试番剧"},
//		{name: "step 0", args: args{step: 0, name: "测试番剧八篇"}, wantNextStep: 1, wantNextName: "测试番剧"},
//		{name: "step 0", args: args{step: 0, name: "测试番剧 第二部"}, wantNextStep: 1, wantNextName: "测试番剧"},
//		{name: "step 0", args: args{step: 0, name: "测试番剧 第2"}, wantNextStep: 2, wantNextName: "测试番剧 第"},
//		{name: "step 0", args: args{step: 0, name: "测试番剧 篇"}, wantNextStep: 3, wantNextName: "测试番剧"},
//
//		{name: "step 1", args: args{step: 1, name: "测试番剧 2nd Season"}, wantNextStep: 2, wantNextName: "测试番剧"},
//		{name: "step 1", args: args{step: 1, name: "测试番剧10thSeason"}, wantNextStep: 2, wantNextName: "测试番剧"},
//		{name: "step 1", args: args{step: 1, name: "测试番剧Season 3"}, wantNextStep: 2, wantNextName: "测试番剧"},
//		{name: "step 1", args: args{step: 1, name: "测试番剧 Season 10"}, wantNextStep: 2, wantNextName: "测试番剧"},
//		{name: "step 1", args: args{step: 1, name: "测试番剧 2dn Season"}, wantNextStep: -1, wantNextName: "测试番剧"},
//		{name: "step 1", args: args{step: 1, name: "测试番剧 Season"}, wantNextStep: -1, wantNextName: "测试番剧"},
//
//		{name: "step 2", args: args{step: 2, name: "魔法使いの嫁 詩篇.75 稲妻ジャックと妖精事件"}, wantNextStep: 3, wantNextName: "魔法使いの嫁"},
//		{name: "step 2", args: args{step: 2, name: "蟲師 特別篇 日蝕む翳"}, wantNextStep: 3, wantNextName: "蟲師"},
//		{name: "step 2", args: args{step: 2, name: "めぞん一刻 完結篇"}, wantNextStep: 3, wantNextName: "めぞん一刻"},
//		{name: "step 2", args: args{step: 2, name: "宇宙戦艦ヤマト2199 第二章「太陽圏の死闘」"}, wantNextStep: 3, wantNextName: "宇宙戦艦ヤマト2199"},
//		{name: "step 2", args: args{step: 2, name: "Re:ゼロから始める異世界生活 第四章 聖域と強欲の魔女"}, wantNextStep: 3, wantNextName: "Re:ゼロから始める異世界生活"},
//		{name: "step 2", args: args{step: 2, name: "明星志願3：甜蜜樂章"}, wantNextStep: -1, wantNextName: "明星志願"},
//		{name: "step 2", args: args{step: 2, name: "水浒传之聚义篇"}, wantNextStep: -1, wantNextName: "水浒传之聚义篇"},
//		{name: "step 2", args: args{step: 2, name: "EUREKA/交響詩篇エウレカセブン ハイエボリューション"}, wantNextStep: -1, wantNextName: "EUREKA/交響詩篇エウレカセブン"},
//		{name: "step 2", args: args{step: 2, name: "幻魔大戦 -神話前夜の章-"}, wantNextStep: 3, wantNextName: "幻魔大戦"},
//
//		{name: "step 3", args: args{step: 3, name: "天外魔境II 卍MARU"}, wantNextStep: -1, wantNextName: "天外魔境"},
//		{name: "step 3", args: args{step: 3, name: "Baldur's Gate II: Shadow of Amn"}, wantNextStep: -1, wantNextName: "Baldur's Gate"},
//		{name: "step 3", args: args{step: 3, name: "グローランサーIV ~Wayfarer of the time~"}, wantNextStep: -1, wantNextName: "グローランサー"},
//		{name: "step 3", args: args{step: 3, name: "提督の決断IV"}, wantNextStep: -1, wantNextName: "提督の決断"},
//		{name: "step 3", args: args{step: 3, name: "Baldur's Gate 3: Shadow of Amn"}, wantNextStep: -1, wantNextName: "Baldur's Gate"},
//		{name: "step 3", args: args{step: 3, name: "グローランサー4 ~Wayfarer of the time~"}, wantNextStep: -1, wantNextName: "グローランサー"},
//		{name: "step 3", args: args{step: 3, name: "提督の決断4"}, wantNextStep: -1, wantNextName: "提督の決断"},
//		{name: "step 3", args: args{step: 3, name: "明星志願3：甜蜜樂章"}, wantNextStep: -1, wantNextName: "明星志願"},
//		{name: "step 3", args: args{step: 3, name: "カードファイト!! ヴァンガード will+Dress"}, wantNextStep: -1, wantNextName: "カードファイト!! ヴァンガード"},
//		{name: "step 3", args: args{step: 3, name: "オーバーロードIV"}, wantNextStep: -1, wantNextName: "オーバーロード"},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			gotNextName, gotNextStep, err := RemoveNameSuffix(tt.args.name, tt.args.step)
//			if err != nil {
//				fmt.Printf("%v RemoveNameSuffix(%v, %v)", err, tt.args.name, tt.args.step)
//			}
//			assert.Equalf(t, tt.wantNextName, gotNextName, "RemoveNameSuffix(%v, %v)", tt.args.name, tt.args.step)
//			assert.Equalf(t, tt.wantNextStep, gotNextStep, "RemoveNameSuffix(%v, %v)", tt.args.name, tt.args.step)
//		})
//	}
//}
