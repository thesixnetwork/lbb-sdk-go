package account

import (
	"crypto/ecdsa"
	"errors"

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
	ValidateMnemonic(mnemonic string) bool
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
		return nil
	}

	cosmosAddress, err := GetBech32AccountFromMnemonic(ctx.GetKeyring(), accountName, mnemonic, password)
	if err != nil {
		return nil
	}

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
