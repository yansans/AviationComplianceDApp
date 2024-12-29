package main

import (
	"aviation-compliance-dapp-wallet/wallet"
	"fmt"
	"log"
)

func main() {
	certPath := "./key/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem"
	keyPath := "./key/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/keystore/priv_sk"
	msp := "Org1MSP"

	identity, err := wallet.LoadIdentityFromFiles(msp, certPath, keyPath)
	if err != nil {
		log.Fatalf("Error loading identity: %v", err)
	}

	store := &wallet.FileWalletStore{}

	wallet, err := wallet.NewWallet(identity, store)
	if err != nil {
		log.Fatalf("Error creating wallet: %v", err)
	}

	labels, err := wallet.List()
	if err != nil {
		log.Fatalf("Error listing identities: %v", err)
	}

	fmt.Println("Identities in the wallet:")
	for _, label := range labels {
		fmt.Println(label)
	}

	retrievedIdentity, err := wallet.Get("user_identity")
	if err != nil {
		log.Fatalf("Error retrieving identity: %v", err)
	}

	fmt.Println("Retrieved identity:")
	fmt.Println(retrievedIdentity)
}
