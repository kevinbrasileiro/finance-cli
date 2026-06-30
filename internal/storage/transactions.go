package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kevinbrasileiro/finance-cli/internal/models"
)

func transactionsFilePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".fin", "transactions.tsv")
}

func AddTransaction(transaction models.Transaction) error {
	path := transactionsFilePath()

	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	line := fmt.Sprintf(
		"%s\t%s\t%d\t%s\t%s\t%s\t%s\n",
		transaction.ID,
		transaction.Date.Format("2006-01-02"),

		transaction.Amount,
		transaction.Title,

		transaction.CategoryID,
		transaction.AccountId,
		transaction.Description,
	)

	_, err = file.WriteString(line)
	return err
}

func ListTransactions() error {
	path := transactionsFilePath()

	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	fmt.Print(string(file))
	return nil
}
