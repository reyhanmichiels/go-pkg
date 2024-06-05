package null

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInt64(t *testing.T) {
	type args struct {
		i     int64
		valid bool
	}
	tests := []struct {
		name string
		args args
		want Int64
	}{
		{
			name: "set invalid",
			args: args{
				i:     1234,
				valid: false,
			},
			want: Int64{Valid: false, Int64: 1234},
		},
		{
			name: "set valid",
			args: args{
				i:     1234,
				valid: true,
			},
			want: Int64{Valid: true, Int64: 1234},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewInt64(tt.args.i, tt.args.valid)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInt64From(t *testing.T) {
	type args struct {
		i int64
	}
	tests := []struct {
		name string
		args args
		want Int64
	}{
		{
			name: "zero",
			args: args{i: 0},
			want: Int64{Valid: true, Int64: 0},
		},
		{
			name: "ok",
			args: args{i: 1234},
			want: Int64{Valid: true, Int64: 1234},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Int64From(tt.args.i)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInt64_Scan(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    Int64
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
			name:    "scan zero",
			args:    args{value: 0},
			wantErr: false,
			want:    Int64{Valid: true, Int64: 0},
		},
		{
			name:    "scan int",
			args:    args{value: 123},
			wantErr: false,
			want:    Int64{Valid: true, Int64: 123},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ni := Int64{}
			if err := ni.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Int64.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, ni)
		})
	}
}

func TestInt64_Value(t *testing.T) {
	tests := []struct {
		name    string
		fields  Int64
		want    driver.Value
		wantErr bool
	}{
		{
			name:    "value of nil",
			wantErr: false,
			want:    nil,
		},
		{
			name:    "value of zero",
			fields:  Int64{Valid: true, Int64: 0},
			wantErr: false,
			want:    int64(0),
		},
		{
			name:    "value of int",
			fields:  Int64{Valid: true, Int64: 123},
			wantErr: false,
			want:    int64(123),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Int64.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInt64_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		fields  Int64
		want    []byte
		wantErr bool
	}{
		{
			name:    "marshal invalid",
			wantErr: false,
			want:    []byte("null"),
		},
		{
			name:    "marshal zero",
			fields:  Int64{Valid: true, Int64: 0},
			wantErr: false,
			want:    []byte("0"),
		},
		{
			name:    "marshal int",
			fields:  Int64{Valid: true, Int64: 123},
			wantErr: false,
			want:    []byte("123"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Int64.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInt64_UnmarshalJSON(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    Int64
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
			wantErr: false,
			want:    Int64{Valid: true, Int64: 1234},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ni := Int64{}
			if err := ni.UnmarshalJSON(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("Int64.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, ni)
		})
	}
}

func TestInt64_IsNullOrZero(t *testing.T) {
	tests := []struct {
		name   string
		fields Int64
		want   bool
	}{
		{
			name: "null or invalid",
			want: true,
		},
		{
			name:   "valid zero",
			fields: Int64{Valid: true, Int64: 0},
			want:   true,
		},
		{
			name:   "valid value",
			fields: Int64{Valid: true, Int64: 69},
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

func TestInt64_Equal(t *testing.T) {
	type args struct {
		other Int64
	}
	tests := []struct {
		name   string
		fields Int64
		args   args
		want   bool
	}{
		{
			name:   "both invalid",
			fields: Int64{Valid: false},
			args:   args{Int64{Valid: false}},
			want:   true,
		},
		{
			name:   "both valid and both same value",
			fields: Int64{Valid: true, Int64: 69},
			args:   args{Int64{Valid: true, Int64: 69}},
			want:   true,
		},
		{
			name:   "valid and invalid but same value",
			fields: Int64{Valid: true, Int64: 69},
			args:   args{Int64{Valid: false, Int64: 69}},
			want:   false,
		},
		{
			name:   "both valid but different value",
			fields: Int64{Valid: true, Int64: 69},
			args:   args{Int64{Valid: true, Int64: 96}},
			want:   false,
		},
		{
			name:   "both valid and both zero",
			fields: Int64{Valid: true, Int64: 0},
			args:   args{Int64{Valid: true, Int64: 0}},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fields.Equal(tt.args.other)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInt64_Is(t *testing.T) {
	type args struct {
		other int64
	}
	tests := []struct {
		name   string
		fields Int64
		args   args
		want   bool
	}{
		{
			name:   "origin invalid and not equal",
			fields: Int64{Valid: false, Int64: 69},
			args:   args{other: 46},
			want:   false,
		},
		{
			name:   "origin invalid and equal",
			fields: Int64{Valid: false, Int64: 69},
			args:   args{other: 69},
			want:   false,
		},
		{
			name:   "origin valid and not equal",
			fields: Int64{Valid: true, Int64: 69},
			args:   args{other: 46},
			want:   false,
		},
		{
			name:   "origin valid and equal",
			fields: Int64{Valid: true, Int64: 69},
			args:   args{other: 69},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.Is(tt.args.other); got != tt.want {
				t.Errorf("Int64.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}
