package bangumi

import (
	"GoBangumi/models"
)

type Bangumi interface {
	Parse(opt *models.BangumiParseOptions) *models.Bangumi
}
