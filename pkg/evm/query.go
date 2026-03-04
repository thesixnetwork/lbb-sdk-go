package evm

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/thesixnetwork/lbb-sdk-go/pkg/evm/assets"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/logger"
)

func (e *EVMClient) TokenOwner(contractAddress common.Address, tokenID uint64) common.Address {
	goCtx := e.GetClient().GetContext()
	ethClient := e.GetClient().GetETHClient()

	var currentOwner common.Address

	stringABI, err := assets.GetContractABIString()
	if err != nil {
		return currentOwner
	}

	contractABI, err := abi.JSON(strings.NewReader(stringABI))
	if err != nil {
		return currentOwner
	}

	// Pack the function call
	data, err := contractABI.Pack("ownerOf", big.NewInt(int64(tokenID)))
	if err != nil {
		logger.Error("Failed to pack data: %v", err)
		return currentOwner
	}

	result, err := ethClient.CallContract(goCtx, ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}, nil)
	if err != nil {
		logger.Error("Failed to call contract: %v", err)
		return currentOwner
	}

	var addressOutput common.Address

	err = contractABI.UnpackIntoInterface(&addressOutput, "ownerOf", result)
	if err != nil {
		logger.Error("Failed to unpack result: %v", err)
		return currentOwner
	}

	return addressOutput
}

func (e *EVMClient) IsMinted(contractAddress common.Address, tokenID uint64) (bool, error) {
	goCtx := e.GetClient().GetContext()
	ethClient := e.GetClient().GetETHClient()

	stringABI, err := assets.GetContractABIString()
	if err != nil {
		return false, fmt.Errorf("failed to get ABI: %w", err)
	}

	contractABI, err := abi.JSON(strings.NewReader(stringABI))
	if err != nil {
		return false, fmt.Errorf("failed to parse ABI: %w", err)
	}

	// 1. Pack the "ownerOf" function call
	data, err := contractABI.Pack("ownerOf", big.NewInt(int64(tokenID)))
	if err != nil {
		return false, fmt.Errorf("failed to pack ownerOf: %w", err)
	}

	result, err := ethClient.CallContract(goCtx, ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}, nil)
	if err != nil {
		return false, nil
	}

	unpacked, err := contractABI.Unpack("ownerOf", result)
	if err != nil {
		return false, fmt.Errorf("failed to unpack result: %w", err)
	}

	owner := *abi.ConvertType(unpacked[0], new(common.Address)).(*common.Address)
	return owner != common.Address{}, nil
}
