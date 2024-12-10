package goliblogging

import (
	"context"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	log = NewLogger()
)

func AddHook(hook Hook) {
	log.AddHook(hook)
}

func StandardLogger() *logrus.Logger {
	return log.StandardLogger()
}

func SetLevel(level logrus.Level) {
	log.SetLevel(level)
}

func Debug(s string, fs ...interface{}) {
	log.Debug(s, fs...)
}

func DebugD(ctx context.Context, s string, f Fields) {
	log.DebugD(ctx, s, f)
}

func Info(s string, fs ...interface{}) {
	log.Info(s, fs...)
}

func InfoD(ctx context.Context, s string, f Fields) {
	log.InfoD(ctx, s, f)
}

func Warn(s string, fs ...interface{}) {
	log.Warn(s, fs...)
}

func WarnD(ctx context.Context, s string, f Fields) {
	log.WarnD(ctx, s, f)
}

func Error(s string, fs ...interface{}) {
	log.Error(s, fs...)
}

func ErrorD(ctx context.Context, s string, err error, f Fields) {
	log.ErrorD(ctx, s, err, f)
}

func Panic(fs ...interface{}) {
	log.Panic(fs...)
}

func PanicD(ctx context.Context, s string, f Fields) {
	log.PanicD(ctx, s, f)
}

func Fatal(fs ...interface{}) {
	log.Fatal(fs...)
}

func FatalD(ctx context.Context, s string, f Fields) {
	log.FatalD(ctx, s, f)
}

func SetupLogging(env string, level string) {
	log.SetEnv(env)
	switch level {
	case "DEBUG":
		log.SetLevel(logrus.DebugLevel)
	case "INFO":
		log.SetLevel(logrus.InfoLevel)
	default:
		log.SetLevel(logrus.ErrorLevel)
	}
}

func GetLogger() Logger {
	return log
}

// PrettifyStack takes the stacktrace produced by debug.Stack() and converts it to
// easily conusmable data
func PrettifyStack(stack string) string {
	lines := strings.Split(strings.TrimSpace(stack), "\n")

	// Strip the first "goroutine" line
	if len(lines) > 0 {
		if first := lines[0]; strings.HasPrefix(first, "goroutine ") && strings.HasSuffix(first, ":") {
			lines = lines[1:]
		}
	}

	sb := strings.Builder{}

	for _, line := range lines {
		// Indented lines are source locations
		if strings.HasPrefix(line, "\t") {
			line = line[1:]
			if offset := strings.LastIndex(line, " +0x"); offset != -1 {
				line = line[:offset]
			}
			sb.WriteString(" (")
			sb.WriteString(line)
			sb.WriteString(")")
			continue
		}

		// Other lines are function calls
		if sb.Len() > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(line)
	}

	return sb.String()
}
