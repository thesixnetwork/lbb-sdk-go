package account

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	client "github.com/thesixnetwork/lbb-sdk-go/client"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bip39 "github.com/cosmos/go-bip39"
)

const (
	mnemonicEntropySize = 256
)

type AccountI interface {
	AccountService
}

type AccountService interface {
	ValidateMnemonic(mnemonic string) bool
	CreateBech32AccountFromMnemonic(mnemonic, password string) (sdk.AccAddress, error)
	CreateEVMAccountFromMnemonic(mnemonic, password string) (common.Address, error)
	CreateEVMAccountFromPrivateKey(pk, password string) (string, error)
	GetPrivateKeyFromMnemonic(mnemonic, password string) (string, error)
}

type Account struct {
	ctx                  client.Context
	accountAddressPrefix string
	accountName          string
	Keyring              keyring.Keyring
	privateKey           *ecdsa.PrivateKey
}

var _ AccountI = (*Account)(nil)

func NewAccountService(ctx client.Context, accountName string, accountAddressPrefix string) AccountI {
	return &Account{
		ctx:                  ctx,
		accountName:          accountName,
		accountAddressPrefix: accountAddressPrefix,
	}
}

func GenerateMnemonic() (string, error) {
	// read entropy seed straight from tmcrypto.Rand and convert to mnemonic
	entropy, err := bip39.NewEntropy(mnemonicEntropySize)
	if err != nil {
		return "", err
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}

func (a *Account) ValidateMnemonic(mnemonic string) bool {
	if !bip39.IsMnemonicValid(mnemonic) {
		return false
	}
	return true
}

func (a *Account) CreateBech32AccountFromMnemonic(mnemonic, password string) (sdk.AccAddress, error) {
	if !a.ValidateMnemonic(mnemonic) {
		return sdk.AccAddress{}, errors.New("invalid mnemonic")
	}
	kb := keyring.NewInMemory(a.ctx.Codec)
	path := GetFullBIP44Path()

	kr, err := kb.NewAccount(a.accountName, mnemonic, password, path, hd.Secp256k1)
	if err != nil {
		return sdk.AccAddress{}, err
	}

	a.Keyring = kb

	account, err := kr.GetAddress()
	if err != nil {
		return sdk.AccAddress{}, err
	}

	return account, nil
}

func (a *Account) CreateEVMAccountFromMnemonic(mnemonic string, password string) (common.Address, error) {
	if !a.ValidateMnemonic(mnemonic) {
		return common.Address{}, errors.New("invalid mnemonic")
	}
	seed := bip39.NewSeed(mnemonic, password)

	privateKey, err := crypto.ToECDSA(seed[:32])
	if err != nil {
		return common.Address{}, err
	}

	a.privateKey = privateKey

	pubkey := privateKey.PublicKey
	return crypto.PubkeyToAddress(pubkey), nil
}

func (a *Account) GetPrivateKeyFromMnemonic(mnemonic, password string) (string, error) {
	if !a.ValidateMnemonic(mnemonic) {
		return "", errors.New("invalid mnemonic")
	}
	// Implementation here
	seed := bip39.NewSeed(mnemonic, password)

	privateKey, err := crypto.ToECDSA(seed[:32])
	if err != nil {
		return "", err
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := hex.EncodeToString(privateKeyBytes)

	return privateKeyHex, nil
}

func (a *Account) CreateEVMAccountFromPrivateKey(pk string, password string) (string, error) {
	if strings.HasPrefix(pk, "0x") {
		pk = strings.TrimPrefix(pk, "0x")
	}

	privateKeyBytes, err := hex.DecodeString(pk)
	if err != nil {
		return "", errors.New("invalid private key format")
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return "", err
	}

	a.privateKey = privateKey

	pubkey := privateKey.PublicKey
	address := crypto.PubkeyToAddress(pubkey)

	return address.Hex(), nil
}
