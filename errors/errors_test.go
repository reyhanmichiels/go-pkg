package errors

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/reyhanmichiels/go-pkg/codes"
	"github.com/reyhanmichiels/go-pkg/language"
	"github.com/stretchr/testify/assert"
)

func Test_App_Error(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "OK",
			want: "invalid format",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &App{
				sys: NewWithCode(codes.CodeBadRequest, "invalid format"),
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("App.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_errors_Compile(t *testing.T) {
	type args struct {
		err  error
		lang string
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 App
	}{
		{
			name: "code exist",
			args: args{err: NewWithCode(codes.CodeBadRequest, "bad request"), lang: language.English},
			want: codes.ErrMsgBadRequest.StatusCode,
			want1: App{
				Code:  codes.CodeBadRequest,
				Title: codes.ErrMsgBadRequest.TitleEN,
				Body:  codes.ErrMsgBadRequest.BodyEN,
				sys:   NewWithCode(codes.CodeBadRequest, "bad request"),
			},
		},
		{
			name: "code not exist",
			args: args{err: NewWithCode(codes.NoCode, "bad request"), lang: language.English},
			want: http.StatusInternalServerError,
			want1: App{
				Code:  codes.NoCode,
				Title: "Service Error Not Defined",
				Body:  "Unknown error. Please contact admin",
				sys:   NewWithCode(codes.CodeBadRequest, "bad request"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Compile(tt.args.err, tt.args.lang)
			if got != tt.want {
				t.Errorf("Compile() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1.Code, tt.want1.Code) {
				t.Errorf("Compile() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_error_GetCaller(t *testing.T) {
	pwd, _ := os.Getwd()
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   int
		want2   string
		wantErr bool
	}{
		{
			name:    "ok",
			args:    args{err: NewWithCode(codes.CodeBadRequest, "bad request")},
			want:    pwd + "/errors_test.go",
			want1:   99,
			want2:   "bad request",
			wantErr: false,
		},
		{
			name:    "not ok",
			args:    args{err: fmt.Errorf("")},
			want:    "",
			want1:   0,
			want2:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := GetCaller(tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCaller() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCaller() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetCaller() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("GetCaller() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_error_GetCallerString(t *testing.T) {
	pwd, _ := os.Getwd()
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "ok",
			args:    args{err: NewWithCode(codes.CodeBadRequest, "bad request")},
			want:    fmt.Sprintf("%s:%#v --- %s", pwd+"/errors_test.go", 147, "bad request"),
			wantErr: false,
		},
		{
			name:    "not ok",
			args:    args{err: fmt.Errorf("")},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCallerString(tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCaller() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
