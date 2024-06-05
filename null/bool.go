package null

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
)

var nullBytes = []byte("null")

// nullable boolean type.
// null will set valid value to false.
// any boolean value will be considered null when valid is set to false
type Bool struct {
	Bool  bool
	Valid bool
}

// create new nullable boolean
func NewBool(b bool, valid bool) Bool {
	return Bool{
		Bool:  b,
		Valid: valid,
	}
}

// create valid nullable boolean
func BoolFrom(b bool) Bool {
	return NewBool(b, true)
}

func (bo *Bool) Scan(value interface{}) error {
	var sqlb sql.NullBool
	if err := sqlb.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*bo = Bool{sqlb.Bool, false}
	} else {
		*bo = Bool{sqlb.Bool, true}
	}
	return nil
}

func (bo Bool) Value() (driver.Value, error) {
	if !bo.Valid {
		return nil, nil
	}
	return bo.Bool, nil
}

func (bo *Bool) MarshalJSON() ([]byte, error) {
	if !bo.Valid {
		return nullBytes, nil
	}
	return json.Marshal(bo.Bool)
}

func (bo *Bool) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, nullBytes) {
		return nil
	}
	err := json.Unmarshal(b, &bo.Bool)
	bo.Valid = (err == nil)
	return err
}

// will return true if invalid or value is empty
func (bo Bool) IsNullOrZero() bool {
	return !bo.Valid || !bo.Bool
}

// returns true if both invalid or both have same value
func (bo Bool) Equal(other Bool) bool {
	return bo.Valid == other.Valid && (!bo.Valid || bo.Bool == other.Bool)
}

// returns true if valid and both have same value
func (bo Bool) Is(other bool) bool {
	return bo.Equal(BoolFrom(other))
}
