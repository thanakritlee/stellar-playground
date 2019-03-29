package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loding .env file")
	}

	sourceSeed := os.Getenv("STELLAR_SOURCE_SEED")
	destinsationSeed := os.Getenv("STELLAR_DESTINATION_SEED")
	url := os.Getenv("URL")

	horizonClient := &horizon.Client{
		URL:  url,
		HTTP: http.DefaultClient,
	}

	// Keys for accounts to issue and receive the new asset.
	issuer, err := keypair.Parse(sourceSeed)
	if err != nil {
		log.Fatal(err)
	}
	recipient, err := keypair.Parse(destinsationSeed)
	if err != nil {
		log.Fatal(err)
	}

	// Create an object to represent the new asset.
	corruptionBaht := build.CreditAsset("CorruptBaht", issuer.Address())

	// First, the receiving account must trust the asset.
	trustTx, err := build.Transaction(
		build.SourceAccount{
			AddressOrSeed: recipient.Address(),
		},
		build.AutoSequence{
			SequenceProvider: horizonClient,
		},
		build.TestNetwork,
		build.Trust(
			corruptionBaht.Code,
			corruptionBaht.Issuer,
			build.Limit("100.25"),
		),
	)
	fmt.Println("After build transaction")
	if err != nil {
		fmt.Println("Error at build transaction")
		log.Fatal(err)
	}

	trustTxe, err := trustTx.Sign(destinsationSeed)
	if err != nil {
		log.Fatal(err)
	}

	trustTxeB64, err := trustTxe.Base64()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Submitting transaction")
	res, err := horizonClient.SubmitTransaction(trustTxeB64)
	fmt.Println(res)
	if err != nil {
		fmt.Println("Error from submitting transaction")
		log.Fatal(err)
	}

	// Second, the issuing account actually sends a payment usng the asset.
	paymentTx, err := build.Transaction(
		build.SourceAccount{
			AddressOrSeed: issuer.Address(),
		},
		build.TestNetwork,
		build.AutoSequence{
			SequenceProvider: horizonClient,
		},
		build.Payment(
			build.Destination{
				AddressOrSeed: recipient.Address(),
			},
			build.CreditAmount{
				Code:   "CorruptBaht",
				Issuer: issuer.Address(),
				Amount: "10",
			},
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	paymentTxe, err := paymentTx.Sign(sourceSeed)
	fmt.Println("signe transaction:")
	if err != nil {
		log.Fatal(err)
	}

	paymentTxeB64, err := paymentTxe.Base64()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Submitting second transaction")
	res, err = horizonClient.SubmitTransaction(paymentTxeB64)
	fmt.Println(res)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Created custom asset and submitted transaction.")

}
