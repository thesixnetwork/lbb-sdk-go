package main

import (
	"context"
	"fmt"

	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/client"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/metadata"
)

// This example demonstrates how to freeze and unfreeze metadata.
// Freezing prevents modifications to the certificate metadata,
// while unfreezing allows modifications again.
//
// Usage:
//   go run 08_freeze_metadata.go
//
// What this script does:
// 1. Connects to the network and creates an account
// 2. Creates a metadata client
// 3. Freezes a certificate metadata
// 4. Unfreezes the certificate metadata
//
// Prerequisites:
// - Schema must be deployed (see 03_deploy_schema.go)
// - Metadata must exist (created during schema deployment)
// - Account must be the owner/creator of the metadata
// - Account must have tokens for transaction fees

const (
	// Schema name - MUST match the schema deployed in step 03
	// Format: {ORGNAME}.{SchemaCode}
	schemaName = "myorg.lbbv01"

	// Token ID to freeze/unfreeze
	// This should match the token ID created in step 03
	tokenId = "1"

	// For this example, we use the test mnemonic
	exampleMnemonic = account.TestMnemonic
)

func main() {
	fmt.Println("=== Step 8: Freeze and Unfreeze Metadata ===")
	fmt.Println()

	// Step 1: Setup client and account
	fmt.Println("Setting up connection...")
	ctx := context.Background()
	client, err := client.NewClient(ctx, false)
	if err != nil {
		panic(fmt.Sprintf("Failed to create client: %v", err))
	}

	acc, err := account.NewAccount(client, "myaccount", exampleMnemonic, "mypassword")
	if err != nil {
		panic(fmt.Sprintf("Failed to create account: %v", err))
	}

	fmt.Printf("Connected with account: %s\n", acc.GetCosmosAddress().String())
	fmt.Println()

	// Step 2: Create metadata client
	fmt.Println("Initializing metadata client...")

	// NewMetadataMsg creates a builder for metadata operations
	meta, err := metadata.NewMetadataMsg(*acc, schemaName)
	if err != nil {
		panic(fmt.Sprintf("Failed to create metadata client: %v", err))
	}

	fmt.Println("Metadata client initialized")
	fmt.Println()

	// Step 3: Freeze the certificate
	fmt.Println("Freezing certificate metadata...")
	fmt.Printf("   Schema: %s\n", schemaName)
	fmt.Printf("   Token ID: %s\n", tokenId)
	fmt.Println()

	// FreezeCertificate locks the metadata, preventing modifications
	// Parameter: token ID of the metadata to freeze
	// Returns:
	// - res: transaction response
	// - err: error if any
	freezeRes, err := meta.FreezeCertificate(tokenId)
	if err != nil {
		panic(fmt.Sprintf("Failed to freeze certificate: %v", err))
	}

	fmt.Printf("Freeze transaction submitted: %s\n", freezeRes.TxHash)
	fmt.Println()

	// Wait for freeze transaction to be confirmed
	fmt.Println("Waiting for freeze transaction confirmation...")
	err = client.WaitForTransaction(freezeRes.TxHash)
	if err != nil {
		panic(fmt.Sprintf("Error waiting for freeze transaction: %v", err))
	}

	fmt.Println("Certificate metadata frozen successfully!")
	fmt.Println()
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("Freeze Information:")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Printf("Schema Code:       %s\n", schemaName)
	fmt.Printf("Token ID:          %s\n", tokenId)
	fmt.Printf("Transaction Hash:  %s\n", freezeRes.TxHash)
	fmt.Printf("Status:            FROZEN\n")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println()

	fmt.Println("Freeze effect:")
	fmt.Println("  • Metadata cannot be modified while frozen")
	fmt.Println("  • This ensures data integrity and immutability")
	fmt.Println("  • Useful for finalized certificates or credentials")
	fmt.Println()

	// Step 4: Unfreeze the certificate
	fmt.Println("Unfreezing certificate metadata...")
	fmt.Printf("   Schema: %s\n", schemaName)
	fmt.Printf("   Token ID: %s\n", tokenId)
	fmt.Println()

	// UnfreezeCertificate unlocks the metadata, allowing modifications
	// Parameter: token ID of the metadata to unfreeze
	// Returns:
	// - res: transaction response
	// - err: error if any
	unfreezeRes, err := meta.UnfreezeCertificate(tokenId)
	if err != nil {
		panic(fmt.Sprintf("Failed to unfreeze certificate: %v", err))
	}

	fmt.Printf("Unfreeze transaction submitted: %s\n", unfreezeRes.TxHash)
	fmt.Println()

	// Wait for unfreeze transaction to be confirmed
	fmt.Println("Waiting for unfreeze transaction confirmation...")
	err = client.WaitForTransaction(unfreezeRes.TxHash)
	if err != nil {
		panic(fmt.Sprintf("Error waiting for unfreeze transaction: %v", err))
	}

	fmt.Println("Certificate metadata unfrozen successfully!")
	fmt.Println()
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("Unfreeze Information:")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Printf("Schema Code:       %s\n", schemaName)
	fmt.Printf("Token ID:          %s\n", tokenId)
	fmt.Printf("Transaction Hash:  %s\n", unfreezeRes.TxHash)
	fmt.Printf("Status:            UNFROZEN\n")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println()

	fmt.Println("Unfreeze effect:")
	fmt.Println("  • Metadata can now be modified again")
	fmt.Println("  • Allows updates to certificate data if needed")
	fmt.Println("  • Owner can freeze it again at any time")
	fmt.Println()

	// Summary
	fmt.Println("=== Summary ===")
	fmt.Printf("Schema Code:        %s\n", schemaName)
	fmt.Printf("Token ID:           %s\n", tokenId)
	fmt.Printf("Freeze Tx:          %s\n", freezeRes.TxHash)
	fmt.Printf("Unfreeze Tx:        %s\n", unfreezeRes.TxHash)
	fmt.Printf("Final Status:       UNFROZEN\n")
	fmt.Println()

	fmt.Println("Use cases for freeze/unfreeze:")
	fmt.Println("  • Freeze: Lock finalized academic certificates")
	fmt.Println("  • Freeze: Prevent tampering with issued credentials")
	fmt.Println("  • Unfreeze: Allow corrections to metadata")
	fmt.Println("  • Unfreeze: Update certificate information if needed")
	fmt.Println()

	fmt.Println("Next steps:")
	fmt.Println("  • Query NFT and metadata information (09_query_nft.go)")
	fmt.Println("  • Implement freeze logic in your application workflow")
	fmt.Println()
}
