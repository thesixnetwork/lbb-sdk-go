package main

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/client"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/evm"
)

// This example demonstrates how to mint an NFT from a deployed contract.
// Minting creates a new token with a unique ID.
//
// Usage:
//   go run 06_mint_nft.go
//
// What this script does:
// 1. Connects to the network and creates an account
// 2. Creates an EVM client
// 3. Mints a certificate NFT with a specific token ID
// 4. Waits for the mint transaction to complete
//
// Prerequisites:
// - EVM contract must be deployed first (see 05_deploy_contract.go)
// - Account must have tokens for gas fees
// - Contract address must be valid

const (
	// IMPORTANT: Replace this with your actual deployed contract address
	// You get this from the output of 05_deploy_contract.go
	contractAddress = "0x0000000000000000000000000000000000000000"

	// Token ID to mint (must be unique)
	tokenId = uint64(1)

	// For this example, we use the test mnemonic
	exampleMnemonic = account.TestMnemonic
)

func main() {
	fmt.Println("=== Step 6: Mint Certificate NFT ===")
	fmt.Println()

	// Validate contract address
	if contractAddress == "0x0000000000000000000000000000000000000000" {
		fmt.Println("ERROR: Please update the contractAddress constant")
		fmt.Println("   Use the contract address from step 05_deploy_contract.go")
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

	fmt.Printf("Connected with account: %s\n", acc.GetEVMAddress().Hex())
	fmt.Println()

	// Step 2: Create EVM client
	fmt.Println("Initializing EVM client...")
	evmClient := evm.NewEVMClient(*acc)
	fmt.Println("EVM client initialized")
	fmt.Println()

	// Step 3: Mint NFT
	fmt.Println("Minting certificate NFT...")
	fmt.Printf("   Contract: %s\n", contractAddress)
	fmt.Printf("   Token ID: %d\n", tokenId)
	fmt.Println()

	// Convert string address to common.Address type
	contractAddr := common.HexToAddress(contractAddress)

	// MintCertificateNFT mints a new NFT
	// Parameters:
	// - contractAddress: the deployed contract address
	// - tokenId: unique identifier for this NFT (must not exist)
	// Returns:
	// - tx: transaction object
	// - err: error if any
	tx, err := evmClient.MintCertificateNFT(contractAddr, tokenId)
	if err != nil {
		panic(fmt.Sprintf("Failed to mint NFT: %v", err))
	}

	fmt.Printf("Mint transaction submitted: %s\n", tx.Hash().Hex())
	fmt.Printf("   Nonce: %d\n", tx.Nonce())
	fmt.Println()

	// Step 4: Wait for transaction confirmation
	fmt.Println("Waiting for transaction to be mined...")

	// WaitForEVMTransaction waits until the transaction is confirmed
	receipt, err := client.WaitForEVMTransaction(tx.Hash())
	if err != nil {
		panic(fmt.Sprintf("Error waiting for mint transaction: %v", err))
	}

	// Step 5: Verify ownership
	fmt.Println("NFT minted successfully!")
	fmt.Println()
	fmt.Println("Verifying ownership...")

	// TokenOwner retrieves the current owner of a token
	owner := evmClient.TokenOwner(contractAddr, tokenId)

	// Step 6: Display results
	fmt.Println()
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("Mint Information:")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Printf("Contract Address:  %s\n", contractAddress)
	fmt.Printf("Token ID:          %d\n", tokenId)
	fmt.Printf("Owner Address:     %s\n", owner.Hex())
	fmt.Printf("Minter Address:    %s\n", acc.GetEVMAddress().Hex())
	fmt.Printf("Transaction Hash:  %s\n", tx.Hash().Hex())
	fmt.Printf("Block Number:      %d\n", receipt.BlockNumber)
	fmt.Printf("Gas Used:          %d\n", receipt.GasUsed)
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println()

	// Explanation
	fmt.Println("What just happened:")
	fmt.Println("  • Created a new NFT with token ID", tokenId)
	fmt.Println("  • The NFT was minted to your address")
	fmt.Println("  • The NFT is linked to the metadata schema")
	fmt.Println("  • You now own this certificate NFT")
	fmt.Println()

	fmt.Println("Next steps:")
	fmt.Println("  • Transfer the NFT to another address (07_transfer_nft.go)")
	fmt.Println("  • Freeze/unfreeze the metadata (08_freeze_metadata.go)")
	fmt.Printf("  • Token ID to use: %d\n", tokenId)
	fmt.Println()
}
