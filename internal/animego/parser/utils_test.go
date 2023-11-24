package parser_test

import (
	"github.com/wetor/AnimeGo/internal/animego/parser"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				name: "[orion origin] Benriya Saitou-san, Isekai ni Iku [04] [1080p] [H265 AAC] [CHT].mp4",
			},
			want: 4,
		},
		{
			name: "2",
			args: args{
				name: "[orion origin] Benriya Saitou-san, Isekai ni Iku [11] [1080p] [H265 AAC] [CHT].mp4",
			},
			want: 11,
		},
		//{
		//	name: "3",
		//	args: args{
		//		name: "[Nekomoe kissaten] Seijo no Maryoku wa Bannou Desu 06 [WebRip 1080p HEVC-10bit AAC ASSx2].mkv",
		//	},
		//	want: 6,
		//},
		{
			name: "5",
			args: args{
				name: "[Nekomoe kissaten&LoliHouse] World Dai Star - 11 [WebRip 1080p HEVC-10bit AAC ASSx2].mkv",
			},
			want: 11,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parser.ParseEp(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
