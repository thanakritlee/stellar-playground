package main

import (
	"fmt"
	"log"
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

	destination := os.Getenv("STELLAR_DESTINATION")
	source := os.Getenv("STELLAR_SOURCE")
	sourceSecretSeed := os.Getenv("STELLAR_SOURCE_SEED")

	// Make sure destination acount exists.
	if _, err := horizon.DefaultTestNetClient.LoadAccount(destination); err != nil {
		panic(err)
	}

	// Make sure source acount exists.
	if _, err := horizon.DefaultTestNetClient.LoadAccount(source); err != nil {
		panic(err)
	}

	tx, err := build.Transaction(
		build.TestNetwork,
		build.SourceAccount{
			AddressOrSeed: source,
		},
		build.AutoSequence{
			SequenceProvider: horizon.DefaultTestNetClient,
		},
		build.Payment(
			build.Destination{
				AddressOrSeed: destination,
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
	txe, err := tx.Sign(sourceSecretSeed)
	if err != nil {
		panic(err)
	}

	txeB64, err := txe.Base64()
	if err != nil {
		panic(err)
	}

	// Send it off to Stellar.
	resp, err := horizon.DefaultTestNetClient.SubmitTransaction(txeB64)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successful Transaction:")
	fmt.Println("Ledger:", resp.Ledger)
	fmt.Println("Hash:", resp.Hash)

}
