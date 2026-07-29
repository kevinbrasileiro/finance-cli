package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/kevinbrasileiro/finance-cli/internal/service"
	"github.com/kevinbrasileiro/finance-cli/internal/storage"
)

func main() {
	switch os.Args[1] {
	case "income":
		args := os.Args[2:]
		if len(args) < 2 {
			fmt.Println("usage: fin income <amount> <title>")
			os.Exit(1)
		}
		amount, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Printf("invalid amount: %q\n", args[0])
			os.Exit(1)
		}
		saved, err := service.RecordIncome(amount, args[1])
		if err != nil {
			panic(err)
		}
		fmt.Printf("added transaction #%d\n", saved.ID)

	case "expense":
		args := os.Args[2:]
		if len(args) < 2 {
			fmt.Println("usage: fin expense <amount> <title>")
			os.Exit(1)
		}
		amount, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Printf("invalid amount: %q\n", args[0])
			os.Exit(1)
		}
		saved, err := service.RecordExpense(amount, args[1])
		if err != nil {
			panic(err)
		}
		fmt.Printf("added transaction #%d\n", saved.ID)

	case "list":
		transactions, err := storage.GetAllTransactions()
		if err != nil {
			panic(err)
		}
		for _, transaction := range transactions {
			fmt.Printf("%+v\n", transaction)
		}

	default:
		fmt.Printf("unknown command: %q\n", os.Args[1])
		os.Exit(1)
	}
}
