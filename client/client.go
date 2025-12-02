package client

import (
	"context"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/thesixnetwork/lbb-sdk-go/config"

	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	// TESTNET DEFAULT
	TestnetRPC     = "https://rpc1.fivenet.sixprotocol.net"
	TestnetAPI     = "https://api1.fivenet.sixprotocol.net"
	TestnetEVMRPC  = "https://rpc-evm.fivenet.sixprotocol.net"
	TestnetChainID = "fivenet"

	// MINNET DEFAULT
	MainnetRPC     = "https://sixnet-rpc.sixprotocol.net/"
	MainnetAPI     = "https://sixnet-api.sixprotocol.net"
	MainnetEVMRPC  = "https://sixnet-rpc.sixprotocol.net"
	MainnetChainID = "sixnet"
)

type Client struct {
	context.Context
	ETHClient         *ethclient.Client
	CosmosClientCTX   client.Context
	Codec             codec.Codec
	InterfaceRegistry codectypes.InterfaceRegistry
	LegacyAmino       *codec.LegacyAmino
	RPCClient         string
	EVMRPCCleint      string
	APIClient         string
	ChainID           string
}

type ClientI interface {
	GetClientCTX() client.Context
	GetKeyring() keyring.Keyring
	GetETHClient() *ethclient.Client
}

var _ ClientI = (*Client)(nil)

// NewClient creates a new Client Context with properly initialized codecs
func NewClient(ctx context.Context, testnet bool) (Client, error) {
	var rpcURL, apiURL, evmrpcURL, chainID string
	if testnet {
		rpcURL = TestnetRPC
		apiURL = TestnetAPI
		evmrpcURL = TestnetEVMRPC
		chainID = TestnetChainID
	} else {
		rpcURL = MainnetRPC
		apiURL = MainnetAPI
		evmrpcURL = MainnetEVMRPC
		chainID = MainnetChainID
	}

	encodingConfig := config.MakeConfig()
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
		WithBroadcastMode(flags.BroadcastAsync).
		WithKeyring(kr).
		WithNodeURI(rpcURL).
		WithClient(rpcclient).
		WithChainID(chainID)

	evmClient, err := ethclient.Dial(evmrpcURL)
	if err != nil {
		return Client{}, nil
	}

	return Client{
		Context:           ctx,
		CosmosClientCTX:   cosmosClientCTX,
		ETHClient:         evmClient,
		Codec:             encodingConfig.Codec,
		InterfaceRegistry: encodingConfig.InterfaceRegistry,
		LegacyAmino:       encodingConfig.Amino,
		RPCClient:         rpcURL,
		EVMRPCCleint:      evmrpcURL,
		APIClient:         apiURL,
		ChainID:           chainID,
	}, nil
}

func NewCustomClient(ctx context.Context, rpcURL, apiURL, evmRPC, chainID string) (Client, error) {
	encodingConfig := config.MakeConfig()
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
		WithClient(rpcclient).
		WithChainID(chainID)

	evmClient, err := ethclient.Dial(evmRPC)
	if err != nil {
		return Client{}, nil
	}
	return Client{
		Context:           ctx,
		ETHClient:         evmClient,
		CosmosClientCTX:   cosmosClientCTX,
		Codec:             encodingConfig.Codec,
		InterfaceRegistry: encodingConfig.InterfaceRegistry,
		LegacyAmino:       encodingConfig.Amino,
		RPCClient:         rpcURL,
		EVMRPCCleint:      evmRPC,
		APIClient:         apiURL,
		ChainID:           chainID,
	}, nil
}

func (c *Client) GetRPCClient() string {
	return c.RPCClient
}

func (c *Client) GetETHClient() *ethclient.Client {
	return c.ETHClient
}

func (c *Client) GetAPIClient() string {
	return c.APIClient
}

func (c *Client) GetEVMRPCClient() string {
	return c.EVMRPCCleint
}

func (c *Client) GetChainID() string {
	return c.ChainID
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

func (c Client) WithFrom(from string) Client {
	c.CosmosClientCTX.From = from
	return c
}

func (c Client) WithFromName(fromName string) Client {
	c.CosmosClientCTX.FromName = fromName
	return c
}

// NewClientFromNode sets up Client implementation that communicates with a CometBFT node over
// JSON RPC and WebSockets
func NewClientFromNode(nodeURI string) (*rpchttp.HTTP, error) {
	return rpchttp.New(nodeURI, "/websocket")
}
