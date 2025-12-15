package main

import (
	"context"
	"fmt"

	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/client"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/metadata"
)

// This example demonstrates how to query certificate schemas and metadata.
// You can retrieve schema information, certificate data, and executor permissions.
//
// Usage:
//   go run 11_query_metadata.go
//
// What this script does:
// 1. Connects to the network and creates an account
// 2. Queries NFT schema information
// 3. Queries certificate metadata for a specific token
// 4. Checks executor permissions for the schema
// 5. Lists all executors for the schema
//
// Prerequisites:
// - Schema must be deployed (see 03_deploy_schema.go)
// - Metadata must be created (see 04_mint_metadata.go)

const (
	// IMPORTANT: Replace with your schema name
	schemaName = "myorg.lbbv01"

	// Token ID to query
	tokenId = "1"

	// For this example, we use the test mnemonic
	exampleMnemonic = account.TestMnemonic
)

func main() {
	fmt.Println("=== Step 11: Query Metadata and Schema ===")
	fmt.Println()

	// Validate configuration
	if schemaName == "" {
		fmt.Println("ERROR: Please update the schemaName constant")
		fmt.Println("Use the schema name from step 03_deploy_schema.go")
		return
	}

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

	// Step 2: Create metadata client for queries
	fmt.Println("Initializing metadata client...")
	metaClient := metadata.NewMetadata(*acc)
	fmt.Println("Metadata client initialized")
	fmt.Println()

	// Step 3: Query NFT Schema
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("Querying NFT Schema Information")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println()

	// GetNFTSchema retrieves the schema definition
	// Parameters:
	// - nftSchemaCode: the schema code (format: {ORGNAME}.{SCHEMACODE})
	// Returns:
	// - NFTSchemaQueryResult with schema details, origin data, and metadata
	schema, err := metaClient.GetNFTSchema(schemaName)
	if err != nil {
		fmt.Printf("ERROR: Failed to get schema: %v\n", err)
		fmt.Println()
		fmt.Println("This could mean:")
		fmt.Println("• Schema doesn't exist - run 03_deploy_schema.go first")
		fmt.Println("• Wrong schema name - check your schema code")
		fmt.Println("• Network connectivity issue")
		return
	}

	fmt.Printf("Schema Code: %s\n", schema.Code)
	fmt.Printf("Name: %s\n", schema.Name)
	fmt.Printf("Owner: %s\n", schema.Owner)
	fmt.Printf("Origin Contract Address: %s\n", schema.OriginData.OriginContractAddress)
	fmt.Printf("Origin Chain: %s\n", schema.OriginData.OriginChain)
	fmt.Printf("URI Retrieval Method: %s\n", schema.OriginData.UriRetrievalMethod)
	fmt.Printf("Metadata Format: %s\n", schema.OriginData.MetadataFormat)
	fmt.Println()

	// Step 4: Query Certificate Metadata
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("Querying Certificate Metadata")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println()

	// GetNFTMetadata retrieves certificate data for a specific token
	// Parameters:
	// - nftSchemaCode: the schema code
	// - tokenID: the token ID to query
	// Returns:
	// - NftData with token metadata, owner, and status
	nftData, err := metaClient.GetNFTMetadata(schemaName, tokenId)
	if err != nil {
		fmt.Printf("ERROR: Failed to get metadata: %v\n", err)
		fmt.Println()
		fmt.Println("This could mean:")
		fmt.Println("• Metadata doesn't exist for this token - run 04_mint_metadata.go first")
		fmt.Println("• Wrong token ID - check your token ID")
		fmt.Println("• Schema exists but no metadata created yet")
		return
	}

	fmt.Println()
	fmt.Printf("Metadata: %+v\n", nftData)
	fmt.Println()

	fmt.Println("What you learned:")
	fmt.Println("• How to query schema definitions and attributes")
	fmt.Println("• How to retrieve certificate metadata for a token")
	fmt.Println()

	fmt.Println("Use cases for these queries:")
	fmt.Println("• Verify certificate authenticity")
	fmt.Println("• Display certificate information in UI")
	fmt.Println("• Check permissions before operations")
	fmt.Println("• Audit schema configuration")
	fmt.Println("• Build certificate verification systems")
	fmt.Println()

	fmt.Println("Next steps:")
	fmt.Println("• Query multiple tokens to see different certificates")
	fmt.Println("• Build a frontend to display certificate data")
	fmt.Println("• Integrate queries into your application")
	fmt.Println("• Use these queries for verification workflows")
	fmt.Println()
}
