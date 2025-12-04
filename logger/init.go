package logger

import (
	"context"
	"errors"
	"github.com/duke-git/lancet/v2/slice"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

type Level zapcore.Level

const (
	DebugLevel = Level(zap.DebugLevel)
	InfoLevel  = Level(zap.InfoLevel)
	WarnLevel  = Level(zap.WarnLevel)
	ErrorLevel = Level(zap.ErrorLevel)
)

// 不导出类型，防止外部包覆盖
type contextKey struct{ name string }

var (
	traceIDKey = contextKey{"trace_id"}
	spanIDKey  = contextKey{"span_id"}
)

var (
	isDebugMode        = false
	logPath            string
	logLevel           = InfoLevel
	logJson            = false
	DefaultLogger      *zap.Logger
	DefaultSugarLogger *zap.SugaredLogger
)

// EnableDebug
// 1. print log to console for debug purpose
// 2. set logger level as DEBUG
func enableDebug() {
	isDebugMode = true
	logLevel = DebugLevel
}

func InitWithConfig(config Config) error {
	if debug, ok := os.LookupEnv("DEBUG"); ok {
		if slice.Contain([]string{"1", "true", "on", "enalbe"}, strings.ToLower(debug)) {
			enableDebug()
		}
	}
	if config.AppName == "" {
		return errors.New("empty app name")
	}

	homeDir, err := os.UserHomeDir()
	logRoot := ""
	if err == nil {
		logRoot = path.Join(homeDir, "logs")
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			logRoot = path.Join(cwd, "logs")
		}
	}
	logPath = path.Join(logRoot, config.AppName, "main.log")

	initMainLogger()
	return nil
}

func initMainLogger() {
	// 自定义输出：https://cloud.tencent.com/developer/article/1811437
	encoder := getEncoder()
	writeSyncer := getWriteSyncer(logPath)

	core := zapcore.NewCore(encoder, writeSyncer, zapcore.Level(logLevel))

	DefaultLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	DefaultSugarLogger = DefaultLogger.Sugar()
}

func getEncoder() zapcore.Encoder {
	var config zapcore.EncoderConfig
	if isDebugMode {
		config = zap.NewDevelopmentEncoderConfig()
	} else {
		config = zap.NewProductionEncoderConfig()
	}
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	if logJson {
		return zapcore.NewJSONEncoder(config)
	} else {
		return zapcore.NewConsoleEncoder(config)
	}
}

func getWriteSyncer(path string) (syncer zapcore.WriteSyncer) {
	writer := getWriter(path)
	if isDebugMode {
		mw := io.MultiWriter(os.Stdout, writer)
		syncer = zapcore.AddSync(mw)
	} else {
		syncer = zapcore.AddSync(writer)
	}
	return
}

func getWriter(path string) io.Writer {
	writer, err := rotatelogs.New(
		path+".%Y-%m-%d",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(-1),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
		rotatelogs.WithRotationCount(3))
	if err != nil {
		panic("create rotatelogs failed, error=" + err.Error())
	}
	return writer
}

func Debugf(fmt string, args ...interface{}) {
	DefaultSugarLogger.Debugf(fmt, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	DefaultSugarLogger.Debugw(msg, keysAndValues...)
}

func Infof(fmt string, args ...interface{}) {
	DefaultSugarLogger.Infof(fmt, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	DefaultSugarLogger.Infow(msg, keysAndValues...)
}

func Warnf(fmt string, args ...interface{}) {
	DefaultSugarLogger.Warnf(fmt, args...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	DefaultSugarLogger.Warnw(msg, keysAndValues...)
}

func Errorf(fmt string, args ...interface{}) {
	DefaultSugarLogger.Errorf(fmt, args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	DefaultSugarLogger.Errorw(msg, keysAndValues...)
}

// GetTraceID 从 context 中获取 trace_id，如果不存在返回空字符串
func GetTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if traceID, ok := ctx.Value(traceIDKey).(string); ok {
		return traceID
	}
	return ""
}

// GetSpanID 从 context 中获取 span_id，如果不存在返回空字符串
func GetSpanID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if spanID, ok := ctx.Value(spanIDKey).(string); ok {
		return spanID
	}
	return ""
}

// WithTraceID 将 trace_id 添加到 context 中
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

// WithSpanID 将 span_id 添加到 context 中
func WithSpanID(ctx context.Context, spanID string) context.Context {
	return context.WithValue(ctx, spanIDKey, spanID)
}

// WithTrace 同时设置 trace_id 和 span_id 到 context 中
func WithTrace(ctx context.Context, traceID, spanID string) context.Context {
	ctx = context.WithValue(ctx, traceIDKey, traceID)
	ctx = context.WithValue(ctx, spanIDKey, spanID)
	return ctx
}

// DebugfCtx 使用 context 记录 Debug 级别的格式化日志，自动添加 trace_id 和 span_id
func DebugfCtx(ctx context.Context, format string, args ...interface{}) {
	traceID := GetTraceID(ctx)
	spanID := GetSpanID(ctx)
	DefaultSugarLogger.Debugw(format, "trace_id", traceID, "span_id", spanID, "args", args)
}

// DebugwCtx 使用 context 记录 Debug 级别的结构化日志，自动添加 trace_id 和 span_id
func DebugwCtx(ctx context.Context, msg string, keysAndValues ...interface{}) {
	traceID := GetTraceID(ctx)
	spanID := GetSpanID(ctx)
	kvs := append([]interface{}{"trace_id", traceID, "span_id", spanID}, keysAndValues...)
	DefaultSugarLogger.Debugw(msg, kvs...)
}

// InfofCtx 使用 context 记录 Info 级别的格式化日志，自动添加 trace_id 和 span_id
func InfofCtx(ctx context.Context, format string, args ...interface{}) {
	traceID := GetTraceID(ctx)
	spanID := GetSpanID(ctx)
	DefaultSugarLogger.Infow(format, "trace_id", traceID, "span_id", spanID, "args", args)
}

// InfowCtx 使用 context 记录 Info 级别的结构化日志，自动添加 trace_id 和 span_id
func InfowCtx(ctx context.Context, msg string, keysAndValues ...interface{}) {
	traceID := GetTraceID(ctx)
	spanID := GetSpanID(ctx)
	kvs := append([]interface{}{"trace_id", traceID, "span_id", spanID}, keysAndValues...)
	DefaultSugarLogger.Infow(msg, kvs...)
}

// WarnfCtx 使用 context 记录 Warn 级别的格式化日志，自动添加 trace_id 和 span_id
func WarnfCtx(ctx context.Context, format string, args ...interface{}) {
	traceID := GetTraceID(ctx)
	spanID := GetSpanID(ctx)
	DefaultSugarLogger.Warnw(format, "trace_id", traceID, "span_id", spanID, "args", args)
}

// WarnwCtx 使用 context 记录 Warn 级别的结构化日志，自动添加 trace_id 和 span_id
func WarnwCtx(ctx context.Context, msg string, keysAndValues ...interface{}) {
	traceID := GetTraceID(ctx)
	spanID := GetSpanID(ctx)
	kvs := append([]interface{}{"trace_id", traceID, "span_id", spanID}, keysAndValues...)
	DefaultSugarLogger.Warnw(msg, kvs...)
}

// ErrorfCtx 使用 context 记录 Error 级别的格式化日志，自动添加 trace_id 和 span_id
func ErrorfCtx(ctx context.Context, format string, args ...interface{}) {
	traceID := GetTraceID(ctx)
	spanID := GetSpanID(ctx)
	DefaultSugarLogger.Errorw(format, "trace_id", traceID, "span_id", spanID, "args", args)
}

// ErrorwCtx 使用 context 记录 Error 级别的结构化日志，自动添加 trace_id 和 span_id
func ErrorwCtx(ctx context.Context, msg string, keysAndValues ...interface{}) {
	traceID := GetTraceID(ctx)
	spanID := GetSpanID(ctx)
	kvs := append([]interface{}{"trace_id", traceID, "span_id", spanID}, keysAndValues...)
	DefaultSugarLogger.Errorw(msg, kvs...)
}
