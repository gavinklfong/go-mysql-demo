package model

import (
	"time"
)

type ForexRateBooking struct {
	ID                 string
	Timestamp          time.Time
	BaseCurrency       string
	CounterCurrency    string
	Rate               float32
	TradeAction        string
	BaseCurrencyAmount float32
	BookingRef         string
	ExpiryTime         time.Time
	CustomerID         string
}
