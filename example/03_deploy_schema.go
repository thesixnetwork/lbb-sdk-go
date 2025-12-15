package main

import (
	"context"
	"fmt"

	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/client"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/metadata"
)

// This example demonstrates how to deploy a certificate schema.
// A schema defines the structure and rules for your certificates/NFTs.
//
// Usage:
//   go run 03_deploy_schema.go
//
// What this script does:
// 1. Connects to the network and creates an account
// 2. Creates a metadata schema message
// 3. Deploys the schema to the blockchain
//
// Prerequisites:
// - Account must have tokens for transactionfees
// - Schema name must follow format: {ORGNAME}.{SchemaCode}

const (
	// Schema name format: {ORGNAME}.{SchemaCode}
	// Example: "mycompany.certificate01"
	schemaName = "myorg.lbbv01"

	// For this example, we use the test mnemonic
	exampleMnemonic = account.TestMnemonic
)

func main() {
	fmt.Println("=== Step 3: Deploy Certificate Schema ===")
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
	fmt.Println("Creating metadata schema...")

	// NewMetadataMsg creates a builder for metadata operations
	// Parameters:
	// - account: the account that will own the schema
	// - schemaName: unique identifier for this schema
	meta, err := metadata.NewMetadataMsg(*acc, schemaName)
	if err != nil {
		panic(fmt.Sprintf("Failed to create metadata message: %v", err))
	}

	// Step 3: Build deploy message
	fmt.Println("Building deployment message...")

	// BuildDeployMsg creates the message to deploy the schema on-chain
	msgDeploySchema, err := meta.BuildDeployMsg()
	if err != nil {
		panic(fmt.Sprintf("Failed to build deploy message: %v", err))
	}

	// Step 4: Broadcast transaction
	fmt.Println("Broadcasting schema deployment transaction...")

	// BroadcastTxAndWait sends the transaction and waits for confirmation
	res, err := meta.BroadcastTxAndWait(msgDeploySchema)
	if err != nil {
		panic(fmt.Sprintf("Failed to broadcast transaction: %v", err))
	}

	// Step 5: Display results
	fmt.Println("Schema deployed successfully!")
	fmt.Println()
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("Deployment Information:")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Printf("Schema Code:       %s\n", schemaName)
	fmt.Printf("Transaction Hash:  %s\n", res.TxHash)
	fmt.Printf("Deployer:          %s\n", acc.GetCosmosAddress().String())
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println()

	// Explanation
	fmt.Println("What just happened:")
	fmt.Println("  • Created a new certificate schema on the blockchain")
	fmt.Println("  • This schema defines the structure for your certificates")
	fmt.Println("  • The schema can now be used to mint metadata instances")
	fmt.Println("  • EVM contracts can reference this schema")
	fmt.Println()

	fmt.Println("Next steps:")
	fmt.Println("  • Mint metadata instances (04_mint_metadata.go)")
	fmt.Println("  • Deploy an EVM NFT contract linked to this schema")
	fmt.Printf("  • Schema name to use: %s\n", schemaName)
	fmt.Println()
}
