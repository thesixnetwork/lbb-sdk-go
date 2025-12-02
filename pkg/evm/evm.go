package evm

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/evm/assets"
)

type EVMClient struct {
	account.Account
}

func NewEVMClient(a account.Account) *EVMClient {
	return &EVMClient{
		Account: a,
	}
}

func (e *EVMClient) GasPrice() (*big.Int, error) {
	gasPrice, err := e.ETHClient.SuggestGasPrice(e.GetContext())
	if err != nil {
		return gasPrice, err
	}
	return gasPrice, nil
}

func (e *EVMClient) GetNonce() (nonce uint64, err error) {
	nonce, err = e.ETHClient.PendingNonceAt(e.GetContext(), e.GetEVMAddress())
	if err != nil {
		return nonce, err
	}

	return nonce, nil
}

func (e *EVMClient) DeployCertificateContract() (common.Address, *types.Transaction, error) {
	stringABI, err := assets.GetContractABIString()
	if err != nil {
		return common.Address{}, &types.Transaction{}, err
	}

	contractABI, err := abi.JSON(strings.NewReader(stringABI))
	if err != nil {
		return common.Address{}, &types.Transaction{}, err
	}

	stringBIN, err := assets.GetContractBINString()
	if err != nil {
		return common.Address{}, &types.Transaction{}, err
	}

	nonce, err := e.GetNonce()
	if err != nil {
		return common.Address{}, &types.Transaction{}, err
	}

	gasPrice, err := e.GasPrice()
	if err != nil {
		return common.Address{}, &types.Transaction{}, err
	}

	auth := e.GetTransactionOps()
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	address, tx, _, err := bind.DeployContract(auth, contractABI, common.FromHex(stringBIN), e.ETHClient, "NFT", "NFT", "URL", "URL", e.GetEVMAddress())
	if err != nil {
		return common.Address{}, &types.Transaction{}, err
	}

	return address, tx, nil
}

func (e *EVMClient) MintCertificateNFT(contractAdderss common.Address, tokenID string) (common.Address, *types.Transaction, error) {

}
