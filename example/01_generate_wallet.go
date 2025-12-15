package main

import (
	"fmt"

	"github.com/thesixnetwork/lbb-sdk-go/account"
)

// This example demonstrates how to generate a new wallet with mnemonic phrase.
// The mnemonic is a 24-word phrase that can be used to recover your wallet.
//
// Usage:
//   go run 01_generate_wallet.go
//
// What this script does:
// 1. Generates a new random mnemonic phrase (24 words)
// 2. Displays the mnemonic for backup
//
// IMPORTANT: Save the mnemonic phrase in a secure location!
// This is the ONLY way to recover your account.

func main() {
	fmt.Println("=== Step 1: Generate New Wallet ===")
	fmt.Println()

	// Generate a new mnemonic phrase
	// This creates a random 24-word phrase following BIP-39 standard
	mnemonic, err := account.GenerateMnemonic()
	if err != nil {
		panic(fmt.Sprintf("Failed to generate mnemonic: %v", err))
	}

	// Display the generated mnemonic
	fmt.Println("Mnemonic generated successfully!")
	fmt.Println()
	fmt.Println("IMPORTANT: Write this mnemonic phrase in a safe place.")
	fmt.Println("This is the ONLY way to recover your account if you ever forget your password.")
	fmt.Println()
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Printf("Mnemonic: %s\n", mnemonic)
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  • Save this mnemonic in a secure location")
	fmt.Println("  • Use this mnemonic in the next example to create an account")
	fmt.Println("  • NEVER share your mnemonic with anyone")
	fmt.Println()
}
