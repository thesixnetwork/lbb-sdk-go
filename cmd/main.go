package main

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/client"

	"github.com/thesixnetwork/lbb-sdk-go/pkg/balance"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/evm"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/metadata"
)

const (
	BobAddress    = "6x13g50hqdqsjk85fmgqz2h5xdxq49lsmjdwlemsp"
	BobEVMAddres = "0x8a28fb81A084Ac7A276800957a19a6054BF86E4D"
	nftSchemaName = "sixnetwork.lbbv01" // {ORGNAME}.{Schemacode}
	contractName = "MyNFTCert"
	contractSymbol = "Cert"
)

func init(){
	mnemonic, _ := account.GenerateMnemonic()
	fmt.Println("-----------------------------------------------------")
	fmt.Println()
	fmt.Println()
	fmt.Printf("THIS IS JUST DEMO HOW TO GEN MNEMONIC \n: %+v \n",mnemonic)
	fmt.Println()
	fmt.Println()
	fmt.Println("-----------------------------------------------------")
}

func main() {
	// Initialize Client
	// Testing local or run replicate node use CustomClient
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

	/*
		   NOTE:: For test in official testnet or mainnet use NewClient
			client, err := client.NewClient(
				context.Background(),
				false,
			)
	*/
	a := account.NewAccount(client, "alice", account.TestMnemonic, "")
	if a == nil {
		panic("ERROR CREATE ACCOUNT: NewAccount returned nil - check mnemonic and keyring initialization")
	}

	balanceClient := balance.NewBalanceMsg(*a)

	sendAmount := sdk.Coin{
		Amount: math.NewInt(20),
		Denom:  "usix",
	}

	res, err := balanceClient.SendBalance(BobAddress, sdk.NewCoins(sendAmount))
	if err != nil {
		fmt.Printf("Send error: %v\n", err)
		return
	}

	// Wait for SendBalance transaction to be confirmed
	err = client.WaitForTransaction(res.TxHash)
	if err != nil {
		fmt.Printf("Error waiting for SendBalance: %v\n", err)
		return
	}
	meta := metadata.NewMetadataMsg(*a, nftSchemaName)

	msgDeploySchema, err := meta.BuildDeployMsg()
		if err != nil {
		fmt.Printf("Deploy error: %v\n", err)
		return
	}

	msgCreateMetadata, err := meta.BuildMintMetadataMsg("1")
	if err != nil {
		fmt.Printf("Mint error: %v\n", err)
		return
	}

	var msgs []sdk.Msg

	msgs = append(msgs,msgDeploySchema,msgCreateMetadata)

	res, err = meta.BroadcastTx(msgs...)
	if err != nil {
		fmt.Printf("Mint error: %v\n", err)
	}
	err = client.WaitForTransaction(res.TxHash)
	if err != nil {
		fmt.Printf("Error waiting for deployment: %v\n", err)
		return
	}

	res, err = meta.FreezeCertificate("1")
	if err != nil {
		fmt.Printf("Freeze error: %v\n", err)
		return
	}

	err = client.WaitForTransaction(res.TxHash)
	if err != nil {
		fmt.Printf("Error waiting for deployment: %v\n", err)
		return
	}

	res, err = meta.UnfreezeCertificate("1")
	if err != nil {
		fmt.Printf("Unfreeze error: %v\n", err)
		return
	}
	fmt.Printf("Unfreeze response: %v\n", res)

	evm := evm.NewEVMClient(*a)
	address, tx, err := evm.DeployCertificateContract(contractName, contractSymbol, nftSchemaName)
	if err != nil {
		fmt.Printf("EVM error: %v\n", err)
		return
	}

	fmt.Printf("Deploy at: %v\n", tx.Hash())
	fmt.Printf("Deploy at Nonce: %v\n", tx.Nonce())
	fmt.Printf("Contract Address: %v\n", address)

	// Wait for deployment transaction to be mined
	_, err = client.WaitForEVMTransaction(tx.Hash())
	if err != nil {
		fmt.Printf("Error waiting for deployment: %v\n", err)
		return
	}

	tx, err = evm.MintCertificateNFT(address,1)
	if err != nil {
		fmt.Printf("EVM error: %v\n", err)
		return
	}
	fmt.Printf("Mint Tx: %+v \n", tx.Hash())
	fmt.Printf("Mint at Nonce: %v\n", tx.Nonce())

	_, err = client.WaitForEVMTransaction(tx.Hash())
	if err != nil {
		fmt.Printf("Error waiting for deployment: %v\n", err)
		return
	}

	tx, err = evm.TransferCertificateNFT(address, common.HexToAddress(BobEVMAddres), 1)
	if err != nil {
		fmt.Printf("EVM error: %v\n", err)
		return
	}
	fmt.Printf("Transfer Tx: %+v \n", tx.Hash())
	fmt.Printf("Transfer at Nonce: %v\n", tx.Nonce())

	_, err = client.WaitForEVMTransaction(tx.Hash())
	if err != nil {
		fmt.Printf("Error waiting for deployment: %v\n", err)
		return
	}

	currentOwner := evm.TokenOwner(address, 1)
	fmt.Printf("Current Owner: %+v \n", currentOwner)
}
