package main

import (
	"context"
	"fmt"

	_ "cosmossdk.io/math"

	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/client"
)

func main() {
	mnemonic, err := account.GenerateMnemonic()
	if err != nil {
		panic(fmt.Sprintf("Failed to generate mnemonic: %v", err))
	}

	fmt.Println("Mnemonic generated")
	fmt.Println("*Important** write this mnemonic phrase in a safe place.")
	fmt.Println("It is the only way to recover your account if you ever forget your password.")
	fmt.Printf("\nMnemonic: %s\n\n", mnemonic)
	fmt.Println("-----------------------------------------------------")

	// Step 2: Initialize client (fivenet = testnet)
	fmt.Println("Step 2: Connecting to network...")
	ctx := context.Background()
	client, err := client.NewClient(ctx, false)
	if err != nil {
		panic(fmt.Sprintf("Failed to create client: %v", err))
	}
	fmt.Println("Connected to fivenet (testnet)")
	fmt.Println()

	fmt.Println("Step 3: Creating account...")
	acc, err := account.NewAccount(client, "my-account", mnemonic, "password")
	if err != nil {
		panic("ERROR CREATE ACCOUNT: NewAccount returned nil - check mnemonic and keyring initialization")
	}

	fmt.Printf("Account created\n")
	fmt.Printf("  EVM Address: %s\n", acc.GetEVMAddress().Hex())
	fmt.Printf("  Cosmos Address: %s\n\n", acc.GetCosmosAddress().String())
}
