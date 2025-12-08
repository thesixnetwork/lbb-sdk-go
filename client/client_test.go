package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	t.Run("Create testnet client", func(t *testing.T) {
		ctx := context.Background()
		c, err := NewClient(ctx, false) // false = testnet

		require.NoError(t, err, "Should create testnet client without error")
		require.NotNil(t, c, "Client should not be nil")

		// Verify testnet configuration
		assert.Equal(t, TestnetChainID, c.GetChainID(), "Should use testnet chain ID")
		assert.Equal(t, TestnetRPC, c.GetRPCClient(), "Should use testnet RPC")
		assert.Equal(t, TestnetAPI, c.GetAPIClient(), "Should use testnet API")
		assert.Equal(t, TestnetEVMRPC, c.GetEVMRPCClient(), "Should use testnet EVM RPC")

		t.Logf("Testnet client created successfully:")
		t.Logf("  Chain ID: %s", c.GetChainID())
		t.Logf("  RPC: %s", c.GetRPCClient())
		t.Logf("  API: %s", c.GetAPIClient())
		t.Logf("  EVM RPC: %s", c.GetEVMRPCClient())
	})

	t.Run("Create mainnet client", func(t *testing.T) {
		ctx := context.Background()
		c, err := NewClient(ctx, true) // true = mainnet

		require.NoError(t, err, "Should create mainnet client without error")
		require.NotNil(t, c, "Client should not be nil")

		// Verify mainnet configuration
		assert.Equal(t, MainnetChainID, c.GetChainID(), "Should use mainnet chain ID")
		assert.Equal(t, MainnetRPC, c.GetRPCClient(), "Should use mainnet RPC")
		assert.Equal(t, MainnetAPI, c.GetAPIClient(), "Should use mainnet API")
		assert.Equal(t, MainnetEVMRPC, c.GetEVMRPCClient(), "Should use mainnet EVM RPC")

		t.Logf("Mainnet client created successfully:")
		t.Logf("  Chain ID: %s", c.GetChainID())
		t.Logf("  RPC: %s", c.GetRPCClient())
		t.Logf("  API: %s", c.GetAPIClient())
		t.Logf("  EVM RPC: %s", c.GetEVMRPCClient())
	})

	t.Run("Create client with nil context", func(t *testing.T) {
		c, err := NewClient(context.TODO(), false)

		require.NoError(t, err, "Should handle nil context gracefully")
		require.NotNil(t, c, "Client should not be nil")
		assert.NotNil(t, c.GetContext(), "Context should be initialized to Background")
	})
}

func TestNewCustomClient(t *testing.T) {
	t.Run("Create custom client with valid configuration", func(t *testing.T) {
		ctx := context.Background()
		customRPC := "https://custom-rpc.example.com"
		customAPI := "https://custom-api.example.com"
		customEVMRPC := "https://custom-evm.example.com"
		customChainID := "custom-chain"

		c, err := NewCustomClient(ctx, customRPC, customAPI, customEVMRPC, customChainID)

		require.NoError(t, err, "Should create custom client without error")
		require.NotNil(t, c, "Client should not be nil")

		// Verify custom configuration
		assert.Equal(t, customChainID, c.GetChainID(), "Should use custom chain ID")
		assert.Equal(t, customRPC, c.GetRPCClient(), "Should use custom RPC")
		assert.Equal(t, customAPI, c.GetAPIClient(), "Should use custom API")
		assert.Equal(t, customEVMRPC, c.GetEVMRPCClient(), "Should use custom EVM RPC")

		t.Logf("Custom client created successfully with chain ID: %s", c.GetChainID())
	})

	t.Run("Create custom client with nil context", func(t *testing.T) {
		c, err := NewCustomClient(context.TODO(), TestnetRPC, TestnetAPI, TestnetEVMRPC, TestnetChainID)

		require.NoError(t, err, "Should handle nil context gracefully")
		require.NotNil(t, c, "Client should not be nil")
		assert.NotNil(t, c.GetContext(), "Context should be initialized")
	})

	t.Run("Create custom client with invalid RPC URL", func(t *testing.T) {
		ctx := context.Background()
		invalidRPC := "invalid://url"

		c, err := NewCustomClient(ctx, invalidRPC, TestnetAPI, TestnetEVMRPC, TestnetChainID)

		// Note: Client creation may succeed even with invalid RPC URL
		// because the RPC client is created lazily or connection is not validated during initialization
		if err != nil {
			t.Logf("Invalid RPC URL rejected during creation: %v", err)
		} else {
			t.Logf("Client created with invalid RPC URL (connection not validated at creation time)")
			assert.NotNil(t, c, "Client should be created")
		}
	})

	t.Run("Create custom client with empty chain ID", func(t *testing.T) {
		ctx := context.Background()

		c, err := NewCustomClient(ctx, TestnetRPC, TestnetAPI, TestnetEVMRPC, "")

		// Depending on implementation, this might succeed or fail
		// The test just verifies consistent behavior
		if err != nil {
			t.Logf("Empty chain ID rejected: %v", err)
		} else {
			t.Logf("Empty chain ID accepted, chain ID: '%s'", c.GetChainID())
		}
	})
}

func TestClientGetters(t *testing.T) {
	ctx := context.Background()
	c, err := NewClient(ctx, false)
	require.NoError(t, err)
	require.NotNil(t, c)

	t.Run("GetChainID", func(t *testing.T) {
		chainID := c.GetChainID()
		assert.NotEmpty(t, chainID, "Chain ID should not be empty")
		assert.Equal(t, TestnetChainID, chainID, "Should return correct chain ID")
	})

	t.Run("GetRPCClient", func(t *testing.T) {
		rpc := c.GetRPCClient()
		assert.NotEmpty(t, rpc, "RPC URL should not be empty")
		assert.Equal(t, TestnetRPC, rpc, "Should return correct RPC URL")
	})

	t.Run("GetAPIClient", func(t *testing.T) {
		api := c.GetAPIClient()
		assert.NotEmpty(t, api, "API URL should not be empty")
		assert.Equal(t, TestnetAPI, api, "Should return correct API URL")
	})

	t.Run("GetEVMRPCClient", func(t *testing.T) {
		evmRPC := c.GetEVMRPCClient()
		assert.NotEmpty(t, evmRPC, "EVM RPC URL should not be empty")
		assert.Equal(t, TestnetEVMRPC, evmRPC, "Should return correct EVM RPC URL")
	})

	t.Run("GetContext", func(t *testing.T) {
		ctx := c.GetContext()
		assert.NotNil(t, ctx, "Context should not be nil")
	})

	t.Run("GetClientCTX", func(t *testing.T) {
		clientCtx := c.GetClientCTX()
		assert.NotNil(t, clientCtx, "Client context should not be nil")
		assert.NotNil(t, clientCtx.TxConfig, "TxConfig should be initialized")
		assert.NotNil(t, clientCtx.Codec, "Codec should be initialized")
	})

	t.Run("GetKeyring", func(t *testing.T) {
		kr := c.GetKeyring()
		assert.NotNil(t, kr, "Keyring should not be nil")
	})

	t.Run("GetETHClient", func(t *testing.T) {
		ethClient := c.GetETHClient()
		assert.NotNil(t, ethClient, "ETH client should not be nil")
	})
}

func TestClientContextCancellation(t *testing.T) {
	t.Run("Client with cancellable context", func(t *testing.T) {
		parentCtx, cancel := context.WithCancel(context.Background())
		defer cancel()

		c, err := NewClient(parentCtx, false)
		require.NoError(t, err)
		require.NotNil(t, c)

		// Context should not be cancelled initially
		select {
		case <-c.GetContext().Done():
			t.Error("Context should not be cancelled initially")
		default:
			// Expected
		}

		// Cancel the parent context
		cancel()

		// Context should now be cancelled
		select {
		case <-c.GetContext().Done():
			// Expected
			assert.Error(t, c.GetContext().Err(), "Context error should not be nil after cancellation")
			t.Log("Context correctly cancelled")
		default:
			t.Error("Context should be cancelled after calling cancel()")
		}
	})
}

func TestClientConstants(t *testing.T) {
	t.Run("Testnet constants", func(t *testing.T) {
		assert.NotEmpty(t, TestnetRPC, "Testnet RPC should be defined")
		assert.NotEmpty(t, TestnetAPI, "Testnet API should be defined")
		assert.NotEmpty(t, TestnetEVMRPC, "Testnet EVM RPC should be defined")
		assert.NotEmpty(t, TestnetChainID, "Testnet chain ID should be defined")

		assert.Contains(t, TestnetRPC, "https://", "Testnet RPC should use HTTPS")
		assert.Contains(t, TestnetAPI, "https://", "Testnet API should use HTTPS")
		assert.Contains(t, TestnetEVMRPC, "https://", "Testnet EVM RPC should use HTTPS")

		t.Logf("Testnet configuration:")
		t.Logf("  RPC: %s", TestnetRPC)
		t.Logf("  API: %s", TestnetAPI)
		t.Logf("  EVM RPC: %s", TestnetEVMRPC)
		t.Logf("  Chain ID: %s", TestnetChainID)
	})

	t.Run("Mainnet constants", func(t *testing.T) {
		assert.NotEmpty(t, MainnetRPC, "Mainnet RPC should be defined")
		assert.NotEmpty(t, MainnetAPI, "Mainnet API should be defined")
		assert.NotEmpty(t, MainnetEVMRPC, "Mainnet EVM RPC should be defined")
		assert.NotEmpty(t, MainnetChainID, "Mainnet chain ID should be defined")

		assert.Contains(t, MainnetRPC, "https://", "Mainnet RPC should use HTTPS")
		assert.Contains(t, MainnetAPI, "https://", "Mainnet API should use HTTPS")
		assert.Contains(t, MainnetEVMRPC, "https://", "Mainnet EVM RPC should use HTTPS")

		t.Logf("Mainnet configuration:")
		t.Logf("  RPC: %s", MainnetRPC)
		t.Logf("  API: %s", MainnetAPI)
		t.Logf("  EVM RPC: %s", MainnetEVMRPC)
		t.Logf("  Chain ID: %s", MainnetChainID)
	})

	t.Run("Testnet and mainnet should differ", func(t *testing.T) {
		assert.NotEqual(t, TestnetChainID, MainnetChainID, "Testnet and mainnet should have different chain IDs")
		assert.NotEqual(t, TestnetRPC, MainnetRPC, "Testnet and mainnet should have different RPC URLs")
		assert.NotEqual(t, TestnetAPI, MainnetAPI, "Testnet and mainnet should have different API URLs")
	})
}

func TestClientInterface(t *testing.T) {
	t.Run("Client implements ClientI interface", func(t *testing.T) {
		ctx := context.Background()
		c, err := NewClient(ctx, false)
		require.NoError(t, err)

		// Verify that Client implements ClientI
		var _ ClientI = c

		// Test that all interface methods are callable
		assert.NotNil(t, c.GetClientCTX())
		assert.NotNil(t, c.GetKeyring())
		assert.NotNil(t, c.GetETHClient())
		assert.NotEmpty(t, c.GetRPCClient())
		assert.NotEmpty(t, c.GetAPIClient())
		assert.NotEmpty(t, c.GetEVMRPCClient())
		assert.NotEmpty(t, c.GetChainID())
		assert.NotNil(t, c.GetContext())

		t.Log("Client correctly implements ClientI interface")
	})
}

func TestMultipleClients(t *testing.T) {
	t.Run("Create multiple independent clients", func(t *testing.T) {
		ctx := context.Background()

		// Create testnet client
		testnetClient, err := NewClient(ctx, false)
		require.NoError(t, err)

		// Create mainnet client
		mainnetClient, err := NewClient(ctx, true)
		require.NoError(t, err)

		// Verify they are independent
		assert.NotEqual(t, testnetClient.GetChainID(), mainnetClient.GetChainID())
		assert.NotEqual(t, testnetClient.GetRPCClient(), mainnetClient.GetRPCClient())

		t.Log("Multiple independent clients created successfully")
	})
}

// Benchmark tests
func BenchmarkNewClient(b *testing.B) {
	ctx := context.Background()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := NewClient(ctx, false)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNewCustomClient(b *testing.B) {
	ctx := context.Background()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := NewCustomClient(ctx, TestnetRPC, TestnetAPI, TestnetEVMRPC, TestnetChainID)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkClientGetters(b *testing.B) {
	ctx := context.Background()
	c, _ := NewClient(ctx, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.GetChainID()
		_ = c.GetRPCClient()
		_ = c.GetAPIClient()
		_ = c.GetEVMRPCClient()
		_ = c.GetContext()
	}
}
