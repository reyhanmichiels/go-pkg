package null

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
)

// nullable int64 type.
// null will set valid value to false.
// any int64 value will be considered null when valid is set to false
// SqlNull is for updating SQL DB value to null
type Int64 struct {
	Int64   int64
	Valid   bool
	SqlNull bool
}

// create new nullable int64
func NewInt64(i int64, valid bool) Int64 {
	return Int64{
		Int64: i,
		Valid: valid,
	}
}

// create valid nullable int64
func Int64From(i int64) Int64 {
	return NewInt64(i, true)
}

func (i *Int64) Scan(value interface{}) error {
	var sqli sql.NullInt64
	if err := sqli.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*i = Int64{sqli.Int64, false, false}
	} else {
		*i = Int64{sqli.Int64, true, false}
	}
	return nil
}

func (i Int64) Value() (driver.Value, error) {
	if !i.Valid {
		return nil, nil
	}
	return i.Int64, nil
}

func (i *Int64) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return nullBytes, nil
	}
	return json.Marshal(i.Int64)
}

func (i *Int64) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, nullBytes) {
		return nil
	}
	err := json.Unmarshal(b, &i.Int64)
	i.Valid = (err == nil)
	return err
}

// will return true if invalid or value is zero
func (i Int64) IsNullOrZero() bool {
	return !i.Valid || i.Int64 == 0
}

// returns true if both invalid or both have same value
func (i Int64) Equal(other Int64) bool {
	return i.Valid == other.Valid && (!i.Valid || i.Int64 == other.Int64)
}

// returns true if valid and both have same value
func (i Int64) Is(other int64) bool {
	return i.Equal(Int64From(other))
}
