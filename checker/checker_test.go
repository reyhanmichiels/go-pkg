package checker

import (
	"reflect"
	"testing"
)

func TestArrayInt64Contains(t *testing.T) {
	type args struct {
		slice []int64
		x     int64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "array contain x",
			args: args{
				slice: []int64{1, 2},
				x:     2,
			},
			want: true,
		},
		{
			name: "array does not contain x",
			args: args{
				slice: []int64{1, 2},
				x:     3,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayInt64Contains(tt.args.slice, tt.args.x); got != tt.want {
				t.Errorf("ArrayInt64Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayContains(t *testing.T) {
	type argsInt64 struct {
		slice []int64
		x     int64
	}
	type argsString struct {
		slice []string
		x     string
	}
	type argsFloat64 struct {
		slice []float64
		x     float64
	}
	type args struct {
		argsInt64
		argsString
		argsFloat64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "array of int64 contain x",
			args: args{
				argsInt64: argsInt64{
					slice: []int64{1, 2},
					x:     2,
				},
			},
			want: true,
		},
		{
			name: "array of int64 does not contain x",
			args: args{
				argsInt64: argsInt64{
					slice: []int64{1, 2},
					x:     3,
				},
			},
			want: false,
		},
		{
			name: "array of string contain x",
			args: args{
				argsString: argsString{
					slice: []string{"a", "b"},
					x:     "b",
				},
			},
			want: true,
		},
		{
			name: "array of string does not contain x",
			args: args{
				argsString: argsString{
					slice: []string{"a", "b"},
					x:     "c",
				},
			},
			want: false,
		},

		{
			name: "array of float64 contain x",
			args: args{
				argsFloat64: argsFloat64{
					slice: []float64{1.0, 2.0},
					x:     2.0,
				},
			},
			want: true,
		},
		{
			name: "array of float64 does not contain x",
			args: args{
				argsFloat64: argsFloat64{
					slice: []float64{1.0, 2.0},
					x:     3.0,
				},
			},
			want: false,
		},
	}

	for i, tt := range tests {
		if i < 2 {
			t.Run(tt.name, func(t *testing.T) {
				if got := ArrayContains(tt.args.argsInt64.slice, tt.args.argsInt64.x); got != tt.want {
					t.Errorf("ArrayContains() = %v, want %v", got, tt.want)
				}
			})
		} else if i < 4 {
			t.Run(tt.name, func(t *testing.T) {
				if got := ArrayContains(tt.args.argsString.slice, tt.args.argsString.x); got != tt.want {
					t.Errorf("ArrayContains() = %v, want %v", got, tt.want)
				}
			})
		} else {
			t.Run(tt.name, func(t *testing.T) {
				if got := ArrayContains(tt.args.argsFloat64.slice, tt.args.argsFloat64.x); got != tt.want {
					t.Errorf("ArrayContains() = %v, want %v", got, tt.want)
				}
			})
		}
	}
}

func TestArrayDeduplicate(t *testing.T) {
	type argsInt64 struct {
		slice []int64
	}
	type argsString struct {
		slice []string
	}
	type argsFloat64 struct {
		slice []float64
	}
	type args struct {
		argsInt64
		argsString
		argsFloat64
	}
	tests := []struct {
		name  string
		args  args
		want  []int64
		want2 []string
		want3 []float64
	}{
		{
			name: "array of int64 contain duplicate",
			args: args{
				argsInt64: argsInt64{
					slice: []int64{1, 2, 1},
				},
			},
			want: []int64{1, 2},
		},
		{
			name: "array of int64 does not contain duplicate",
			args: args{
				argsInt64: argsInt64{
					slice: []int64{1, 2, 3},
				},
			},
			want: []int64{1, 2, 3},
		},
		{
			name: "array of string contain duplicate",
			args: args{
				argsString: argsString{
					slice: []string{"a", "b", "a"},
				},
			},
			want2: []string{"a", "b"},
		},
		{
			name: "array of string does not contain duplicate",
			args: args{
				argsString: argsString{
					slice: []string{"a", "b"},
				},
			},
			want2: []string{"a", "b"},
		},

		{
			name: "array of float64 contain duplicate",
			args: args{
				argsFloat64: argsFloat64{
					slice: []float64{1.0, 2.0, 1.0},
				},
			},
			want3: []float64{1.0, 2.0},
		},
		{
			name: "array of float64 does not contain duplicate",
			args: args{
				argsFloat64: argsFloat64{
					slice: []float64{1.0, 2.0},
				},
			},
			want3: []float64{1.0, 2.0},
		},
	}

	for i, tt := range tests {
		if i < 2 {
			t.Run(tt.name, func(t *testing.T) {
				got := ArrayDeduplicate(tt.args.argsInt64.slice)
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("ArrayDumpDuplicate() = %v, want %v", got, tt.want)
				}
			})
		} else if i < 4 {
			t.Run(tt.name, func(t *testing.T) {
				got := ArrayDeduplicate(tt.args.argsString.slice)
				if !reflect.DeepEqual(got, tt.want2) {
					t.Errorf("ArrayDumpDuplicate() = %v, want %v", got, tt.want2)
				}
			})
		} else {
			t.Run(tt.name, func(t *testing.T) {
				got := ArrayDeduplicate(tt.args.argsFloat64.slice)
				if !reflect.DeepEqual(got, tt.want3) {
					t.Errorf("ArrayDumpDuplicate() = %v, want %v", got, tt.want3)
				}
			})
		}
	}
}

func TestIsPhoneNumber(t *testing.T) {
	type args struct {
		phone string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "false",
			args: args{
				phone: "1082736464",
			},
			want: false,
		},
		{
			name: "true",
			args: args{
				phone: "0812345678",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPhoneNumber(tt.args.phone); got != tt.want {
				t.Errorf("IsPhoneNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsEmail(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "false",
			args: args{
				email: "abcaaa",
			},
			want: false,
		},
		{
			name: "true",
			args: args{
				email: "a@mail.com",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmail(tt.args.email); got != tt.want {
				t.Errorf("IsEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
