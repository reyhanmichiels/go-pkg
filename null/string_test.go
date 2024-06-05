package null

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewString(t *testing.T) {
	type args struct {
		s     string
		valid bool
	}
	tests := []struct {
		name string
		args args
		want String
	}{
		{
			name: "set invalid",
			args: args{
				s:     "ABCD",
				valid: false,
			},
			want: String{Valid: false, String: "ABCD"},
		},
		{
			name: "set valid",
			args: args{
				s:     "ABCD",
				valid: true,
			},
			want: String{Valid: true, String: "ABCD"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewString(tt.args.s, tt.args.valid)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStringFrom(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want String
	}{
		{
			name: "empty",
			args: args{s: ""},
			want: String{Valid: true, String: ""},
		},
		{
			name: "ok",
			args: args{s: "ABCD"},
			want: String{Valid: true, String: "ABCD"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringFrom(tt.args.s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestString_Scan(t *testing.T) {
	getStringPtr := func(s string) *string { return &s }
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    String
	}{
		{
			name:    "scan nil",
			args:    args{value: nil},
			wantErr: false,
		},
		{
			name:    "scan string",
			args:    args{value: "abc"},
			wantErr: false,
			want:    String{Valid: true, String: "abc"},
		},
		{
			name:    "scan empty",
			args:    args{value: ""},
			wantErr: false,
			want:    String{Valid: true, String: ""},
		},
		{
			name:    "scan int",
			args:    args{value: 123},
			wantErr: false,
			want:    String{Valid: true, String: "123"},
		},
		{
			name:    "scan boolean",
			args:    args{value: true},
			wantErr: false,
			want:    String{Valid: true, String: "true"},
		},
		{
			name:    "scan pointer of string",
			args:    args{value: getStringPtr("abc")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ns := String{}
			if err := ns.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("String.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, ns)
		})
	}
}

func TestString_Value(t *testing.T) {
	tests := []struct {
		name    string
		fields  String
		want    driver.Value
		wantErr bool
	}{
		{
			name:    "value of nil",
			wantErr: false,
			want:    nil,
		},
		{
			name:    "value of empty",
			fields:  String{Valid: true, String: ""},
			wantErr: false,
			want:    "",
		},
		{
			name:    "value of string",
			fields:  String{Valid: true, String: "abc"},
			wantErr: false,
			want:    "abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("String.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestString_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		fields  String
		want    []byte
		wantErr bool
	}{
		{
			name:    "marshal invalid",
			wantErr: false,
			want:    []byte("null"),
		},
		{
			name:    "marshal empty",
			fields:  String{Valid: true, String: ""},
			wantErr: false,
			want:    []byte("\"\""),
		},
		{
			name:    "marshal string",
			fields:  String{Valid: true, String: "abc"},
			wantErr: false,
			want:    []byte("\"abc\""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("String.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestString_UnmarshalJSON(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    String
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
			wantErr: false,
			want:    String{Valid: true, String: "abcd"},
		},
		{
			name: "unmarshal string of number",
			args: args{
				b: []byte("\"1234\""),
			},
			wantErr: false,
			want:    String{Valid: true, String: "1234"},
		},
		{
			name: "unmarshal number",
			args: args{
				b: []byte("1234"),
			},
			wantErr: true,
		},
		{
			name: "unmarshal string with whitespaces",
			args: args{
				b: []byte("\"          abc 123              \""),
			},
			wantErr: false,
			want:    String{Valid: true, String: "abc 123"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ns := String{}
			if err := ns.UnmarshalJSON(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("String.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, ns)
		})
	}
}

func TestString_IsNullOrZero(t *testing.T) {
	tests := []struct {
		name   string
		fields String
		want   bool
	}{
		{
			name: "null or invalid",
			want: true,
		},
		{
			name:   "valid empty",
			fields: String{Valid: true, String: ""},
			want:   true,
		},
		{
			name:   "valid value",
			fields: String{Valid: true, String: "abc"},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.IsNullOrZero(); got != tt.want {
				t.Errorf("String.IsNullOrZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_Equal(t *testing.T) {
	type args struct {
		other String
	}
	tests := []struct {
		name   string
		fields String
		args   args
		want   bool
	}{
		{
			name:   "both invalid",
			fields: String{Valid: false},
			args:   args{String{Valid: false}},
			want:   true,
		},
		{
			name:   "both valid and both same value",
			fields: String{Valid: true, String: "abc"},
			args:   args{String{Valid: true, String: "abc"}},
			want:   true,
		},
		{
			name:   "valid and invalid but same value",
			fields: String{Valid: true, String: "abc"},
			args:   args{String{Valid: false, String: "abc"}},
			want:   false,
		},
		{
			name:   "both valid but different value",
			fields: String{Valid: true, String: "abc"},
			args:   args{String{Valid: true, String: "bca"}},
			want:   false,
		},
		{
			name:   "both valid and both zero",
			fields: String{Valid: true, String: ""},
			args:   args{String{Valid: true, String: ""}},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.Equal(tt.args.other); got != tt.want {
				t.Errorf("String.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_Is(t *testing.T) {
	type args struct {
		other string
	}
	tests := []struct {
		name   string
		fields String
		args   args
		want   bool
	}{
		{
			name:   "origin invalid and not equal",
			fields: String{Valid: false, String: "kucing"},
			args:   args{other: "cat"},
			want:   false,
		},
		{
			name:   "origin invalid and equal",
			fields: String{Valid: false, String: "kucing"},
			args:   args{other: "kucing"},
			want:   false,
		},
		{
			name:   "origin valid and not equal",
			fields: String{Valid: true, String: "kucing"},
			args:   args{other: "cat"},
			want:   false,
		},
		{
			name:   "origin valid and equal",
			fields: String{Valid: true, String: "kucing"},
			args:   args{other: "kucing"},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.Is(tt.args.other); got != tt.want {
				t.Errorf("String.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}
