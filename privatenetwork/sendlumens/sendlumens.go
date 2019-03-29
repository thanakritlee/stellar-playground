package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
)

func main() {

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loding .env file")
	}

	account0Seed := os.Getenv("ACCOUNT_0_SEED")
	account0Address := os.Getenv("ACCOUNT_0_ADDRESS")
	account1Address := os.Getenv("ACCOUNT_1_ADDRESS")
	url := os.Getenv("URL")
	networkPassPhrase := os.Getenv("NETWORK_PASSPHRASE")

	horizonClient := &horizon.Client{
		URL:  url,
		HTTP: http.DefaultClient,
	}

	// Make sure destination acount exists.
	if _, err := horizonClient.LoadAccount(account1Address); err != nil {
		panic(err)
	}

	// Make sure source acount exists.
	if _, err := horizonClient.LoadAccount(account0Address); err != nil {
		panic(err)
	}

	tx, err := build.Transaction(
		build.Network{Passphrase: networkPassPhrase},
		build.SourceAccount{
			AddressOrSeed: account0Address,
		},
		build.AutoSequence{
			SequenceProvider: horizonClient,
		},
		build.Payment(
			build.Destination{
				AddressOrSeed: account1Address,
			},
			build.NativeAmount{
				Amount: "10",
			},
		),
	)

	if err != nil {
		panic(err)
	}

	// Sign the transaction to prove you are actually the person sending it.
	// Cryptographically sign it using the source secret seed.
	txe, err := tx.Sign(account0Seed)
	if err != nil {
		panic(err)
	}

	txeB64, err := txe.Base64()
	if err != nil {
		panic(err)
	}

	// Send it off to Stellar.
	resp, err := horizonClient.SubmitTransaction(txeB64)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successful Transaction:")
	fmt.Println("Ledger:", resp.Ledger)
	fmt.Println("Hash:", resp.Hash)

}
