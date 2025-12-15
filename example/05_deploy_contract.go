package main

import (
	"context"
	"fmt"

	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/client"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/evm"
)

// This example demonstrates how to deploy an EVM NFT contract.
// The contract is linked to a previously deployed schema.
//
// Usage:
//   go run 04_deploy_contract.go
//
// What this script does:
// 1. Connects to the network and creates an account
// 2. Creates an EVM client
// 3. Deploys a Certificate NFT contract
// 4. Waits for deployment confirmation
//
// Prerequisites:
// - Schema must be deployed first (see 03_deploy_schema.go)
// - Metadata must be minted (see 04_mint_metadata.go)
// - Account must have tokens for gas fees
// - Schema name must match the deployed schema

const (
	// Contract configuration
	contractName   = "MyCertificate"
	contractSymbol = "CERT"

	// Schema name - MUST match the schema deployed in step 03
	// Format: {ORGNAME}.{SchemaCode}
	schemaName = "myorg.lbbv01"

	// For this example, we use the test mnemonic
	exampleMnemonic = account.TestMnemonic
)

func main() {
	fmt.Println("=== Step 5: Deploy EVM NFT Contract ===")
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

	fmt.Printf("Connected with account: %s\n", acc.GetEVMAddress().Hex())
	fmt.Println()

	// Step 2: Create EVM client
	fmt.Println("Initializing EVM client...")

	// NewEVMClient creates a client for EVM operations
	// Parameter: account with signing capabilities
	evmClient := evm.NewEVMClient(*acc)

	fmt.Println("EVM client initialized")
	fmt.Println()

	// Step 3: Deploy Certificate Contract
	fmt.Println("Deploying Certificate NFT contract...")
	fmt.Printf("   Contract Name: %s\n", contractName)
	fmt.Printf("   Symbol: %s\n", contractSymbol)
	fmt.Printf("   Linked Schema: %s\n", schemaName)
	fmt.Println()

	// DeployCertificateContract deploys the NFT contract
	// Parameters:
	// - contractName: human-readable name for the NFT collection
	// - contractSymbol: short symbol (e.g., "CERT", "NFT")
	// - schemaName: the schema code this contract is linked to
	// Returns:
	// - contractAddress: deployed contract address
	// - tx: transaction object
	// - err: error if any
	contractAddress, tx, err := evmClient.DeployCertificateContract(contractName, contractSymbol, schemaName)
	if err != nil {
		panic(fmt.Sprintf("Failed to deploy contract: %v", err))
	}

	fmt.Printf("Transaction submitted: %s\n", tx.Hash().Hex())
	fmt.Printf("   Nonce: %d\n", tx.Nonce())
	fmt.Println()

	// Step 4: Wait for deployment confirmation
	fmt.Println("Waiting for transaction to be mined...")

	// WaitForEVMTransaction waits until the transaction is confirmed
	_, err = client.WaitForEVMTransaction(tx.Hash())
	if err != nil {
		panic(fmt.Sprintf("Error waiting for deployment: %v", err))
	}

	// Step 5: Display results
	fmt.Println("Contract deployed successfully!")
	fmt.Println()
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("Contract Information:")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Printf("Contract Address:  %s\n", contractAddress.Hex())
	fmt.Printf("Contract Name:     %s\n", contractName)
	fmt.Printf("Symbol:            %s\n", contractSymbol)
	fmt.Printf("Linked Schema:     %s\n", schemaName)
	fmt.Printf("Transaction Hash:  %s\n", tx.Hash().Hex())
	fmt.Printf("Deployer (EVM):    %s\n", acc.GetEVMAddress().Hex())
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println()

	// Explanation
	fmt.Println("What just happened:")
	fmt.Println("  • Deployed an ERC-721 compatible NFT contract")
	fmt.Println("  • The contract is linked to your metadata schema")
	fmt.Println("  • You can now mint NFTs using this contract")
	fmt.Println("  • Each NFT will reference metadata from the schema")
	fmt.Println()

	fmt.Println("Next steps:")
	fmt.Println("  • Mint NFTs using this contract (06_mint_nft.go)")
	fmt.Println("  • Save the contract address for future operations")
	fmt.Printf("  • Contract address: %s\n", contractAddress.Hex())
	fmt.Println()
}
