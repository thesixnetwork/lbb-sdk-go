package account

import (
	"crypto/ecdsa"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bip39 "github.com/cosmos/go-bip39"
)

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

func CreatePrivateKeyFromMnemonic(mnemonic string, password string) (*ecdsa.PrivateKey, error) {
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

func CreateAccountFromPrivateKey(hexprivatekey string) (*ecdsa.PrivateKey, error) {
	privateKey, err := crypto.HexToECDSA(hexprivatekey)
	if err != nil {
		return &ecdsa.PrivateKey{}, err
	}
	return privateKey, nil
}

func GetBech32AccountFromMnemonic(keyring keyring.Keyring, accountName, mnemonic, password string) (sdk.AccAddress, error) {
	if bip39.IsMnemonicValid(mnemonic) {
		return sdk.AccAddress{}, errors.New("invalid mnemonic")
	}

	path := GetFullBIP44Path()

	kr, err := keyring.NewAccount(accountName, mnemonic, password, path, hd.Secp256k1)
	if err != nil {
		return sdk.AccAddress{}, err
	}

	account, err := kr.GetAddress()
	if err != nil {
		return sdk.AccAddress{}, err
	}

	return account, nil
}

func GetAddressFromMnemonic(mnemonic, password string) (common.Address, error) {
	if bip39.IsMnemonicValid(mnemonic) {
		return common.Address{}, errors.New("invalid mnemonic")
	}

	// Implementation here
	seed := bip39.NewSeed(mnemonic, password)

	privateKey, err := crypto.ToECDSA(seed[:32])
	if err != nil {
		return common.Address{}, err
	}

	pubkey := privateKey.PublicKey
	address := crypto.PubkeyToAddress(pubkey)

	return address, nil
}

func GetAddressFromPrivateKeyHex(privateKey string) (common.Address, error) {
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return common.Address{}, err
	}

	pubkey := key.PublicKey
	address := crypto.PubkeyToAddress(pubkey)

	return address, nil
}

func GetAddressFromPrivateKey(privatekey *ecdsa.PrivateKey) (common.Address, error) {
	pubkey := privatekey.PublicKey
	address := crypto.PubkeyToAddress(pubkey)

	return address, nil
}
