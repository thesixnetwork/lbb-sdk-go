package client

import (
	"context"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/thesixnetwork/six-protocol/v4/app"
)

type Context struct {
	Context           context.Context
	Codec             codec.Codec
	InterfaceRegistry codectypes.InterfaceRegistry
	LegacyAmino       *codec.LegacyAmino
	RPCClient         string
	EVMRPCCleint      string
	APIClient         string
	Keyring           keyring.Keyring
}

// NewContext creates a new Context with properly initialized codecs
func NewContext(ctx context.Context, rpcClient, evmRPCCleint, apiClient string) Context {
	encodingConfig := app.MakeEncodingConfig()
	return Context{
		Context:           ctx,
		Codec:             encodingConfig.Codec,
		InterfaceRegistry: encodingConfig.InterfaceRegistry,
		LegacyAmino:       encodingConfig.Amino,
		RPCClient:         rpcClient,
		EVMRPCCleint:      evmRPCCleint,
		APIClient:         apiClient,
		Keyring:           keyring.NewInMemory(encodingConfig.Codec),
	}
}

func (c *Context) GetKeyring() keyring.Keyring {
	return c.Keyring
}

func (c *Context) GetRPCClient() string {
	return c.RPCClient
}

func (c *Context) GetAPIClient() string {
	return c.APIClient
}

func (c *Context) GetEVMRPCClient() string {
	return c.EVMRPCCleint
}
