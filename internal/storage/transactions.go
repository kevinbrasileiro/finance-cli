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

const transactionFields = 7

func transactionsFilePath() (string, error) {
	dir, err := dataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "transactions.tsv"), nil
}

func transactionToRecord(t models.Transaction) []string {
	return []string{
		strconv.FormatInt(t.ID, 10),
		t.Date.Format(dateLayout),
		strconv.FormatInt(t.Amount, 10),
		t.Title,
		t.CategoryID,
		t.AccountID,
		t.Description,
	}
}

func recordToTransaction(record []string) (models.Transaction, error) {
	if len(record) != transactionFields {
		return models.Transaction{}, fmt.Errorf(
			"malformed record: expected %d fields, got %d", transactionFields, len(record),
		)
	}

	id, err := strconv.ParseInt(record[0], 10, 64)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("parsing id %q: %w", record[0], err)
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
		ID:          id,
		Date:        date,
		Amount:      amount,
		Title:       record[3],
		CategoryID:  record[4],
		AccountID:   record[5],
		Description: record[6],
	}, nil
}

func AddTransaction(transaction models.Transaction) (models.Transaction, error) {
	id, err := nextID("transactions")
	if err != nil {
		return models.Transaction{}, err
	}
	transaction.ID = id

	path, err := transactionsFilePath()
	if err != nil {
		return models.Transaction{}, err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return models.Transaction{}, err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = '\t'

	if err := writer.Write(transactionToRecord(transaction)); err != nil {
		return models.Transaction{}, fmt.Errorf("writing transaction: %w", err)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return models.Transaction{}, err
	}

	return transaction, nil
}

func GetAllTransactions() ([]models.Transaction, error) {
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
