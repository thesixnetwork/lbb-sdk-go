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
	BobAddress       = "6x13g50hqdqsjk85fmgqz2h5xdxq49lsmjdwlemsp"
	BobEVMAddres     = "0x8a28fb81A084Ac7A276800957a19a6054BF86E4D"
	ChalieEVMAddress = "0xde609F435E82D1D5f71105CED56d06dDADB148B3"
	nftSchemaName    = "sixnetwork.lbbv05" // {ORGNAME}.{Schemacode}
	contractName     = "MyNFTCert"
	contractSymbol   = "Cert"
)

func main() {
	fmt.Println("\n========== Welcome to the LBB SDK-Go Quick Start Guide ==========")
	fmt.Println()

	// Step 1: Generate new wallet
	fmt.Println("\n========================================")
	fmt.Println("[Step 1] Generating a new wallet")
	fmt.Println("========================================")
	mnemonic, err := account.GenerateMnemonic()
	if err != nil {
		panic(fmt.Sprintf("ERROR: Failed to generate mnemonic: %v", err))
	}

	fmt.Println("SUCCESS: Mnemonic generated successfully")
	fmt.Println("\nIMPORTANT: Write this mnemonic phrase in a safe place.")
	fmt.Println("It is the only way to recover your account if you ever forget your password.")
	fmt.Printf("\nMnemonic: %s\n", mnemonic)
	fmt.Println("========================================")

	// Step 2: Initialize client (fivenet = testnet)
	fmt.Println("\n========================================")
	fmt.Println("[Step 2] Connecting to network")
	fmt.Println("========================================")
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
		panic(fmt.Sprintf("ERROR: Failed to create client: %v", err))
	}
	fmt.Println("SUCCESS: Connected to testnet (localhost)")
	fmt.Println("  RPC: http://localhost:26657")
	fmt.Println("  REST: http://localhost:1317")
	fmt.Println("  EVM: http://localhost:8545")
	fmt.Println("========================================")

	// Step 3: Create account from mnemonic
	fmt.Println("\n========================================")
	fmt.Println("[Step 3] Creating account from mnemonic")
	fmt.Println("========================================")
	acc, err := account.NewAccount(client, "alice", account.TestMnemonic, "")
	if err != nil {
		panic(fmt.Sprintf("ERROR: Failed to create account: %v", err))
	}

	fmt.Println("SUCCESS: Account created successfully")
	fmt.Printf("  Account Name: alice\n")
	fmt.Printf("  EVM Address: %s\n", acc.GetEVMAddress().Hex())
	fmt.Printf("  Cosmos Address: %s\n", acc.GetCosmosAddress().String())
	fmt.Println("========================================")

	// Step 4: Deploy Certificate Schema
	fmt.Println("\n========================================")
	fmt.Println("[Step 4] Deploying certificate schema")
	fmt.Println("========================================")

	meta, err := metadata.NewMetadataMsg(*acc, nftSchemaName)
	if err != nil {
		fmt.Printf("ERROR: Failed to create metadata message: %v\n", err)
		return
	}

	fmt.Println("Building deployment messages...")
	msgDeploySchema, err := meta.BuildDeployMsg()
	if err != nil {
		fmt.Printf("ERROR: Failed to build deploy message: %v\n", err)
		return
	}
	fmt.Println("  Deploy schema message built")

	msgCreateMetadata, err := meta.BuildMintMetadataMsg("1")
	if err != nil {
		fmt.Printf("ERROR: Failed to build metadata #1: %v\n", err)
		return
	}
	fmt.Println("  Mint metadata #1 message built")

	msgCreateMetadataWithInfo, err := meta.BuildMintMetadataWithInfoMsg("2", metadata.CertificateInfo{
		Status:       metadata.CertStatusType_ACTIVE,
		GoldStandard: "LBI",
		Weight:       "2000g",
		CertNumber:   "LBB_V1_01",
		CustomerID:   "LBB_V1_USER_01",
		IssueDate:    "Mon Dec 15 16:12:28 2025",
	})
	if err != nil {
		fmt.Printf("ERROR: Failed to build metadata #2 with info: %v\n", err)
		return
	}
	fmt.Println("  Mint metadata #2 with certificate info built")

	var msgs []sdk.Msg
	msgs = append(msgs, msgDeploySchema, msgCreateMetadata, msgCreateMetadataWithInfo)

	fmt.Println("\nBroadcasting transaction to blockchain...")
	res, err := meta.BroadcastTxAndWait(msgs...)
	if err != nil {
		fmt.Printf("ERROR: Failed to broadcast transaction: %v\n", err)
		return
	}

	fmt.Println("SUCCESS: Schema deployed successfully")
	fmt.Printf("  Schema Code: %s\n", nftSchemaName)
	fmt.Printf("  Transaction Hash: %s\n", res.TxHash)
	fmt.Println("========================================")
	// Step 5: Deploy EVM NFT Contract
	fmt.Println("\n========================================")
	fmt.Println("[Step 5] Deploying EVM NFT contract")
	fmt.Println("========================================")
	evmClient := evm.NewEVMClient(*acc)

	fmt.Printf("Contract Details:\n")
	fmt.Printf("  Name: %s\n", contractName)
	fmt.Printf("  Symbol: %s\n", contractSymbol)
	fmt.Printf("  Schema: %s\n", nftSchemaName)

	fmt.Println("\nDeploying contract...")
	contractAddress, tx, err := evmClient.DeployCertificateContract(contractName, contractSymbol, nftSchemaName)
	if err != nil {
		fmt.Printf("ERROR: Failed to deploy certificate contract: %v\n", err)
		return
	}

	fmt.Println("Waiting for transaction to be mined...")
	_, err = client.WaitForEVMTransaction(tx.Hash())
	if err != nil {
		fmt.Printf("ERROR: Error waiting for deployment: %v\n", err)
		return
	}

	fmt.Println("SUCCESS: Contract deployed successfully")
	fmt.Printf("  Contract Address: %s\n", contractAddress.Hex())
	fmt.Printf("  Transaction Hash: %s\n", tx.Hash().Hex())
	fmt.Println("========================================")

	// Step 6: Mint Certificate NFT
	fmt.Println("\n========================================")
	fmt.Println("[Step 6] Minting certificate NFTs")
	fmt.Println("========================================")

	// Mint NFT #1 to self
	tokenID := uint64(1)
	fmt.Printf("Minting NFT #%d to self...\n", tokenID)
	tx, err = evmClient.MintCertificateNFT(contractAddress, tokenID)
	if err != nil {
		fmt.Printf("ERROR: Failed to mint NFT: %v\n", err)
		return
	}
	fmt.Printf("  Transaction Hash: %s\n", tx.Hash().Hex())
	fmt.Printf("  Nonce: %v\n", tx.Nonce())

	fmt.Println("Waiting for transaction to be mined...")
	_, err = client.WaitForEVMTransaction(tx.Hash())
	if err != nil {
		fmt.Printf("ERROR: Error waiting for mint transaction: %v\n", err)
		return
	}

	fmt.Printf("SUCCESS: NFT #%d minted successfully to owner\n\n", tokenID)

	// Mint NFT #2 to Bob
	fmt.Printf("Minting NFT #%d to Bob's address...\n", tokenID+1)
	fmt.Printf("  Destination: %s\n", BobEVMAddres)
	tx, err = evmClient.MintCertificateNFTToDestination(contractAddress, tokenID+1, common.HexToAddress(BobEVMAddres))
	if err != nil {
		fmt.Printf("ERROR: Failed to mint NFT to destination: %v\n", err)
		return
	}
	fmt.Printf("  Transaction Hash: %s\n", tx.Hash().Hex())
	fmt.Printf("  Nonce: %v\n", tx.Nonce())

	fmt.Println("Waiting for transaction to be mined...")
	_, err = client.WaitForEVMTransaction(tx.Hash())
	if err != nil {
		fmt.Printf("ERROR: Error waiting for mint transaction: %v\n", err)
		return
	}

	fmt.Printf("SUCCESS: NFT #%d minted successfully to Bob\n", tokenID+1)
	fmt.Println("========================================")

	// Step 7: Change certificate state (Freeze/Unfreeze)
	fmt.Println("\n========================================")
	fmt.Println("[Step 7] Managing certificate state")
	fmt.Println("========================================")

	fmt.Println("Freezing certificate #1...")
	res, err = meta.FreezeCertificate("1")
	if err != nil {
		fmt.Printf("ERROR: Failed to freeze certificate: %v\n", err)
		return
	}

	fmt.Println("Waiting for freeze transaction to be confirmed...")
	err = client.WaitForTransaction(res.TxHash)
	if err != nil {
		fmt.Printf("ERROR: Error waiting for freeze transaction: %v\n", err)
		return
	}

	fmt.Println("SUCCESS: Certificate #1 frozen successfully")
	fmt.Printf("  Transaction Hash: %s\n\n", res.TxHash)

	fmt.Println("Unfreezing certificate #1...")
	res, err = meta.UnfreezeCertificate("1")
	if err != nil {
		fmt.Printf("ERROR: Failed to unfreeze certificate: %v\n", err)
		return
	}

	fmt.Println("Waiting for unfreeze transaction to be confirmed...")
	err = client.WaitForTransaction(res.TxHash)
	if err != nil {
		fmt.Printf("ERROR: Error waiting for unfreeze transaction: %v\n", err)
		return
	}

	fmt.Println("SUCCESS: Certificate #1 unfrozen successfully")
	fmt.Printf("  Transaction Hash: %s\n", res.TxHash)
	fmt.Println("========================================")

	// Step 8: Transfer NFT
	fmt.Println("\n========================================")
	fmt.Println("[Step 8] Transferring NFT directly")
	fmt.Println("========================================")

	fmt.Printf("Transferring NFT #%d to Bob...\n", tokenID)
	fmt.Printf("  From: %s\n", acc.GetEVMAddress().Hex())
	fmt.Printf("  To: %s\n", BobEVMAddres)

	tx, err = evmClient.TransferCertificateNFT(contractAddress, common.HexToAddress(BobEVMAddres), tokenID)
	if err != nil {
		fmt.Printf("ERROR: Failed to transfer NFT: %v\n", err)
		return
	}
	fmt.Printf("  Transaction Hash: %s\n", tx.Hash().Hex())
	fmt.Printf("  Nonce: %v\n", tx.Nonce())

	fmt.Println("Waiting for transfer transaction to be mined...")
	_, err = client.WaitForEVMTransaction(tx.Hash())
	if err != nil {
		fmt.Printf("ERROR: Error waiting for transfer: %v\n", err)
		return
	}

	fmt.Println("SUCCESS: NFT transferred successfully")
	fmt.Println("========================================")

	// Step 8.1: Meta-transaction (gasless transfer)
	fmt.Println("\n========================================")
	fmt.Println("[Step 8.1] Gasless Transfer Demo")
	fmt.Println("Admin pays gas for user's NFT transfer")
	fmt.Println("========================================")

	fmt.Println("\nProcess Overview:")
	fmt.Println("  1. Create new user account")
	fmt.Println("  2. Mint NFT to new user")
	fmt.Println("  3. User signs transfer offline (no gas)")
	fmt.Println("  4. Admin broadcasts signed transaction")
	fmt.Println()

	// Step 8.1.1: Create new account
	fmt.Println("┌─ [Step 8.1.1] Creating new user account")
	accFromGenMnemonic, err := account.NewAccount(client, "new_gen", mnemonic, "")
	if err != nil {
		fmt.Printf("   ERROR: Failed to create account: %v\n", err)
		return
	}
	fmt.Println("└─ SUCCESS: Account created")
	fmt.Printf("      Name: new_gen\n")
	fmt.Printf("      Address: %s\n\n", accFromGenMnemonic.GetEVMAddress().Hex())

	// Step 8.1.2: Mint NFT to new account
	fmt.Printf("┌─ [Step 8.1.2] Minting NFT #%d to new account\n", tokenID+2)
	tx, err = evmClient.MintCertificateNFTToDestination(contractAddress, tokenID+2, accFromGenMnemonic.GetEVMAddress())
	if err != nil {
		fmt.Printf("   ERROR: Failed to mint NFT: %v\n", err)
		return
	}
	fmt.Printf("   Transaction Hash: %s\n", tx.Hash().Hex())
	fmt.Printf("   Recipient: %s\n", accFromGenMnemonic.GetEVMAddress().Hex())

	fmt.Println("   Waiting for mint transaction...")
	_, err = client.WaitForEVMTransaction(tx.Hash())
	if err != nil {
		fmt.Printf("   ERROR: Error waiting for mint: %v\n", err)
		return
	}
	fmt.Printf("└─ SUCCESS: NFT #%d minted to new user\n\n", tokenID+2)

	// Step 8.1.3: User signs transaction offline
	fmt.Println("┌─ [Step 8.1.3] User signs transfer offline (no gas required)")
	offlineEVMClient := evm.NewEVMClient(*accFromGenMnemonic)

	fmt.Printf("   Signing transfer of NFT #%d\n", tokenID+2)
	fmt.Printf("   Destination: %s\n", ChalieEVMAddress)
	signedTx, err := offlineEVMClient.SignTransferNFT(contractAddress, common.HexToAddress(ChalieEVMAddress), tokenID+2)
	if err != nil {
		fmt.Printf("   ERROR: Failed to sign offline transaction: %v\n", err)
		return
	}
	fmt.Println("└─ SUCCESS: Transaction signed by user")
	fmt.Println()

	// Step 8.1.4: Admin broadcasts transaction
	fmt.Println("┌─ [Step 8.1.4] Admin broadcasts signed transaction")
	fmt.Println("   NOTE: Admin pays for gas")
	err = evmClient.SendTransaction(signedTx)
	if err != nil {
		fmt.Printf("   ERROR: Failed to send transaction: %v\n", err)
		return
	}

	fmt.Println("   Waiting for transaction confirmation...")
	_, err = client.WaitForEVMTransaction(signedTx.Hash())
	if err != nil {
		fmt.Printf("   ERROR: Error waiting for transfer: %v\n", err)
		return
	}

	fmt.Println("└─ SUCCESS: Gasless transfer completed successfully")
	fmt.Printf("\n   Transfer Summary:\n")
	fmt.Printf("      NFT #%d transferred\n", tokenID+2)
	fmt.Printf("      From (User): %s\n", accFromGenMnemonic.GetEVMAddress().Hex())
	fmt.Printf("      To (Charlie): %s\n", ChalieEVMAddress)
	fmt.Printf("      Gas paid by: Admin (%s)\n", acc.GetEVMAddress().Hex())
	fmt.Printf("      Transaction Hash: %s\n", signedTx.Hash().Hex())
	fmt.Println("========================================")

	// Step 9: Verify ownership
	fmt.Println("\n========================================")
	fmt.Println("[Step 9] Verifying final ownership")
	fmt.Println("========================================")

	fmt.Printf("Checking owner of NFT #%d...\n", tokenID+2)
	currentOwner := evmClient.TokenOwner(contractAddress, tokenID+2)
	fmt.Printf("Current owner: %s\n", currentOwner.Hex())
	fmt.Printf("  (Should be Charlie: %s)\n", ChalieEVMAddress)
	fmt.Println("========================================")

	// Summary
	fmt.Println("\n╔════════════════════════════════════════╗")
	fmt.Println("║          EXECUTION SUMMARY             ║")
	fmt.Println("╚════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("Schema Information:\n")
	fmt.Printf("  Schema Code: %s\n\n", nftSchemaName)

	fmt.Printf("Contract Information:\n")
	fmt.Printf("  Address: %s\n", contractAddress.Hex())
	fmt.Printf("  Name: %s\n", contractName)
	fmt.Printf("  Symbol: %s\n\n", contractSymbol)

	fmt.Printf("NFTs Created:\n")
	fmt.Printf("  NFT #1: Transferred to Bob\n")
	fmt.Printf("  NFT #2: Minted to Bob\n")
	fmt.Printf("  NFT #3: Gasless transfer to Charlie\n\n")

	fmt.Printf("Final Owner of NFT #%d:\n", tokenID+2)
	fmt.Printf("  Address: %s\n", currentOwner.Hex())
	fmt.Println()
	fmt.Println("Quick start completed successfully!")
	fmt.Println("╔════════════════════════════════════════╗")
}
