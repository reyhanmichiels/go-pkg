package log

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/reyhanmichiels/go-pkg/appcontext"
)

var (
	once sync.Once
	now  = time.Now
)

const (
	levelDebug = "debug"
	levelInfo  = "info"
	levelWarn  = "warn"
	levelError = "error"
)

type Interface interface {
	Info(ctx context.Context, obj any)
	Debug(ctx context.Context, obj any)
	Warn(ctx context.Context, obj any)
	Error(ctx context.Context, obj any)
}

type Config struct {
	Level string
}

// TODO: implement fatal logging
type logger struct {
	log *slog.Logger
}

func Init(cfg Config) Interface {
	var slogLogger *slog.Logger

	once.Do(func() {
		level, err := parsingLogLevel(cfg.Level)
		if err != nil {
			log.Fatal(err)
		}

		slogLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level,
			AddSource: true,
		}))
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

// TODO: improve error caller
func (l *logger) getCaller(obj any) string {
	switch tr := obj.(type) {
	case error:
		return tr.Error()
	case string:
		return tr
	default:
		return fmt.Sprintf("%#v", tr)
	}

}
