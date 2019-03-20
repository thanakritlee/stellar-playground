package main

import (
	"fmt"

	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
)

func main() {
	source := "GDFQBJPSBH3VMH6TZVNWURFB636O4OSK6SCKQBMQX5YZYTZ5VWXW2DYU"
	sourceSecretSeed := "SCX5EA6BKJ3R64CNSTBXCQ36NSNG2ERKSHU47YKDR7DKGAEO36YNROTH"
	destination := "GAXGIUCKOTCNF2RHHVOLK23I2U6YDP3INTW6EB6GRXASWVYDO7RCSEI3"

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
