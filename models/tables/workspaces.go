package tables

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"time"
)

type (
	NullTime   sql.NullTime
	Workspaces struct {
		ID          int       `db:"id"`
		Name        string    `db:"name"`
		Description string    `db:"description"`
		CreatedAt   time.Time `db:"created_at"`
		UpdatedAt   *NullTime `db:"updated_at"`
	}
)

func (ns *NullTime) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.Time)
}

func (ns *NullTime) Scan(value interface{}) error {
	var s sql.NullTime
	if err := s.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		ns.Valid = false
	} else {
		ns.Valid = true
		ns.Time = s.Time
	}

	return nil
}
