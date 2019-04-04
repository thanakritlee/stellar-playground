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
	"github.com/thanakritlee/stellar-playground/testnetwork/genseed"
)

func main() {

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loding .env file")
	}

	rootSeed := os.Getenv("ROOT_SEED")
	url := os.Getenv("URL")
	networkPassPhrase := os.Getenv("NETWORK_PASSPHRASE")

	horizonClient := &horizon.Client{
		URL:  url,
		HTTP: http.DefaultClient,
	}

	rootPair, err := keypair.Parse(rootSeed)
	if err != nil {
		fmt.Println("Can't parse root seed")
		fmt.Println(err)
	}

	fmt.Println("Root Address:")
	fmt.Println(rootPair.Address())

	// Generate key pair of first Stellar account to create on the private network.
	pair := genseed.GenPair()
	address := pair.Address()
	seed := pair.Seed()

	fmt.Println("New Account Seed:")
	fmt.Printf("%s\n", seed)
	fmt.Println("New Account Address:")
	fmt.Printf("%s\n", address)

	tx, err := build.Transaction(
		build.SourceAccount{AddressOrSeed: rootSeed},
		// build.Sequence{Sequence: 10101},
		build.AutoSequence{
			SequenceProvider: horizonClient,
		},
		build.Network{Passphrase: networkPassPhrase},
		build.CreateAccount(
			build.Destination{AddressOrSeed: address},
			build.NativeAmount{Amount: "100"},
		),
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	txe, err := tx.Sign(rootSeed)
	if err != nil {
		fmt.Println(err)
		return
	}

	txeB64, err := txe.Base64()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("tx base64: %s\n", txeB64)

	res, err := horizonClient.SubmitTransaction(txeB64)
	if err != nil {
		fmt.Println("Error from submitting transaction")
		log.Fatal(err)
	}

	fmt.Println("tx submitted response:")
	fmt.Println(res)
}
