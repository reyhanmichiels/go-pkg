package null

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFloat64(t *testing.T) {
	type args struct {
		f     float64
		valid bool
	}
	tests := []struct {
		name string
		args args
		want Float64
	}{
		{
			name: "set invalid",
			args: args{
				f:     12.34,
				valid: false,
			},
			want: Float64{Valid: false, Float64: 12.34},
		},
		{
			name: "set valid",
			args: args{
				f:     12.34,
				valid: true,
			},
			want: Float64{Valid: true, Float64: 12.34},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFloat64(tt.args.f, tt.args.valid)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFloat64From(t *testing.T) {
	type args struct {
		f float64
	}
	tests := []struct {
		name string
		args args
		want Float64
	}{
		{
			name: "zero",
			args: args{f: 0},
			want: Float64{Valid: true, Float64: 0},
		},
		{
			name: "ok",
			args: args{f: 12.34},
			want: Float64{Valid: true, Float64: 12.34},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Float64From(tt.args.f)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFloat64_Scan(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    Float64
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
			want:    Float64{Valid: true, Float64: 0},
		},
		{
			name:    "scan float",
			args:    args{value: 12.34},
			wantErr: false,
			want:    Float64{Valid: true, Float64: 12.34},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nf := Float64{}
			if err := nf.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Float64.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, nf)
		})
	}
}

func TestFloat64_Value(t *testing.T) {
	tests := []struct {
		name    string
		fields  Float64
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
			fields:  Float64{Valid: true, Float64: 0},
			wantErr: false,
			want:    float64(0),
		},
		{
			name:    "value of float",
			fields:  Float64{Valid: true, Float64: 12.34},
			wantErr: false,
			want:    float64(12.34),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Float64.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFloat64_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		fields  Float64
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
			fields:  Float64{Valid: true, Float64: 0},
			wantErr: false,
			want:    []byte("0"),
		},
		{
			name:    "marshal float",
			fields:  Float64{Valid: true, Float64: 12.34},
			wantErr: false,
			want:    []byte("12.34"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Float64.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFloat64_UnmarshalJSON(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    Float64
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
			name: "unmarshal string of float",
			args: args{
				b: []byte("\"12.34\""),
			},
			wantErr: true,
		},
		{
			name: "unmarshal float",
			args: args{
				b: []byte("12.34"),
			},
			wantErr: false,
			want:    Float64{Valid: true, Float64: 12.34},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nf := Float64{}
			if err := nf.UnmarshalJSON(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("Float64.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, nf)
		})
	}
}

func TestFloat64_IsNullOrZero(t *testing.T) {
	tests := []struct {
		name   string
		fields Float64
		want   bool
	}{
		{
			name: "null or invalid",
			want: true,
		},
		{
			name:   "valid zero",
			fields: Float64{Valid: true, Float64: 0},
			want:   true,
		},
		{
			name:   "valid value",
			fields: Float64{Valid: true, Float64: 69.69},
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

func TestFloat64_Equal(t *testing.T) {
	type args struct {
		other Float64
	}
	tests := []struct {
		name   string
		fields Float64
		args   args
		want   bool
	}{
		{
			name:   "both invalid",
			fields: Float64{Valid: false},
			args:   args{Float64{Valid: false}},
			want:   true,
		},
		{
			name:   "both valid and both same value",
			fields: Float64{Valid: true, Float64: 69.69},
			args:   args{Float64{Valid: true, Float64: 69.69}},
			want:   true,
		},
		{
			name:   "valid and invalid but same value",
			fields: Float64{Valid: true, Float64: 69.69},
			args:   args{Float64{Valid: false, Float64: 69.69}},
			want:   false,
		},
		{
			name:   "both valid but different value",
			fields: Float64{Valid: true, Float64: 69.69},
			args:   args{Float64{Valid: true, Float64: 96.96}},
			want:   false,
		},
		{
			name:   "both valid and both zero",
			fields: Float64{Valid: true, Float64: 0},
			args:   args{Float64{Valid: true, Float64: 0}},
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

func TestFloat64_Is(t *testing.T) {
	type args struct {
		other float64
	}
	tests := []struct {
		name   string
		fields Float64
		args   args
		want   bool
	}{
		{
			name:   "origin invalid and not equal",
			fields: Float64{Valid: false, Float64: 69.69},
			args:   args{other: 46.46},
			want:   false,
		},
		{
			name:   "origin invalid and equal",
			fields: Float64{Valid: false, Float64: 69.69},
			args:   args{other: 69.69},
			want:   false,
		},
		{
			name:   "origin valid and not equal",
			fields: Float64{Valid: true, Float64: 69.69},
			args:   args{other: 46.46},
			want:   false,
		},
		{
			name:   "origin valid and equal",
			fields: Float64{Valid: true, Float64: 69.69},
			args:   args{other: 69.69},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.Is(tt.args.other); got != tt.want {
				t.Errorf("Float64.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}
