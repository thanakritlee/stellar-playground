package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loding .env file")
	}

	sourceSeed := os.Getenv("STELLAR_SOURCE_SEED")
	destinsationSeed := os.Getenv("STELLAR_DESTINATION_SEED")

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
			SequenceProvider: horizon.DefaultTestNetClient,
		},
		build.TestNetwork,
		build.Trust(
			corruptionBaht.Code,
			corruptionBaht.Issuer,
			build.Limit("100.25"),
		),
	)
	if err != nil {
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

	_, err = horizon.DefaultTestNetClient.SubmitTransaction(trustTxeB64)
	if err != nil {
		log.Fatal(err)
	}

	// Second, the issuing account actually sends a payment usng the asset.
	paymentTx, err := build.Transaction(
		build.SourceAccount{
			AddressOrSeed: issuer.Address(),
		},
		build.TestNetwork,
		build.AutoSequence{
			SequenceProvider: horizon.DefaultTestNetClient,
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
	if err != nil {
		log.Fatal(err)
	}

	paymentTxeB64, err := paymentTxe.Base64()
	if err != nil {
		log.Fatal(err)
	}

	_, err = horizon.DefaultTestNetClient.SubmitTransaction(paymentTxeB64)
	if err != nil {
		log.Fatal(err)
	}

}
