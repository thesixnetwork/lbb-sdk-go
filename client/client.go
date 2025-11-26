package client

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/thesixnetwork/six-protocol/v4/app"
)

type Client struct {
	context.Context
	CosmosClientCTX   client.Context
	Codec             codec.Codec
	InterfaceRegistry codectypes.InterfaceRegistry
	LegacyAmino       *codec.LegacyAmino
	RPCClient         string
	EVMRPCCleint      string
	APIClient         string
}

type ClientI interface {
	GetCosmosClientCTX() client.Context
	GetKeyring() keyring.Keyring
}

var _ ClientI = (*Client)(nil)

// NewClient creates a new Client Context with properly initialized codecs
func NewClient(ctx context.Context, rpcClient, evmRPCCleint, apiClient string) Client {
	encodingConfig := app.MakeEncodingConfig()
	txConfig := authtx.NewTxConfig(encodingConfig.Codec, authtx.DefaultSignModes)
	kr := keyring.NewInMemory(encodingConfig.Codec)
	cosmosClientCTX := client.Context{}.
		WithCodec(encodingConfig.Codec).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(txConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastSync).
		WithKeyring(kr)

	return Client{
		Context:           ctx,
		CosmosClientCTX:   cosmosClientCTX,
		Codec:             encodingConfig.Codec,
		InterfaceRegistry: encodingConfig.InterfaceRegistry,
		LegacyAmino:       encodingConfig.Amino,
		RPCClient:         rpcClient,
		EVMRPCCleint:      evmRPCCleint,
		APIClient:         apiClient,
	}
}

func (c *Client) GetRPCClient() string {
	return c.RPCClient
}

func (c *Client) GetAPIClient() string {
	return c.APIClient
}

func (c *Client) GetEVMRPCClient() string {
	return c.EVMRPCCleint
}

func (c *Client) GetCosmosClientCTX() client.Context {
	return c.CosmosClientCTX
}

func (c *Client) GetKeyring() keyring.Keyring {
	return c.CosmosClientCTX.Keyring
}
