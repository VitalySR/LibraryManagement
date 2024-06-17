package repository

import (
	"database/sql"
	"encoding/json"
)

type StringNull struct {
	sql.NullString
}
type Int32Null struct {
	sql.NullInt32
}

func (s *StringNull) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(s.String)
}

func (s *StringNull) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &s.String)
	s.Valid = err == nil
	return err
}

func (v *Int32Null) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return json.Marshal(0)
	}
	return json.Marshal(v.Int32)
}

func (v *Int32Null) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &v.Int32)
	v.Valid = err == nil
	return err
}

func StringWithNull(s *string) sql.NullString {
	if s == nil || len(*s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: *s,
		Valid:  true,
	}
}

func Int32WithNull(i *int32) sql.NullInt32 {
	if i == nil || *i == 0 {
		return sql.NullInt32{}
	}
	return sql.NullInt32{
		Int32: *i,
		Valid: true,
	}
}
