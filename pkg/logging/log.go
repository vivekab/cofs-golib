package goliblogging

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/debug"
	"strings"

	lr "github.com/sirupsen/logrus"
	golibconstants "github.com/vivekab/golib/pkg/constants"
	"google.golang.org/grpc/metadata"
)

const (
	stackKey                 = "stack"
	reportLocationKey        = "reportLocation"
	errorInfoKey             = "errorInfo"
	requestIdKey             = "request_id"
	applicationCodeLocation1 = "/app/cmd"
	applicationCodeLocation2 = "/app/internal"
	interceptorLocation      = "/app/vendor/github.com/vivekab/golib/grpc/interceptor.go"
	qldbLocation             = "/app/vendor/github.com/vivekab/golib/qldb"
)

func getRequestIdFromContext(ctx context.Context) string {
	data, ok := metadata.FromIncomingContext(ctx)
	requestId := ""
	if ok {
		if v, ok := data[golibconstants.HeaderRequestID]; ok {
			requestId = v[0]
		}
	} else {
		data, ok = metadata.FromOutgoingContext(ctx)
		if ok {
			if v, ok := data[golibconstants.HeaderRequestID]; ok {
				requestId = v[0]
			}
		}
	}
	return requestId
}

type Logger interface {
	Debug(string, ...interface{})
	DebugD(context.Context, string, Fields)
	Info(string, ...interface{})
	InfoD(context.Context, string, Fields)
	Warn(string, ...interface{})
	WarnD(context.Context, string, Fields)
	Error(string, ...interface{})
	ErrorD(context.Context, string, error, Fields)
	Panic(...interface{})
	PanicD(context.Context, string, Fields)
	Fatal(...interface{})
	FatalD(context.Context, string, Fields)
	Log(string, ...interface{})
	SetLevel(level lr.Level)
	AddHook(hook Hook)
	SetEnv(env string)
	StandardLogger() *lr.Logger
}

type logger struct {
	l   *lr.Logger
	env string
}

type Fields lr.Fields
type Hook lr.Hook

// NewLogger returns a logger interface
func NewLogger(opts ...Option) Logger {
	//Default args
	la := &loggerOpts{
		output: os.Stdout,
		format: &lr.JSONFormatter{},
	}

	for _, opt := range opts {
		opt(la)
	}

	l := lr.New()

	l.SetOutput(la.output)
	l.SetFormatter(la.format)

	return &logger{
		l: l,
	}
}

type loggerOpts struct {
	output io.Writer
	format lr.Formatter
}

// Option optional funcs passed into NewLogger
type Option func(*loggerOpts)

// SetOutput is optionally passed into NewLogger it's used to set the log output
func SetOutput(i io.Writer) Option {
	return func(opts *loggerOpts) {
		opts.output = i
	}
}

// SetFormat is optionally passed into NewLogger it's used to set the log format
func SetFormat(f lr.Formatter) Option {
	return func(opts *loggerOpts) {
		opts.format = f
	}
}

func (l *logger) SetLevel(level lr.Level) {
	l.l.SetLevel(level)
}
func (l *logger) SetEnv(env string) {
	l.env = env
}

func (l *logger) AddHook(hook Hook) {
	l.l.AddHook(hook)
}

// Debug will log debug level logs, will not log in production
func (l *logger) Debug(s string, fs ...interface{}) {
	f := getFields(fs)
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).Debug(s, fs)
}

// DebugD will log debug level logs with Fields, will not log in production
func (l *logger) DebugD(ctx context.Context, s string, f Fields) {
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).WithFields(lr.Fields{
		requestIdKey: getRequestIdFromContext(ctx),
	}).Debug(s)
}

// Info will log info level logs, will appear in production logs
func (l *logger) Info(s string, fs ...interface{}) {
	f := getFields(fs)
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).Info(s, fs)
}
func (l *logger) Infof(ctx context.Context, s string, args ...interface{}) {
	f := Fields{}
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).WithFields(lr.Fields{
		requestIdKey: getRequestIdFromContext(ctx),
	}).Info(s, args)
}

func (l *logger) InfoD(ctx context.Context, s string, f Fields) {
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).WithFields(lr.Fields{
		requestIdKey: getRequestIdFromContext(ctx),
	}).Info(s)
}

// Warn will log warn level logs and a stacktrace, will appear in production logs
func (l *logger) Warn(s string, fs ...interface{}) {
	f := getFields(fs)
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).Warn(s, fs)
}

func (l *logger) WarnD(ctx context.Context, s string, f Fields) {
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).WithFields(lr.Fields{
		requestIdKey: getRequestIdFromContext(ctx),
	}).Warn(s)
}

// Error will log error level logs and a stacktrace, will appear in production logs
func (l *logger) Error(s string, fs ...interface{}) {
	f := getFields(fs)
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).Error(s, fs)
}

func (l *logger) Errorf(ctx context.Context, s string, args ...interface{}) {
	f := Fields{}
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).WithFields(lr.Fields{
		requestIdKey: getRequestIdFromContext(ctx),
	}).Error(s, args)
}

func (l *logger) ErrorD(ctx context.Context, s string, err error, f Fields) {
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).WithFields(lr.Fields{
		requestIdKey: getRequestIdFromContext(ctx),
		errorInfoKey: err,
	}).Error(s)
}

// Panic will log panic level logs and a stacktrace, will appear in production logs
func (l *logger) Panic(fs ...interface{}) {
	f := getFields(fs)
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).Panic(fs...)
}
func (l *logger) PanicD(ctx context.Context, s string, f Fields) {
	f = l.appendStack(f)
	f = l.appendReportLocation(f)
	f[requestIdKey] = getRequestIdFromContext(ctx)
	l.l.WithFields(f.format()).Panic(s)
}

func (l *logger) FatalD(ctx context.Context, s string, f Fields) {
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).WithFields(lr.Fields{
		requestIdKey: getRequestIdFromContext(ctx),
	}).Fatal(s)
}

func (l *logger) Fatal(fs ...interface{}) {
	f := getFields(fs)
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).Fatal(fs...)
}

func (l *logger) Log(s string, fs ...interface{}) {
	f := getFields(fs)
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).Infof(s, fs...)
}

func (l *logger) StandardLogger() *lr.Logger {
	return l.l
}

func (l *logger) appendStack(f Fields) Fields {
	f[stackKey] = string(debug.Stack())

	return f
}

func (l *logger) appendReportLocation(f Fields) Fields {
	if l.env == golibconstants.EnvProd {
		return f
	}
	hash := map[string]string{}

	// we traverse the stack frames and add any relevant application code location to the hash
	for i := 1; i < 15; i++ {
		pc, file, line, ok := runtime.Caller(i)
		addFrame := ok && (strings.HasPrefix(file, applicationCodeLocation1) || strings.HasPrefix(file, applicationCodeLocation2) || strings.HasPrefix(file, interceptorLocation) || strings.HasPrefix(file, qldbLocation))
		if addFrame {
			key := fmt.Sprintf("%s:%d:%d", file, line, i)
			hash[key] = runtime.FuncForPC(pc).Name()
		}
	}
	f[reportLocationKey] = hash
	return f
}

// Format to logrus formatted fields
// We could also put any log data sanitisazation in here
func (f Fields) format() lr.Fields {
	return lr.Fields(f)
}

// fs passed in will be in the form of "method", req.Method, so we only want an even number of fs, otherwise just skip it
func getFields(vfs ...interface{}) Fields {
	fields := Fields{}
	if len(vfs) == 0 {
		return fields
	}
	fs := vfs[0].([]interface{})
	if len(fs) > 0 {
		if f, ok := fs[0].([]interface{}); ok {
			if len(f) > 0 {
				if val, ok := f[0].(Fields); ok {
					fields = val
				}
			}
		}
	}
	return fields
}
