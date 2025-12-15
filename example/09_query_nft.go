package main

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/client"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/evm"
)

// This example demonstrates how to query NFT information.
// You can check ownership, token details, and verify NFT status.
//
// Usage:
//   go run 09_query_nft.go
//
// What this script does:
// 1. Connects to the network and creates an account
// 2. Creates an EVM client
// 3. Queries NFT ownership
// 4. Displays token information
//
// Prerequisites:
// - EVM contract must be deployed (see 05_deploy_contract.go)
// - NFT must be minted (see 06_mint_nft.go)
// - Contract address must be valid

const (
	// IMPORTANT: Replace with your deployed contract address
	contractAddress = "0x0000000000000000000000000000000000000000"

	// Token ID to query
	tokenId = uint64(1)

	// For this example, we use the test mnemonic
	exampleMnemonic = account.TestMnemonic
)

func main() {
	fmt.Println("=== Step 9: Query NFT Information ===")
	fmt.Println()

	// Validate configuration
	if contractAddress == "0x0000000000000000000000000000000000000000" {
		fmt.Println("ERROR: Please update the contractAddress constant")
		fmt.Println("Use the contract address from step 05_deploy_contract.go")
		return
	}

	// Step 1: Setup client and account
	fmt.Println("Setting up connection...")
	ctx := context.Background()
	client, err := client.NewClient(ctx, false)
	if err != nil {
		panic(fmt.Sprintf("Failed to create client: %v", err))
	}

	acc, err := account.NewAccount(client, "myaccout", exampleMnemonic, "mypassword")
	if err != nil {
		panic(fmt.Sprintf("Failed to create account: %v", err))
	}

	fmt.Printf("Connected with account: %s\n", acc.GetEVMAddress().Hex())
	fmt.Println()

	// Step 2: Create EVM client
	fmt.Println("Initializing EVM client...")
	evmClient := evm.NewEVMClient(*acc)
	fmt.Println("EVM client initialized")
	fmt.Println()

	// Step 3: Query NFT ownership
	contractAddr := common.HexToAddress(contractAddress)

	fmt.Println("Querying NFT information...")
	fmt.Printf("   Contract: %s\n", contractAddress)
	fmt.Printf("   Token ID: %d\n", tokenId)
	fmt.Println()

	// TokenOwner retrieves the current owner of a token
	// Parameters:
	// - contractAddress: the deployed contract address
	// - tokenId: the NFT token ID to query
	// Returns:
	// - owner: address of the current owner
	owner := evmClient.TokenOwner(contractAddr, tokenId)

	// Step 4: Display results
	fmt.Println("Query completed!")
	fmt.Println()
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("NFT Information:")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Printf("Contract Address:  %s\n", contractAddress)
	fmt.Printf("Token ID:          %d\n", tokenId)
	fmt.Printf("Current Owner:     %s\n", owner.Hex())
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println()

	// Step 5: Check if querier is the owner
	fmt.Println("Ownership verification:")
	if owner.Hex() == acc.GetEVMAddress().Hex() {
		fmt.Printf("   You own this NFT\n")
		fmt.Printf("Your address: %s\n", acc.GetEVMAddress().Hex())
	} else {
		fmt.Printf("You don't own this NFT\n")
		fmt.Printf("Your address:  %s\n", acc.GetEVMAddress().Hex())
		fmt.Printf("Owner address: %s\n", owner.Hex())
	}
	fmt.Println()

	// Explanation
	fmt.Println("What this query tells us:")
	fmt.Println("• Who currently owns the NFT")
	fmt.Println("• The NFT exists and has been minted")
	fmt.Println("• The contract is deployed and accessible")
	fmt.Println("• Ownership can be verified on-chain")
	fmt.Println()

	// Additional information
	fmt.Println("Additional query capabilities:")
	fmt.Println("• You can query any token ID on the contract")
	fmt.Println("• Queries are read-only and don't cost gas")
	fmt.Println("• Anyone can query NFT ownership")
	fmt.Println("• This is useful for verification and auditing")
	fmt.Println()

	fmt.Println("Next steps:")
	fmt.Println("• Query multiple token IDs to see different owners")
	fmt.Println("• Build a frontend to display NFT information")
	fmt.Println("• Integrate queries into your application logic")
	fmt.Println()
}
