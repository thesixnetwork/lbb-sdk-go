package main

import (
	"context"
	"fmt"

	// "cosmossdk.io/math"
	// sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/client"

	"github.com/thesixnetwork/lbb-sdk-go/pkg/evm"
	// "github.com/thesixnetwork/lbb-sdk-go/pkg/metadata"
	//"github.com/thesixnetwork/lbb-sdk-go/pkg/balance"
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
	address, tx, err := evm.DeployCertificateContract()
	if err != nil {
		fmt.Printf("EVM error: %v\n", err)
		return
	}

	fmt.Printf("Deployed Tx: %+v, Address: %v \n", tx.Hash(), address)
}
