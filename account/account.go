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
}

type Account struct {
	client.Client
	auth          *bind.TransactOpts
	mnemonic      string
	privateKey    *ecdsa.PrivateKey
	evmAddress    common.Address
	cosmosAddress sdk.AccAddress
	accountName   string
}

var _ AccountI = (*Account)(nil)

func NewAccount(ctx client.Client, accountName, mnemonic, password string) *Account {
	evmAddress, err := GetAddressFromMnemonic(mnemonic, password)
	if err != nil {
		fmt.Printf("ERROR: Failed to get EVM address from mnemonic for account '%s': %v\n", accountName, err)
		return nil
	}

	cosmosAddress, err := GetBech32AccountFromMnemonic(ctx.GetKeyring(), accountName, mnemonic, password)
	if err != nil {
		fmt.Printf("ERROR: Failed to get Bech32 Cosmos address from mnemonic for account '%s': %v\n", accountName, err)
		fmt.Printf("  - Account name: %s\n", accountName)
		fmt.Printf("  - Mnemonic valid: %v\n", bip39.IsMnemonicValid(mnemonic))
		fmt.Printf("  - Password length: %d\n", len(password))
		return nil
	}

	privateKey, err := CreatePrivateKeyFromMnemonic(mnemonic, password)
	if err != nil {
		fmt.Printf("ERROR: Failed to generate PrivateKey from mnemonic '%s': %v\n", accountName, err)
		return nil
	}

	authz, err := bind.NewKeyedTransactorWithChainID(privateKey, ChainIDMapping[ctx.ChainID])
	if err != nil {
		fmt.Printf("ERROR: Failed to bind account '%s': %v\n", accountName, err)
		return nil
	}

	fmt.Printf("Account created successfully: %s (Cosmos: %s, EVM: %s)\n", accountName, cosmosAddress.String(), evmAddress.Hex())

	return &Account{
		Client:        ctx,
		auth:          authz,
		privateKey:    privateKey,
		mnemonic:      mnemonic,
		evmAddress:    evmAddress,
		cosmosAddress: cosmosAddress,
		accountName:   accountName,
	}
}

func (a *Account) ValidateMnemonic(mnemonic string) bool {
	return bip39.IsMnemonicValid(mnemonic)
}

func (a *Account) GetPrivateKey() *ecdsa.PrivateKey {
	return a.privateKey
}

func (a *Account) GetTransactionOps() *bind.TransactOpts {
	return a.auth
}

func (a *Account) GetCosmosAddress() sdk.AccAddress {
	return a.cosmosAddress
}

func (a *Account) GetEVMAddress() common.Address {
	return a.evmAddress
}
