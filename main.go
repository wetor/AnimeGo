package main

import (
	"GoBangumi/utils/logger"
)

func main() {
	logger.Init()
	defer logger.Flush()
	zap.S().Debug("hello world")
}
