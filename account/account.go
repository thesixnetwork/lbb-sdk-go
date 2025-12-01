package account

import (
	"crypto/ecdsa"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	client "github.com/thesixnetwork/lbb-sdk-go/client"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bip39 "github.com/cosmos/go-bip39"
)

const (
	mnemonicEntropySize = 256
)

type AccountI interface {
	GetCosmosAddress() sdk.AccAddress
	GetEVMAddress() common.Address
	GetPrivateKey(ctx client.Client, mnemonic string, password string) (*ecdsa.PrivateKey, error)
}

type Account struct {
	client.Client
	mnemonic      string
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

	fmt.Printf("Account created successfully: %s (Cosmos: %s, EVM: %s)\n", accountName, cosmosAddress.String(), evmAddress.Hex())

	return &Account{
		Client:        ctx,
		mnemonic:      mnemonic,
		evmAddress:    evmAddress,
		cosmosAddress: cosmosAddress,
		accountName:   accountName,
	}
}

func (a *Account) ValidateMnemonic(mnemonic string) bool {
	return bip39.IsMnemonicValid(mnemonic)
}

func (*Account) GetPrivateKey(ctx client.Client, mnemonic string, password string) (*ecdsa.PrivateKey, error) {
	if bip39.IsMnemonicValid(mnemonic) {
		return &ecdsa.PrivateKey{}, errors.New("invalid mnemonic")
	}

	seed := bip39.NewSeed(mnemonic, password)

	privateKey, err := crypto.ToECDSA(seed[:32])
	if err != nil {
		return &ecdsa.PrivateKey{}, err
	}

	return privateKey, nil
}

func (a *Account) GetCosmosAddress() sdk.AccAddress {
	return a.cosmosAddress
}

func (a *Account) GetEVMAddress() common.Address {
	return a.evmAddress
}
