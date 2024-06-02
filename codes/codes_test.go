package codes

import (
	"testing"

	"github.com/reyhanmichiels/go-pkg/language"
	"github.com/stretchr/testify/assert"
)

func Test_codes_Compile(t *testing.T) {
	type args struct {
		code Code
		lang string
	}

	mockResult := DisplayMessage{
		StatusCode: SuccessDefault.StatusCode,
		Title:      SuccessDefault.TitleEN,
		Body:       SuccessDefault.BodyEN,
	}

	tests := []struct {
		name string
		args
		want DisplayMessage
	}{
		{
			name: "code exist",
			args: args{
				code: CodeSuccess,
				lang: language.English,
			},
			want: mockResult,
		},
		{
			name: "code doesn't exist",
			args: args{
				code: CodeBadRequest,
				lang: language.English,
			},
			want: mockResult,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Compile(tt.args.code, tt.args.lang)
			assert.Equal(t, tt.want, got)
		})
	}
}
