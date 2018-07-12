package main

import (
	"fmt"
	"os"

	"github.com/republicprotocol/atom-go/adapters/keystore"
	"github.com/republicprotocol/atom-go/utils"
)

func main() {

	ksPath := os.Getenv("GOPATH") + "/src/github.com/republicprotocol/atom-go/secrets/keystore.json"
	ks := keystore.NewKeystore(ksPath)

	key, err := ks.LoadKeypair("bitcoin")
	if err != nil {
		panic(err)
	}

	wallet := utils.NewWallet(key, "mainnet")
	addr, err := wallet.GetWIF()

	if err != nil {
		panic(err)
	}

	fmt.Println(addr)
}
