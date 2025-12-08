package account

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bip39 "github.com/cosmos/go-bip39"

	client "github.com/thesixnetwork/lbb-sdk-go/client"
)

const (
	mnemonicEntropySize = 256
)

type AccountI interface {
	GetCosmosAddress() sdk.AccAddress
	GetEVMAddress() common.Address
	GetAccountName() string
	GetPrivateKey() *ecdsa.PrivateKey
	GetTransactOpts() *bind.TransactOpts
	GetClient() client.ClientI
}

// Account represents a blockchain account with both Cosmos and EVM capabilities
type Account struct {
	client        client.ClientI
	auth          *bind.TransactOpts
	mnemonic      string
	privateKey    *ecdsa.PrivateKey
	evmAddress    common.Address
	cosmosAddress sdk.AccAddress
	accountName   string
}

var _ AccountI = (*Account)(nil)

// NewAccount creates a new Account instance from a mnemonic and password
// Returns an error if any step in the account creation process fails
func NewAccount(ctx client.ClientI, accountName, mnemonic, password string) (*Account, error) {
	if ctx == nil {
		return nil, fmt.Errorf("client cannot be nil")
	}

	if accountName == "" {
		return nil, fmt.Errorf("account name cannot be empty")
	}

	if !ValidateMnemonic(mnemonic) {
		return nil, fmt.Errorf("invalid mnemonic provided")
	}

	evmAddress, err := GetAddressFromMnemonic(mnemonic, password)
	if err != nil {
		return nil, fmt.Errorf("failed to get EVM address from mnemonic for account '%s': %w", accountName, err)
	}

	cosmosAddress, err := GetBech32AccountFromMnemonic(ctx.GetKeyring(), accountName, mnemonic, password)
	if err != nil {
		return nil, fmt.Errorf("failed to get Bech32 Cosmos address from mnemonic for account '%s': %w", accountName, err)
	}

	privateKey, err := CreatePrivateKeyFromMnemonic(mnemonic, password)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key from mnemonic for account '%s': %w", accountName, err)
	}

	// Get chain ID for EVM operations
	chainIDBigInt, ok := ChainIDMapping[ctx.GetChainID()]
	if !ok {
		return nil, fmt.Errorf("chain ID '%s' not found in mapping", ctx.GetChainID())
	}

	// Create transaction options for EVM operations
	authz, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIDBigInt)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor for account '%s': %w", accountName, err)
	}

	fmt.Printf("Account '%s' created successfully\n", accountName)
	fmt.Printf("  Cosmos Address: %s\n", cosmosAddress.String())
	fmt.Printf("  EVM Address: %s\n", evmAddress.Hex())

	return &Account{
		client:        ctx,
		auth:          authz,
		privateKey:    privateKey,
		mnemonic:      mnemonic,
		evmAddress:    evmAddress,
		cosmosAddress: cosmosAddress,
		accountName:   accountName,
	}, nil
}

// GetCosmosAddress returns the Cosmos Bech32 address
func (a *Account) GetCosmosAddress() sdk.AccAddress {
	return a.cosmosAddress
}

// GetEVMAddress returns the EVM (Ethereum) address
func (a *Account) GetEVMAddress() common.Address {
	return a.evmAddress
}

// GetAccountName returns the account name
func (a *Account) GetAccountName() string {
	return a.accountName
}

// GetPrivateKey returns the ECDSA private key
func (a *Account) GetPrivateKey() *ecdsa.PrivateKey {
	return a.privateKey
}

// GetTransactOpts returns the transaction options for EVM operations
func (a *Account) GetTransactOpts() *bind.TransactOpts {
	return a.auth
}

// GetClient returns the underlying client
func (a *Account) GetClient() client.ClientI {
	return a.client
}

// GetMnemonic returns the mnemonic phrase (use with caution)
// This should only be used for backup purposes and the result should be kept secure
func (a *Account) GetMnemonic() string {
	return a.mnemonic
}

// ValidateMnemonic validates a BIP39 mnemonic phrase
// This is a package-level function as it doesn't require account state
func ValidateMnemonic(mnemonic string) bool {
	return bip39.IsMnemonicValid(mnemonic)
}

// GenerateNewMnemonic generates a new BIP39 mnemonic with the default entropy size
func GenerateNewMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(mnemonicEntropySize)
	if err != nil {
		return "", fmt.Errorf("failed to generate entropy: %w", err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("failed to generate mnemonic: %w", err)
	}

	return mnemonic, nil
}

// String returns a string representation of the account (for debugging)
func (a *Account) String() string {
	return fmt.Sprintf("Account{name: %s, cosmos: %s, evm: %s}",
		a.accountName,
		a.cosmosAddress.String(),
		a.evmAddress.Hex())
}

// NewAccountFromPrivateKey creates a new Account instance from an existing private key
// This is useful when you want to create an account without using a mnemonic
// NOTE:: Use this function for creating isolate EVM client ONLY. And prevent initial packages of cosmos such as metadata and bank
func NewAccountFromPrivateKey(ctx client.ClientI, accountName string, privateKey *ecdsa.PrivateKey) (*Account, error) {
	if ctx == nil {
		return nil, fmt.Errorf("client cannot be nil")
	}

	if accountName == "" {
		return nil, fmt.Errorf("account name cannot be empty")
	}

	if privateKey == nil {
		return nil, fmt.Errorf("private key cannot be nil")
	}

	// Derive EVM address from private key
	evmAddress, err := GetAddressFromPrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to derive EVM address from private key: %w", err)
	}

	// Get chain ID for EVM operations
	chainIDBigInt, ok := ChainIDMapping[ctx.GetChainID()]
	if !ok {
		return nil, fmt.Errorf("chain ID '%s' not found in mapping", ctx.GetChainID())
	}

	// Create transaction options for EVM operations
	authz, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIDBigInt)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor for account '%s': %w", accountName, err)
	}

	return &Account{
		client:      ctx,
		auth:        authz,
		privateKey:  privateKey,
		evmAddress:  evmAddress,
		accountName: accountName,
	}, nil
}
