package main

import (
	"time"

	"github.com/kevinbrasileiro/finance-cli/internal/models"
	"github.com/kevinbrasileiro/finance-cli/internal/storage"
)

func main() {
	tx := models.Transaction{
		ID:          "1",
		Date:        time.Now(),
		Amount:      1250,
		Title:       "Pizza",
		CategoryID:  "food",
		AccountId:   "cash",
		Description: "test",
	}

	err := storage.AddTransaction(tx)
	if err != nil {
		panic(err)
	}

	wow := storage.ListTransactions()
	if wow != nil {
		panic(wow)
	}
}
