package language

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_language_HTTPStatusText(t *testing.T) {
	type args struct {
		lang string
		code int
	}

	tests := []struct {
		name string
		args
		want string
	}{
		{
			name: "language id",
			args: args{
				lang: Indonesian,
				code: http.StatusOK,
			},
			want: statusTextId[http.StatusOK],
		},
		{
			name: "language en",
			args: args{
				lang: English,
				code: http.StatusOK,
			},
			want: statusTextId[http.StatusOK],
		},
	}
	for _, tt := range tests {
		got := HTTPStatusText(tt.args.lang, tt.args.code)
		assert.Equal(t, tt.want, got)
	}
}
