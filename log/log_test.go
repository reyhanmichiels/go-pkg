package log

import (
	"context"
	"errors"
	"testing"

	"github.com/reyhanmichiels/go-pkg/v2/appcontext"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	myerr "github.com/reyhanmichiels/go-pkg/v2/errors"
)

func Test_log_Info(t *testing.T) {
	type args struct {
		ctx context.Context
		obj any
	}

	mockLog := Init(Config{Level: levelInfo})

	mockCtx := context.Background()
	mockCtx = appcontext.SetRequestId(mockCtx, "the request id")
	mockCtx = appcontext.SetUserAgent(mockCtx, "the user agent")
	mockCtx = appcontext.SetUserId(mockCtx, 1)
	mockCtx = appcontext.SetServiceVersion(mockCtx, "the service version")
	mockCtx = appcontext.SetAppResponseCode(mockCtx, codes.Code(100))
	mockCtx = appcontext.SetAppErrorMessage(mockCtx, "the error message")
	mockCtx = appcontext.SetRequestStartTime(mockCtx, now())

	tests := []struct {
		name string
		args
		mockFunc func(mockLogger Interface, arg args)
	}{
		{
			name: "info",
			args: args{
				ctx: mockCtx,
				obj: "test log info",
			},
			mockFunc: func(mockLogger Interface, arg args) {
				mockLogger.Info(mockCtx, arg.obj)
			},
		},
		{
			name: "debug",
			args: args{
				ctx: mockCtx,
				obj: struct {
					Test string
				}{
					"test log debug",
				},
			},
			mockFunc: func(mockLogger Interface, arg args) {
				mockLogger.Debug(mockCtx, arg.obj)
			},
		},
		{
			name: "warn",
			args: args{
				ctx: mockCtx,
				obj: 1,
			},
			mockFunc: func(mockLogger Interface, arg args) {
				mockLogger.Warn(mockCtx, arg.obj)
			},
		},
		{
			name: "error",
			args: args{
				ctx: mockCtx,
				obj: errors.New("test error"),
			},
			mockFunc: func(mockLogger Interface, arg args) {
				mockLogger.Error(mockCtx, arg.obj)
			},
		}, {
			name: "error with codes",
			args: args{
				ctx: mockCtx,
				obj: myerr.NewWithCode(codes.CodeBadRequest, "test error with code"),
			},
			mockFunc: func(mockLogger Interface, arg args) {
				mockLogger.Error(mockCtx, arg.obj)
			},
		},
		{
			name: "fatal",
			args: args{
				ctx: mockCtx,
				obj: "test log fatal",
			},
			mockFunc: func(mockLogger Interface, arg args) {
				defer func() {
					if err := recover(); err != nil {
						mockLogger.Panic(err)
					}
				}()
				panic(myerr.NewWithCode(codes.CodeBadRequest, "test panic"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mockLog, tt.args)
		})
	}
}
