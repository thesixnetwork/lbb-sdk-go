package client

import (
	"context"
	"fmt"
	"time"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/thesixnetwork/lbb-sdk-go/config"
)

const (
	// TESTNET DEFAULT
	TestnetRPC     = "https://rpc1.fivenet.sixprotocol.net"
	TestnetAPI     = "https://api1.fivenet.sixprotocol.net"
	TestnetEVMRPC  = "https://rpc-evm.fivenet.sixprotocol.net"
	TestnetChainID = "fivenet"

	// MAINNET DEFAULT
	MainnetRPC     = "https://sixnet-rpc.sixprotocol.net/"
	MainnetAPI     = "https://sixnet-api.sixprotocol.net"
	MainnetEVMRPC  = "https://sixnet-rpc.sixprotocol.net"
	MainnetChainID = "sixnet"

	// Transaction timeout settings
	transactionPollInterval = 1 * time.Second
	transactionTimeout      = 20 * time.Second
)

// ClientI defines the interface for blockchain client operations
type ClientI interface {
	GetClientCTX() client.Context
	GetKeyring() keyring.Keyring
	GetETHClient() *ethclient.Client
	GetRPCClient() string
	GetAPIClient() string
	GetEVMRPCClient() string
	GetChainID() string
	GetContext() context.Context
	WaitForTransaction(txhash string) error
	WaitForEVMTransaction(txHash common.Hash) (*types.Receipt, error)
}

type Client struct {
	ctx               context.Context
	ethClient         *ethclient.Client
	cosmosClientCTX   client.Context
	codec             codec.Codec
	interfaceRegistry codectypes.InterfaceRegistry
	legacyAmino       *codec.LegacyAmino
	rpcClient         string
	evmRPCClient      string
	apiClient         string
	chainID           string
}

var _ ClientI = (*Client)(nil)

// NewClient creates a new Client instance with default mainnet or testnet configuration
func NewClient(ctx context.Context, mainnet bool) (*Client, error) {
	var rpcURL, apiURL, evmRPCURL, chainID string

	if mainnet {
		rpcURL = MainnetRPC
		apiURL = MainnetAPI
		evmRPCURL = MainnetEVMRPC
		chainID = MainnetChainID
	} else {
		rpcURL = TestnetRPC
		apiURL = TestnetAPI
		evmRPCURL = TestnetEVMRPC
		chainID = TestnetChainID
	}

	return NewCustomClient(ctx, rpcURL, apiURL, evmRPCURL, chainID)
}

// NewCustomClient creates a new Client instance with custom configuration
func NewCustomClient(ctx context.Context, rpcURL, apiURL, evmRPC, chainID string) (*Client, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	encodingConfig := config.MakeConfig()
	kr := keyring.NewInMemory(encodingConfig.Codec)
	rpcClient, err := newClientFromNode(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create RPC client for %s: %w", rpcURL, err)
	}
	txConfig := authtx.NewTxConfig(encodingConfig.Codec, authtx.DefaultSignModes)
	cosmosClientCTX := client.Context{}.
		WithCodec(encodingConfig.Codec).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(txConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastSync).
		WithKeyring(kr).
		WithNodeURI(rpcURL).
		WithClient(rpcClient).
		WithChainID(chainID)
	evmClient, err := ethclient.Dial(evmRPC)
	if err != nil {
		return nil, fmt.Errorf("failed to create EVM client for %s: %w", evmRPC, err)
	}

	return &Client{
		ctx:               ctx,
		ethClient:         evmClient,
		cosmosClientCTX:   cosmosClientCTX,
		codec:             encodingConfig.Codec,
		interfaceRegistry: encodingConfig.InterfaceRegistry,
		legacyAmino:       encodingConfig.Amino,
		rpcClient:         rpcURL,
		evmRPCClient:      evmRPC,
		apiClient:         apiURL,
		chainID:           chainID,
	}, nil
}

// GetRPCClient returns the RPC client URL
func (c *Client) GetRPCClient() string {
	return c.rpcClient
}

// GetETHClient returns the Ethereum client instance
func (c *Client) GetETHClient() *ethclient.Client {
	return c.ethClient
}

// GetAPIClient returns the API client URL
func (c *Client) GetAPIClient() string {
	return c.apiClient
}

// GetEVMRPCClient returns the EVM RPC client URL
func (c *Client) GetEVMRPCClient() string {
	return c.evmRPCClient
}

// GetChainID returns the chain ID
func (c *Client) GetChainID() string {
	return c.chainID
}

// GetClientCTX returns the Cosmos client context
func (c *Client) GetClientCTX() client.Context {
	return c.cosmosClientCTX
}

// GetKeyring returns the keyring from the Cosmos client context
func (c *Client) GetKeyring() keyring.Keyring {
	return c.cosmosClientCTX.Keyring
}

// GetContext returns the context
func (c *Client) GetContext() context.Context {
	return c.ctx
}

// GetCodec returns the codec
func (c *Client) GetCodec() codec.Codec {
	return c.codec
}

// GetInterfaceRegistry returns the interface registry
func (c *Client) GetInterfaceRegistry() codectypes.InterfaceRegistry {
	return c.interfaceRegistry
}

// GetLegacyAmino returns the legacy amino codec
func (c *Client) GetLegacyAmino() *codec.LegacyAmino {
	return c.legacyAmino
}

// WaitForTransaction waits for a Cosmos transaction to be mined and returns an error if it fails
// The timeout is set to 20 seconds (approximately 3 blocks at 6.3s block time)
func (c *Client) WaitForTransaction(txHash string) error {
	if txHash == "" {
		return fmt.Errorf("transaction hash cannot be empty")
	}

	fmt.Printf("Waiting for transaction %s to be mined...\n", txHash)

	ticker := time.NewTicker(transactionPollInterval)
	defer ticker.Stop()

	timeout := time.After(transactionTimeout)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for transaction %s to be mined", txHash)
		case <-ticker.C:
			output, err := authtx.QueryTx(c.cosmosClientCTX, txHash)
			if err != nil {
				// Transaction not yet available, continue waiting
				continue
			}

			if output.Empty() {
				return fmt.Errorf("no transaction found with hash %s", txHash)
			}

			if output.Code != 0 {
				return fmt.Errorf("transaction %s failed with code %d: %s", txHash, output.Code, output.RawLog)
			}

			fmt.Printf("Transaction %s successfully mined in block %d\n", txHash, output.Height)
			return nil
		}
	}
}

// WaitForEVMTransaction waits for an EVM transaction to be mined and returns the receipt
// The timeout is set to 20 seconds (approximately 3 blocks at 6.3s block time)
func (c *Client) WaitForEVMTransaction(txHash common.Hash) (*types.Receipt, error) {
	if txHash == (common.Hash{}) {
		return nil, fmt.Errorf("transaction hash cannot be empty")
	}

	fmt.Printf("Waiting for EVM transaction %s to be mined...\n", txHash.Hex())

	ticker := time.NewTicker(transactionPollInterval)
	defer ticker.Stop()

	timeout := time.After(transactionTimeout)

	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("timeout waiting for transaction %s to be mined", txHash.Hex())
		case <-ticker.C:
			receipt, err := c.ethClient.TransactionReceipt(c.ctx, txHash)
			if err != nil {
				// Transaction not yet mined, continue waiting
				continue
			}

			if receipt.Status == 0 {
				return receipt, fmt.Errorf("transaction %s failed", txHash.Hex())
			}

			fmt.Printf("Transaction %s successfully mined in block %d\n", txHash.Hex(), receipt.BlockNumber.Uint64())
			return receipt, nil
		}
	}
}

// WithFrom returns a new Client with the specified from address
func (c *Client) WithFrom(from string) *Client {
	newClient := *c
	newClient.cosmosClientCTX = newClient.cosmosClientCTX.WithFrom(from)
	return &newClient
}

// WithFromName returns a new Client with the specified from name
func (c *Client) WithFromName(fromName string) *Client {
	newClient := *c
	newClient.cosmosClientCTX = newClient.cosmosClientCTX.WithFromName(fromName)
	return &newClient
}

// WithBroadcastMode returns a new Client with the specified broadcast mode
func (c *Client) WithBroadcastMode(mode string) *Client {
	newClient := *c
	newClient.cosmosClientCTX = newClient.cosmosClientCTX.WithBroadcastMode(mode)
	return &newClient
}

// newClientFromNode creates an RPC client that communicates with a CometBFT node
// over JSON RPC and WebSockets
func newClientFromNode(nodeURI string) (*rpchttp.HTTP, error) {
	if nodeURI == "" {
		return nil, fmt.Errorf("node URI cannot be empty")
	}
	return rpchttp.New(nodeURI, "/websocket")
}
