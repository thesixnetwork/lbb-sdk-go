package main

import (
	"context"
	"fmt"

	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/client"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/metadata"
)

// This example demonstrates how to mint certificate metadata.
// Metadata instances are created after a schema has been deployed.
//
// Usage:
//   go run 04_mint_metadata.go
//
// What this script does:
// 1. Connects to the network and creates an account
// 2. Creates a metadata message builder
// 3. Mints a new metadata instance with a specific token ID
// 4. Waits for transaction confirmation
//
// Prerequisites:
// - Schema must be deployed first (see 03_deploy_schema.go)
// - Account must have tokens for transaction fees
// - Schema name must match the deployed schema

const (
	// Schema name - MUST match the schema deployed in step 03
	// Format: {ORGNAME}.{SchemaCode}
	schemaName = "myorg.lbbv01"

	// Token ID to mint
	// Each metadata instance needs a unique token ID
	tokenId = "1"

	// For this example, we use the test mnemonic
	exampleMnemonic = account.TestMnemonic
)

func main() {
	fmt.Println("=== Step 4: Mint Certificate Metadata ===")
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

	// Step 2: Create metadata message builder
	fmt.Println("Initializing metadata client...")

	// NewMetadataMsg creates a builder for metadata operations
	// Parameters:
	// - account: the account that will mint the metadata
	// - schemaName: the schema code (must already be deployed)
	meta, err := metadata.NewMetadataMsg(*acc, schemaName)
	if err != nil {
		panic(fmt.Sprintf("Failed to create metadata message: %v", err))
	}

	fmt.Println("Metadata client initialized")
	fmt.Println()

	// Step 3: Build mint metadata message
	fmt.Println("Building mint metadata message...")
	fmt.Printf("   Schema: %s\n", schemaName)
	fmt.Printf("   Token ID: %s\n", tokenId)
	fmt.Println()

	// BuildMintMetadataMsg creates a new metadata instance
	// Parameter: token ID for this metadata instance (must be unique)
	msgMintMetadata, err := meta.BuildMintMetadataWithInfoMsg(tokenId, metadata.CertificateInfo{
		Status:       "TCI",
		GoldStandard: "LBI",
		Weight:       "2000g",
		CertNumber:   "LBB_V1_01",
		CustomerID:   "LBB_V1_USER_01",
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to build mint message: %v", err))
	}

	// Step 4: Broadcast transaction
	fmt.Println("Broadcasting mint transaction...")

	// BroadcastTxAndWait sends the transaction and waits for confirmation
	res, err := meta.BroadcastTxAndWait(msgMintMetadata)
	if err != nil {
		panic(fmt.Sprintf("Failed to broadcast transaction: %v", err))
	}

	// Step 5: Display results
	fmt.Println("Metadata minted successfully!")
	fmt.Println()
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("Mint Information:")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Printf("Schema Code:       %s\n", schemaName)
	fmt.Printf("Token ID:          %s\n", tokenId)
	fmt.Printf("Transaction Hash:  %s\n", res.TxHash)
	fmt.Printf("Minter:            %s\n", acc.GetCosmosAddress().String())
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println()

	// Explanation
	fmt.Println("What just happened:")
	fmt.Println("  • Created a new metadata instance on the blockchain")
	fmt.Println("  • This metadata is linked to the deployed schema")
	fmt.Println("  • The metadata can be referenced by NFTs on the EVM layer")
	fmt.Println()

	fmt.Println("Next steps:")
	fmt.Println("  • Deploy an EVM NFT contract (05_deploy_contract.go)")
	fmt.Println("  • Link the contract to this schema")
	fmt.Println("  • Mint NFTs that reference this metadata")
	fmt.Printf("  • Token ID to reference: %s\n", tokenId)
	fmt.Println()
}
