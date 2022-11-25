package models

import (
	"time"
)

const (
	BtcusdtTmeFormat = "2006-01-02 15:04:05"
)

//easyjson:json
type HistoryBtcusdt struct {
	Total   uint32     `json:"total"`
	History []*Btcusdt `json:"history"`
}

//easyjson:json
type Btcusdt struct {
	Time         *time.Time `db:"timestamp"`
	Buy          float64    `db:"buy"`
	Sell         float64    `db:"sell"`
	High         float64    `db:"high"`
	Low          float64    `db:"low"`
	Last         float64    `db:"last"`
	AveragePrice float64    `db:"average_price"`
}

//easyjson:json
type HistoryBtc struct {
	Total   uint32                       `json:"total"`
	History map[string]map[string]*BtcDB `json:"history"`
}

//easyjson:json
type BtcDB struct {
	Time         *time.Time `db:"timestamp"       json:"-"`
	CodeStr      string     `db:"code"            json:"-"`
	CodeID       uint32     `db:"code_id"         json:"-"`
	Buy          float64    `db:"buy"             json:"buy"`
	Sell         float64    `db:"sell"            json:"sell"`
	High         float64    `db:"high"            json:"high"`
	Low          float64    `db:"low"             json:"low"`
	Last         float64    `db:"last"            json:"last"`
	AveragePrice float64    `db:"average_price"   json:"average_price"`
}
