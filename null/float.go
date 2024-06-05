package null

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
)

// nullable float64 type.
// null will set valid value to false.
// any float64 value will be considered null when valid is set to false
type Float64 struct {
	Float64 float64
	Valid   bool
}

// create new nullable float64
func NewFloat64(f float64, valid bool) Float64 {
	return Float64{
		Float64: f,
		Valid:   valid,
	}
}

// create valid nullable float64
func Float64From(f float64) Float64 {
	return NewFloat64(f, true)
}

func (f *Float64) Scan(value interface{}) error {
	var sqlf sql.NullFloat64
	if err := sqlf.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*f = Float64{sqlf.Float64, false}
	} else {
		*f = Float64{sqlf.Float64, true}
	}

	return nil
}

func (f Float64) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.Float64, nil
}

func (f *Float64) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return nullBytes, nil
	}
	return json.Marshal(f.Float64)
}

func (f *Float64) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, nullBytes) {
		return nil
	}
	err := json.Unmarshal(b, &f.Float64)
	f.Valid = (err == nil)
	return err
}

// will return true if invalid or value is empty
func (f Float64) IsNullOrZero() bool {
	return !f.Valid || f.Float64 == 0
}

// returns true if both invalid or both have same value
func (f Float64) Equal(other Float64) bool {
	return f.Valid == other.Valid && (!f.Valid || f.Float64 == other.Float64)
}

// returns true if valid and both have same value
func (f Float64) Is(other float64) bool {
	return f.Equal(Float64From(other))
}
