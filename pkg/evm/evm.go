package evm

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
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

func (e *EVMClient) GasLimit(callMsg ethereum.CallMsg) (uint64, error) {
	gasLimit, err := e.ETHClient.EstimateGas(e.GetContext(), callMsg)
	if err != nil {
		return gasLimit, err
	}
	return gasLimit, nil
}

func (e *EVMClient) ChainID() (*big.Int, error) {
	chainID, err := e.ETHClient.NetworkID(e.GetContext())
	if err != nil {
		return chainID, err
	}
	return chainID, err
}

func (e *EVMClient) GetNonce() (nonce uint64, err error) {
	nonce, err = e.ETHClient.PendingNonceAt(e.GetContext(), e.GetEVMAddress())
	if err != nil {
		return nonce, err
	}

	return nonce, nil
}

func (e *EVMClient) DynamicABI(contractAddress common.Address, functionName string, args interface{}) (tx *types.Transaction, err error) {
	stringABI, err := assets.GetContractABIString()
	if err != nil {
		return  &types.Transaction{}, err
	}

	contractABI, err := abi.JSON(strings.NewReader(stringABI))
	if err != nil {
		return &types.Transaction{}, err
	}

	// Pack the function call
	data, err := contractABI.Pack(functionName, args)
	if err != nil {
		return nil, fmt.Errorf("failed to pack data: %w", err)
	}

	gasLimit, err := e.GasLimit(ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	})
	if err != nil {
		return &types.Transaction{}, err
	}

	nonce, err := e.GetNonce()
	if err != nil {
		return &types.Transaction{}, err
	}

	gasPrice, err := e.GasPrice()
	if err != nil {
		return &types.Transaction{}, err
	}

	tx = types.NewTransaction(nonce, contractAddress, big.NewInt(0), gasLimit, gasPrice, data)

	chainID, err := e.ChainID()
	if err != nil {
		return &types.Transaction{}, err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), e.GetPrivateKey())
	if err != nil {
		return &types.Transaction{}, err
	}

	err = e.ETHClient.SendTransaction(e.GetContext(), signedTx)
	if err != nil {
		return &types.Transaction{}, err
	}

	return signedTx, nil
}

