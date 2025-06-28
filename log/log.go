package log

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/reyhanmichiels/go-pkg/v2/appcontext"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	once             sync.Once
	now              = time.Now
	CustomLevelNames = map[slog.Leveler]string{
		LevelPanic: "PANIC",
		LevelFatal: "FATAL",
	}
)

const (
	levelDebug = "debug"
	levelInfo  = "info"
	levelWarn  = "warn"
	levelError = "error"
	levelFatal = "fatal"
	levelPanic = "panic"

	// customize slog level
	LevelFatal = slog.Level(10)
	LevelPanic = slog.Level(12)

	OutputFile = "file"
)

type Interface interface {
	Info(ctx context.Context, obj any)
	Debug(ctx context.Context, obj any)
	Warn(ctx context.Context, obj any)
	Error(ctx context.Context, obj any)
	Fatal(ctx context.Context, obj any)
	Panic(obj any)
}

type Config struct {
	Level  string
	Output string

	LumberjackConfig LumbejackConfig
}

type LumbejackConfig struct {
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

type logger struct {
	log *slog.Logger
}

func DefaultLogger() Interface {
	return &logger{
		log: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:       slog.LevelDebug,
			AddSource:   true,
			ReplaceAttr: getCustomLevelName,
		})),
	}
}

func Init(cfg Config) Interface {
	var slogLogger *slog.Logger

	once.Do(func() {
		level, err := parsingLogLevel(cfg.Level)
		if err != nil {
			log.Panic(err)
		}

		switch cfg.Output {
		case OutputFile:
			if cfg.LumberjackConfig.Filename == "" {
				log.Panic("filename cannot be empty")
			}

			logFile := lumberjack.Logger{
				Filename:   cfg.LumberjackConfig.Filename,
				MaxSize:    cfg.LumberjackConfig.MaxSize,
				MaxBackups: cfg.LumberjackConfig.MaxBackups,
				MaxAge:     cfg.LumberjackConfig.MaxAge,
				Compress:   cfg.LumberjackConfig.Compress,
			}

			slogLogger = slog.New(slog.NewJSONHandler(&logFile, &slog.HandlerOptions{
				Level:       level,
				AddSource:   true,
				ReplaceAttr: getCustomLevelName,
			}))
		default:
			slogLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level:       level,
				AddSource:   true,
				ReplaceAttr: getCustomLevelName,
			}))
		}
	})

	return &logger{
		log: slogLogger,
	}
}

func (l *logger) Info(ctx context.Context, obj any) {
	l.log.LogAttrs(
		ctx,
		slog.LevelInfo,
		l.getCaller(obj),
		l.getFieldsFromContext(ctx)...,
	)
}

func (l *logger) Debug(ctx context.Context, obj any) {
	l.log.LogAttrs(
		ctx,
		slog.LevelDebug,
		l.getCaller(obj),
		l.getFieldsFromContext(ctx)...,
	)
}

func (l *logger) Warn(ctx context.Context, obj any) {
	l.log.LogAttrs(
		ctx,
		slog.LevelWarn,
		l.getCaller(obj),
		l.getFieldsFromContext(ctx)...,
	)
}

func (l *logger) Error(ctx context.Context, obj any) {
	l.log.LogAttrs(
		ctx,
		slog.LevelError,
		l.getCaller(obj),
		l.getFieldsFromContext(ctx)...,
	)
}

func (l *logger) Fatal(ctx context.Context, obj any) {
	l.log.LogAttrs(
		ctx,
		LevelFatal,
		l.getCaller(obj),
		l.getFieldsFromContext(ctx)...,
	)

	os.Exit(1)
}

func (l *logger) Panic(obj any) {
	l.log.LogAttrs(
		context.Background(),
		LevelPanic,
		l.getCaller(obj),
		l.getPanicStacktrace(),
	)
}

func parsingLogLevel(text string) (slog.Level, error) {
	level := strings.ToLower(text)

	switch level {
	case levelInfo:
		return slog.LevelInfo, nil
	case levelError:
		return slog.LevelError, nil
	case levelDebug:
		return slog.LevelDebug, nil
	case levelWarn:
		return slog.LevelWarn, nil
	case levelPanic:
		return LevelPanic, nil
	case levelFatal:
		return LevelFatal, nil
	}

	return slog.Level(-1), fmt.Errorf("invalid log level %s", text)
}

func (l *logger) getFieldsFromContext(ctx context.Context) []slog.Attr {
	reqStart := appcontext.GetRequestStartTime(ctx)
	appRespCode := appcontext.GetAppResponseCodeInt(ctx)
	appErrMsg := appcontext.GetAppErrorMessage(ctx)
	timeElapsed := "0ms"
	if !time.Time.IsZero(reqStart) {
		timeElapsed = fmt.Sprintf("%dms", int64(now().Sub(reqStart)/time.Millisecond))
	}

	fields := []slog.Attr{
		slog.String("request_id", appcontext.GetRequestId(ctx)),
		slog.String("user_agent", appcontext.GetUserAgent(ctx)),
		slog.Int("user_id", appcontext.GetUserId(ctx)),
		slog.String("service_version", appcontext.GetServiceVersion(ctx)),
		slog.String("time_elapsed", timeElapsed),
	}

	if appRespCode > 0 {
		fields = append(fields, slog.Int("app_resp_code", appRespCode))
	}

	if appErrMsg != "" {
		fields = append(fields, slog.String("app_err_msg", appErrMsg))
	}

	return fields
}

func (l *logger) getCaller(obj any) string {
	switch tr := obj.(type) {
	case error:
		caller, err := errors.GetCallerString(tr)
		if err != nil {
			return tr.Error()
		}

		return caller
	case string:
		return tr
	default:
		return fmt.Sprintf("%#v", tr)
	}
}

func getCustomLevelName(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey {
		level := a.Value.Any().(slog.Level)

		levelLabel, exists := CustomLevelNames[level]
		if !exists {
			levelLabel = level.String()
		}

		a.Value = slog.StringValue(levelLabel)
	}

	return a
}

func (l *logger) getPanicStacktrace() slog.Attr {
	errStackAttr := []any{}
	errStack := strings.Split(strings.ReplaceAll(string(debug.Stack()), "\t", ""), "\n")

	for i, v := range errStack {
		errStackAttr = append(errStackAttr, slog.String(fmt.Sprintf("stack - %v", i), v))
	}

	return slog.Group("stack_trace", errStackAttr...)
}
