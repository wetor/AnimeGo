package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	logTmFmt = "2006-01-02 15:04:05"
)

func GetLogger(opt *Options) {
	level := zapcore.DebugLevel
	if !opt.Debug {
		level = zapcore.InfoLevel
	}
	newCore := zapcore.NewTee(
		zapcore.NewCore(GetEncoder(true, opt.Debug), GetWriteSyncer(opt.File), level), // 写入文件
		zapcore.NewCore(GetEncoder(false, opt.Debug), zapcore.Lock(os.Stdout), level), // 写入控制台
	)
	logger := zap.New(newCore, zap.AddCaller())
	zap.ReplaceGlobals(logger)
}

// GetEncoder 自定义的Encoder
func GetEncoder(file, debug bool) zapcore.Encoder {
	conf := zapcore.EncoderConfig{
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
	}
	if file {
		conf.EncodeLevel = cEncodeLevel
	} else {
		conf.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	if debug {
		conf.EncodeCaller = cEncodeCaller
	} else {
		conf.EncodeCaller = nil
	}

	return zapcore.NewConsoleEncoder(conf)
}

// GetWriteSyncer 自定义的WriteSyncer
func GetWriteSyncer(file string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file,
		MaxSize:    16,
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
	enc.AppendString("[" + caller.TrimmedPath() + "]")
}
