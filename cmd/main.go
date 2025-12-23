package main

import (
	"context"
	"fmt"
	"math/big"
	"time"

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
	nftSchemaName    = "sixnetwork.lbbv01" // {ORGNAME}.{Schemacode}
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

	// Step 8.1: EIP-2612 Permit (gasless transfer)
	fmt.Println("\n========================================")
	fmt.Println("[Step 8.1] Gasless Transfer Demo (EIP-2612 Permit)")
	fmt.Println("Admin pays gas for user's NFT transfer")
	fmt.Println("========================================")

	fmt.Println("\nProcess Overview:")
	fmt.Println("  1. Create new user account (no funds needed!)")
	fmt.Println("  2. Mint NFT to new user")
	fmt.Println("  3. User signs EIP-712 permit message offline (no gas, no blockchain interaction)")
	fmt.Println("  4. Admin broadcasts transfer with permit (admin pays all gas)")
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
	fmt.Printf("      Address: %s\n", accFromGenMnemonic.GetEVMAddress().Hex())
	fmt.Printf("      Balance: 0 (no funds needed!)\n\n")

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

	// Step 8.1.3: User signs EIP-712 permit offline
	fmt.Println("┌─ [Step 8.1.3] User signs EIP-712 permit message (completely offline)")
	offlineEVMClient := evm.NewEVMClient(*accFromGenMnemonic)

	fmt.Printf("   User signs permit for NFT #%d transfer\n", tokenID+2)
	fmt.Printf("   From (User): %s\n", accFromGenMnemonic.GetEVMAddress().Hex())
	fmt.Printf("   To (Charlie): %s\n", ChalieEVMAddress)
	fmt.Println("   NOTE: This is just a signature, NO transaction, NO gas needed!")

	// Set deadline to 1 hour from now (Unix timestamp)
	deadline := big.NewInt(time.Now().Unix() + 3600)
	permitSig, err := offlineEVMClient.SignPermit(
		contractName,
		contractAddress,
		acc.GetEVMAddress(), // Spender (admin/relay)
		big.NewInt(int64(tokenID+2)),
		deadline,
	)
	if err != nil {
		fmt.Printf("   ERROR: Failed to sign permit: %v\n", err)
		return
	}
	fmt.Println("└─ SUCCESS: Permit signed offline (user never touched blockchain!)")
	fmt.Println()

	// Step 8.1.4: Admin broadcasts transfer with permit
	fmt.Println("┌─ [Step 8.1.4] Admin broadcasts transfer using user's permit")
	fmt.Println("   NOTE: Admin pays ALL gas, user pays NOTHING")
	fmt.Printf("   Transferring NFT #%d with permit\n", tokenID+2)

	transferTx, err := evmClient.TransferWithPermit(
		contractAddress,
		accFromGenMnemonic.GetEVMAddress(),    // From
		common.HexToAddress(ChalieEVMAddress), // To
		big.NewInt(int64(tokenID+2)),
		permitSig,
	)
	if err != nil {
		fmt.Printf("   ERROR: Failed to execute transfer with permit: %v\n", err)
		return
	}

	fmt.Printf("   Transaction Hash: %s\n", transferTx.Hash().Hex())
	fmt.Println("   Waiting for confirmation...")
	_, err = client.WaitForEVMTransaction(transferTx.Hash())
	if err != nil {
		fmt.Printf("   ERROR: Error waiting for transfer: %v\n", err)
		return
	}

	fmt.Println("└─ SUCCESS: Gasless transfer completed!")
	fmt.Printf("\n   Transfer Summary:\n")
	fmt.Printf("      NFT #%d transferred\n", tokenID+2)
	fmt.Printf("      Owner (User): %s (paid 0 gas! 🎉)\n", accFromGenMnemonic.GetEVMAddress().Hex())
	fmt.Printf("      Destination: %s\n", ChalieEVMAddress)
	fmt.Printf("      Gas paid by: Admin (%s)\n", acc.GetEVMAddress().Hex())
	fmt.Printf("      Transaction Hash: %s\n", transferTx.Hash().Hex())
	fmt.Printf("      Method: EIP-2612 Permit (transferWithPermit)\n")
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

	// Step 10: Burn NFT directly
	fmt.Println("\n========================================")
	fmt.Println("[Step 10] Burning NFT directly")
	fmt.Println("========================================")

	// First mint a new NFT to burn
	burnTokenID := tokenID + 3
	fmt.Printf("Minting NFT #%d to burn...\n", burnTokenID)
	tx, err = evmClient.MintCertificateNFT(contractAddress, burnTokenID)
	if err != nil {
		fmt.Printf("ERROR: Failed to mint NFT for burning: %v\n", err)
		return
	}
	fmt.Printf("  Transaction Hash: %s\n", tx.Hash().Hex())

	fmt.Println("Waiting for mint transaction...")
	_, err = client.WaitForEVMTransaction(tx.Hash())
	if err != nil {
		fmt.Printf("ERROR: Error waiting for mint: %v\n", err)
		return
	}
	fmt.Printf("SUCCESS: NFT #%d minted\n\n", burnTokenID)

	// Now burn it
	fmt.Printf("Burning NFT #%d...\n", burnTokenID)
	fmt.Printf("  Owner: %s\n", acc.GetEVMAddress().Hex())

	burnTx, err := evmClient.BurnCertificateNFT(contractAddress, burnTokenID)
	if err != nil {
		fmt.Printf("ERROR: Failed to burn NFT: %v\n", err)
		return
	}
	fmt.Printf("  Transaction Hash: %s\n", burnTx.Hash().Hex())
	fmt.Printf("  Nonce: %v\n", burnTx.Nonce())

	fmt.Println("Waiting for burn transaction to be mined...")
	_, err = client.WaitForEVMTransaction(burnTx.Hash())
	if err != nil {
		fmt.Printf("ERROR: Error waiting for burn: %v\n", err)
		return
	}

	fmt.Printf("SUCCESS: NFT #%d burned successfully\n", burnTokenID)

	// Verify the token was burned by checking owner
	fmt.Println("\nVerifying burn...")
	fmt.Printf("Querying owner of burned NFT #%d...\n", burnTokenID)
	burnedOwner := evmClient.TokenOwner(contractAddress, burnTokenID)
	fmt.Printf("  Owner address: %s\n", burnedOwner.Hex())

	zeroAddress := common.HexToAddress("0x0000000000000000000000000000000000000000")
	if burnedOwner == zeroAddress {
		fmt.Println("  ✓ VERIFIED: Token burned successfully (owner is zero address)")
	} else {
		fmt.Printf("  ⚠ WARNING: Token still has owner: %s\n", burnedOwner.Hex())
	}
	fmt.Println("========================================")

	// Step 11: Gasless Burn with Permit
	fmt.Println("\n========================================")
	fmt.Println("[Step 11] Gasless Burn Demo (Burn with Permit)")
	fmt.Println("Admin pays gas for user's NFT burn")
	fmt.Println("========================================")

	fmt.Println("\nProcess Overview:")
	fmt.Println("  1. Mint NFT to user account (the one we created earlier)")
	fmt.Println("  2. User signs EIP-712 permit message offline (no gas)")
	fmt.Println("  3. Admin broadcasts burn with permit (admin pays all gas)")
	fmt.Println()

	// Step 11.1: Mint NFT to user account
	gaslessBurnTokenID := burnTokenID + 1
	fmt.Printf("┌─ [Step 11.1] Minting NFT #%d to user account for gasless burn\n", gaslessBurnTokenID)
	tx, err = evmClient.MintCertificateNFTToDestination(contractAddress, gaslessBurnTokenID, accFromGenMnemonic.GetEVMAddress())
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
	fmt.Printf("└─ SUCCESS: NFT #%d minted to user\n\n", gaslessBurnTokenID)

	// Step 11.2: User signs EIP-712 permit for burn
	fmt.Println("┌─ [Step 11.2] User signs EIP-712 permit for burn (completely offline)")
	offlineEVMClientForBurn := evm.NewEVMClient(*accFromGenMnemonic)

	fmt.Printf("   User signs permit for NFT #%d burn\n", gaslessBurnTokenID)
	fmt.Printf("   Owner (User): %s\n", accFromGenMnemonic.GetEVMAddress().Hex())
	fmt.Println("   NOTE: This is just a signature, NO transaction, NO gas needed!")

	// Set deadline to 1 hour from now
	burnDeadline := big.NewInt(time.Now().Unix() + 3600)
	burnPermitSig, err := offlineEVMClientForBurn.SignPermit(
		contractName,
		contractAddress,
		acc.GetEVMAddress(), // Spender (admin/relay)
		big.NewInt(int64(gaslessBurnTokenID)),
		burnDeadline,
	)
	if err != nil {
		fmt.Printf("   ERROR: Failed to sign burn permit: %v\n", err)
		return
	}
	fmt.Println("└─ SUCCESS: Burn permit signed offline")
	fmt.Println()

	// Step 11.3: Admin broadcasts burn with permit
	fmt.Println("┌─ [Step 11.3] Admin broadcasts burn using user's permit")
	fmt.Println("   NOTE: Admin pays ALL gas, user pays NOTHING")
	fmt.Printf("   Burning NFT #%d with permit\n", gaslessBurnTokenID)

	burnWithPermitTx, err := evmClient.BurnWithPermit(
		contractAddress,
		accFromGenMnemonic.GetEVMAddress(), // From (owner)
		big.NewInt(int64(gaslessBurnTokenID)),
		burnPermitSig,
	)
	if err != nil {
		fmt.Printf("   ERROR: Failed to execute burn with permit: %v\n", err)
		return
	}

	fmt.Printf("   Transaction Hash: %s\n", burnWithPermitTx.Hash().Hex())
	fmt.Println("   Waiting for confirmation...")
	_, err = client.WaitForEVMTransaction(burnWithPermitTx.Hash())
	if err != nil {
		fmt.Printf("   ERROR: Error waiting for burn: %v\n", err)
		return
	}

	fmt.Println("└─ SUCCESS: Gasless burn completed!")
	fmt.Printf("\n   Burn Summary:\n")
	fmt.Printf("      NFT #%d burned\n", gaslessBurnTokenID)
	fmt.Printf("      Owner (User): %s (paid 0 gas! 🎉)\n", accFromGenMnemonic.GetEVMAddress().Hex())
	fmt.Printf("      Gas paid by: Admin (%s)\n", acc.GetEVMAddress().Hex())
	fmt.Printf("      Transaction Hash: %s\n", burnWithPermitTx.Hash().Hex())
	fmt.Printf("      Method: EIP-2612 Permit (burnWithPermit)\n")

	// Verify the token was burned by checking owner
	fmt.Println("\n   Verifying burn...")
	fmt.Printf("   Querying owner of burned NFT #%d...\n", gaslessBurnTokenID)
	gaslessBurnedOwner := evmClient.TokenOwner(contractAddress, gaslessBurnTokenID)
	fmt.Printf("   Owner address: %s\n", gaslessBurnedOwner.Hex())

	zeroAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")
	if gaslessBurnedOwner == zeroAddress {
		fmt.Println("   ✓ VERIFIED: Token burned successfully (owner is zero address)")
	} else {
		fmt.Printf("   ⚠ WARNING: Token still has owner: %s\n", gaslessBurnedOwner.Hex())
	}
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
	fmt.Printf("  NFT #3: Gasless transfer to Charlie\n")
	fmt.Printf("  NFT #4: Burned directly by owner\n")
	fmt.Printf("  NFT #5: Gasless burn via permit\n\n")

	fmt.Printf("Final Owner of NFT #%d:\n", tokenID+2)
	fmt.Printf("  Address: %s\n", currentOwner.Hex())
	fmt.Println()
	fmt.Println("Quick start completed successfully!")
	fmt.Println("╔════════════════════════════════════════╗")
}
