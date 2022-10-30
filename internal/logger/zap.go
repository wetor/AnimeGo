package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

const (
	logTmFmt = "2006-01-02 15:04:05"
)

var (
	debug bool
)

func GetLogger(opt *InitOptions) {
	debug = opt.Debug
	Encoder := GetEncoder()
	level := zapcore.DebugLevel
	if !debug {
		level = zapcore.InfoLevel
	}
	newCore := zapcore.NewTee(
		zapcore.NewCore(Encoder, GetWriteSyncer(opt.File), zapcore.DebugLevel), // 写入文件
		zapcore.NewCore(Encoder, zapcore.Lock(os.Stdout), level),               // 写入控制台
	)
	logger := zap.New(newCore, zap.AddCaller())
	zap.ReplaceGlobals(logger)
}

// GetEncoder 自定义的Encoder
func GetEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(
		zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller_line",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    cEncodeLevel,
			EncodeTime:     cEncodeTime,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   cEncodeCaller,
		})
}

// GetWriteSyncer 自定义的WriteSyncer
func GetWriteSyncer(file string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file,
		MaxSize:    50,
		MaxBackups: 10,
		MaxAge:     30,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// cEncodeLevel 自定义日志级别显示
func cEncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

// cEncodeTime 自定义时间格式显示
func cEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format(logTmFmt) + "]")
}

// cEncodeCaller 自定义行号显示
func cEncodeCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	if debug {
		enc.AppendString("[" + caller.TrimmedPath() + "]")
	}

}
