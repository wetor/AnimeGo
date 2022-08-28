package logger

import (
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestGetLogger(t *testing.T) {
	Init()
	defer Flush()
	zap.S().Infow("failed to fetch URL",
		"url", "http://123.com",
		"attempt", 3,
		"backoff", time.Second,
	)
}
