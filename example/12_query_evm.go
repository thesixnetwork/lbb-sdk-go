package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/client"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/evm"
)

// This example demonstrates how to query EVM layer information.
// You can check NFT ownership, gas prices, chain ID, nonce, and transaction receipts.
//
// Usage:
//   go run 12_query_evm.go
//
// What this script does:
// 1. Connects to the network and creates an account
// 2. Queries current gas price
// 3. Queries chain ID
// 4. Queries account nonce
// 5. Queries NFT token ownership
// 6. Demonstrates transaction receipt checking
//
// Prerequisites:
// - Network connection to testnet or mainnet
// - For NFT queries: contract must be deployed and NFT minted

const (
	// IMPORTANT: Replace with your deployed contract address
	contractAddress = "0x0000000000000000000000000000000000000000"

	// Token ID to query
	tokenId = uint64(1)

	// For this example, we use the test mnemonic
	exampleMnemonic = account.TestMnemonic
)

func main() {
	fmt.Println("=== Step 12: Query EVM Information ===")
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

	fmt.Printf("Connected with account\n")
	fmt.Printf("  Cosmos Address: %s\n", acc.GetCosmosAddress().String())
	fmt.Printf("  EVM Address: %s\n", acc.GetEVMAddress().Hex())
	fmt.Println()

	// Step 2: Create EVM client
	fmt.Println("Initializing EVM client...")
	evmClient := evm.NewEVMClient(*acc)
	fmt.Println("EVM client initialized")
	fmt.Println()

	// Step 3: Query Gas Price
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("Querying Current Gas Price")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println()

	// GasPrice retrieves the current suggested gas price from the network
	// This is useful for estimating transaction costs
	// Returns: *big.Int (gas price in wei)
	gasPrice, err := evmClient.GasPrice()
	if err != nil {
		fmt.Printf("Warning: Failed to get gas price: %v\n", err)
	} else {
		fmt.Printf("Current Gas Price: %s wei\n", gasPrice.String())

		// Convert to Gwei for readability (1 Gwei = 1,000,000,000 wei)
		gwei := new(big.Int).Div(gasPrice, big.NewInt(1000000000))
		fmt.Printf("  (≈ %s Gwei)\n", gwei.String())

		// Estimate cost for a typical transaction (21,000 gas)
		standardGas := big.NewInt(21000)
		estimatedCost := new(big.Int).Mul(gasPrice, standardGas)
		fmt.Printf("\nEstimated cost for standard transaction (21,000 gas):\n")
		fmt.Printf("  %s wei\n", estimatedCost.String())
	}
	fmt.Println()

	// Step 4: Query Chain ID
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("Querying Chain ID")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println()

	// ChainID retrieves the network's chain ID
	// This is used for transaction signing (EIP-155)
	// Returns: *big.Int (chain ID)
	chainID, err := evmClient.ChainID()
	if err != nil {
		fmt.Printf("Warning: Failed to get chain ID: %v\n", err)
	} else {
		fmt.Printf("Chain ID: %s\n", chainID.String())

		// Identify network
		switch chainID.String() {
		case "97":
			fmt.Println("Network: Fivenet (Testnet)")
		case "98":
			fmt.Println("Network: Sixnet (Mainnet)")
		default:
			fmt.Printf("Network: Unknown (Chain ID: %s)\n", chainID.String())
		}
	}
	fmt.Println()

	// Step 5: Query Account Nonce
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("Querying Account Nonce")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println()

	// GetNonce retrieves the next nonce for the account
	// The nonce is used to order transactions and prevent replay attacks
	// Returns: uint64 (next nonce to use)
	nonce, err := evmClient.GetNonce()
	if err != nil {
		fmt.Printf("Warning: Failed to get nonce: %v\n", err)
	} else {
		fmt.Printf("Current Nonce: %d\n", nonce)
		fmt.Printf("  This is the transaction number for your next EVM transaction\n")

		if nonce == 0 {
			fmt.Println("  (No EVM transactions sent yet from this account)")
		} else {
			fmt.Printf("  (You have sent %d EVM transaction(s) so far)\n", nonce)
		}
	}
	fmt.Println()

	// Step 6: Query NFT Token Ownership (if contract address is provided)
	if contractAddress != "0x0000000000000000000000000000000000000000" {
		fmt.Println("────────────────────────────────────────────────────────────────")
		fmt.Println("Querying NFT Token Ownership")
		fmt.Println("────────────────────────────────────────────────────────────────")
		fmt.Println()

		contractAddr := common.HexToAddress(contractAddress)

		fmt.Printf("Contract: %s\n", contractAddress)
		fmt.Printf("Token ID: %d\n", tokenId)
		fmt.Println()

		// TokenOwner retrieves the current owner of an NFT token
		// Parameters:
		// - contractAddress: the NFT contract address
		// - tokenID: the token ID to query
		// Returns: common.Address (owner's address)
		owner := evmClient.TokenOwner(contractAddr, tokenId)

		fmt.Printf("Token Owner: %s\n", owner.Hex())
		fmt.Println()

		// Verify ownership
		if owner.Hex() == acc.GetEVMAddress().Hex() {
			fmt.Println("✓ You own this NFT")
		} else if owner.Hex() == "0x0000000000000000000000000000000000000000" {
			fmt.Println("⚠ Token doesn't exist or hasn't been minted yet")
		} else {
			fmt.Println("✗ This NFT is owned by someone else")
			fmt.Printf("Your address: %s\n", acc.GetEVMAddress().Hex())
		}
		fmt.Println()
	} else {
		fmt.Println("────────────────────────────────────────────────────────────────")
		fmt.Println("NFT Query Skipped")
		fmt.Println("────────────────────────────────────────────────────────────────")
		fmt.Println()
		fmt.Println("To query NFT ownership:")
		fmt.Println("1. Update the contractAddress constant with your contract")
		fmt.Println("2. Make sure the contract is deployed (see 05_deploy_contract.go)")
		fmt.Println("3. Make sure NFTs are minted (see 06_mint_nft.go)")
		fmt.Println()
	}

	// Step 7: Transaction Receipt Example
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("Checking Transaction Receipt")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println()
	fmt.Println("To check a transaction receipt, use:")
	fmt.Println()
	fmt.Println("  txHash := common.HexToHash(\"0x...\")")
	fmt.Println("  err := evmClient.CheckTransactionReceipt(txHash)")
	fmt.Println("  if err != nil {")
	fmt.Println("      fmt.Printf(\"Transaction failed: %v\\n\", err)")
	fmt.Println("  } else {")
	fmt.Println("      fmt.Println(\"Transaction successful!\")")
	fmt.Println("  }")
	fmt.Println()
	fmt.Println("This will display:")
	fmt.Println("  • Transaction hash")
	fmt.Println("  • Block number")
	fmt.Println("  • Gas used")
	fmt.Println("  • Success/failure status")
	fmt.Println("  • Contract address (if contract deployment)")
	fmt.Println()

	// Summary
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Println("EVM Query Summary")
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Printf("Your EVM Address: %s\n", acc.GetEVMAddress().Hex())
	if nonce > 0 {
		fmt.Printf("Transaction Count: %d\n", nonce)
	}
	if gasPrice != nil {
		fmt.Printf("Current Gas Price: %s wei\n", gasPrice.String())
	}
	if chainID != nil {
		fmt.Printf("Chain ID: %s\n", chainID.String())
	}
	fmt.Println()

	fmt.Println("What you learned:")
	fmt.Println("• How to query current gas price for cost estimation")
	fmt.Println("• How to get the network chain ID")
	fmt.Println("• How to check account nonce for transaction ordering")
	fmt.Println("• How to query NFT token ownership")
	fmt.Println("• How to verify transaction receipts")
	fmt.Println()

	fmt.Println("Why these queries matter:")
	fmt.Println("• Gas Price: Estimate transaction costs before sending")
	fmt.Println("• Chain ID: Ensure you're on the correct network")
	fmt.Println("• Nonce: Track transaction count and prevent issues")
	fmt.Println("• Token Owner: Verify NFT ownership on-chain")
	fmt.Println("• Receipt: Confirm transaction success and get details")
	fmt.Println()

	fmt.Println("Best practices:")
	fmt.Println("• Always check gas price before expensive operations")
	fmt.Println("• Verify chain ID matches your intended network")
	fmt.Println("• Use nonce to track transaction history")
	fmt.Println("• Query ownership before transfer operations")
	fmt.Println("• Check receipts to confirm transaction success")
	fmt.Println()

	fmt.Println("Next steps:")
	fmt.Println("• Use these queries to build transaction monitoring")
	fmt.Println("• Integrate into your application for real-time data")
	fmt.Println("• Build cost estimation features")
	fmt.Println("• Create ownership verification systems")
	fmt.Println()
}
