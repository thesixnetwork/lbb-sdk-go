package main

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/client"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/evm"
)

const (
	// IMPORTANT: Replace with your deployed contract address
	contractAddress = "0x0000000000000000000000000000000000000000"

	// Token ID to transfer (must exist and be owned by you)
	tokenId = uint64(1)

	// Recipient address (who will receive the NFT)
	// Replace with the actual recipient's EVM address
	recipientAddress = "0x8a28fb81A084Ac7A276800957a19a6054BF86E4D"

	// For this example, we use the test mnemonic
	exampleMnemonic = account.TestMnemonic
)

func main(){

}
