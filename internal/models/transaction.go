package models

import "time"

type Transaction struct {
	ID   string
	Date time.Time

	Amount int64
	Title  string

	CategoryID  string
	AccountId   string
	Description string
}
