package main

import (
	"context"
	"fmt"

	// "cosmossdk.io/math"
	//sdk "github.com/cosmos/cosmos-sdk/types"
	// "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/client"

	"github.com/thesixnetwork/lbb-sdk-go/pkg/evm"
	// "github.com/thesixnetwork/lbb-sdk-go/pkg/metadata"
	// "github.com/thesixnetwork/lbb-sdk-go/pkg/balance"
)

const (
	BobAddress = "6x13g50hqdqsjk85fmgqz2h5xdxq49lsmjdwlemsp"
)

func main() {
	// Initialize Client
	client, err := client.NewCustomClient(
		context.Background(),
		"http://localhost:26657",
		"http://localhost:1317",
		"http://localhost:8545",
		"testnet",
	)
	//client, err := client.NewClient(
	//	context.Background(),
	//	true,
	//)
	if err != nil {
		panic(fmt.Sprintf("ERROR CREATE CLEINT %v", err))
	}

	a := account.NewAccount(client, "alice", account.TestMnemonic, "")
	if a == nil {
		panic("ERROR CREATE ACCOUNT: NewAccount returned nil - check mnemonic and keyring initialization")
	}

	// balanceClient := balance.NewBalanceMsg(*a)

	// sendAmount := sdk.Coin{
	// 	Amount: math.NewInt(20),
	// 	Denom:  "usix",
	// }

	// res, err := balanceClient.SendBalance(BobAddress, sdk.NewCoins(sendAmount))
	// if err != nil {
	// 	fmt.Printf("Send error: %v\n", err)
	// 	return
	// }

	// meta := metadata.NewMetadataMsg(*a, "sixnetwork.hamdee")
	// msgCreateMetadata2, err := meta.BuildMintMetadataMsg("3")
	// if err != nil {
	// 	fmt.Printf("Mint error: %v\n", err)
	// 	return
	// }
	//
	// var msgs []sdk.Msg

	// msgs = append(msgs, msgCreateMetadata2)

	// res, err := meta.BroadcastTx(msgs...)
	// if err != nil {
	// 	fmt.Printf("Mint error: %v\n", err)
	// }

	// fmt.Printf("Freeze response: %v\n", res)
	//res, err := meta.FreezeCertificate("1")
	//if err != nil {
	//	fmt.Printf("Freeze error: %v\n", err)
	//	return
	//}
	//fmt.Printf("Freeze response: %v\n", res)

	//res, err := meta.UnfreezeCertificate("1")
	//if err != nil {
	//	fmt.Printf("Unfreeze error: %v\n", err)
	//	return
	//}
	//fmt.Printf("Unfreeze response: %v\n", res)
	evm := evm.NewEVMClient(*a)
	address, tx, err := evm.DeployCertificateContract("NFT", "NFT", "sixnetwork.hamdee")
	if err != nil {
		fmt.Printf("EVM error: %v\n", err)
		return
	}

	fmt.Printf("Deploy at: %v\n", tx.Hash())
	fmt.Printf("Deploy at Nonce: %v\n", tx.Nonce())
	fmt.Printf("Contract Address: %v\n", address)

	// Wait for deployment transaction to be mined
	_, err = evm.WaitForTransaction(tx.Hash())
	if err != nil {
		fmt.Printf("Error waiting for deployment: %v\n", err)
		return
	}

	tx, err = evm.MintCertificateNFT(common.HexToAddress("0x67b18d8d5B82c7D8633a37d2909f6c82b7aCD6e7"), "1")
	if err != nil {
		fmt.Printf("EVM error: %v\n", err)
		return
	}
	fmt.Printf("Mint Tx: %+v \n", tx.Hash())
	fmt.Printf("Mint at Nonce: %v\n", tx.Nonce())

	//tx, err = evm.TransferCertificateNFT(address, common.HexToAddress("0xd907f36f7D83344057a619b6D83A45B3288c3c21"), "2")
	//if err != nil {
	//	fmt.Printf("EVM error: %v\n", err)
	//	return
	//}
	//fmt.Printf("Transfer Tx: %+v \n", tx.Hash())
	//fmt.Printf("Transfer at Nonce: %v\n", tx.Nonce())
	//Check your transactions:
	// err = evm.CheckTransactionReceipt()
	// if err != nil {
	// 	fmt.Printf("Mint transaction check: %v\n", err)
	// }

	//currentOwner := evm.CurrentOwner(common.HexToAddress("0x3224E227969A7a661798B59aF92fD250e9983dB6"), "2")
	//fmt.Printf("Current Owner: %+v \n", currentOwner)
}
