package evm

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/evm/assets"
)

func (e *EVMClient) CurrentOwner(contractAddress common.Address, tokenID string) common.Address {
	var currentOwner common.Address

	fmt.Printf("Input data: %v, %v\n", contractAddress, tokenID)

	stringABI, err := assets.GetContractABIString()
	if err != nil {
		return currentOwner
	}

	contractABI, err := abi.JSON(strings.NewReader(stringABI))
	if err != nil {
		return currentOwner
	}

	tokenIDInt, err := strconv.Atoi(tokenID)
	if err != nil {
		fmt.Printf("failed to : %w \n", err)
		return currentOwner
	}

	tokenIdptr := big.NewInt(int64(tokenIDInt))

	// Pack the function call
	data, err := contractABI.Pack("ownerOf", tokenIdptr)
	if err != nil {
		fmt.Printf("failed to pack data: %w\n", err)
		return currentOwner
	}

	fmt.Printf("Data Pack: %v \n", data)

	result, err := e.ETHClient.CallContract(e.GetContext(), ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}, nil)
	if err != nil {
		fmt.Errorf("failed to call contract: %w\n", err)
		return currentOwner
	}

	fmt.Printf("Pack result: %v \n", result)

	var addressOutpu common.Address

	err = contractABI.UnpackIntoInterface(&addressOutpu,"ownerOf", result)
	if err != nil {
		fmt.Printf("failed to call contract: %w\n", err)
		return currentOwner
	}

	return addressOutpu
}
