package parser

import "GoBangumi/model"

type Parser interface {
	ParseBangumiName(opt *model.ParseBangumiNameOptions) *model.BangumiName
}
