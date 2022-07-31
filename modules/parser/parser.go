package parser

import "GoBangumi/models"

type Parser interface {
	Parse(opt *models.ParseOptions) *models.ParseResult
}
