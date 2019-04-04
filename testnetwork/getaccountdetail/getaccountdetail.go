package main

import (
	"fmt"
	"log"

	"github.com/stellar/go/clients/horizon"
	"github.com/thanakritlee/stellar-playground/testnetwork/createtestaccount"
)

func main() {

	pair := createtestaccount.CreateAccount()

	account, err := horizon.DefaultTestNetClient.LoadAccount(pair.Address())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Balances for account:", pair.Address())
	fmt.Println("Account secret seed: ", pair.Seed())

	// Accounts can carry multiple balances - one for each type of
	// currency they hold.
	for _, balance := range account.Balances {
		log.Println(balance)
	}
}
