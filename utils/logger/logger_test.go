package logger

import (
	"go.uber.org/zap"
	"testing"
)

func TestGetLogger(t *testing.T) {
	Init()
	defer Flush()
	zap.S().Info("test", "hello world1111")
}
