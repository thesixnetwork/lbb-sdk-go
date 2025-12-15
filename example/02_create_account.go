package main

import (
	"context"
	"fmt"

	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/client"
)

// This example demonstrates how to connect to the network and create an account.
// An account allows you to interact with both Cosmos and EVM features.
//
// Usage:
//   go run 02_create_account.go
//
// What this script does:
// 1. Connects to the fivenet testnet
// 2. Creates an account from a mnemonic phrase
// 3. Displays both Cosmos and EVM addresses
//
// Prerequisites:
// - You need a mnemonic phrase (generate one using 01_generate_wallet.go)

const (
	// For this example, we use the test mnemonic
	// In production, replace this with your own mnemonic from step 01
	exampleMnemonic = account.TestMnemonic

	// Account name (can be any identifier you prefer)
	accountName = "my-account"

	// Password for the keyring (optional, can be empty string)
	accountPassword = ""
)

func main() {
	fmt.Println("=== Step 2: Create Account ===")
	fmt.Println()

	// Step 1: Initialize client connection
	fmt.Println("Connecting to network...")
	ctx := context.Background()

	// NewClient parameters:
	// - ctx: context for the connection
	// - isMainnet: false = fivenet (testnet), true = sixnet (mainnet)
	client, err := client.NewClient(ctx, false)
	if err != nil {
		panic(fmt.Sprintf("Failed to create client: %v", err))
	}
	fmt.Println("Connected to fivenet (testnet)")
	fmt.Println()

	// Step 2: Create account from mnemonic
	fmt.Println("Creating account from mnemonic...")

	// NewAccount parameters:
	// - client: the network client
	// - name: account identifier in keyring
	// - mnemonic: your 24-word phrase
	// - password: optional password for keyring security
	acc, err := account.NewAccount(client, accountName, exampleMnemonic, accountPassword)
	if err != nil {
		panic(fmt.Sprintf("Failed to create account: %v", err))
	}

	fmt.Println("Account created successfully!")
	fmt.Println()

	// Step 3: Display account information
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("Account Information:")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Printf("Account Name:    %s\n", accountName)
	fmt.Printf("EVM Address:     %s\n", acc.GetEVMAddress().Hex())
	fmt.Printf("Cosmos Address:  %s\n", acc.GetCosmosAddress().String())
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println()

	// Explanation of addresses
	fmt.Println("Understanding Addresses:")
	fmt.Println("  • EVM Address: Used for Ethereum-compatible interactions (0x...)")
	fmt.Println("  • Cosmos Address: Used for Cosmos SDK operations (6x...)")
	fmt.Println("  • Both addresses are derived from the same private key")
	fmt.Println()

	fmt.Println("Next steps:")
	fmt.Println("  • Use this account to deploy schemas (03_deploy_schema.go)")
	fmt.Println("  • Make sure you have tokens in your account for transactions")
	fmt.Println()
}
