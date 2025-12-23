package main

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/client"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/evm"
)

// This example demonstrates how to transfer an NFT to another address.
// Transfer changes the ownership of an NFT from one address to another.
//
// Usage:
//   go run 07_transfer_nft.go
//
// What this script does:
// 1. Connects to the network and creates an account
// 2. Creates an EVM client
// 3. Transfers the NFT to a recipient address
// 4. Verifies the new owner
//
// Prerequisites:
// - NFT must be minted first (see 06_mint_nft.go)
// - Account must own the NFT being transferred
// - Account must have tokens for gas fees
// - Recipient address must be valid

const (
	// IMPORTANT: Replace with your deployed contract address
	contractAddress = "0x0000000000000000000000000000000000000000"

	// Token ID to transfer (must exist and be owned by you)
	tokenId = uint64(1)

	// Recipient address (who will receive the NFT)
	// Replace with the actual recipient's EVM address
	recipientAddress = "0x8a28fb81A084Ac7A276800957a19a6054BF86E4D"

	// For this example, we use the test mnemonic
	exampleMnemonic = account.TestMnemonic
)

func main() {
	fmt.Println("=== Step 7: Transfer NFT ===")
	fmt.Println()

	// Validate configuration
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

	// Step 3: Verify current ownership
	contractAddr := common.HexToAddress(contractAddress)

	fmt.Println("Verifying current ownership...")
	currentOwner := evmClient.TokenOwner(contractAddr, tokenId)
	fmt.Printf("   Current owner: %s\n", currentOwner.Hex())
	fmt.Printf("   Your address:  %s\n", acc.GetEVMAddress().Hex())

	if currentOwner.Hex() != acc.GetEVMAddress().Hex() {
		fmt.Println()
		fmt.Println("ERROR: You don't own this NFT!")
		fmt.Printf("   Token ID %d is owned by: %s\n", tokenId, currentOwner.Hex())
		fmt.Printf("   But you are: %s\n", acc.GetEVMAddress().Hex())
		return
	}
	fmt.Println("Ownership verified")
	fmt.Println()

	// Step 4: Transfer NFT
	fmt.Println("Transferring NFT...")
	fmt.Printf("   From:     %s\n", acc.GetEVMAddress().Hex())
	fmt.Printf("   To:       %s\n", recipientAddress)
	fmt.Printf("   Token ID: %d\n", tokenId)
	fmt.Println()

	// Convert recipient string to common.Address
	recipientAddr := common.HexToAddress(recipientAddress)

	// TransferCertificateNFT transfers the NFT to a new owner
	// Parameters:
	// - contractAddress: the deployed contract address
	// - to: recipient's address
	// - tokenId: the NFT token ID to transfer
	// Returns:
	// - tx: transaction object
	// - err: error if any
	tx, err := evmClient.TransferCertificateNFT(contractAddr, recipientAddr, tokenId)
	if err != nil {
		panic(fmt.Sprintf("Failed to transfer NFT: %v", err))
	}

	fmt.Printf("Transfer transaction submitted: %s\n", tx.Hash().Hex())
	fmt.Printf("   Nonce: %d\n", tx.Nonce())
	fmt.Println()

	// Step 5: Wait for transaction confirmation
	fmt.Println("Waiting for transaction to be mined...")

	receipt, err := client.WaitForEVMTransaction(tx.Hash())
	if err != nil {
		panic(fmt.Sprintf("Error waiting for transfer transaction: %v", err))
	}

	fmt.Println("Transfer completed!")
	fmt.Println()

	// Step 6: Verify new ownership
	fmt.Println("Verifying new ownership...")
	newOwner := evmClient.TokenOwner(contractAddr, tokenId)

	// Step 7: Display results
	fmt.Println()
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("Transfer Information:")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Printf("Contract Address:  %s\n", contractAddress)
	fmt.Printf("Token ID:          %d\n", tokenId)
	fmt.Printf("Previous Owner:    %s\n", acc.GetEVMAddress().Hex())
	fmt.Printf("New Owner:         %s\n", newOwner.Hex())
	fmt.Printf("Recipient:         %s\n", recipientAddress)
	fmt.Printf("Transaction Hash:  %s\n", tx.Hash().Hex())
	fmt.Printf("Block Number:      %d\n", receipt.BlockNumber)
	fmt.Printf("Gas Used:          %d\n", receipt.GasUsed)
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println()

	// Verify the transfer was successful
	if newOwner.Hex() == recipientAddr.Hex() {
		fmt.Println("Transfer verified successfully!")
		fmt.Printf("   Token ID %d is now owned by %s\n", tokenId, newOwner.Hex())
	} else {
		fmt.Println("WARNING: Owner verification mismatch!")
		fmt.Printf("   Expected: %s\n", recipientAddr.Hex())
		fmt.Printf("   Actual:   %s\n", newOwner.Hex())
	}
	fmt.Println()

	// Explanation
	fmt.Println("What just happened:")
	fmt.Println("  • Transferred ownership of the NFT")
	fmt.Println("  • The recipient now owns the certificate")
	fmt.Println("  • The transfer is recorded on the blockchain")
	fmt.Println("  • Only the new owner can transfer it again")
	fmt.Println()

	fmt.Println("Next steps:")
	fmt.Println("  • Freeze/unfreeze the metadata (08_freeze_metadata.go)")
	fmt.Println("  • Query NFT information (09_query_nft.go)")
	fmt.Println("  • The recipient can now transfer the NFT to others")
	fmt.Println()
}
