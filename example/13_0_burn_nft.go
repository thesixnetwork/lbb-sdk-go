package main

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/client"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/evm"
)

// This example demonstrates how to burn an NFT.
// Burning permanently destroys the NFT, making it non-existent.
// After burning, the token owner becomes the zero address.
//
// Usage:
//   go run 13_0_burn_nft.go
//
// What this script does:
// 1. Connects to the network and creates an account
// 2. Creates an EVM client
// 3. Verifies ownership of the NFT
// 4. Burns the NFT (permanently destroys it)
// 5. Verifies the burn by checking the owner is now the zero address
//
// Prerequisites:
// - NFT must be minted first (see 06_mint_nft.go)
// - Account must own the NFT being burned
// - Account must have tokens for gas fees
//
// Important Notes:
// - Burning is IRREVERSIBLE - the NFT cannot be recovered
// - After burning, the token ID cannot be used again (depending on contract implementation)
// - The owner after burn will be the zero address (0x0000000000000000000000000000000000000000)

const (
	// IMPORTANT: Replace with your deployed contract address
	contractAddress = "0x0000000000000000000000000000000000000000"

	// Token ID to burn (must exist and be owned by you)
	// WARNING: This NFT will be permanently destroyed!
	tokenId = uint64(1)

	// For this example, we use the test mnemonic
	exampleMnemonic = account.TestMnemonic
)

func main() {
	fmt.Println("=== Step 13.0: Burn NFT ===")
	fmt.Println()
	fmt.Println("⚠️  WARNING: Burning permanently destroys the NFT!")
	fmt.Println("This action cannot be undone.")
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

	acc, err := account.NewAccount(client, "myaccount", exampleMnemonic, "")
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
		fmt.Println()
		fmt.Println("Only the owner can burn their NFT.")
		return
	}
	fmt.Println("✓ Ownership verified")
	fmt.Println()

	// Step 4: Burn NFT
	fmt.Println("Burning NFT...")
	fmt.Printf("   Owner:    %s\n", acc.GetEVMAddress().Hex())
	fmt.Printf("   Token ID: %d\n", tokenId)
	fmt.Println()
	fmt.Println("⚠️  Last chance to cancel! The NFT will be permanently destroyed.")
	fmt.Println()

	// BurnCertificateNFT burns the NFT permanently
	// Parameters:
	// - contractAddress: the deployed contract address
	// - tokenId: the NFT token ID to burn
	// Returns:
	// - tx: transaction object
	// - err: error if any
	tx, err := evmClient.BurnCertificateNFT(contractAddr, tokenId)
	if err != nil {
		panic(fmt.Sprintf("Failed to burn NFT: %v", err))
	}

	fmt.Printf("Burn transaction submitted: %s\n", tx.Hash().Hex())
	fmt.Printf("   Nonce: %d\n", tx.Nonce())
	fmt.Println()

	// Step 5: Wait for transaction confirmation
	fmt.Println("Waiting for transaction to be mined...")

	receipt, err := client.WaitForEVMTransaction(tx.Hash())
	if err != nil {
		panic(fmt.Sprintf("Error waiting for burn transaction: %v", err))
	}

	fmt.Println("Burn completed!")
	fmt.Println()

	// Step 6: Verify the NFT was burned
	fmt.Println("Verifying burn...")
	burnedOwner := evmClient.TokenOwner(contractAddr, tokenId)
	zeroAddress := common.HexToAddress("0x0000000000000000000000000000000000000000")

	// Step 7: Display results
	fmt.Println()
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("Burn Information:")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Printf("Contract Address:  %s\n", contractAddress)
	fmt.Printf("Token ID:          %d\n", tokenId)
	fmt.Printf("Previous Owner:    %s\n", acc.GetEVMAddress().Hex())
	fmt.Printf("Current Owner:     %s\n", burnedOwner.Hex())
	fmt.Printf("Zero Address:      %s\n", zeroAddress.Hex())
	fmt.Printf("Transaction Hash:  %s\n", tx.Hash().Hex())
	fmt.Printf("Block Number:      %d\n", receipt.BlockNumber)
	fmt.Printf("Gas Used:          %d\n", receipt.GasUsed)
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println()

	// Verify the burn was successful
	if burnedOwner == zeroAddress {
		fmt.Println("✓ NFT successfully burned!")
		fmt.Printf("   Token ID %d no longer exists\n", tokenId)
		fmt.Printf("   Owner is now the zero address: %s\n", burnedOwner.Hex())
		fmt.Println()
		fmt.Println("♨️  The NFT has been permanently destroyed and cannot be recovered.")
	} else {
		fmt.Println("⚠️  WARNING: Burn verification failed!")
		fmt.Printf("   Expected owner: %s (zero address)\n", zeroAddress.Hex())
		fmt.Printf("   Actual owner:   %s\n", burnedOwner.Hex())
		fmt.Println()
		fmt.Println("The burn transaction completed but the owner is not the zero address.")
		fmt.Println("This might indicate an issue with the contract or burn implementation.")
	}
	fmt.Println()

	// Explanation
	fmt.Println("What just happened:")
	fmt.Println("  • The NFT was permanently destroyed")
	fmt.Println("  • The token owner is now the zero address")
	fmt.Println("  • The token ID cannot be transferred anymore")
	fmt.Println("  • This action is irreversible - the NFT is gone forever")
	fmt.Println()

	fmt.Println("Use cases for burning NFTs:")
	fmt.Println("  • Removing expired or invalid certificates")
	fmt.Println("  • Destroying test or demo NFTs")
	fmt.Println("  • Implementing token deflationary mechanics")
	fmt.Println("  • Certificate revocation (permanent invalidation)")
	fmt.Println("  • Cleaning up unused tokens")
	fmt.Println()

	fmt.Println("Burn vs Transfer:")
	fmt.Println("  • Burn: Permanently destroys the NFT (cannot be undone)")
	fmt.Println("  • Transfer: Changes ownership (can be transferred again)")
	fmt.Println("  • Both can be done with or without gas (see gasless examples)")
	fmt.Println()

	fmt.Println("Next steps:")
	fmt.Println("  • Try gasless burning (13_gasless_burn.go)")
	fmt.Println("  • Query metadata to see certificate status")
	fmt.Println("  • Mint a new NFT to replace the burned one")
	fmt.Println()
}
