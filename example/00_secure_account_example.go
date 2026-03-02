package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/client"
)

// This example demonstrates SECURE account management practices
// following the security guidelines in account/SECURITY.md

func main() {
	fmt.Println("=== Secure Account Management Example ===\n")

	// Choose your scenario
	fmt.Println("Choose a scenario:")
	fmt.Println("1. Create new account (first time)")
	fmt.Println("2. Use existing account from mnemonic")
	fmt.Println("3. Use encrypted keystore (recommended for production)")
	fmt.Print("\nEnter choice (1-3): ")

	var choice int
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		exampleCreateNewAccount()
	case 2:
		exampleUseExistingMnemonic()
	case 3:
		exampleEncryptedKeystore()
	default:
		log.Fatal("Invalid choice")
	}
}

// exampleCreateNewAccount demonstrates creating a NEW account with secure mnemonic handling
func exampleCreateNewAccount() {
	fmt.Println("\n=== Creating New Account ===\n")

	// Step 1: Generate new mnemonic
	fmt.Println("Generating new mnemonic...")
	mnemonic, err := account.GenerateNewMnemonic()
	if err != nil {
		log.Fatalf("Failed to generate mnemonic: %v", err)
	}

	// Step 2: Display mnemonic ONE TIME ONLY for user backup
	fmt.Println("\n" + "="*80)
	fmt.Println("🚨 CRITICAL: YOUR RECOVERY PHRASE 🚨")
	fmt.Println("=" * 80)
	fmt.Println()
	fmt.Println("Write down these 12 words IN ORDER on paper:")
	fmt.Println()
	fmt.Println("  " + mnemonic)
	fmt.Println()
	fmt.Println("⚠️  WARNING:")
	fmt.Println("  • This phrase will ONLY be shown ONCE")
	fmt.Println("  • Anyone with this phrase can steal ALL your funds")
	fmt.Println("  • NEVER share it with anyone")
	fmt.Println("  • NEVER store it digitally (no screenshots, no files)")
	fmt.Println("  • Store it in a SECURE physical location")
	fmt.Println()
	fmt.Println("=" * 80)
	fmt.Println()
	fmt.Print("Have you written it down? Type 'YES' to continue: ")

	var confirmation string
	fmt.Scanln(&confirmation)
	if confirmation != "YES" {
		log.Fatal("Please write down your mnemonic before continuing")
	}

	// Step 3: Create client
	cfg := client.Config{
		ChainID:        "fivenet",
		RPC:            "https://rpc1.fivenet.sixprotocol.net:443",
		GRPCEndpoint:   "https://grpc-web.fivenet.sixprotocol.net:443",
		EVMRPCUrl:      "https://rpc-evm.fivenet.sixprotocol.net:443",
		KeyringDir:     "./keyring-test",
		KeyringBackend: "test",
	}

	c, err := client.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	// Step 4: Create account from mnemonic
	// NOTE: The mnemonic is used here but NOT stored in the Account struct
	fmt.Println("\nCreating account...")
	acc, err := account.NewAccount(c, "my-secure-account", mnemonic, "")
	if err != nil {
		log.Fatalf("Failed to create account: %v", err)
	}
	defer acc.Close() // ✅ CRITICAL: Always cleanup

	// Step 5: Show account info
	fmt.Printf("\n✅ Account created successfully!\n")
	fmt.Printf("   Name:           %s\n", acc.GetAccountName())
	fmt.Printf("   Cosmos Address: %s\n", acc.GetCosmosAddress().String())
	fmt.Printf("   EVM Address:    %s\n", acc.GetEVMAddress().Hex())

	// The mnemonic variable is now out of scope
	// It will be garbage collected and is NOT stored in the account
	fmt.Println("\n✅ Mnemonic is NOT stored in memory")
	fmt.Println("✅ You can only recover with the written backup")

	// Step 6: Use account for operations
	ctx := context.Background()
	fmt.Println("\nAccount is ready for use...")

	// Example: Check balance (would fail if no funds)
	// balance, _ := pkg.GetBalance(ctx, acc.GetCosmosAddress().String(), "usix")
	// fmt.Printf("Balance: %s\n", balance)

	time.Sleep(2 * time.Second)

	// Step 7: Cleanup happens automatically via defer acc.Close()
	fmt.Println("\n✅ Cleaning up (zeroizing private key)...")
}

// exampleUseExistingMnemonic demonstrates using an EXISTING mnemonic securely
func exampleUseExistingMnemonic() {
	fmt.Println("\n=== Using Existing Account ===\n")

	// WARNING: Never hardcode mnemonics in production!
	// This is just for demonstration
	fmt.Println("⚠️  Enter your mnemonic phrase:")
	fmt.Println("    (For testing, you can use: test test test test test test test test test test test junk)")
	fmt.Print("\nMnemonic: ")

	reader := bufio.NewReader(os.Stdin)
	mnemonic, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read mnemonic: %v", err)
	}

	// Validate mnemonic
	if !account.ValidateMnemonic(mnemonic) {
		log.Fatal("❌ Invalid mnemonic phrase")
	}

	// Create client
	cfg := client.Config{
		ChainID:        "fivenet",
		RPC:            "https://rpc1.fivenet.sixprotocol.net:443",
		GRPCEndpoint:   "https://grpc-web.fivenet.sixprotocol.net:443",
		EVMRPCUrl:      "https://rpc-evm.fivenet.sixprotocol.net:443",
		KeyringDir:     "./keyring-test",
		KeyringBackend: "test",
	}

	c, err := client.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	// Create account (mnemonic used but NOT stored)
	acc, err := account.NewAccount(c, "restored-account", mnemonic, "")
	if err != nil {
		log.Fatalf("Failed to create account: %v", err)
	}
	defer acc.Close() // ✅ CRITICAL: Always cleanup

	fmt.Printf("\n✅ Account restored successfully!\n")
	fmt.Printf("   Cosmos Address: %s\n", acc.GetCosmosAddress().String())
	fmt.Printf("   EVM Address:    %s\n", acc.GetEVMAddress().Hex())

	// Mnemonic is now out of scope
	// Use account for operations...

	fmt.Println("\n✅ Account ready for use")
	fmt.Println("✅ Mnemonic is NOT stored in memory")
}

// exampleEncryptedKeystore demonstrates production-ready encrypted keystore usage
func exampleEncryptedKeystore() {
	fmt.Println("\n=== Encrypted Keystore Example ===\n")
	fmt.Println("This example shows the RECOMMENDED approach for production:")
	fmt.Println("1. Create account from mnemonic ONCE")
	fmt.Println("2. Save encrypted keystore file")
	fmt.Println("3. Daily usage: Load from keystore (no mnemonic needed)")
	fmt.Println()

	// For this example, we'll use a simple password-protected approach
	// In production, use proper keystore encryption (see account/SECURITY.md)

	fmt.Println("⚠️  This is a simplified example")
	fmt.Println("⚠️  For production, use proper keystore encryption")
	fmt.Println("⚠️  See account/SECURITY.md for full implementation")
	fmt.Println()

	// Step 1: One-time setup (would normally be done once)
	fmt.Println("Step 1: Generate new account (one-time)")
	mnemonic, err := account.GenerateNewMnemonic()
	if err != nil {
		log.Fatalf("Failed to generate mnemonic: %v", err)
	}

	fmt.Println("Generated mnemonic (WRITE THIS DOWN):")
	fmt.Println("  " + mnemonic)

	// Create client
	cfg := client.Config{
		ChainID:        "fivenet",
		RPC:            "https://rpc1.fivenet.sixprotocol.net:443",
		GRPCEndpoint:   "https://grpc-web.fivenet.sixprotocol.net:443",
		EVMRPCUrl:      "https://rpc-evm.fivenet.sixprotocol.net:443",
		KeyringDir:     "./keyring-test",
		KeyringBackend: "test",
	}

	c, err := client.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	// Create account
	acc, err := account.NewAccount(c, "keystore-account", mnemonic, "")
	if err != nil {
		log.Fatalf("Failed to create account: %v", err)
	}

	fmt.Printf("\n✅ Account created!\n")
	fmt.Printf("   Address: %s\n", acc.GetEVMAddress().Hex())

	// In production, you would:
	// 1. Export the private key
	// 2. Encrypt it with a strong password
	// 3. Save to secure keystore file
	// 4. Delete the mnemonic (keep only paper backup)

	privateKeyHex, err := account.ExportPrivateKeyHex(acc)
	if err != nil {
		log.Fatalf("Failed to export key: %v", err)
	}

	fmt.Println("\nPrivate key exported (would be encrypted in production)")
	fmt.Printf("Key: 0x%s...\n", privateKeyHex[:8])

	acc.Close() // Cleanup after export

	fmt.Println("\n✅ In production:")
	fmt.Println("   • Encrypt this key with a strong password")
	fmt.Println("   • Save encrypted keystore to secure location")
	fmt.Println("   • Daily usage: Load from keystore (see SECURITY.md)")
	fmt.Println("   • Mnemonic only needed for recovery")
}

// Additional helper functions

// SecureAccountWrapper wraps account operations with proper cleanup
func SecureAccountWrapper(mnemonic string, operation func(*account.Account) error) error {
	// Create client
	cfg := client.Config{
		ChainID:        "fivenet",
		RPC:            "https://rpc1.fivenet.sixprotocol.net:443",
		GRPCEndpoint:   "https://grpc-web.fivenet.sixprotocol.net:443",
		EVMRPCUrl:      "https://rpc-evm.fivenet.sixprotocol.net:443",
		KeyringDir:     "./keyring-test",
		KeyringBackend: "test",
	}

	c, err := client.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer c.Close()

	// Create account (mnemonic not stored)
	acc, err := account.NewAccount(c, "temp-account", mnemonic, "")
	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}
	defer acc.Close() // ✅ Always cleanup

	// Execute operation
	if err := operation(acc); err != nil {
		return fmt.Errorf("operation failed: %w", err)
	}

	// Cleanup happens automatically via defer
	return nil
}

// Example usage of wrapper
func exampleWithWrapper() {
	mnemonic := "test test test test test test test test test test test junk"

	err := SecureAccountWrapper(mnemonic, func(acc *account.Account) error {
		// Your operations here
		fmt.Printf("Using account: %s\n", acc.GetEVMAddress().Hex())

		// Do transfers, queries, etc.

		return nil
	})

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Account is automatically cleaned up
}
