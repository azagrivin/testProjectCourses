package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

const (
	CurrencyDateFormat = "2006-01-02"
)

//easyjson:json
type Currency struct {
	Time    time.Time `db:"time"`
	CodeID  uint32    `db:"code_id"`
	CodeStr string    `db:"code_str"`
	Value   float64   `db:"value"`
}

//easyjson:json
type HistoryCurrencies struct {
	Total   uint32    `json:"total"`
	History []History `json:"history"`
}

//easyjson:json
type History map[string]interface{}

func (h *History) Scan(val interface{}) error {
	switch v := val.(type) {
	case string:
		return json.Unmarshal([]byte(v), h)
	case []byte:
		return json.Unmarshal(v, h)
	default:
		return fmt.Errorf("scan pkg.NullJson error, unsupported type %T", v)
	}
}

//easyjson:json
func (h *History) Value() (driver.Value, error) {
	bytes, err := json.Marshal(h)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
