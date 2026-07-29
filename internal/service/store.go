package service

import (
	"time"

	"github.com/kevinbrasileiro/finance-cli/internal/models"
	"github.com/kevinbrasileiro/finance-cli/internal/storage"
)

func RecordIncome(amount int64, title string) (models.Transaction, error) {
	tx := models.Transaction{
		Date:   time.Now(),
		Amount: amount,
		Title:  title,
	}
	return storage.AddTransaction(tx)
}

func RecordExpense(amount int64, title string) (models.Transaction, error) {
	tx := models.Transaction{
		Date:   time.Now(),
		Amount: -amount,
		Title:  title,
	}
	return storage.AddTransaction(tx)
}
