package utils_test

import (
	"reflect"
	"testing"

	"github.com/wetor/AnimeGo/internal/animego/parser/utils"
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.ParseEp(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
