package main

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/client"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/evm"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/metadata"
)

const (
	// Configuration
	contractName   = "MyCertificate"
	contractSymbol = "CERT"
	schemaName     = "myorg.lbbv01" // Format: {ORGNAME}.{SchemaCode}

	// Recipient addresses for testing
	recipientCosmos = "6x13g50hqdqsjk85fmgqz2h5xdxq49lsmjdwlemsp"
	recipientEVM    = "0x8a28fb81A084Ac7A276800957a19a6054BF86E4D"
)

func main() {
	fmt.Println("=== LBB SDK Go - Quick Start Example ===\n")

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
	client, err := client.NewClient(ctx, false)
	if err != nil {
		panic(fmt.Sprintf("Failed to create client: %v", err))
	}
	fmt.Println("Connected to fivenet (testnet)")
	fmt.Println()

	// Step 3: Create account from mnemonic
	fmt.Println("Step 3: Creating account...")
	acc := account.NewAccount(client, "quickstart", mnemonic, "")
	if acc == nil {
		panic("Failed to create account")
	}

	fmt.Printf("Account created\n")
	fmt.Printf("  EVM Address: %s\n", acc.GetEVMAddress().Hex())
	fmt.Printf("  Cosmos Address: %s\n\n", acc.GetCosmosAddress().String())

	// Step 4: Deploy Certificate Schema
	fmt.Println("Step 4: Deploying certificate schema...")
	meta := metadata.NewMetadataMsg(*acc, schemaName)
	msgDeploySchema, err := meta.BuildDeployMsg()
	if err != nil {
		panic(fmt.Sprintf("Failed to build deploy message: %v", err))
	}

	res, err := meta.BroadcastTx(msgDeploySchema)
	if err != nil {
		panic(fmt.Sprintf("Failed to deploy schema: %v", err))
	}

	err = client.WaitForTransaction(res.TxHash)
	if err != nil {
		panic(fmt.Sprintf("Error waiting for schema deployment: %v", err))
	}

	fmt.Printf("Schema deployed\n")
	fmt.Printf("  Schema Code: %s\n", schemaName)
	fmt.Printf("  Transaction: %s\n\n", res.TxHash)

	// Step 5: Deploy EVM NFT Contract
	fmt.Println("Step 5: Deploying EVM NFT contract...")
	evmClient := evm.NewEVMClient(*acc)

	contractAddress, tx, err := evmClient.DeployCertificateContract(
		contractName,
		contractSymbol,
		schemaName,
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to deploy contract: %v", err))
	}

	_, err = client.WaitForEVMTransaction(tx.Hash())
	if err != nil {
		panic(fmt.Sprintf("Error waiting for contract deployment: %v", err))
	}

	fmt.Printf("Contract deployed\n")
	fmt.Printf("  Contract Address: %s\n", contractAddress.Hex())
	fmt.Printf("  Transaction: %s\n\n", tx.Hash().Hex())

	// Step 6: Mint Certificate NFT
	fmt.Println("Step 6: Minting certificate NFT...")
	tokenId := uint64(1)

	mintTx, err := evmClient.MintCertificateNFT(contractAddress, tokenId)
	if err != nil {
		panic(fmt.Sprintf("Failed to mint NFT: %v", err))
	}

	_, err = client.WaitForEVMTransaction(mintTx.Hash())
	if err != nil {
		panic(fmt.Sprintf("Error waiting for mint: %v", err))
	}

	fmt.Printf("NFT minted\n")
	fmt.Printf("  Token ID: %d\n", tokenId)
	fmt.Printf("  Transaction: %s\n\n", mintTx.Hash().Hex())

	// Step 7: Create Certificate Metadata
	fmt.Println("Step 7: Creating certificate metadata...")
	tokenIdStr := "1"

	msgCreateMetadata, err := meta.BuildMintMetadataMsg(tokenIdStr)
	if err != nil {
		panic(fmt.Sprintf("Failed to build mint metadata message: %v", err))
	}

	res, err = meta.BroadcastTx(msgCreateMetadata)
	if err != nil {
		panic(fmt.Sprintf("Failed to create metadata: %v", err))
	}

	err = client.WaitForTransaction(res.TxHash)
	if err != nil {
		panic(fmt.Sprintf("Error waiting for metadata creation: %v", err))
	}

	fmt.Printf("Metadata created\n")
	fmt.Printf("  Token ID: %s\n", tokenIdStr)
	fmt.Printf("  Transaction: %s\n\n", res.TxHash)

	// Step 8: Transfer NFT
	fmt.Println("Step 8: Transferring NFT to recipient...")
	transferTx, err := evmClient.TransferCertificateNFT(
		contractAddress,
		common.HexToAddress(recipientEVM),
		tokenId,
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to transfer NFT: %v", err))
	}

	_, err = client.WaitForEVMTransaction(transferTx.Hash())
	if err != nil {
		panic(fmt.Sprintf("Error waiting for transfer: %v", err))
	}

	fmt.Printf("NFT transferred\n")
	fmt.Printf("  To: %s\n", recipientEVM)
	fmt.Printf("  Transaction: %s\n\n", transferTx.Hash().Hex())

	// Step 9: Verify ownership
	fmt.Println("Step 9: Verifying new owner...")
	currentOwner := evmClient.TokenOwner(contractAddress, tokenId)
	fmt.Printf("Current owner: %s\n\n", currentOwner.Hex())

	// Summary
	fmt.Println("=== Summary ===")
	fmt.Printf("Schema Code: %s\n", schemaName)
	fmt.Printf("Contract Address: %s\n", contractAddress.Hex())
	fmt.Printf("Token ID: %d\n", tokenId)
	fmt.Printf("Current Owner: %s\n", currentOwner.Hex())
	fmt.Println("\nQuick start completed successfully!")
}
