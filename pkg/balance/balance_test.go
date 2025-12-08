package balance_test

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/client"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/balance"
)

func Example_queryOnly() {
	// Setup client and account
	ctx := context.Background()
	c, _ := client.NewClient(ctx, false)
	acc, _ := account.NewAccount(c, "myaccount", "your mnemonic here...", "password")

	// Create Balance instance for queries (lightweight, no tx overhead)
	bal := balance.NewBalance(*acc)

	// Query all balances
	coins, err := bal.GetBalance()
	if err != nil {
		fmt.Printf("Error getting balance: %v\n", err)
		return
	}
	fmt.Printf("All balances: %v\n", coins)

	cosmosBalance, _ := bal.GetCosmosBalance()
	fmt.Printf("Cosmos balance: %v\n", cosmosBalance)

	evmBalance, _ := bal.GetEVMBalance()
	fmt.Printf("EVM balance: %v\n", evmBalance)

	// Query custom denom
	customBalance, _ := bal.GetBalanceByDenom("usix")
	fmt.Printf("Custom denom balance: %v\n", customBalance)
}

func Example_transactions() {
	ctx := context.Background()
	c, _ := client.NewClient(ctx, false)
	acc, _ := account.NewAccount(c, "myaccount", "your mnemonic here...", "password")

	// Create BalanceTx instance for both queries and transactions
	balTx, err := balance.NewBalanceMsg(*acc)
	if err != nil {
		fmt.Printf("Error creating BalanceTx: %v\n", err)
		return
	}

	coins, _ := balTx.GetBalance()
	fmt.Printf("Current balance: %v\n", coins)

	// Send tokens with default settings
	amount := sdk.NewCoins(sdk.NewInt64Coin("usix", 1000000))
	res, err := balTx.SendBalance("six1recipient...", amount)
	if err != nil {
		fmt.Printf("Error sending balance: %v\n", err)
		return
	}
	fmt.Printf("Transaction hash: %s\n", res.TxHash)
}

func Example_transactionWithOptions() {
	ctx := context.Background()
	c, _ := client.NewClient(ctx, false)
	acc, _ := account.NewAccount(c, "myaccount", "your mnemonic here...", "password")

	balTx, _ := balance.NewBalanceMsg(*acc)

	// Configure transaction with custom gas, memo, etc.
	amount := sdk.NewCoins(sdk.NewInt64Coin("usix", 1000000))
	res, err := balTx.
		WithGas(500000).
		WithGasAdjustment(1.5).
		WithMemo("Payment for services").
		SendBalance("six1recipient...", amount)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Transaction sent with hash: %s\n", res.TxHash)
	fmt.Printf("Gas used: %d\n", res.GasUsed)
}

// Example_sendAndWait demonstrates sending a transaction and waiting for confirmation
func Example_sendAndWait() {
	// Setup client and account
	ctx := context.Background()
	c, _ := client.NewClient(ctx, false)
	acc, _ := account.NewAccount(c, "myaccount", "your mnemonic here...", "password")

	balTx, _ := balance.NewBalanceMsg(*acc)

	amount := sdk.NewCoins(sdk.NewInt64Coin("usix", 1000000))

	// Send and wait for confirmation
	res, err := balTx.
		WithMemo("Payment").
		SendBalanceAndWait("six1recipient...", amount)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Transaction confirmed! Hash: %s\n", res.TxHash)
}

// Example_buildMessage demonstrates building a message without broadcasting
// Useful for batch operations or custom transaction handling
func Example_buildMessage() {
	ctx := context.Background()
	c, _ := client.NewClient(ctx, false)
	acc, _ := account.NewAccount(c, "myaccount", "your mnemonic here...", "password")

	balTx, _ := balance.NewBalanceMsg(*acc)

	amount := sdk.NewCoins(sdk.NewInt64Coin("usix", 1000000))

	// Build message without broadcasting
	msg, err := balTx.BuildSendMsg("6xsomething...", amount)
	if err != nil {
		fmt.Printf("Error building message: %v\n", err)
		return
	}

	fmt.Printf("Message built: %v\n", msg)

	// Example: broadcast multiple messages in one transaction
	msg2, _ := balTx.BuildSendMsg("six1another...", amount)
	res, err := balTx.BroadcastTx(msg, msg2)
	if err != nil {
		fmt.Printf("Error broadcasting: %v\n", err)
		return
	}

	fmt.Printf("Batch transaction hash: %s\n", res.TxHash)
}
