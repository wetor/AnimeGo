package logger

import (
	"strings"
	"time"

	"github.com/wetor/AnimeGo/pkg/log"
)

type CronLoggerAdapter struct {
}

func NewCronLoggerAdapter() *CronLoggerAdapter {
	return &CronLoggerAdapter{}
}

// Info logs routine messages about cron's operation.
func (c CronLoggerAdapter) Info(msg string, keysAndValues ...interface{}) {
	keysAndValues = formatTimes(keysAndValues)
	log.Debugf(formatString(len(keysAndValues)),
		append([]interface{}{msg}, keysAndValues...)...)
}

// Error logs an error condition.
func (c CronLoggerAdapter) Error(err error, msg string, keysAndValues ...interface{}) {
	keysAndValues = formatTimes(keysAndValues)
	log.Debugf(formatString(len(keysAndValues)+2),
		append([]interface{}{msg, "error", err}, keysAndValues...)...)
}

// formatString returns a logfmt-like format string for the number of
// key/values.
func formatString(numKeysAndValues int) string {
	var sb strings.Builder
	sb.WriteString("%s")
	if numKeysAndValues > 0 {
		sb.WriteString(", ")
	}
	for i := 0; i < numKeysAndValues/2; i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("%v=%v")
	}
	return sb.String()
}

// formatTimes formats any time.Time values as RFC3339.
func formatTimes(keysAndValues []interface{}) []interface{} {
	var formattedArgs []interface{}
	for _, arg := range keysAndValues {
		if t, ok := arg.(time.Time); ok {
			arg = t.Format(time.RFC3339)
		}
		formattedArgs = append(formattedArgs, arg)
	}
	return formattedArgs
}
