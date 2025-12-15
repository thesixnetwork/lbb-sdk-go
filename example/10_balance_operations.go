package main

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/client"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/balance"
)

// This example demonstrates how to query and transfer balances.
// You can check both Cosmos layer (usix) and EVM layer (asix) balances,
// and transfer tokens between addresses.
//
// Usage:
//   go run 10_balance_operations.go
//
// What this script does:
// 1. Connects to the network and creates an account
// 2. Queries all balances for the account
// 3. Queries Cosmos layer balance (usix)
// 4. Queries EVM layer balance (asix)
// 5. Demonstrates how to transfer tokens
//
// Prerequisites:
// - Account must have some balance for queries to show results
// - Account must have sufficient balance for transfers

const (
	// For this example, we use the test mnemonic
	exampleMnemonic = account.TestMnemonic

	// IMPORTANT: Replace with actual recipient address for transfers
	recipientAddress = "6x13g50hqdqsjk85fmgqz2h5xdxq49lsmjdwlemsp"
)

func main() {
	fmt.Println("=== Step 10: Balance Operations ===")
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

	// Step 2: Create balance client for queries
	fmt.Println("Initializing balance client...")
	bal := balance.NewBalance(*acc)
	fmt.Println("Balance client initialized")
	fmt.Println()

	// Step 3: Query all balances
	fmt.Println("Querying all balances...")

	// GetBalance retrieves all token balances for the account
	// This includes both Cosmos layer (usix) and any other tokens
	// Returns: sdk.Coins (a list of all balances)
	allBalances, err := bal.GetBalance()
	if err != nil {
		panic(fmt.Sprintf("Failed to get balance: %v", err))
	}

	fmt.Println("All balances:")
	if allBalances.IsZero() {
		fmt.Println("  No balances found")
		fmt.Println("  Request test tokens from the faucet to see balances")
	} else {
		for _, coin := range allBalances {
			fmt.Printf("  %s\n", coin.String())
		}
	}
	fmt.Println()

	// Step 4: Query Cosmos layer balance (usix)
	fmt.Println("Querying Cosmos layer balance (usix)...")

	// GetCosmosBalance retrieves only the Cosmos native token balance
	// Denomination: usix (1 six = 1,000,000 usix)
	// Returns: sdk.Coin (amount and denomination)
	cosmosBalance, err := bal.GetCosmosBalance()
	if err != nil {
		panic(fmt.Sprintf("Failed to get cosmos balance: %v", err))
	}

	fmt.Printf("Cosmos balance: %s\n", cosmosBalance.String())
	if !cosmosBalance.IsZero() {
		// Convert usix to six for readability
		sixAmount := cosmosBalance.Amount.Quo(sdkmath.NewInt(1000000))
		fmt.Printf("  (≈ %s SIX)\n", sixAmount.String())
	}
	fmt.Println()

	// Step 5: Query EVM layer balance (asix)
	fmt.Println("Querying EVM layer balance (asix)...")

	// GetEVMBalance retrieves the EVM layer token balance
	// Denomination: asix (used for EVM transactions)
	// Returns: sdk.Coin (amount and denomination)
	evmBalance, err := bal.GetEVMBalance()
	if err != nil {
		panic(fmt.Sprintf("Failed to get EVM balance: %v", err))
	}

	fmt.Printf("EVM balance: %s\n", evmBalance.String())
	if !evmBalance.IsZero() {
		// Convert asix to six for readability
		sixAmount := evmBalance.Amount.Quo(sdkmath.NewInt(1000000000000000000))
		fmt.Printf("  (≈ %s SIX)\n", sixAmount.String())
	}
	fmt.Println()

	// Step 6: Demonstrate balance transfer (commented out by default)
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("Balance Transfer Example")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println()
	fmt.Println("To send tokens, use the following code:")
	fmt.Println()
	fmt.Println("  // Create balance message client")
	fmt.Println("  balMsg, err := balance.NewBalanceMsg(*acc)")
	fmt.Println("  if err != nil {")
	fmt.Println("      panic(fmt.Sprintf(\"Failed to create balance msg: %v\", err))")
	fmt.Println("  }")
	fmt.Println()
	fmt.Println("  // Define amount to send (1 SIX = 1,000,000 usix)")
	fmt.Println("  amount := sdk.NewCoins(sdk.NewInt64Coin(\"usix\", 1000000))")
	fmt.Println()
	fmt.Println("  // Send balance (returns immediately)")
	fmt.Println("  res, err := balMsg.SendBalance(recipientAddress, amount)")
	fmt.Println("  if err != nil {")
	fmt.Println("      panic(fmt.Sprintf(\"Failed to send balance: %v\", err))")
	fmt.Println("  }")
	fmt.Println("  fmt.Printf(\"Transfer tx: %s\\n\", res.TxHash)")
	fmt.Println()
	fmt.Println("  // OR send and wait for confirmation")
	fmt.Println("  res, err := balMsg.SendBalanceAndWait(recipientAddress, amount)")
	fmt.Println("  if err != nil {")
	fmt.Println("      panic(fmt.Sprintf(\"Failed to send balance: %v\", err))")
	fmt.Println("  }")
	fmt.Println("  fmt.Printf(\"Transfer confirmed: %s\\n\", res.TxHash)")
	fmt.Println()
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println()

	// Uncomment below to actually send tokens
	/*
		// Create balance message client for transactions
		balMsg, err := balance.NewBalanceMsg(*acc)
		if err != nil {
			panic(fmt.Sprintf("Failed to create balance msg: %v", err))
		}

		// Define amount to send (1 SIX = 1,000,000 usix)
		amount := sdk.NewCoins(sdk.NewInt64Coin("usix", 1000000))

		fmt.Println("Sending tokens...")
		fmt.Printf("  From: %s\n", acc.GetCosmosAddress().String())
		fmt.Printf("  To: %s\n", recipientAddress)
		fmt.Printf("  Amount: %s\n", amount.String())
		fmt.Println()

		// SendBalanceAndWait sends tokens and waits for confirmation
		// Parameters:
		// - dest: recipient address (Cosmos format: 6x...)
		// - amount: coins to send
		// Returns:
		// - TxResponse with transaction hash and result
		res, err := balMsg.SendBalanceAndWait(recipientAddress, amount)
		if err != nil {
			panic(fmt.Sprintf("Failed to send balance: %v", err))
		}

		fmt.Println("Transfer completed!")
		fmt.Printf("  Transaction Hash: %s\n", res.TxHash)
		fmt.Printf("  Height: %d\n", res.Height)
		fmt.Println()
	*/

	// Summary
	fmt.Println("Balance operations completed!")
	fmt.Println()
	fmt.Println("What you learned:")
	fmt.Println("• How to query all balances for an account")
	fmt.Println("• How to query Cosmos layer balance (usix)")
	fmt.Println("• How to query EVM layer balance (asix)")
	fmt.Println("• How to transfer tokens between addresses")
	fmt.Println()
	fmt.Println("Key differences:")
	fmt.Println("• Cosmos layer (usix): Used for Cosmos transactions (metadata, schemas)")
	fmt.Println("• EVM layer (asix): Used for EVM transactions (contract deployment, NFT minting)")
	fmt.Println("• Both use the same underlying token but different accounting")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("• Check your balance before performing transactions")
	fmt.Println("• Ensure sufficient balance for gas fees")
	fmt.Println("• Use SendBalanceAndWait() for confirmed transfers")
	fmt.Println()
}
