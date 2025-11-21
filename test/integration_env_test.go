package main_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Integration test that uses actual environment variables
func TestIntegrationWithActualEnvVars(t *testing.T) {
	// Get actual environment variables
	evmRPC := os.Getenv("EVM_FIVENET_RPC")
	fivenetRPC := os.Getenv("FIVENET_RPC")
	fivenetAPI := os.Getenv("FIVENET_API")

	t.Run("Test with actual environment variables", func(t *testing.T) {
		// Skip test if environment variables are not set
		if evmRPC == "" || fivenetRPC == "" || fivenetAPI == "" {
			t.Skip("Environment variables EVM_FIVENET_RPC, FIVENET_RPC, or FIVENET_API not set")
		}

		// Create context using actual environment variables
		ctx := context.Background()
		clientCtx := NewTestContext(ctx, fivenetRPC, evmRPC, fivenetAPI)

		// Verify context creation with actual values
		require.NotNil(t, clientCtx)
		assert.Equal(t, ctx, clientCtx.Context)
		assert.Equal(t, fivenetRPC, clientCtx.RPCClient)
		assert.Equal(t, evmRPC, clientCtx.EVMRPCCleint)
		assert.Equal(t, fivenetAPI, clientCtx.APIClient)
		assert.NotNil(t, clientCtx.Codec)
		assert.NotNil(t, clientCtx.InterfaceRegistry)

		// Create client
		client := NewTestClient(clientCtx)
		require.NotNil(t, client)
		assert.Equal(t, clientCtx, client.TestContext)

		// Verify client has access to actual environment variable values
		assert.Equal(t, fivenetRPC, client.TestContext.RPCClient)
		assert.Equal(t, evmRPC, client.TestContext.EVMRPCCleint)
		assert.Equal(t, fivenetAPI, client.TestContext.APIClient)

		// Log the actual values being tested (for debugging)
		t.Logf("Testing with FIVENET_RPC: %s", fivenetRPC)
		t.Logf("Testing with EVM_FIVENET_RPC: %s", evmRPC)
		t.Logf("Testing with FIVENET_API: %s", fivenetAPI)
	})

	t.Run("Test with expected environment variable values", func(t *testing.T) {
		// Expected values based on provided environment variables
		expectedEVMRPC := "https://rpc-evm.fivenet.sixprotocol.net:443"
		expectedFivenetRPC := "https://rpc1.fivenet.sixprotocol.net:443"
		expectedFivenetAPI := "https://api1.fivenet.sixprotocol.net:443"

		// Only run this test if environment variables are set
		if evmRPC == "" || fivenetRPC == "" || fivenetAPI == "" {
			t.Skip("Environment variables not set")
		}

		// Verify the actual environment variables match expected values
		assert.Equal(t, expectedEVMRPC, evmRPC, "EVM_FIVENET_RPC should match expected value")
		assert.Equal(t, expectedFivenetRPC, fivenetRPC, "FIVENET_RPC should match expected value")
		assert.Equal(t, expectedFivenetAPI, fivenetAPI, "FIVENET_API should match expected value")

		// Create context and client with verified values
		ctx := context.Background()
		clientCtx := NewTestContext(ctx, fivenetRPC, evmRPC, fivenetAPI)
		client := NewTestClient(clientCtx)

		// Verify exact values in client
		assert.Equal(t, expectedFivenetRPC, client.TestContext.RPCClient)
		assert.Equal(t, expectedEVMRPC, client.TestContext.EVMRPCCleint)
		assert.Equal(t, expectedFivenetAPI, client.TestContext.APIClient)
	})
}

// Test demonstrating how to use environment variables in a real application
func TestRealWorldEnvVarUsage(t *testing.T) {
	t.Run("Application startup with environment variables", func(t *testing.T) {
		// Simulate application startup that reads from environment
		evmRPC := os.Getenv("EVM_FIVENET_RPC")
		fivenetRPC := os.Getenv("FIVENET_RPC")
		fivenetAPI := os.Getenv("FIVENET_API")

		// Use fallback values if environment variables are not set (for testing)
		if evmRPC == "" {
			evmRPC = "https://rpc-evm.fivenet.sixprotocol.net:443"
			t.Log("Using fallback value for EVM_FIVENET_RPC")
		}
		if fivenetRPC == "" {
			fivenetRPC = "https://rpc1.fivenet.sixprotocol.net:443"
			t.Log("Using fallback value for FIVENET_RPC")
		}
		if fivenetAPI == "" {
			fivenetAPI = "https://api1.fivenet.sixprotocol.net:443"
			t.Log("Using fallback value for FIVENET_API")
		}

		// Initialize application context
		ctx := context.Background()
		clientCtx := NewTestContext(ctx, fivenetRPC, evmRPC, fivenetAPI)
		client := NewTestClient(clientCtx)

		// Verify client is properly initialized for application use
		require.NotNil(t, client)
		require.NotNil(t, client.TestContext.Codec)
		require.NotNil(t, client.TestContext.InterfaceRegistry)

		// Verify network endpoints are configured
		assert.NotEmpty(t, client.TestContext.RPCClient)
		assert.NotEmpty(t, client.TestContext.EVMRPCCleint)
		assert.NotEmpty(t, client.TestContext.APIClient)

		// Verify URLs are valid format (basic validation)
		assert.Contains(t, client.TestContext.RPCClient, "https://")
		assert.Contains(t, client.TestContext.EVMRPCCleint, "https://")
		assert.Contains(t, client.TestContext.APIClient, "https://")
	})
}

// Test with context timeout using environment variables
func TestContextWithTimeoutAndEnvVars(t *testing.T) {
	evmRPC := os.Getenv("EVM_FIVENET_RPC")
	fivenetRPC := os.Getenv("FIVENET_RPC")
	fivenetAPI := os.Getenv("FIVENET_API")

	// Use fallback values if environment variables are not set
	if evmRPC == "" {
		evmRPC = "https://rpc-evm.fivenet.sixprotocol.net:443"
	}
	if fivenetRPC == "" {
		fivenetRPC = "https://rpc1.fivenet.sixprotocol.net:443"
	}
	if fivenetAPI == "" {
		fivenetAPI = "https://api1.fivenet.sixprotocol.net:443"
	}

	t.Run("Context with timeout and environment variables", func(t *testing.T) {
		// Create context with timeout for real-world scenario
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Create client context with environment variables
		clientCtx := NewTestContext(ctx, fivenetRPC, evmRPC, fivenetAPI)
		client := NewTestClient(clientCtx)

		// Verify context with timeout is properly stored
		require.NotNil(t, client.TestContext.Context)

		// Check that context has deadline
		deadline, ok := client.TestContext.Context.Deadline()
		assert.True(t, ok, "Context should have a deadline")
		assert.False(t, deadline.IsZero(), "Deadline should not be zero")

		// Verify timeout is reasonable (should be around 30 seconds from now)
		assert.True(t, time.Until(deadline) <= 30*time.Second)
		assert.True(t, time.Until(deadline) > 29*time.Second)
	})
}

// Test environment variable validation
func TestEnvironmentVariableValidation(t *testing.T) {
	t.Run("Validate environment variable format", func(t *testing.T) {
		evmRPC := os.Getenv("EVM_FIVENET_RPC")
		fivenetRPC := os.Getenv("FIVENET_RPC")
		fivenetAPI := os.Getenv("FIVENET_API")

		// Skip validation if environment variables are not set
		if evmRPC == "" || fivenetRPC == "" || fivenetAPI == "" {
			t.Skip("Environment variables not set, skipping validation")
		}

		// Validate URL format
		assert.Contains(t, evmRPC, "https://", "EVM_FIVENET_RPC should use HTTPS")
		assert.Contains(t, fivenetRPC, "https://", "FIVENET_RPC should use HTTPS")
		assert.Contains(t, fivenetAPI, "https://", "FIVENET_API should use HTTPS")

		// Validate domain
		assert.Contains(t, evmRPC, "fivenet.sixprotocol.net", "EVM_FIVENET_RPC should use fivenet.sixprotocol.net domain")
		assert.Contains(t, fivenetRPC, "fivenet.sixprotocol.net", "FIVENET_RPC should use fivenet.sixprotocol.net domain")
		assert.Contains(t, fivenetAPI, "fivenet.sixprotocol.net", "FIVENET_API should use fivenet.sixprotocol.net domain")

		// Validate port
		assert.Contains(t, evmRPC, ":443", "EVM_FIVENET_RPC should use port 443")
		assert.Contains(t, fivenetRPC, ":443", "FIVENET_RPC should use port 443")
		assert.Contains(t, fivenetAPI, ":443", "FIVENET_API should use port 443")
	})
}

// Benchmark with actual environment variables
func BenchmarkClientCreationWithActualEnvVars(b *testing.B) {
	evmRPC := os.Getenv("EVM_FIVENET_RPC")
	fivenetRPC := os.Getenv("FIVENET_RPC")
	fivenetAPI := os.Getenv("FIVENET_API")

	// Use fallback values if environment variables are not set
	if evmRPC == "" {
		evmRPC = "https://rpc-evm.fivenet.sixprotocol.net:443"
	}
	if fivenetRPC == "" {
		fivenetRPC = "https://rpc1.fivenet.sixprotocol.net:443"
	}
	if fivenetAPI == "" {
		fivenetAPI = "https://api1.fivenet.sixprotocol.net:443"
	}

	ctx := context.Background()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		clientCtx := NewTestContext(ctx, fivenetRPC, evmRPC, fivenetAPI)
		NewTestClient(clientCtx)
	}
}

// Test error handling when environment variables are missing or invalid
func TestEnvironmentVariableErrorHandling(t *testing.T) {
	// Temporarily clear environment variables
	originalEVMRPC := os.Getenv("EVM_FIVENET_RPC")
	originalFivenetRPC := os.Getenv("FIVENET_RPC")
	originalFivenetAPI := os.Getenv("FIVENET_API")

	os.Unsetenv("EVM_FIVENET_RPC")
	os.Unsetenv("FIVENET_RPC")
	os.Unsetenv("FIVENET_API")

	// Cleanup after test
	defer func() {
		if originalEVMRPC != "" {
			os.Setenv("EVM_FIVENET_RPC", originalEVMRPC)
		}
		if originalFivenetRPC != "" {
			os.Setenv("FIVENET_RPC", originalFivenetRPC)
		}
		if originalFivenetAPI != "" {
			os.Setenv("FIVENET_API", originalFivenetAPI)
		}
	}()

	t.Run("Handle missing environment variables gracefully", func(t *testing.T) {
		evmRPC := os.Getenv("EVM_FIVENET_RPC")
		fivenetRPC := os.Getenv("FIVENET_RPC")
		fivenetAPI := os.Getenv("FIVENET_API")

		// All should be empty now
		assert.Empty(t, evmRPC)
		assert.Empty(t, fivenetRPC)
		assert.Empty(t, fivenetAPI)

		// Should still be able to create context and client with empty values
		ctx := context.Background()
		clientCtx := NewTestContext(ctx, fivenetRPC, evmRPC, fivenetAPI)
		client := NewTestClient(clientCtx)

		// Client should be created but with empty endpoint values
		require.NotNil(t, client)
		assert.Empty(t, client.TestContext.RPCClient)
		assert.Empty(t, client.TestContext.EVMRPCCleint)
		assert.Empty(t, client.TestContext.APIClient)

		// Codec and interface registry should still be initialized
		assert.NotNil(t, client.TestContext.Codec)
		assert.NotNil(t, client.TestContext.InterfaceRegistry)
	})
}
