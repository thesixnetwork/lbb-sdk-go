package client

import (
	"context"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/thesixnetwork/six-protocol/v4/app"
)

type Context struct {
	Context           context.Context
	Codec             codec.Codec
	InterfaceRegistry codectypes.InterfaceRegistry
	RPCClient         string
	EVMRPCCleint      string
	APIClient         string
}

type Client struct {
	Context
}

// NewClient creates a new Client with the provided Context
func NewClient(ctx Context) *Client {
	return &Client{
		Context: ctx,
	}
}

// NewContext creates a new Context with properly initialized codecs
func NewContext(ctx context.Context, rpcClient, evmRPCCleint, apiClient string) Context {
	encodingConfig := app.MakeEncodingConfig()
	return Context{
		Context:           ctx,
		Codec:             encodingConfig.Codec,
		InterfaceRegistry: encodingConfig.InterfaceRegistry,
		RPCClient:         rpcClient,
		EVMRPCCleint:      evmRPCCleint,
		APIClient:         apiClient,
	}
}

// NewContextWithCodec creates a new Context with custom codec configuration
func NewContextWithCodec(ctx context.Context, codec codec.Codec, interfaceRegistry codectypes.InterfaceRegistry, rpcClient, evmRPCCleint, apiClient string) Context {
	return Context{
		Context:           ctx,
		Codec:             codec,
		InterfaceRegistry: interfaceRegistry,
		RPCClient:         rpcClient,
		EVMRPCCleint:      evmRPCCleint,
		APIClient:         apiClient,
	}
}
