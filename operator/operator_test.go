package operator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_operator_Ternary(t *testing.T) {
	type args struct {
		condition bool
		a         string
		b         string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "condition true",
			args: args{
				condition: true,
				a:         "string A",
				b:         "string B",
			},
			want: "string A",
		},
		{
			name: "condition false",
			args: args{
				condition: false,
				a:         "string A",
				b:         "string B",
			},
			want: "string B",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Ternary(tt.args.condition, tt.args.a, tt.args.b)
			assert.Equal(t, tt.want, got)
		})
	}
}
