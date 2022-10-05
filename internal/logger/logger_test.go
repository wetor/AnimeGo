package logger

import (
	"context"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestGetLogger(t *testing.T) {
	Init(&InitOptions{
		File:    "debug.log",
		Debug:   true,
		Context: context.Background(),
	})
	defer Flush()
	zap.S().Infow("failed to fetch URL",
		"url", "http://123.com",
		"attempt", 3,
		"backoff", time.Second,
	)
}
