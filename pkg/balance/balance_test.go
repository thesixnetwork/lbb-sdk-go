package balance_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/client"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/balance"
)

// Test mnemonic (DO NOT use in production)
const (
	testMnemonic = "test test test test test test test test test test test junk"
	testPassword = ""
)

func TestNewBalance(t *testing.T) {
	t.Run("Create Balance instance", func(t *testing.T) {
		ctx := context.Background()
		c, err := client.NewClient(ctx, false)
		require.NoError(t, err, "Should create client without error")

		acc, err := account.NewAccount(c, "testaccount", testMnemonic, testPassword)
		require.NoError(t, err, "Should create account without error")

		bal := balance.NewBalance(*acc)
		require.NotNil(t, bal, "Balance should not be nil")

		// Verify we can access the account
		retrievedAcc := bal.GetAccount()
		assert.Equal(t, acc.GetAccountName(), retrievedAcc.GetAccountName())
	})
}

func TestNewBalanceMsg(t *testing.T) {
	t.Run("Create BalanceMsg instance", func(t *testing.T) {
		ctx := context.Background()
		c, err := client.NewClient(ctx, false)
		require.NoError(t, err, "Should create client without error")

		acc, err := account.NewAccount(c, "testaccount", testMnemonic, testPassword)
		require.NoError(t, err, "Should create account without error")

		balMsg, err := balance.NewBalanceMsg(*acc)
		require.NoError(t, err, "Should create BalanceMsg without error")
		require.NotNil(t, balMsg, "BalanceMsg should not be nil")

		// Verify we can access the account (inherited from Balance)
		retrievedAcc := balMsg.GetAccount()
		assert.Equal(t, acc.GetAccountName(), retrievedAcc.GetAccountName())
	})
}

func TestBalanceInterface(t *testing.T) {
	t.Run("Balance implements BalanceI", func(t *testing.T) {
		ctx := context.Background()
		c, err := client.NewClient(ctx, false)
		require.NoError(t, err)

		acc, err := account.NewAccount(c, "testaccount", testMnemonic, testPassword)
		require.NoError(t, err)

		bal := balance.NewBalance(*acc)

		// Verify Balance implements BalanceI
		var _ balance.BalanceI = bal

		t.Log("Balance correctly implements BalanceI interface")
	})
}

func TestBalanceMsgInterface(t *testing.T) {
	t.Run("BalanceMsg implements BalanceMsgI", func(t *testing.T) {
		ctx := context.Background()
		c, err := client.NewClient(ctx, false)
		require.NoError(t, err)

		acc, err := account.NewAccount(c, "testaccount", testMnemonic, testPassword)
		require.NoError(t, err)

		balMsg, err := balance.NewBalanceMsg(*acc)
		require.NoError(t, err)

		// Verify BalanceMsg implements BalanceMsgI
		var _ balance.BalanceMsgI = balMsg

		t.Log("BalanceMsg correctly implements BalanceMsgI interface")
	})
}

func TestBuildSendMsg(t *testing.T) {
	t.Run("Build send message", func(t *testing.T) {
		ctx := context.Background()
		c, err := client.NewClient(ctx, false)
		require.NoError(t, err)

		acc, err := account.NewAccount(c, "testaccount", testMnemonic, testPassword)
		require.NoError(t, err)

		balMsg, err := balance.NewBalanceMsg(*acc)
		require.NoError(t, err)

		amount := sdk.NewCoins(sdk.NewInt64Coin("usix", 1000000))
		destAddr := "six1qyqszqgpqyqszqgpqyqszqgpqyqszqgpqyqszqgpqyqszqgpqyk9rgcg"

		msg, err := balMsg.BuildSendMsg(destAddr, amount)
		require.NoError(t, err, "Should build message without error")
		require.NotNil(t, msg, "Message should not be nil")

		assert.Equal(t, acc.GetCosmosAddress().String(), msg.FromAddress)
		assert.Equal(t, destAddr, msg.ToAddress)
		assert.Equal(t, amount, msg.Amount)

		t.Logf("Built send message from %s to %s", msg.FromAddress, msg.ToAddress)
	})

	t.Run("Build multiple send messages", func(t *testing.T) {
		ctx := context.Background()
		c, err := client.NewClient(ctx, false)
		require.NoError(t, err)

		acc, err := account.NewAccount(c, "testaccount", testMnemonic, testPassword)
		require.NoError(t, err)

		balMsg, err := balance.NewBalanceMsg(*acc)
		require.NoError(t, err)

		amount := sdk.NewCoins(sdk.NewInt64Coin("usix", 1000000))
		dest1 := "six1qyqszqgpqyqszqgpqyqszqgpqyqszqgpqyqszqgpqyqszqgpqyk9rgcg"
		dest2 := "six1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"

		msg1, err := balMsg.BuildSendMsg(dest1, amount)
		require.NoError(t, err)

		msg2, err := balMsg.BuildSendMsg(dest2, amount)
		require.NoError(t, err)

		assert.NotEqual(t, msg1.ToAddress, msg2.ToAddress)
		assert.Equal(t, msg1.FromAddress, msg2.FromAddress)

		t.Log("Successfully built multiple messages for batch transaction")
	})
}

func TestBalanceMsgConfiguration(t *testing.T) {
	t.Run("Configure transaction with fluent API", func(t *testing.T) {
		ctx := context.Background()
		c, err := client.NewClient(ctx, false)
		require.NoError(t, err)

		acc, err := account.NewAccount(c, "testaccount", testMnemonic, testPassword)
		require.NoError(t, err)

		balMsg, err := balance.NewBalanceMsg(*acc)
		require.NoError(t, err)

		// Test fluent API - should return new instances
		balMsg2 := balMsg.WithGas(500000)
		assert.NotNil(t, balMsg2)

		balMsg3 := balMsg2.WithGasAdjustment(1.5)
		assert.NotNil(t, balMsg3)

		balMsg4 := balMsg3.WithMemo("test memo")
		assert.NotNil(t, balMsg4)

		balMsg5 := balMsg4.WithGasPrices("1.25usix")
		assert.NotNil(t, balMsg5)

		balMsg6 := balMsg5.WithFees("1000usix")
		assert.NotNil(t, balMsg6)

		balMsg7 := balMsg6.WithTimeoutHeight(1000000)
		assert.NotNil(t, balMsg7)

		t.Log("Fluent API configuration works correctly")
	})

	t.Run("Chain configuration methods", func(t *testing.T) {
		ctx := context.Background()
		c, err := client.NewClient(ctx, false)
		require.NoError(t, err)

		acc, err := account.NewAccount(c, "testaccount", testMnemonic, testPassword)
		require.NoError(t, err)

		balMsg, err := balance.NewBalanceMsg(*acc)
		require.NoError(t, err)

		// Test method chaining
		configured := balMsg.
			WithGas(500000).
			WithGasAdjustment(1.5).
			WithMemo("chained config")

		assert.NotNil(t, configured)
		t.Log("Method chaining works correctly")
	})
}

func TestBalanceMsgInheritsBalance(t *testing.T) {
	t.Run("BalanceMsg inherits Balance methods", func(t *testing.T) {
		ctx := context.Background()
		c, err := client.NewClient(ctx, false)
		require.NoError(t, err)

		acc, err := account.NewAccount(c, "testaccount", testMnemonic, testPassword)
		require.NoError(t, err)

		balMsg, err := balance.NewBalanceMsg(*acc)
		require.NoError(t, err)

		// These methods are inherited from Balance
		// They should be accessible on BalanceMsg
		retrievedAcc := balMsg.GetAccount()
		assert.Equal(t, acc.GetAccountName(), retrievedAcc.GetAccountName())

		// Note: GetBalance, GetCosmosBalance, GetEVMBalance would make actual network calls
		// so we don't test them here, but they should be accessible

		t.Log("BalanceMsg successfully inherits Balance methods")
	})
}

func TestConstants(t *testing.T) {
	t.Run("Balance constants are defined", func(t *testing.T) {
		assert.Equal(t, "usix", balance.BaseDenom, "Base denom should be usix")
		assert.Equal(t, "asix", balance.EVMDenom, "EVM denom should be asix")

		t.Logf("Base denom: %s", balance.BaseDenom)
		t.Logf("EVM denom: %s", balance.EVMDenom)
	})
}

func TestKeeperPattern(t *testing.T) {
	t.Run("Verify keeper pattern structure", func(t *testing.T) {
		ctx := context.Background()
		c, err := client.NewClient(ctx, false)
		require.NoError(t, err)

		acc, err := account.NewAccount(c, "testaccount", testMnemonic, testPassword)
		require.NoError(t, err)

		// Pattern 1: Use Balance for query-only operations
		bal := balance.NewBalance(*acc)
		assert.NotNil(t, bal)

		// Pattern 2: Use BalanceMsg for both queries and transactions
		balMsg, err := balance.NewBalanceMsg(*acc)
		require.NoError(t, err)
		assert.NotNil(t, balMsg)

		// BalanceMsg should have access to Balance methods (through embedding)
		acc1 := bal.GetAccount()
		acc2 := balMsg.GetAccount()
		assert.Equal(t, acc1.GetAccountName(), acc2.GetAccountName())

		t.Log("Keeper pattern verified:")
		t.Log("  - Balance provides query operations")
		t.Log("  - BalanceMsg embeds Balance and adds transaction operations")
	})
}

// Benchmark tests
func BenchmarkNewBalance(b *testing.B) {
	ctx := context.Background()
	c, _ := client.NewClient(ctx, false)
	acc, _ := account.NewAccount(c, "benchaccount", testMnemonic, testPassword)

	b.ResetTimer()
	for b.Loop() {
		_ = balance.NewBalance(*acc)
	}
}

func BenchmarkNewBalanceMsg(b *testing.B) {
	ctx := context.Background()
	c, _ := client.NewClient(ctx, false)
	acc, _ := account.NewAccount(c, "benchaccount", testMnemonic, testPassword)

	b.ResetTimer()
	for b.Loop() {
		_, _ = balance.NewBalanceMsg(*acc)
	}
}

func BenchmarkBuildSendMsg(b *testing.B) {
	ctx := context.Background()
	c, _ := client.NewClient(ctx, false)
	acc, _ := account.NewAccount(c, "benchaccount", testMnemonic, testPassword)
	balMsg, _ := balance.NewBalanceMsg(*acc)

	amount := sdk.NewCoins(sdk.NewInt64Coin("usix", 1000000))
	destAddr := "six1qyqszqgpqyqszqgpqyqszqgpqyqszqgpqyqszqgpqyqszqgpqyk9rgcg"

	b.ResetTimer()
	for b.Loop() {
		_, _ = balMsg.BuildSendMsg(destAddr, amount)
	}
}

func BenchmarkWithConfiguration(b *testing.B) {
	ctx := context.Background()
	c, _ := client.NewClient(ctx, false)
	acc, _ := account.NewAccount(c, "benchaccount", testMnemonic, testPassword)
	balMsg, _ := balance.NewBalanceMsg(*acc)

	b.ResetTimer()
	for b.Loop() {
		_ = balMsg.
			WithGas(500000).
			WithGasAdjustment(1.5).
			WithMemo("benchmark")
	}
}
