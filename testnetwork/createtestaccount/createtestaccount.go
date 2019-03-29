package createtestaccount

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/stellar/go/keypair"
	"github.com/thanakritlee/stellar-playground/genseed"
)

func main() {
	pair := genseed.GenPair()

	address := pair.Address()
	resp, err := http.Get("https://friendbot.stellar.org/?addr=" + address)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))

}

// CreateAccount creates an account with the Stellar test network.
func CreateAccount() *keypair.Full {
	pair := genseed.GenPair()

	address := pair.Address()
	_, err := http.Get("https://friendbot.stellar.org/?addr=" + address)
	if err != nil {
		log.Fatal(err)
	}

	return pair
}
