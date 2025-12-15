package main

import (
	"context"
	"fmt"

	_ "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/client"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/evm"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/metadata"
)

const (
	BobAddress     = "6x13g50hqdqsjk85fmgqz2h5xdxq49lsmjdwlemsp"
	BobEVMAddres   = "0x8a28fb81A084Ac7A276800957a19a6054BF86E4D"
	nftSchemaName  = "sixnetwork.lbbv05" // {ORGNAME}.{Schemacode}
	contractName   = "MyNFTCert"
	contractSymbol = "Cert"
)

func main() {
	fmt.Println("=== LBB SDK Go - Quick Start Example ===")
	fmt.Println()

	// Step 1: Generate new wallet
	fmt.Println("Step 1: Generating new wallet...")
	fmt.Println("-----------------------------------------------------")
	fmt.Println()
	fmt.Println()
	mnemonic, err := account.GenerateMnemonic()
	if err != nil {
		panic(fmt.Sprintf("Failed to generate mnemonic: %v", err))
	}

	fmt.Println("Mnemonic generated")
	fmt.Println("*Important** write this mnemonic phrase in a safe place.")
	fmt.Println("It is the only way to recover your account if you ever forget your password.")
	fmt.Printf("\nMnemonic: %s\n\n", mnemonic)
	fmt.Println("-----------------------------------------------------")

	// Step 2: Initialize client (fivenet = testnet)
	fmt.Println("Step 2: Connecting to network...")
	ctx := context.Background()
	// client, err := client.NewClient(ctx, false)
	client, err := client.NewCustomClient(
		ctx,
		"http://localhost:26657",
		"http://localhost:1317",
		"http://localhost:8545",
		"testnet",
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to create client: %v", err))
	}
	fmt.Println("Connected to fivenet (testnet)")
	fmt.Println()

	// Step 3: Create account from mnemonic
	fmt.Println("Step 3: Creating account...")
	acc, err := account.NewAccount(client, "alice", account.TestMnemonic, "")
	if err != nil {
		panic("ERROR CREATE ACCOUNT: NewAccount returned nil - check mnemonic and keyring initialization")
	}

	fmt.Printf("Account created\n")
	fmt.Printf("  EVM Address: %s\n", acc.GetEVMAddress().Hex())
	fmt.Printf("  Cosmos Address: %s\n\n", acc.GetCosmosAddress().String())

	// Step 4: Deploy Certificate Schema
	fmt.Println("Step 4: Deploying certificate schema...")

	meta, err := metadata.NewMetadataMsg(*acc, nftSchemaName)
	if err != nil {
		fmt.Printf("NewMetadataMsg error: %v\n", err)
		return
	}

	msgDeploySchema, err := meta.BuildDeployMsg()
	if err != nil {
		fmt.Printf("Failed to build deploy message: %v\n", err)
		return
	}

	msgCreateMetadata, err := meta.BuildMintMetadataMsg("1")
	if err != nil {
		fmt.Printf("Failed to build create metadata: %v\n", err)
		return
	}

	msgCreateMetadataWithInfo, err := meta.BuildMintMetadataWithInfoMsg("2", metadata.CertificateInfo{
		Status:       "TCI",
		GoldStandard: "LBI",
		Weight:       "2000g",
		CertNumber:   "LBB_V1_01",
		CustomerID:   "LBB_V1_USER_01",
		IssueDate:    "Mon Dec 15 16:12:28 2025",
	})
	if err != nil {
		fmt.Printf("Failed to build create metadata: %v\n", err)
		return
	}

	var msgs []sdk.Msg

	msgs = append(msgs, msgDeploySchema, msgCreateMetadata, msgCreateMetadataWithInfo)

	res, err := meta.BroadcastTxAndWait(msgs...)
	if err != nil {
		fmt.Printf("Broadcast Tx error: %v\n", err)
	}

	fmt.Printf("Schema deployed\n")
	fmt.Printf("  Schema Code: %s\n", nftSchemaName)
	fmt.Printf("  Transaction: %s\n\n", res.TxHash)
	// Step 5: Deploy EVM NFT Contract
	fmt.Println("Step 5: Deploying EVM NFT contract...")
	evmClient := evm.NewEVMClient(*acc)
	contractAddress, tx, err := evmClient.DeployCertificateContract(contractName, contractSymbol, nftSchemaName)
	if err != nil {
		fmt.Printf("EVM deploy certificate erro: %v\n", err)
		return
	}

	// Wait for deployment transaction to be mined
	_, err = client.WaitForEVMTransaction(tx.Hash())
	if err != nil {
		fmt.Printf("Error waiting for deployment: %v\n", err)
		return
	}

	fmt.Printf("Contract deployed\n")
	fmt.Printf("  Contract Address: %s\n", contractAddress.Hex())
	fmt.Printf("  Transaction: %s\n\n", tx.Hash().Hex())

	// Step 6: Mint Certificate NFT
	fmt.Println("Step 6: Minting certificate NFT...")
	tokenId := uint64(1)
	tx, err = evmClient.MintCertificateNFT(contractAddress, tokenId)
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

	fmt.Printf("NFT minted\n")
	fmt.Printf("  Token ID: %d\n", tokenId)
	fmt.Printf("  Transaction: %s\n\n", tx.Hash().Hex())

	tx, err = evmClient.MintCertificateNFTToDestination(contractAddress, tokenId+1, common.HexToAddress(BobEVMAddres))
	if err != nil {
		fmt.Printf("EVM error: %v\n", err)
		return
	}
	fmt.Printf("Mint Tx to destinatin address: %+v \n", tx.Hash())
	fmt.Printf("Mint at Nonce: %v\n", tx.Nonce())

	_, err = client.WaitForEVMTransaction(tx.Hash())
	if err != nil {
		fmt.Printf("Error waiting for deployment: %v\n", err)
		return
	}

	fmt.Printf("NFT minted\n")
	fmt.Printf("  Token ID: %d\n", tokenId)
	fmt.Printf("  Transaction: %s\n\n", tx.Hash().Hex())

	// Step 7: Tryto change state of metadata

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

	// Step 8: Transfer NFT
	fmt.Println("Step 8: Transferring NFT to recipient...")

	tx, err = evmClient.TransferCertificateNFT(contractAddress, common.HexToAddress(BobEVMAddres), tokenId)
	if err != nil {
		fmt.Printf("EVM error: %v\n", err)
		return
	}
	fmt.Printf("Transfer Tx: %+v \n", tx.Hash())
	fmt.Printf("Transfer at Nonce: %v\n", tx.Nonce())

	_, err = client.WaitForEVMTransaction(tx.Hash())
	if err != nil {
		fmt.Printf("Error waiting for transfer: %v\n", err)
		return
	}

	fmt.Printf("NFT transferred\n")
	fmt.Printf("  To: %s\n", BobEVMAddres)
	fmt.Printf("  Transaction: %s\n\n", tx.Hash().Hex())

	// Step 9: Verify ownership
	fmt.Println("Step 9: Verifying new owner...")
	currentOwner := evmClient.TokenOwner(contractAddress, tokenId)
	fmt.Printf("Current owner: %s\n\n", currentOwner.Hex())

	// Summary
	fmt.Println("=== Summary ===")
	fmt.Printf("Schema Code: %s\n", nftSchemaName)
	fmt.Printf("Contract Address: %s\n", contractAddress.Hex())
	fmt.Printf("Token ID: %d\n", tokenId)
	fmt.Printf("Current Owner: %s\n", currentOwner.Hex())
	fmt.Println("\nQuick start completed successfully!")
}
