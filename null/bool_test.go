package null

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBool(t *testing.T) {
	type args struct {
		b     bool
		valid bool
	}
	tests := []struct {
		name string
		args args
		want Bool
	}{
		{
			name: "set invalid",
			args: args{
				b:     true,
				valid: false,
			},
			want: Bool{Valid: false, Bool: true},
		},
		{
			name: "set valid",
			args: args{
				b:     true,
				valid: true,
			},
			want: Bool{Valid: true, Bool: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewBool(tt.args.b, tt.args.valid)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBoolFrom(t *testing.T) {
	type args struct {
		b bool
	}
	tests := []struct {
		name string
		args args
		want Bool
	}{
		{
			name: "false",
			args: args{b: false},
			want: Bool{Valid: true, Bool: false},
		},
		{
			name: "true",
			args: args{b: true},
			want: Bool{Valid: true, Bool: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BoolFrom(tt.args.b)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBool_Scan(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    Bool
	}{
		{
			name:    "scan nil",
			args:    args{value: nil},
			wantErr: false,
		},
		{
			name:    "scan string",
			args:    args{value: "abc"},
			wantErr: true,
		},
		{
			name:    "scan 1",
			args:    args{value: 1},
			wantErr: false,
			want:    Bool{Valid: true, Bool: true},
		},
		{
			name:    "scan 0",
			args:    args{value: 0},
			wantErr: false,
			want:    Bool{Valid: true, Bool: false},
		},
		{
			name:    "scan string of true",
			args:    args{value: "true"},
			wantErr: false,
			want:    Bool{Valid: true, Bool: true},
		},
		{
			name:    "scan string of false",
			args:    args{value: "false"},
			wantErr: false,
			want:    Bool{Valid: true, Bool: false},
		},
		{
			name:    "scan true",
			args:    args{value: true},
			wantErr: false,
			want:    Bool{Valid: true, Bool: true},
		},
		{
			name:    "scan false",
			args:    args{value: false},
			wantErr: false,
			want:    Bool{Valid: true, Bool: false},
		},
		{
			name:    "scan int",
			args:    args{value: 123},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nb := Bool{}
			if err := nb.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Bool.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, nb)
		})
	}
}

func TestBool_Value(t *testing.T) {
	tests := []struct {
		name    string
		fields  Bool
		want    driver.Value
		wantErr bool
	}{
		{
			name:    "value of nil",
			wantErr: false,
			want:    nil,
		},
		{
			name:    "value of false",
			fields:  Bool{Valid: true, Bool: false},
			wantErr: false,
			want:    false,
		},
		{
			name:    "value of true",
			fields:  Bool{Valid: true, Bool: true},
			wantErr: false,
			want:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Bool.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBool_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		fields  Bool
		want    []byte
		wantErr bool
	}{
		{
			name:    "marshal invalid",
			wantErr: false,
			want:    []byte("null"),
		},
		{
			name:    "marshal false",
			fields:  Bool{Valid: true, Bool: false},
			wantErr: false,
			want:    []byte("false"),
		},
		{
			name:    "marshal true",
			fields:  Bool{Valid: true, Bool: true},
			wantErr: false,
			want:    []byte("true"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Bool.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBool_UnmarshalJSON(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    Bool
	}{
		{
			name: "unmarshal nil",
			args: args{
				b: []byte("null"),
			},
			wantErr: false,
		},
		{
			name: "unmarshal string",
			args: args{
				b: []byte("\"abcd\""),
			},
			wantErr: true,
		},
		{
			name: "unmarshal string of number",
			args: args{
				b: []byte("\"1234\""),
			},
			wantErr: true,
		},
		{
			name: "unmarshal number",
			args: args{
				b: []byte("1234"),
			},
			wantErr: true,
		},
		{
			name: "unmarshal 0",
			args: args{
				b: []byte("0"),
			},
			wantErr: true,
		},
		{
			name: "unmarshal 1",
			args: args{
				b: []byte("1"),
			},
			wantErr: true,
		},
		{
			name: "unmarshal string of true",
			args: args{
				b: []byte("\"true\""),
			},
			wantErr: true,
		},
		{
			name: "unmarshal true",
			args: args{
				b: []byte("true"),
			},
			wantErr: false,
			want:    Bool{Valid: true, Bool: true},
		},
		{
			name: "unmarshal false",
			args: args{
				b: []byte("false"),
			},
			wantErr: false,
			want:    Bool{Valid: true, Bool: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nb := Bool{}
			if err := nb.UnmarshalJSON(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("Bool.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, nb)
		})
	}
}

func TestBool_IsNullOrZero(t *testing.T) {
	tests := []struct {
		name   string
		fields Bool
		want   bool
	}{
		{
			name: "null or invalid",
			want: true,
		},
		{
			name:   "valid false",
			fields: Bool{Valid: true, Bool: false},
			want:   true,
		},
		{
			name:   "valid true",
			fields: Bool{Valid: true, Bool: true},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fields.IsNullOrZero()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBool_Equal(t *testing.T) {
	type args struct {
		other Bool
	}
	tests := []struct {
		name   string
		fields Bool
		args   args
		want   bool
	}{
		{
			name:   "both invalid",
			fields: Bool{Valid: false},
			args:   args{Bool{Valid: false}},
			want:   true,
		},
		{
			name:   "both valid and both same value",
			fields: Bool{Valid: true, Bool: true},
			args:   args{Bool{Valid: true, Bool: true}},
			want:   true,
		},
		{
			name:   "valid and invalid but same value",
			fields: Bool{Valid: true, Bool: true},
			args:   args{Bool{Valid: false, Bool: true}},
			want:   false,
		},
		{
			name:   "both valid but different value",
			fields: Bool{Valid: true, Bool: true},
			args:   args{Bool{Valid: true, Bool: false}},
			want:   false,
		},
		{
			name:   "both valid and both zero",
			fields: Bool{Valid: true, Bool: false},
			args:   args{Bool{Valid: true, Bool: false}},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.Equal(tt.args.other); got != tt.want {
				t.Errorf("Bool.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBool_Is(t *testing.T) {
	type args struct {
		other bool
	}
	tests := []struct {
		name   string
		fields Bool
		args   args
		want   bool
	}{
		{
			name:   "origin invalid and not equal",
			fields: Bool{Valid: false, Bool: false},
			args:   args{other: true},
			want:   false,
		},
		{
			name:   "origin invalid and equal",
			fields: Bool{Valid: false, Bool: true},
			args:   args{other: true},
			want:   false,
		},
		{
			name:   "origin valid and not equal",
			fields: Bool{Valid: true, Bool: true},
			args:   args{other: false},
			want:   false,
		},
		{
			name:   "origin valid and equal",
			fields: Bool{Valid: true, Bool: true},
			args:   args{other: true},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.Is(tt.args.other); got != tt.want {
				t.Errorf("Bool.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}
