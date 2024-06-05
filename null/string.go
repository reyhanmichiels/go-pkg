package null

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"strings"
)

// nullable string type.
// null will set valid value to false.
// any string value will be considered null when valid is set to false
// SqlNull is for updating SQL DB value to null
type String struct {
	String  string
	Valid   bool
	SqlNull bool
}

// create new nullable string
func NewString(s string, valid bool) String {
	return String{
		String: s,
		Valid:  valid,
	}
}

// create valid nullable string
func StringFrom(s string) String {
	return NewString(s, true)
}

func (s *String) Scan(value interface{}) error {
	var sqls sql.NullString
	if err := sqls.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*s = String{sqls.String, false, false}
	} else {
		*s = String{sqls.String, true, false}
	}

	return nil
}

func (s String) Value() (driver.Value, error) {
	if !s.Valid {
		return nil, nil
	}
	return s.String, nil
}

func (s *String) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return nullBytes, nil
	}
	return json.Marshal(s.String)
}

func (s *String) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, nullBytes) {
		return nil
	}
	err := json.Unmarshal(b, &s.String)
	s.String = strings.TrimSpace(s.String)
	s.Valid = (err == nil)
	return err
}

// will return true if invalid or value is empty
func (s String) IsNullOrZero() bool {
	return !s.Valid || s.String == ""
}

// returns true if both invalid or both have same value
func (s String) Equal(other String) bool {
	return s.Valid == other.Valid && (!s.Valid || s.String == other.String)
}

// returns true if valid and both have same value
func (s String) Is(other string) bool {
	return s.Equal(StringFrom(other))
}
