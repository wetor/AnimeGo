package bangumi

import "GoBangumi/model"

type Bangumi interface {
	Parse(opt *model.BangumiParseOptions) *model.Bangumi
}
