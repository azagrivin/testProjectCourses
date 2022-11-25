package models

import "time"

type FilterBtcusdt struct {
	TimeFrom  *time.Time
	TimeTo    *time.Time
	Limit     uint64
	Page      uint64
	FirstOnly bool
}

type FilterCurrencies struct {
	TimeFrom   *time.Time
	TimeTo     *time.Time
	Currencies []string
	Limit      uint64
	Page       uint64
	FirstOnly  bool
}
