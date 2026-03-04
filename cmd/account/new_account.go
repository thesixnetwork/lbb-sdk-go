package main

import (
	"context"
	"fmt"

	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/client"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/logger"
)

func main() {
	mnemonic, err := account.GenerateMnemonic()
	if err != nil {
		logger.Fatal("Failed to generate mnemonic: %v", err)
	}

	logger.Info("Mnemonic generated")
	logger.Info("*Important** write this mnemonic phrase in a safe place.")
	logger.Info("It is the only way to recover your account if you ever forget your password.")
	fmt.Printf("\nMnemonic: %s\n\n", mnemonic)
	logger.Info("-----------------------------------------------------")

	// Step 2: Initialize client (fivenet = testnet)
	logger.Info("Step 2: Connecting to network...")
	ctx := context.Background()
	client, err := client.NewClient(ctx, false)
	if err != nil {
		logger.Fatal("Failed to create client: %v", err)
	}
	logger.Info("Connected to fivenet (testnet)")
	fmt.Println()

	logger.Info("Step 3: Creating account...")
	acc, err := account.NewAccount(client, "my-account", mnemonic, "password")
	if err != nil {
		logger.Fatal("Failed to create account: %v", err)
	}

	logger.Info("Account created")
	logger.Info("  EVM Address: %s", acc.GetEVMAddress().Hex())
	logger.Info("  Cosmos Address: %s", acc.GetCosmosAddress().String())
	fmt.Println()
}
