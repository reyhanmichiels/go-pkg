package errors

import (
	"testing"

	"github.com/reyhanmichiels/go-pkg/codes"
)

func Test_stacktrace_Error(t *testing.T) {
	type fields struct {
		message  string
		cause    error
		code     codes.Code
		file     string
		function string
		line     int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "ok",
			fields: fields{message: "failed to start"},
			want:   "failed to start",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &stacktrace{
				message:  tt.fields.message,
				cause:    tt.fields.cause,
				code:     tt.fields.code,
				file:     tt.fields.file,
				function: tt.fields.function,
				line:     tt.fields.line,
			}
			if got := st.Error(); got != tt.want {
				t.Errorf("stacktrace.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stacktrace_ExitCode(t *testing.T) {
	type fields struct {
		message  string
		cause    error
		code     codes.Code
		file     string
		function string
		line     int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name:   "no code",
			fields: fields{code: codes.NoCode},
			want:   1,
		},
		{
			name:   "bad request",
			fields: fields{code: codes.CodeBadRequest},
			want:   int(codes.CodeBadRequest),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &stacktrace{
				message:  tt.fields.message,
				cause:    tt.fields.cause,
				code:     tt.fields.code,
				file:     tt.fields.file,
				function: tt.fields.function,
				line:     tt.fields.line,
			}
			if got := st.ExitCode(); got != tt.want {
				t.Errorf("stacktrace.ExitCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
