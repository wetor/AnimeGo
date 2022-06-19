package parser

import "GoBangumi/models"

type Parser interface {
	ParseBangumiName(opt *models.ParseBangumiNameOptions) *models.BangumiName
}
