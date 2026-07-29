package models

import "time"

type Transaction struct {
	ID   int64
	Date time.Time

	Amount int64
	Title  string

	CategoryID  string
	AccountID   string
	Description string
}
