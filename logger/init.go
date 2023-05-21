package logger

import (
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
