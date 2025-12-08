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
		fmt.Printf("failed to pack data: %v\n", err)
		return currentOwner
	}

	result, err := ethClient.CallContract(goCtx, ethereum.CallMsg{
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

