package evm

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/thesixnetwork/lbb-sdk-go/pkg/evm/assets"
	incrementassets "github.com/thesixnetwork/lbb-sdk-go/pkg/evm/assets/increment"
)

const (
	mainnetBaseURIPath = "https://gen2-api.sixprotocol.com/api/nft/metadata/"
	testnetBaseURIPath = "https://gen2-api.fivenet.sixprotocol.com/api/nft/metadata/"
)

func (e *EVMClient) SignTransferNFT(contractAddress common.Address, destAddress common.Address, tokenID uint64) (tx *types.Transaction, err error) {
	stringABI, err := assets.GetContractABIString()
	if err != nil {
		return &types.Transaction{}, err
	}

	contractABI, err := abi.JSON(strings.NewReader(stringABI))
	if err != nil {
		return &types.Transaction{}, err
	}

	// Pack the function call
	data, err := contractABI.Pack("safeTransferFrom", e.GetEVMAddress(), destAddress, big.NewInt(int64(tokenID)))
	if err != nil {
		return nil, fmt.Errorf("failed to pack data: %w", err)
	}

	gasLimit, err := e.GasLimit(ethereum.CallMsg{
		From: e.GetEVMAddress(),
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

	return signedTx, nil
}

func (e *EVMClient) SendTransaction(signedTx *types.Transaction) error {
	goCtx := e.GetClient().GetContext()
	ethClient := e.GetClient().GetETHClient()
	fmt.Printf("Sender: %v\n", e.GetEVMAddress())
	err := ethClient.SendTransaction(goCtx, signedTx)
	if err != nil {
		return err
	}

	return nil
}

func (e *EVMClient) DeployCertificateContract(contractName, symbol, nftSchemaCode string) (common.Address, *types.Transaction, error) {
	goCtx := e.GetClient().GetContext()
	_ = goCtx
	ethClient := e.GetClient().GetETHClient()

	var baseURI string
	if e.GetClient().GetChainID() == "sixnet" {
		baseURI = mainnetBaseURIPath + nftSchemaCode
	} else {
		baseURI = testnetBaseURIPath + nftSchemaCode
	}

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

	var construcArg []interface{}

	construcArg = append(construcArg, contractName, symbol, baseURI, e.GetEVMAddress())

	auth := e.GetTransactOpts()
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei
	gasLimit, err := e.EstimateDeployGas(contractABI, common.FromHex(stringBIN), construcArg...)
	if err != nil {
		return common.Address{}, &types.Transaction{}, err
	}
	auth.GasLimit = gasLimit
	auth.GasPrice = gasPrice

	address, tx, _, err := bind.DeployContract(auth, contractABI, common.FromHex(stringBIN), ethClient, construcArg...)
	if err != nil {
		return common.Address{}, &types.Transaction{}, err
	}

	return address, tx, nil
}

func (e *EVMClient) MintCertificateNFT(contractAddress common.Address, tokenID uint64) (tx *types.Transaction, err error) {
	goCtx := e.GetClient().GetContext()
	ethClient := e.GetClient().GetETHClient()

	stringABI, err := assets.GetContractABIString()
	if err != nil {
		return &types.Transaction{}, err
	}

	contractABI, err := abi.JSON(strings.NewReader(stringABI))
	if err != nil {
		return &types.Transaction{}, err
	}

	// Pack the function call
	data, err := contractABI.Pack("safeMint", e.GetEVMAddress(), big.NewInt(int64(tokenID)))
	if err != nil {
		return nil, fmt.Errorf("failed to pack data: %w", err)
	}

	gasLimit, err := e.GasLimit(ethereum.CallMsg{
		From: e.GetEVMAddress(),
		To:   &contractAddress,
		Data: data,
	})

	fmt.Printf("Mint Gas Limit: %v \n", gasLimit)

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

	err = ethClient.SendTransaction(goCtx, signedTx)
	if err != nil {
		return &types.Transaction{}, err
	}

	return signedTx, nil
}

func (e *EVMClient) MintCertificateNFTToDestination(contractAddress common.Address, tokenID uint64, destAddress common.Address) (tx *types.Transaction, err error) {
	goCtx := e.GetClient().GetContext()
	ethClient := e.GetClient().GetETHClient()

	stringABI, err := assets.GetContractABIString()
	if err != nil {
		return &types.Transaction{}, err
	}

	contractABI, err := abi.JSON(strings.NewReader(stringABI))
	if err != nil {
		return &types.Transaction{}, err
	}

	// Pack the function call
	data, err := contractABI.Pack("safeMint", destAddress, big.NewInt(int64(tokenID)))
	if err != nil {
		return nil, fmt.Errorf("failed to pack data: %w", err)
	}

	gasLimit, err := e.GasLimit(ethereum.CallMsg{
		From: e.GetEVMAddress(),
		To:   &contractAddress,
		Data: data,
	})

	fmt.Printf("Mint Gas Limit: %v \n", gasLimit)

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

	err = ethClient.SendTransaction(goCtx, signedTx)
	if err != nil {
		return &types.Transaction{}, err
	}

	return signedTx, nil
}

func (e *EVMClient) TransferCertificateNFT(contractAddress common.Address, destAddress common.Address, tokenID uint64) (tx *types.Transaction, err error) {
	goCtx := e.GetClient().GetContext()
	ethClient := e.GetClient().GetETHClient()

	signedTx, err := e.SignTransferNFT(contractAddress, destAddress, tokenID)
	if err != nil {
		return nil, err
	}

	err = ethClient.SendTransaction(goCtx, signedTx)
	if err != nil {
		return &types.Transaction{}, err
	}

	return signedTx, nil
}

func (e *EVMClient) DeployCertIDIncrementContract(contractName, symbol, nftSchemaCode string) (common.Address, *types.Transaction, error) {
	ethClient := e.GetClient().GetETHClient()

	var baseURI string
	if e.GetClient().GetChainID() == "sixnet" {
		baseURI = mainnetBaseURIPath + nftSchemaCode
	} else {
		baseURI = testnetBaseURIPath + nftSchemaCode
	}

	stringABI, err := incrementassets.GetContractABIString()
	if err != nil {
		return common.Address{}, &types.Transaction{}, err
	}

	contractABI, err := abi.JSON(strings.NewReader(stringABI))
	if err != nil {
		return common.Address{}, &types.Transaction{}, err
	}

	stringBIN, err := incrementassets.GetContractBINString()
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

	var construcArg []interface{}

	construcArg = append(construcArg, contractName, symbol, baseURI, e.GetEVMAddress())

	auth := e.GetTransactOpts()
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei
	gasLimit, err := e.EstimateDeployGas(contractABI, common.FromHex(stringBIN), construcArg...)
	if err != nil {
		return common.Address{}, &types.Transaction{}, err
	}
	auth.GasLimit = gasLimit
	auth.GasPrice = gasPrice

	address, tx, _, err := bind.DeployContract(auth, contractABI, common.FromHex(stringBIN), ethClient, construcArg...)
	if err != nil {
		return common.Address{}, &types.Transaction{}, err
	}

	return address, tx, nil
}

func (e *EVMClient) MintCertNFT(contractAddress common.Address) (tx *types.Transaction, err error) {
	goCtx := e.GetClient().GetContext()
	ethClient := e.GetClient().GetETHClient()
	stringABI, err := incrementassets.GetContractABIString()
	if err != nil {
		return &types.Transaction{}, err
	}

	contractABI, err := abi.JSON(strings.NewReader(stringABI))
	if err != nil {
		return &types.Transaction{}, err
	}

	// Pack the function call
	data, err := contractABI.Pack("safeMint", e.GetEVMAddress())
	if err != nil {
		return nil, fmt.Errorf("failed to pack data: %w", err)
	}

	gasLimit, err := e.GasLimit(ethereum.CallMsg{
		From: e.GetEVMAddress(),
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

	err = ethClient.SendTransaction(goCtx, signedTx)
	if err != nil {
		return &types.Transaction{}, err
	}

	return signedTx, nil
}
