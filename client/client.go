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
	"github.com/thesixnetwork/six-protocol/v4/encoding"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
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
	GetClientCTX() client.Context
	GetKeyring() keyring.Keyring
}

var _ ClientI = (*Client)(nil)

// NewClient creates a new Client Context with properly initialized codecs
func NewClient(ctx context.Context, rpcURL, evmrpcURL, apiURL string) (Client, error) {
	encodingConfig := encoding.MakeConfig()
	cdc := encodingConfig.Codec
	LegacyAmino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry
	txConfig := authtx.NewTxConfig(encodingConfig.Codec, authtx.DefaultSignModes)
	kr := keyring.NewInMemory(encodingConfig.Codec)
	rpcclient, err := NewClientFromNode(rpcURL)
	if err != nil {
		return Client{}, nil
	}
	cosmosClientCTX := client.Context{}.
		WithCodec(cdc).
		WithInterfaceRegistry(interfaceRegistry).
		WithTxConfig(txConfig).
		WithLegacyAmino(LegacyAmino).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastSync).
		WithKeyring(kr).
		WithNodeURI(rpcURL).
		WithClient(rpcclient)

	return Client{
		Context:           ctx,
		CosmosClientCTX:   cosmosClientCTX,
		Codec:             encodingConfig.Codec,
		InterfaceRegistry: encodingConfig.InterfaceRegistry,
		LegacyAmino:       encodingConfig.Amino,
		RPCClient:         rpcURL,
		EVMRPCCleint:      evmrpcURL,
		APIClient:         apiURL,
	}, nil
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

func (c *Client) GetClientCTX() client.Context {
	return c.CosmosClientCTX
}

func (c *Client) GetKeyring() keyring.Keyring {
	return c.CosmosClientCTX.Keyring
}

func (c *Client) GetContext() context.Context {
	return c.Context
}

// NewClientFromNode sets up Client implementation that communicates with a CometBFT node over
// JSON RPC and WebSockets
func NewClientFromNode(nodeURI string) (*rpchttp.HTTP, error) {
	return rpchttp.New(nodeURI, "/websocket")
}
