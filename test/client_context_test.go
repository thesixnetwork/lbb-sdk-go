package main_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock structs to test the client and context concepts
type MockCodec interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, any) error
}

type MockInterfaceRegistry interface {
	RegisterImplementations(interface{}, ...interface{})
}

type mockCodec struct{}

func (c *mockCodec) Marshal(v interface{}) ([]byte, error) {
	return []byte("mocked"), nil
}

func (c *mockCodec) Unmarshal(data []byte, v interface{}) error {
	return nil
}

type mockInterfaceRegistry struct{}

func (r *mockInterfaceRegistry) RegisterImplementations(iface interface{}, impls ...interface{}) {}

// Context struct similar to the one in client/context.go
type TestContext struct {
	Context           context.Context
	Codec             MockCodec
	InterfaceRegistry MockInterfaceRegistry
	RPCClient         string
	EVMRPCCleint      string
	APIClient         string
}

// Client struct similar to the one in client/context.go
type TestClient struct {
	TestContext
}

// NewTestContext creates a new TestContext
func NewTestContext(ctx context.Context, rpcClient, evmRPCClient, apiClient string) TestContext {
	return TestContext{
		Context:           ctx,
		Codec:             &mockCodec{},
		InterfaceRegistry: &mockInterfaceRegistry{},
		RPCClient:         rpcClient,
		EVMRPCCleint:      evmRPCClient,
		APIClient:         apiClient,
	}
}

// NewTestClient creates a new TestClient with the provided TestContext
func NewTestClient(ctx TestContext) *TestClient {
	return &TestClient{
		TestContext: ctx,
	}
}

// Test constants using the provided environment variables
const (
	TestRPCClient    = "https://rpc1.fivenet.sixprotocol.net:443"
	TestEVMRPCClient = "https://rpc-evm.fivenet.sixprotocol.net:443"
	TestAPIClient    = "https://api1.fivenet.sixprotocol.net:443"
)

func TestNewTestContext(t *testing.T) {
	tests := []struct {
		name         string
		ctx          context.Context
		rpcClient    string
		evmRPCClient string
		apiClient    string
	}{
		{
			name:         "Valid context creation with all endpoints",
			ctx:          context.Background(),
			rpcClient:    TestRPCClient,
			evmRPCClient: TestEVMRPCClient,
			apiClient:    TestAPIClient,
		},
		{
			name:         "Context creation with empty RPC client",
			ctx:          context.Background(),
			rpcClient:    "",
			evmRPCClient: TestEVMRPCClient,
			apiClient:    TestAPIClient,
		},
		{
			name:         "Context creation with empty EVM RPC client",
			ctx:          context.Background(),
			rpcClient:    TestRPCClient,
			evmRPCClient: "",
			apiClient:    TestAPIClient,
		},
		{
			name:         "Context creation with empty API client",
			ctx:          context.Background(),
			rpcClient:    TestRPCClient,
			evmRPCClient: TestEVMRPCClient,
			apiClient:    "",
		},
		{
			name:         "Context creation with all empty clients",
			ctx:          context.Background(),
			rpcClient:    "",
			evmRPCClient: "",
			apiClient:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientCtx := NewTestContext(tt.ctx, tt.rpcClient, tt.evmRPCClient, tt.apiClient)

			// Verify the context was created
			assert.NotNil(t, clientCtx)
			assert.Equal(t, tt.ctx, clientCtx.Context)
			assert.Equal(t, tt.rpcClient, clientCtx.RPCClient)
			assert.Equal(t, tt.evmRPCClient, clientCtx.EVMRPCCleint)
			assert.Equal(t, tt.apiClient, clientCtx.APIClient)

			// Verify codec and interface registry are properly initialized
			assert.NotNil(t, clientCtx.Codec)
			assert.NotNil(t, clientCtx.InterfaceRegistry)
		})
	}
}

func TestNewTestContextWithNilContext(t *testing.T) {
	clientCtx := NewTestContext(nil, TestRPCClient, TestEVMRPCClient, TestAPIClient)

	assert.NotNil(t, clientCtx)
	assert.Nil(t, clientCtx.Context)
	assert.Equal(t, TestRPCClient, clientCtx.RPCClient)
	assert.Equal(t, TestEVMRPCClient, clientCtx.EVMRPCCleint)
	assert.Equal(t, TestAPIClient, clientCtx.APIClient)
	assert.NotNil(t, clientCtx.Codec)
	assert.NotNil(t, clientCtx.InterfaceRegistry)
}

func TestNewTestClient(t *testing.T) {
	tests := []struct {
		name string
		ctx  TestContext
	}{
		{
			name: "Valid client creation with populated context",
			ctx: NewTestContext(
				context.Background(),
				TestRPCClient,
				TestEVMRPCClient,
				TestAPIClient,
			),
		},
		{
			name: "Client creation with empty endpoints",
			ctx: NewTestContext(
				context.Background(),
				"",
				"",
				"",
			),
		},
		{
			name: "Client creation with nil context",
			ctx: NewTestContext(
				nil,
				TestRPCClient,
				TestEVMRPCClient,
				TestAPIClient,
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewTestClient(tt.ctx)

			// Verify client was created
			require.NotNil(t, client)

			// Verify client contains the correct context
			assert.Equal(t, tt.ctx, client.TestContext)

			// Verify embedded context fields are accessible
			assert.Equal(t, tt.ctx.Context, client.TestContext.Context)
			assert.Equal(t, tt.ctx.Codec, client.TestContext.Codec)
			assert.Equal(t, tt.ctx.InterfaceRegistry, client.TestContext.InterfaceRegistry)
			assert.Equal(t, tt.ctx.RPCClient, client.TestContext.RPCClient)
			assert.Equal(t, tt.ctx.EVMRPCCleint, client.TestContext.EVMRPCCleint)
			assert.Equal(t, tt.ctx.APIClient, client.TestContext.APIClient)
		})
	}
}

func TestContextFieldAssignment(t *testing.T) {
	ctx := context.Background()
	clientCtx := NewTestContext(ctx, TestRPCClient, TestEVMRPCClient, TestAPIClient)

	t.Run("Context field assignment validation", func(t *testing.T) {
		// Test that Context field is properly assigned
		assert.Equal(t, ctx, clientCtx.Context)

		// Test endpoint assignments
		assert.Equal(t, TestRPCClient, clientCtx.RPCClient)
		assert.Equal(t, TestEVMRPCClient, clientCtx.EVMRPCCleint)
		assert.Equal(t, TestAPIClient, clientCtx.APIClient)

		// Test that codec-related fields are not nil
		assert.NotNil(t, clientCtx.Codec)
		assert.NotNil(t, clientCtx.InterfaceRegistry)
	})
}

func TestClientEmbeddedContextAccess(t *testing.T) {
	ctx := NewTestContext(
		context.Background(),
		TestRPCClient,
		TestEVMRPCClient,
		TestAPIClient,
	)
	client := NewTestClient(ctx)

	t.Run("Client embedded context field access", func(t *testing.T) {
		// Test that client has direct access to all context fields through embedding
		assert.Equal(t, ctx.Context, client.TestContext.Context)
		assert.Equal(t, ctx.Codec, client.TestContext.Codec)
		assert.Equal(t, ctx.InterfaceRegistry, client.TestContext.InterfaceRegistry)
		assert.Equal(t, ctx.RPCClient, client.TestContext.RPCClient)
		assert.Equal(t, ctx.EVMRPCCleint, client.TestContext.EVMRPCCleint)
		assert.Equal(t, ctx.APIClient, client.TestContext.APIClient)

		// Test that the entire context is properly embedded
		assert.Equal(t, ctx, client.TestContext)
	})
}

func TestContextZeroValues(t *testing.T) {
	t.Run("Context with zero values", func(t *testing.T) {
		var zeroCtx TestContext
		client := NewTestClient(zeroCtx)

		assert.NotNil(t, client)
		assert.Equal(t, zeroCtx, client.TestContext)
		assert.Nil(t, client.TestContext.Context)
		assert.Nil(t, client.TestContext.Codec)
		assert.Nil(t, client.TestContext.InterfaceRegistry)
		assert.Empty(t, client.TestContext.RPCClient)
		assert.Empty(t, client.TestContext.EVMRPCCleint)
		assert.Empty(t, client.TestContext.APIClient)
	})
}

func TestContextWithCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clientCtx := NewTestContext(ctx, TestRPCClient, TestEVMRPCClient, TestAPIClient)
	client := NewTestClient(clientCtx)

	// Verify the cancellable context is properly stored
	assert.Equal(t, ctx, client.TestContext.Context)

	// Cancel the context and verify it's cancelled
	cancel()
	select {
	case <-client.TestContext.Context.Done():
		// Context is properly cancelled
		assert.Error(t, client.TestContext.Context.Err())
	default:
		t.Error("Context should be cancelled")
	}
}

func TestContextWithValues(t *testing.T) {
	type contextKey string
	const testKey contextKey = "testKey"
	const testValue = "testValue"

	ctx := context.WithValue(context.Background(), testKey, testValue)
	clientCtx := NewTestContext(ctx, TestRPCClient, TestEVMRPCClient, TestAPIClient)
	client := NewTestClient(clientCtx)

	// Verify the context value is accessible through the client
	value := client.TestContext.Context.Value(testKey)
	assert.Equal(t, testValue, value)
}

func TestMockCodecFunctionality(t *testing.T) {
	clientCtx := NewTestContext(context.Background(), TestRPCClient, TestEVMRPCClient, TestAPIClient)

	t.Run("Test mock codec marshal", func(t *testing.T) {
		data, err := clientCtx.Codec.Marshal("test")
		assert.NoError(t, err)
		assert.Equal(t, []byte("mocked"), data)
	})

	t.Run("Test mock codec unmarshal", func(t *testing.T) {
		var result string
		err := clientCtx.Codec.Unmarshal([]byte("test"), &result)
		assert.NoError(t, err)
	})
}

// Benchmark tests
func BenchmarkNewTestContext(b *testing.B) {
	ctx := context.Background()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		NewTestContext(ctx, TestRPCClient, TestEVMRPCClient, TestAPIClient)
	}
}

func BenchmarkNewTestClient(b *testing.B) {
	ctx := NewTestContext(
		context.Background(),
		TestRPCClient,
		TestEVMRPCClient,
		TestAPIClient,
	)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		NewTestClient(ctx)
	}
}

// Integration test with environment variables simulation
func TestEnvironmentVariableSimulation(t *testing.T) {
	// Simulated environment variables
	envVars := map[string]string{
		"EVM_FIVENET_RPC": "https://rpc-evm.fivenet.sixprotocol.net:443",
		"FIVENET_RPC":     "https://rpc1.fivenet.sixprotocol.net:443",
		"FIVENET_API":     "https://api1.fivenet.sixprotocol.net:443",
	}

	t.Run("Create context and client with env var values", func(t *testing.T) {
		ctx := context.Background()
		clientCtx := NewTestContext(
			ctx,
			envVars["FIVENET_RPC"],
			envVars["EVM_FIVENET_RPC"],
			envVars["FIVENET_API"],
		)
		client := NewTestClient(clientCtx)

		// Verify the client was created with environment variable values
		assert.Equal(t, envVars["FIVENET_RPC"], client.TestContext.RPCClient)
		assert.Equal(t, envVars["EVM_FIVENET_RPC"], client.TestContext.EVMRPCCleint)
		assert.Equal(t, envVars["FIVENET_API"], client.TestContext.APIClient)
	})
}

// Example usage test
func TestExampleUsage(t *testing.T) {
	t.Run("Example: Create client for blockchain operations", func(t *testing.T) {
		// Step 1: Create a context
		ctx := context.Background()

		// Step 2: Create a client context with network endpoints
		clientCtx := NewTestContext(
			ctx,
			"https://rpc1.fivenet.sixprotocol.net:443",    // RPC endpoint
			"https://rpc-evm.fivenet.sixprotocol.net:443", // EVM RPC endpoint
			"https://api1.fivenet.sixprotocol.net:443",    // API endpoint
		)

		// Step 3: Create the client
		client := NewTestClient(clientCtx)

		// Step 4: Verify client is ready for operations
		assert.NotNil(t, client)
		assert.NotNil(t, client.TestContext.Codec)
		assert.NotNil(t, client.TestContext.InterfaceRegistry)
		assert.NotEmpty(t, client.TestContext.RPCClient)
		assert.NotEmpty(t, client.TestContext.EVMRPCCleint)
		assert.NotEmpty(t, client.TestContext.APIClient)
	})
}
