package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

type Client struct {
	LBBContext      Context
	CosmosClientCTX client.Context
}

type ClientI interface {
	GetCosmosClientCTX() client.Context
	GetLBBClientCTX() Context
}

var _ ClientI = (*Client)(nil)

// NewClient creates a new Client with the provided Context
func NewClient(ctx Context) *Client {
	cdc := ctx.Codec
	txConfig := authtx.NewTxConfig(cdc, authtx.DefaultSignModes)
	cosmosClientCTX := client.Context{}.
		WithCodec(cdc).
		WithInterfaceRegistry(ctx.InterfaceRegistry).
		WithTxConfig(txConfig).
		WithLegacyAmino(ctx.LegacyAmino).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastSync).
		WithKeyring(ctx.Keyring)

	return &Client{
		LBBContext:      ctx,
		CosmosClientCTX: cosmosClientCTX,
	}
}

func (c *Client) GetCosmosClientCTX() client.Context {
	return c.CosmosClientCTX
}

func (c *Client) GetLBBClientCTX() Context {
	return c.LBBContext
}
