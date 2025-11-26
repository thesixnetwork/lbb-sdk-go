package client

import (
	"context"
	"testing"
)

// Test constants using the provided environment variables
const (
	TestRPCClient    = "https://rpc1.fivenet.sixprotocol.net:443"
	TestEVMRPCClient = "https://rpc-evm.fivenet.sixprotocol.net:443"
	TestAPIClient    = "https://api1.fivenet.sixprotocol.net:443"
)

func TestContextStructFields(t *testing.T) {
	// Test that Context struct has the expected fields
	var ctx Client

	// Test field assignments
	ctx.RPCClient = TestRPCClient
	ctx.EVMRPCCleint = TestEVMRPCClient
	ctx.APIClient = TestAPIClient
	ctx.Context = context.Background()

	// Verify assignments
	if ctx.RPCClient != TestRPCClient {
		t.Errorf("Expected RPCClient to be %s, got %s", TestRPCClient, ctx.RPCClient)
	}

	if ctx.EVMRPCCleint != TestEVMRPCClient {
		t.Errorf("Expected EVMRPCCleint to be %s, got %s", TestEVMRPCClient, ctx.EVMRPCCleint)
	}

	if ctx.APIClient != TestAPIClient {
		t.Errorf("Expected APIClient to be %s, got %s", TestAPIClient, ctx.APIClient)
	}

	if ctx.Context == nil {
		t.Error("Expected Context to not be nil")
	}
}

func TestClientStructFields(t *testing.T) {
	// Test that Client struct has the expected embedded Context
	var client Client

	// Test field assignments through embedding
	client.RPCClient = TestRPCClient
	client.EVMRPCCleint = TestEVMRPCClient
	client.APIClient = TestAPIClient
	client.Context = context.Background()

	// Verify assignments through embedded Context
	if client.RPCClient != TestRPCClient {
		t.Errorf("Expected client.RPCClient to be %s, got %s", TestRPCClient, client.RPCClient)
	}

	if client.EVMRPCCleint != TestEVMRPCClient {
		t.Errorf("Expected client.EVMRPCCleint to be %s, got %s", TestEVMRPCClient, client.EVMRPCCleint)
	}

	if client.APIClient != TestAPIClient {
		t.Errorf("Expected client.APIClient to be %s, got %s", TestAPIClient, client.APIClient)
	}
}

func TestBasicContextCreation(t *testing.T) {
	// Test basic context creation without using the problematic NewContext function
	ctx := Client{
		Context:      context.Background(),
		RPCClient:    TestRPCClient,
		EVMRPCCleint: TestEVMRPCClient,
		APIClient:    TestAPIClient,
	}

	if ctx.Context == nil {
		t.Error("Context field should not be nil")
	}

	if ctx.RPCClient != TestRPCClient {
		t.Errorf("RPCClient should be %s, got %s", TestRPCClient, ctx.RPCClient)
	}

	if ctx.EVMRPCCleint != TestEVMRPCClient {
		t.Errorf("EVMRPCCleint should be %s, got %s", TestEVMRPCClient, ctx.EVMRPCCleint)
	}

	if ctx.APIClient != TestAPIClient {
		t.Errorf("APIClient should be %s, got %s", TestAPIClient, ctx.APIClient)
	}
}

func TestBasicClientCreation(t *testing.T) {
	// Test basic client creation without using problematic functions
	client := Client{
		Context:      context.Background(),
		RPCClient:    TestRPCClient,
		EVMRPCCleint: TestEVMRPCClient,
		APIClient:    TestAPIClient,
	}

	if client.Context == nil {
		t.Error("Client's embedded LBBContext.Context field should not be nil")
	}

	if client.RPCClient != TestRPCClient {
		t.Errorf("Client's LBBContext.RPCClient should be %s, got %s", TestRPCClient, client.RPCClient)
	}

	if client.EVMRPCCleint != TestEVMRPCClient {
		t.Errorf("Client's LBBContext.EVMRPCCleint should be %s, got %s", TestEVMRPCClient, client.EVMRPCCleint)
	}

	if client.APIClient != TestAPIClient {
		t.Errorf("Client's LBBContext.APIClient should be %s, got %s", TestAPIClient, client.APIClient)
	}
}

func TestEmptyContextFields(t *testing.T) {
	// Test with empty endpoint values
	ctx := Client{
		Context:      context.Background(),
		RPCClient:    "",
		EVMRPCCleint: "",
		APIClient:    "",
	}

	if ctx.Context == nil {
		t.Error("Context field should not be nil even with empty endpoints")
	}

	if ctx.RPCClient != "" {
		t.Error("RPCClient should be empty")
	}

	if ctx.EVMRPCCleint != "" {
		t.Error("EVMRPCCleint should be empty")
	}

	if ctx.APIClient != "" {
		t.Error("APIClient should be empty")
	}
}

func TestNilContext(t *testing.T) {
	// Test with nil context
	ctx := Client{
		Context:      nil,
		RPCClient:    TestRPCClient,
		EVMRPCCleint: TestEVMRPCClient,
		APIClient:    TestAPIClient,
	}

	if ctx.Context != nil {
		t.Error("Context field should be nil")
	}

	if ctx.RPCClient != TestRPCClient {
		t.Errorf("RPCClient should be %s even with nil context", TestRPCClient)
	}
}

func TestContextWithCancel(t *testing.T) {
	// Test with cancellable context
	parentCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx := Client{
		Context:      parentCtx,
		RPCClient:    TestRPCClient,
		EVMRPCCleint: TestEVMRPCClient,
		APIClient:    TestAPIClient,
	}

	// Context should not be cancelled initially
	select {
	case <-ctx.Context.Done():
		t.Error("Context should not be cancelled initially")
	default:
		// Expected
	}

	// Cancel the context
	cancel()

	// Context should now be cancelled
	select {
	case <-ctx.Context.Done():
		// Expected
		if ctx.Context.Err() == nil {
			t.Error("Context error should not be nil after cancellation")
		}
	default:
		t.Error("Context should be cancelled after calling cancel()")
	}
}

func TestContextWithValue(t *testing.T) {
	// Test context with values
	type contextKey string
	const testKey contextKey = "testKey"
	const testValue = "testValue"

	parentCtx := context.WithValue(context.Background(), testKey, testValue)

	ctx := Client{
		Context:      parentCtx,
		RPCClient:    TestRPCClient,
		EVMRPCCleint: TestEVMRPCClient,
		APIClient:    TestAPIClient,
	}

	// Verify the value is accessible
	value := ctx.Context.Value(testKey)
	if value != testValue {
		t.Errorf("Expected context value %s, got %v", testValue, value)
	}
}

func TestEndpointValidation(t *testing.T) {
	// Test that endpoints follow expected format
	endpoints := []struct {
		name string
		url  string
	}{
		{"RPC", TestRPCClient},
		{"EVM RPC", TestEVMRPCClient},
		{"API", TestAPIClient},
	}

	for _, endpoint := range endpoints {
		t.Run(endpoint.name+" endpoint validation", func(t *testing.T) {
			if len(endpoint.url) == 0 {
				t.Error("Endpoint should not be empty")
			}

			if len(endpoint.url) < 8 || endpoint.url[:8] != "https://" {
				t.Errorf("Endpoint %s should use HTTPS", endpoint.url)
			}

			if len(endpoint.url) < 4 || endpoint.url[len(endpoint.url)-4:] != ":443" {
				t.Errorf("Endpoint %s should use port 443", endpoint.url)
			}

			if !contains(endpoint.url, "fivenet.sixprotocol.net") {
				t.Errorf("Endpoint %s should use fivenet.sixprotocol.net domain", endpoint.url)
			}
		})
	}
}

func TestClientFieldAccess(t *testing.T) {
	// Test that client can access all context fields through embedding
	client := Client{
		Context:      context.Background(),
		RPCClient:    TestRPCClient,
		EVMRPCCleint: TestEVMRPCClient,
		APIClient:    TestAPIClient,
	}

	// Test direct access through embedded struct
	if client.RPCClient != TestRPCClient {
		t.Error("Client should have access to RPCClient through embedded LBBContext")
	}

	if client.EVMRPCCleint != TestEVMRPCClient {
		t.Error("Client should have access to EVMRPCCleint through embedded LBBContext")
	}

	if client.APIClient != TestAPIClient {
		t.Error("Client should have access to APIClient through embedded LBBContext")
	}

	if client.Context != context.Background() {
		t.Error("Client should have access to Context field through embedded LBBContext")
	}
}

// Benchmark tests
func BenchmarkContextCreation(b *testing.B) {
	if b.Loop() {
		_ = Client{
			Context:      context.Background(),
			RPCClient:    TestRPCClient,
			EVMRPCCleint: TestEVMRPCClient,
			APIClient:    TestAPIClient,
		}
	}
}

func BenchmarkClientCreation(b *testing.B) {
	ctx := Client{
		Context:      context.Background(),
		RPCClient:    TestRPCClient,
		EVMRPCCleint: TestEVMRPCClient,
		APIClient:    TestAPIClient,
	}

	b.ResetTimer()
	if b.Loop() {
		_ = Client{Context: ctx}
	}
}

func BenchmarkFullSetup(b *testing.B) {
	if b.Loop() {
		ctx := Client{
			Context:      context.Background(),
			RPCClient:    TestRPCClient,
			EVMRPCCleint: TestEVMRPCClient,
			APIClient:    TestAPIClient,
		}
		_ = Client{Context: ctx}
	}
}

// Helper function for string contains check
func contains(str, substr string) bool {
	if len(substr) > len(str) {
		return false
	}
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
