package null

import (
	"database/sql/driver"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTime(t *testing.T) {
	mockTime := time.Now()
	type args struct {
		t     time.Time
		valid bool
	}
	tests := []struct {
		name string
		args args
		want Time
	}{
		{
			name: "set invalid",
			args: args{
				t:     mockTime,
				valid: false,
			},
			want: Time{Valid: false, Time: mockTime},
		},
		{
			name: "set valid",
			args: args{
				t:     mockTime,
				valid: true,
			},
			want: Time{Valid: true, Time: mockTime},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTime(tt.args.t, tt.args.valid)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTimeFrom(t *testing.T) {
	mockTime := time.Now()
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want Time
	}{
		{
			name: "zero time",
			args: args{t: time.Time{}},
			want: Time{Valid: true, Time: time.Time{}},
		},
		{
			name: "valid time",
			args: args{t: mockTime},
			want: Time{Valid: true, Time: mockTime},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TimeFrom(tt.args.t)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTime_Scan(t *testing.T) {
	mockTime := time.Date(2022, 02, 02, 0, 0, 0, 0, time.UTC)
	mockTimeString := "2022-02-02T00:00:00Z"
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    Time
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
			name:    "scan zero time string",
			args:    args{value: "0001-01-01T00:00:00Z"},
			wantErr: true,
		},
		{
			name:    "scan zero time",
			args:    args{value: time.Time{}},
			wantErr: false,
			want:    Time{Valid: true, Time: time.Time{}},
		},
		{
			name:    "scan int",
			args:    args{value: 123},
			wantErr: true,
		},
		{
			name:    "scan valid time string",
			args:    args{value: mockTimeString},
			wantErr: true,
		},
		{
			name:    "scan valid time",
			args:    args{value: mockTime},
			wantErr: false,
			want:    Time{Valid: true, Time: mockTime},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nt := Time{}
			if err := nt.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Time.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, nt)
		})
	}
}

func TestTime_Value(t *testing.T) {
	mockTime := time.Now()
	tests := []struct {
		name    string
		fields  Time
		want    driver.Value
		wantErr bool
	}{
		{
			name:    "value of nil",
			wantErr: false,
			want:    nil,
		},
		{
			name:    "value of zero time",
			fields:  Time{Valid: true, Time: time.Time{}},
			wantErr: false,
			want:    time.Time{},
		},
		{
			name:    "value of int",
			fields:  Time{Valid: true, Time: mockTime},
			wantErr: false,
			want:    mockTime,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Time.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTime_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		fields  Time
		want    []byte
		wantErr bool
	}{
		{
			name:    "marshal invalid",
			wantErr: false,
			want:    []byte("null"),
		},
		{
			name:    "marshal zero time",
			fields:  Time{Valid: true, Time: time.Time{}},
			wantErr: false,
			want:    []byte("\"0001-01-01T00:00:00Z\""),
		},
		{
			name:    "marshal valid time",
			fields:  Time{Valid: true, Time: time.Date(2022, 02, 02, 0, 0, 0, 0, time.UTC)},
			wantErr: false,
			want:    []byte("\"2022-02-02T00:00:00Z\""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Time.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTime_UnmarshalJSON(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    Time
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
			name: "unmarshal string of zero time",
			args: args{
				b: []byte("\"0001-01-01T00:00:00Z"),
			},
			wantErr: false,
			want:    Time{Valid: true, Time: time.Time{}},
		},
		{
			name: "unmarshal string of valid time",
			args: args{
				b: []byte("\"2022-02-02T00:00:00Z"),
			},
			wantErr: false,
			want:    Time{Valid: true, Time: time.Date(2022, 02, 02, 00, 00, 00, 00, time.UTC)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nt := Time{}
			if err := nt.UnmarshalJSON(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("Time.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, nt)
		})
	}
}

func TestTime_IsNullOrZero(t *testing.T) {
	tests := []struct {
		name   string
		fields Time
		want   bool
	}{
		{
			name: "null or invalid",
			want: true,
		},
		{
			name:   "valid zero time",
			fields: Time{Valid: true, Time: time.Time{}},
			want:   true,
		},
		{
			name:   "valid with valid time",
			fields: Time{Valid: true, Time: time.Now()},
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

func TestTime_Equal(t *testing.T) {
	mockTime := time.Now()
	mockTime2 := mockTime.AddDate(0, 0, 1)
	type args struct {
		other Time
	}
	tests := []struct {
		name   string
		fields Time
		args   args
		want   bool
	}{
		{
			name:   "both invalid",
			fields: Time{Valid: false},
			args:   args{Time{Valid: false}},
			want:   true,
		},
		{
			name:   "both valid and both same value",
			fields: Time{Valid: true, Time: mockTime},
			args:   args{Time{Valid: true, Time: mockTime}},
			want:   true,
		},
		{
			name:   "valid and invalid but same value",
			fields: Time{Valid: true, Time: mockTime},
			args:   args{Time{Valid: false, Time: mockTime}},
			want:   false,
		},
		{
			name:   "both valid but different value",
			fields: Time{Valid: true, Time: mockTime},
			args:   args{Time{Valid: true, Time: mockTime2}},
			want:   false,
		},
		{
			name:   "both valid and both zero",
			fields: Time{Valid: true, Time: time.Time{}},
			args:   args{Time{Valid: true, Time: time.Time{}}},
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

func TestTime_Is(t *testing.T) {
	mockTime1 := time.Now()
	mockTime2 := mockTime1.Add(time.Hour)

	type args struct {
		other time.Time
	}
	tests := []struct {
		name   string
		fields Time
		args   args
		want   bool
	}{
		{
			name:   "origin invalid and not equal",
			fields: Time{Valid: false, Time: mockTime1},
			args:   args{other: mockTime2},
			want:   false,
		},
		{
			name:   "origin invalid and equal",
			fields: Time{Valid: false, Time: mockTime1},
			args:   args{other: mockTime1},
			want:   false,
		},
		{
			name:   "origin valid and not equal",
			fields: Time{Valid: true, Time: mockTime1},
			args:   args{other: mockTime2},
			want:   false,
		},
		{
			name:   "origin valid and equal",
			fields: Time{Valid: true, Time: mockTime1},
			args:   args{other: mockTime1},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.Is(tt.args.other); got != tt.want {
				t.Errorf("Time.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTime_IsBefore(t *testing.T) {
	mockTime1 := time.Now()
	mockTime2 := mockTime1.Add(time.Hour)

	type args struct {
		other time.Time
	}
	tests := []struct {
		name   string
		fields Time
		args   args
		want   bool
	}{
		{
			name:   "origin invalid and other is not before the origin",
			fields: Time{Valid: false, Time: mockTime1},
			args:   args{other: mockTime1},
			want:   false,
		},
		{
			name:   "origin invalid and other is before the origin",
			fields: Time{Valid: false, Time: mockTime2},
			args:   args{other: mockTime1},
			want:   false,
		},
		{
			name:   "origin valid and other is not before the origin",
			fields: Time{Valid: true, Time: mockTime1},
			args:   args{other: mockTime1},
			want:   false,
		},
		{
			name:   "origin valid and other is before the origin",
			fields: Time{Valid: true, Time: mockTime2},
			args:   args{other: mockTime1},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.IsBefore(tt.args.other); got != tt.want {
				t.Errorf("Time.IsBefore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTime_IsAfter(t *testing.T) {
	mockTime1 := time.Now()
	mockTime2 := mockTime1.Add(time.Hour)

	type args struct {
		other time.Time
	}
	tests := []struct {
		name   string
		fields Time
		args   args
		want   bool
	}{
		{
			name:   "origin invalid and other is not after the origin",
			fields: Time{Valid: false, Time: mockTime1},
			args:   args{other: mockTime1},
			want:   false,
		},
		{
			name:   "origin invalid and is after the origin",
			fields: Time{Valid: false, Time: mockTime1},
			args:   args{other: mockTime2},
			want:   false,
		},
		{
			name:   "origin valid and is not after the origin",
			fields: Time{Valid: true, Time: mockTime1},
			args:   args{other: mockTime1},
			want:   false,
		},
		{
			name:   "origin valid and is after the origin",
			fields: Time{Valid: true, Time: mockTime1},
			args:   args{other: mockTime2},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.IsAfter(tt.args.other); got != tt.want {
				t.Errorf("Time.IsAfter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDate(t *testing.T) {
	mockTime := time.Now()
	mockTimeResponse := time.Date(mockTime.Year(), mockTime.Month(), mockTime.Day(), 0, 0, 0, 0, mockTime.Location())
	type args struct {
		t     time.Time
		valid bool
	}
	tests := []struct {
		name string
		args args
		want Date
	}{
		{
			name: "set invalid",
			args: args{
				t:     mockTime,
				valid: false,
			},
			want: Date{Valid: false, Time: mockTimeResponse},
		},
		{
			name: "set valid",
			args: args{
				t:     mockTime,
				valid: true,
			},
			want: Date{Valid: true, Time: mockTimeResponse},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDate(tt.args.t, tt.args.valid)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDateFrom(t *testing.T) {
	mockTime := time.Now()
	mockTimeResponse := time.Date(mockTime.Year(), mockTime.Month(), mockTime.Day(), 0, 0, 0, 0, mockTime.Location())
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want Date
	}{
		{
			name: "zero time",
			args: args{t: time.Time{}},
			want: Date{Valid: true, Time: time.Time{}},
		},
		{
			name: "valid time",
			args: args{t: mockTime},
			want: Date{Valid: true, Time: mockTimeResponse},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DateFrom(tt.args.t)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDate_Scan(t *testing.T) {
	getStrPtr := func(s string) *string { return &s }
	mockTime := time.Date(2022, 02, 02, 0, 0, 0, 0, time.UTC)
	mockTimeString := "2022-02-02T00:00:00Z"
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		fields  Date
		args    args
		wantErr bool
		want    Date
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
			name:    "scan zero time string",
			args:    args{value: "0001-01-01T00:00:00Z"},
			wantErr: false,
			want:    Date{Valid: true, Time: time.Time{}},
		},
		{
			name:    "scan zero time",
			args:    args{value: time.Time{}},
			wantErr: false,
			want:    Date{Valid: true, Time: time.Time{}},
		},
		{
			name:    "scan int",
			args:    args{value: 123},
			wantErr: true,
		},
		{
			name:    "scan valid time string",
			args:    args{value: mockTimeString},
			wantErr: false,
			want:    Date{Valid: true, Time: mockTime},
		},
		{
			name:    "scan valid time string pointer",
			args:    args{value: getStrPtr(mockTimeString)},
			wantErr: true,
		},
		{
			name:    "scan valid time",
			args:    args{value: mockTime},
			wantErr: false,
			want:    Date{Valid: true, Time: mockTime},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nd := Date{}
			if err := nd.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Date.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, nd)
		})
	}
}

func TestDate_Value(t *testing.T) {
	mockTime := time.Now()
	mockTimeResponse := time.Date(mockTime.Year(), mockTime.Month(), mockTime.Day(), 0, 0, 0, 0, mockTime.Location())
	tests := []struct {
		name    string
		fields  Date
		want    driver.Value
		wantErr bool
	}{
		{
			name:    "value of nil",
			wantErr: false,
			want:    nil,
		},
		{
			name:    "value of zero time",
			fields:  Date{Valid: true, Time: time.Time{}},
			wantErr: false,
			want:    time.Time{},
		},
		{
			name:    "value of int",
			fields:  Date{Valid: true, Time: mockTime},
			wantErr: false,
			want:    mockTimeResponse,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Date.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDate_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		fields  Date
		want    []byte
		wantErr bool
	}{
		{
			name:    "marshal invalid",
			wantErr: false,
			want:    []byte("null"),
		},
		{
			name:    "marshal zero time",
			fields:  Date{Valid: true, Time: time.Time{}},
			wantErr: false,
			want:    []byte("\"0001-01-01T00:00:00Z\""),
		},
		{
			name:    "marshal valid time",
			fields:  Date{Valid: true, Time: time.Date(2022, 02, 02, 0, 0, 0, 0, time.UTC)},
			wantErr: false,
			want:    []byte("\"2022-02-02T00:00:00Z\""),
		},
		{
			name:    "marshal valid time with set time to 0",
			fields:  Date{Valid: true, Time: time.Date(2022, 02, 02, 1, 3, 2, 1, time.UTC)},
			wantErr: false,
			want:    []byte("\"2022-02-02T00:00:00Z\""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Date.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDate_UnmarshalJSON(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    Date
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
			name: "unmarshal string of zero time",
			args: args{
				b: []byte("\"0001-01-01T00:00:00Z"),
			},
			wantErr: false,
			want:    Date{Valid: true, Time: time.Time{}},
		},
		{
			name: "unmarshal string of valid time",
			args: args{
				b: []byte("\"2022-02-02T00:00:00Z"),
			},
			wantErr: false,
			want:    Date{Valid: true, Time: time.Date(2022, 02, 02, 00, 00, 00, 00, time.UTC)},
		},
		{
			name: "unmarshal string of valid time with set time to 0",
			args: args{
				b: []byte("\"2022-02-02T10:10:10Z"),
			},
			wantErr: false,
			want:    Date{Valid: true, Time: time.Date(2022, 02, 02, 00, 00, 00, 00, time.UTC)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nd := Date{}
			if err := nd.UnmarshalJSON(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("Date.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, nd)
		})
	}
}

func TestDate_IsNullOrZero(t *testing.T) {
	tests := []struct {
		name   string
		fields Date
		want   bool
	}{
		{
			name: "null or invalid",
			want: true,
		},
		{
			name:   "valid zero time",
			fields: Date{Valid: true, Time: time.Time{}},
			want:   true,
		},
		{
			name:   "valid with valid time",
			fields: Date{Valid: true, Time: time.Now()},
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

func TestDate_Equal(t *testing.T) {
	mockTime := time.Now()
	mockTime2 := mockTime.AddDate(0, 0, 1)
	type args struct {
		other Date
	}
	tests := []struct {
		name   string
		fields Date
		args   args
		want   bool
	}{
		{
			name:   "both invalid",
			fields: Date{Valid: false},
			args:   args{Date{Valid: false}},
			want:   true,
		},
		{
			name:   "both valid and both same value",
			fields: Date{Valid: true, Time: mockTime},
			args:   args{Date{Valid: true, Time: mockTime}},
			want:   true,
		},
		{
			name:   "valid and invalid but same value",
			fields: Date{Valid: true, Time: mockTime},
			args:   args{Date{Valid: false, Time: mockTime}},
			want:   false,
		},
		{
			name:   "both valid but different value",
			fields: Date{Valid: true, Time: mockTime},
			args:   args{Date{Valid: true, Time: mockTime2}},
			want:   false,
		},
		{
			name:   "both valid and both zero",
			fields: Date{Valid: true, Time: time.Time{}},
			args:   args{Date{Valid: true, Time: time.Time{}}},
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

func TestDate_Is(t *testing.T) {
	mockTime1 := time.Now()
	mockTime2 := mockTime1.Add(time.Hour)

	type args struct {
		other time.Time
	}
	tests := []struct {
		name   string
		fields Time
		args   args
		want   bool
	}{
		{
			name:   "origin invalid and origin is not equal",
			fields: Time{Valid: false, Time: mockTime1},
			args:   args{other: mockTime2},
			want:   false,
		},
		{
			name:   "origin invalid and origin is equal",
			fields: Time{Valid: false, Time: mockTime1},
			args:   args{other: mockTime1},
			want:   false,
		},
		{
			name:   "origin valid and origin is not equal",
			fields: Time{Valid: true, Time: mockTime1},
			args:   args{other: mockTime2},
			want:   false,
		},
		{
			name:   "origin valid and origin is equal",
			fields: Time{Valid: true, Time: mockTime1},
			args:   args{other: mockTime1},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.Is(tt.args.other); got != tt.want {
				t.Errorf("Date.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_IsBefore(t *testing.T) {
	mockTime1 := time.Now()
	mockTime2 := mockTime1.Add(time.Hour)

	type args struct {
		other time.Time
	}
	tests := []struct {
		name   string
		fields Time
		args   args
		want   bool
	}{
		{
			name:   "origin invalid and other is not before the origin",
			fields: Time{Valid: false, Time: mockTime1},
			args:   args{other: mockTime1},
			want:   false,
		},
		{
			name:   "origin invalid and other is before the origin",
			fields: Time{Valid: false, Time: mockTime2},
			args:   args{other: mockTime1},
			want:   false,
		},
		{
			name:   "origin valid and other is not before the origin",
			fields: Time{Valid: true, Time: mockTime1},
			args:   args{other: mockTime1},
			want:   false,
		},
		{
			name:   "origin valid and other is before the origin",
			fields: Time{Valid: true, Time: mockTime2},
			args:   args{other: mockTime1},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.IsBefore(tt.args.other); got != tt.want {
				t.Errorf("Date.IsBefore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_IsAfter(t *testing.T) {
	mockTime1 := time.Now()
	mockTime2 := mockTime1.Add(time.Hour)

	type args struct {
		other time.Time
	}
	tests := []struct {
		name   string
		fields Time
		args   args
		want   bool
	}{
		{
			name:   "origin invalid and other is not after the origin",
			fields: Time{Valid: false, Time: mockTime1},
			args:   args{other: mockTime1},
			want:   false,
		},
		{
			name:   "origin invalid and other is after the origin",
			fields: Time{Valid: false, Time: mockTime1},
			args:   args{other: mockTime2},
			want:   false,
		},
		{
			name:   "origin valid and other is not after the origin",
			fields: Time{Valid: true, Time: mockTime1},
			args:   args{other: mockTime1},
			want:   false,
		},
		{
			name:   "origin valid and other is after the origin",
			fields: Time{Valid: true, Time: mockTime1},
			args:   args{other: mockTime2},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.IsAfter(tt.args.other); got != tt.want {
				t.Errorf("Date.IsAfter() = %v, want %v", got, tt.want)
			}
		})
	}
}
