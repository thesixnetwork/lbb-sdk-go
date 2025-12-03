package evm

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/evm/assets"
)

func (e *EVMClient) DeployCertificateContract(contractName, symbol, nftSchemaCode string) (common.Address, *types.Transaction, error) {
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

	address, tx, _, err := bind.DeployContract(auth, contractABI, common.FromHex(stringBIN), e.ETHClient, contractName, symbol, "URL", "URL", e.GetEVMAddress())
	if err != nil {
		return common.Address{}, &types.Transaction{}, err
	}

	return address, tx, nil
}

func (e *EVMClient) MintCertificateNFT(contractAddress common.Address, tokenID string) (tx *types.Transaction, err error) {
	stringABI, err := assets.GetContractABIString()
	if err != nil {
		return  &types.Transaction{}, err
	}

	contractABI, err := abi.JSON(strings.NewReader(stringABI))
	if err != nil {
		return &types.Transaction{}, err
	}

	// Pack the function call
	data, err := contractABI.Pack("safeMint", e.GetEVMAddress(), tokenID)
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

func (e *EVMClient) TransferCertificateNFT(contractAddress common.Address, destAddress common.Address ,tokenID string) (tx *types.Transaction, err error) {
	stringABI, err := assets.GetContractABIString()
	if err != nil {
		return  &types.Transaction{}, err
	}

	contractABI, err := abi.JSON(strings.NewReader(stringABI))
	if err != nil {
		return &types.Transaction{}, err
	}

	tokenIDInt,err := strconv.Atoi(tokenID)
	if err != nil {
		return &types.Transaction{}, err
	}

	tokenIdptr := big.NewInt(int64(tokenIDInt))
	// Pack the function call
	data, err := contractABI.Pack("safeTransferFrom", e.GetEVMAddress(), destAddress, tokenIdptr)
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
