package storage

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/kevinbrasileiro/finance-cli/internal/models"
)

const dateLayout = "2006-01-02"
const transactionFields = 7

func dataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("locating home directory: %w", err)
	}

	dir := filepath.Join(home, ".fin")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("creating data directory %q: %w", dir, err)
	}
	return dir, nil
}

func transactionsFilePath() (string, error) {
	dir, err := dataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "transactions.tsv"), nil
}

func transactionToRecord(t models.Transaction) []string {
	return []string{
		t.ID,
		t.Date.Format(dateLayout),
		strconv.FormatInt(t.Amount, 10),
		t.Title,
		t.CategoryID,
		t.AccountId,
		t.Description,
	}
}

func recordToTransaction(record []string) (models.Transaction, error) {
	if len(record) != transactionFields {
		return models.Transaction{}, fmt.Errorf(
			"malformed record: expected %d fields, got %d", transactionFields, len(record),
		)
	}

	date, err := time.Parse(dateLayout, record[1])
	if err != nil {
		return models.Transaction{}, fmt.Errorf("parsing date %q: %w", record[1], err)
	}

	amount, err := strconv.ParseInt(record[2], 10, 64)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("parsing amount %q: %w", record[2], err)
	}

	return models.Transaction{
		ID:          record[0],
		Date:        date,
		Amount:      amount,
		Title:       record[3],
		CategoryID:  record[4],
		AccountId:   record[5],
		Description: record[6],
	}, nil
}

func AddTransaction(transaction models.Transaction) error {
	path, err := transactionsFilePath()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = '\t'

	if err := writer.Write(transactionToRecord(transaction)); err != nil {
		return fmt.Errorf("writing transaction: %w", err)
	}

	writer.Flush()
	return writer.Error()
}

func ReadTransactions() ([]models.Transaction, error) {
	path, err := transactionsFilePath()
	if err != nil {
		return []models.Transaction{}, err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDONLY, 0o644)
	if err != nil {
		return []models.Transaction{}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	reader.FieldsPerRecord = transactionFields

	var transactions []models.Transaction
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return []models.Transaction{}, fmt.Errorf("reading transactions: %w", err)
		}

		transaction, err := recordToTransaction(record)
		if err != nil {
			return []models.Transaction{}, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
