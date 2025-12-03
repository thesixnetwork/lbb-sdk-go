package evm

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/evm/assets"
)

func (e *EVMClient) CurrentOwner(contractAddress common.Address, tokenID string) common.Address {
	var currentOwner common.Address

	stringABI, err := assets.GetContractABIString()
	if err != nil {
		return currentOwner
	}

	contractABI, err := abi.JSON(strings.NewReader(stringABI))
	if err != nil {
		return currentOwner
	}

	// Get numeric token ID from string token ID
	numericTokenID, err := e.GetNumericTokenID(contractAddress, tokenID)
	if err != nil {
		fmt.Printf("failed to get numeric token ID: %v\n", err)
		return currentOwner
	}

	// Pack the function call
	data, err := contractABI.Pack("ownerOf", numericTokenID)
	if err != nil {
		fmt.Printf("failed to pack data: %v\n", err)
		return currentOwner
	}

	result, err := e.ETHClient.CallContract(e.GetContext(), ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}, nil)
	if err != nil {
		fmt.Printf("failed to call contract: %v\n", err)
		return currentOwner
	}

	var addressOutpu common.Address

	err = contractABI.UnpackIntoInterface(&addressOutpu, "ownerOf", result)
	if err != nil {
		fmt.Printf("failed to call contract: %v\n", err)
		return currentOwner
	}

	return addressOutpu
}

func (e *EVMClient) GetNumericTokenID(contractAddress common.Address, tokenIDString string) (*big.Int, error) {
	stringABI, err := assets.GetContractABIString()
	if err != nil {
		return nil, fmt.Errorf("failed to get ABI: %w", err)
	}

	contractABI, err := abi.JSON(strings.NewReader(stringABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	// Pack the function call for getTokenIdFromString
	data, err := contractABI.Pack("getTokenIdFromString", tokenIDString)
	if err != nil {
		return nil, fmt.Errorf("failed to pack data: %w", err)
	}

	result, err := e.ETHClient.CallContract(e.GetContext(), ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %w", err)
	}

	var numericTokenID *big.Int
	err = contractABI.UnpackIntoInterface(&numericTokenID, "getTokenIdFromString", result)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack result: %w", err)
	}

	if numericTokenID.Cmp(big.NewInt(0)) == 0 {
		return nil, fmt.Errorf("token ID string not found: %s", tokenIDString)
	}

	return numericTokenID, nil
}
