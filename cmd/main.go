package main

import (
	"context"
	"math/big"
	"time"

	_ "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/client"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/evm"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/logger"
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
	// Configure logger
	logger.SetTime(false)
	logger.SetColors(true)

	logger.Info("========== Welcome to the LBB SDK-Go Quick Start Guide ==========")

	// Step 1: Generate new wallet
	logger.Info("========================================")
	logger.Info("[Step 1] Generating a new wallet")
	logger.Info("========================================")
	mnemonic, err := account.GenerateMnemonic()
	if err != nil {
		logger.Fatal("Failed to generate mnemonic: %v", err)
	}

	logger.Info("SUCCESS: Mnemonic generated successfully")
	logger.Info("IMPORTANT: Write this mnemonic phrase in a safe place.")
	logger.Info("It is the only way to recover your account if you ever forget your password.")
	logger.Info("Mnemonic: %s", mnemonic)
	logger.Info("========================================")

	// Step 2: Initialize client (fivenet = testnet)
	logger.Info("========================================")
	logger.Info("[Step 2] Connecting to network")
	logger.Info("========================================")
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
		logger.Fatal("Failed to create client: %v", err)
	}
	logger.Info("SUCCESS: Connected to testnet (localhost)")
	logger.Info("  RPC: http://localhost:26657")
	logger.Info("  REST: http://localhost:1317")
	logger.Info("  EVM: http://localhost:8545")
	logger.Info("========================================")

	// Step 3: Create account from mnemonic
	logger.Info("========================================")
	logger.Info("[Step 3] Creating account from mnemonic")
	logger.Info("========================================")
	acc, err := account.NewAccount(client, "alice", account.TestMnemonic, "")
	if err != nil {
		logger.Fatal("Failed to create account: %v", err)
	}
	defer acc.Close()

	logger.Info("SUCCESS: Account created successfully")
	logger.Info("  Account Name: alice")
	logger.Info("  EVM Address: %s", acc.GetEVMAddress().Hex())
	logger.Info("  Cosmos Address: %s", acc.GetCosmosAddress().String())
	logger.Info("========================================")

	// Step 4: Deploy Certificate Schema
	logger.Info("========================================")
	logger.Info("[Step 4] Deploying certificate schema")
	logger.Info("========================================")

	meta, err := metadata.NewMetadataMsg(*acc, nftSchemaName)
	if err != nil {
		logger.Error("Failed to create metadata message: %v", err)
		return
	}

	logger.Info("Building deployment messages...")
	msgDeploySchema, err := meta.BuildDeployMsg()
	if err != nil {
		logger.Error("Failed to build deploy message: %v", err)
		return
	}
	logger.Info("  Deploy schema message built")

	msgCreateMetadata, err := meta.BuildMintMetadataMsg("1")
	if err != nil {
		logger.Error("Failed to build metadata #1: %v", err)
		return
	}
	logger.Info("  Mint metadata #1 message built")

	msgCreateMetadataWithInfo, err := meta.BuildMintMetadataWithInfoMsg("2", metadata.CertificateInfo{
		Status:       metadata.CertStatusType_ACTIVE,
		GoldStandard: "LBI",
		Weight:       "2000g",
		CertNumber:   "LBB_V1_01",
		CustomerID:   "LBB_V1_USER_01",
		IssueDate:    "Mon Dec 15 16:12:28 2025",
	})
	if err != nil {
		logger.Error("Failed to build metadata #2 with info: %v", err)
		return
	}
	logger.Info("  Mint metadata #2 with certificate info built")

	var msgs []sdk.Msg
	msgs = append(msgs, msgDeploySchema, msgCreateMetadata, msgCreateMetadataWithInfo)

	logger.Info("Broadcasting transaction to blockchain...")
	res, err := meta.BroadcastTxAndWait(msgs...)
	if err != nil {
		logger.Error("Failed to broadcast transaction: %v", err)
		return
	}

	logger.Info("SUCCESS: Schema deployed successfully")
	logger.Info("  Schema Code: %s", nftSchemaName)
	logger.Info("  Transaction Hash: %s", res.TxHash)
	logger.Info("========================================")
	// Step 5: Deploy EVM NFT Contract
	logger.Info("========================================")
	logger.Info("[Step 5] Deploying EVM NFT contract")
	logger.Info("========================================")
	evmClient := evm.NewEVMClient(*acc)

	logger.Info("Contract Details:")
	logger.Info("  Name: %s", contractName)
	logger.Info("  Symbol: %s", contractSymbol)
	logger.Info("  Schema: %s", nftSchemaName)

	logger.Info("Deploying contract...")
	contractAddress, tx, err := evmClient.DeployCertificateContract(contractName, contractSymbol, nftSchemaName)
	if err != nil {
		logger.Error("Failed to deploy certificate contract: %v", err)
		return
	}

	logger.Info("Waiting for transaction to be mined...")
	_, err = client.WaitForEVMTransaction(tx.Hash())
	if err != nil {
		logger.Error("Error waiting for deployment: %v", err)
		return
	}

	logger.Info("SUCCESS: Contract deployed successfully")
	logger.Info("  Contract Address: %s", contractAddress.Hex())
	logger.Info("  Transaction Hash: %s", tx.Hash().Hex())
	logger.Info("========================================")

	// Step 6: Mint Certificate NFT
	logger.Info("========================================")
	logger.Info("[Step ")
	logger.Info("========================================")

	// Mint NFT #1 to self
	tokenID := uint64(1)
	logger.Info("Minting NFT #%d to self...", tokenID)
	tx, err = evmClient.MintCertificateNFT(contractAddress, tokenID)
	if err != nil {
		logger.Error("Failed to mint NFT: %v", err)
		return
	}
	logger.Info("  Transaction Hash: %s", tx.Hash().Hex())
	logger.Info("  Nonce: %v", tx.Nonce())

	logger.Info("Waiting for transaction to be mined...")
	_, err = client.WaitForEVMTransaction(tx.Hash())
	if err != nil {
		logger.Error("Error waiting for mint transaction: %v", err)
		return
	}

	logger.Info("SUCCESS: NFT #%d minted successfully to owner", tokenID)

	// Mint NFT #2 to Bob
	logger.Info("Minting NFT #%d to Bob's address...", tokenID+1)
	logger.Info("  Destination: %s", BobEVMAddres)
	tx, err = evmClient.MintCertificateNFTToDestination(contractAddress, tokenID+1, common.HexToAddress(BobEVMAddres))
	if err != nil {
		logger.Error("Failed to mint NFT to destination: %v", err)
		return
	}
	logger.Info("  Transaction Hash: %s", tx.Hash().Hex())
	logger.Info("  Nonce: %v", tx.Nonce())

	logger.Info("Waiting for transaction to be mined...")
	_, err = client.WaitForEVMTransaction(tx.Hash())
	if err != nil {
		logger.Error("Error waiting for mint transaction: %v", err)
		return
	}

	logger.Info("SUCCESS: NFT #%d minted successfully to Bob", tokenID+1)
	logger.Info("========================================")

	// Step 7: Change certificate state (Freeze/Unfreeze)
	logger.Info("========================================")
	logger.Info("[Step ")
	logger.Info("========================================")

	logger.Info("Freezing certificate #1...")
	res, err = meta.FreezeCertificate("1")
	if err != nil {
		logger.Error("Failed to freeze certificate: %v", err)
		return
	}

	logger.Info("Waiting for freeze transaction to be confirmed...")
	err = client.WaitForTransaction(res.TxHash)
	if err != nil {
		logger.Error("Error waiting for freeze transaction: %v", err)
		return
	}

	logger.Info("SUCCESS")
	logger.Info("  Transaction Hash: %s\n", res.TxHash)

	logger.Info("Unfreezing certificate #1...")
	res, err = meta.UnfreezeCertificate("1")
	if err != nil {
		logger.Error("Failed to unfreeze certificate: %v", err)
		return
	}

	logger.Info("Waiting for unfreeze transaction to be confirmed...")
	err = client.WaitForTransaction(res.TxHash)
	if err != nil {
		logger.Error("Error waiting for unfreeze transaction: %v", err)
		return
	}

	logger.Info("SUCCESS")
	logger.Info("  Transaction Hash: %s", res.TxHash)
	logger.Info("========================================")

	// Step 8: Transfer NFT
	logger.Info("========================================")
	logger.Info("[Step ")
	logger.Info("========================================")

	logger.Info("Transferring NFT #%d to Bob...", tokenID)
	logger.Info("  From: %s", acc.GetEVMAddress().Hex())
	logger.Info("  To: %s", BobEVMAddres)

	tx, err = evmClient.TransferCertificateNFT(contractAddress, common.HexToAddress(BobEVMAddres), tokenID)
	if err != nil {
		logger.Error("Failed to transfer NFT: %v", err)
		return
	}
	logger.Info("  Transaction Hash: %s", tx.Hash().Hex())
	logger.Info("  Nonce: %v", tx.Nonce())

	logger.Info("Waiting for transfer transaction to be mined...")
	_, err = client.WaitForEVMTransaction(tx.Hash())
	if err != nil {
		logger.Error("Error waiting for transfer: %v", err)
		return
	}

	logger.Info("SUCCESS")
	logger.Info("========================================")

	// Step 8.1: EIP-2612 Permit (gasless transfer)
	logger.Info("========================================")
	logger.Info("[Step ")
	logger.Info("Admin pays gas for user's NFT transfer")
	logger.Info("========================================")

	logger.Info("Process Overview:")
	logger.Info("  1. Create new user account (no funds needed!)")
	logger.Info("  2. Mint NFT to new user")
	logger.Info("  3. User signs EIP-712 permit message offline (no gas, no blockchain interaction)")
	logger.Info("  4. Admin broadcasts transfer with permit (admin pays all gas)")

	// Step 8.1.1: Create new account
	logger.Info("┌─ [Step ")
	accFromGenMnemonic, err := account.NewAccount(client, "new_gen", mnemonic, "")
	if err != nil {
		logger.Info("   ERROR: Failed to create account: %v", err)
		return
	}
	logger.Info("└─ SUCCESS")
	logger.Info("      Name: new_gen")
	logger.Info("      Address: %s", accFromGenMnemonic.GetEVMAddress().Hex())
	logger.Info("      Balance: 0 (no funds needed!)")

	// Step 8.1.2: Mint NFT to new account
	logger.Info("┌─ [Step 8.1.2] Minting NFT #%d to new account", tokenID+2)
	tx, err = evmClient.MintCertificateNFTToDestination(contractAddress, tokenID+2, accFromGenMnemonic.GetEVMAddress())
	if err != nil {
		logger.Info("   ERROR: Failed to mint NFT: %v", err)
		return
	}
	logger.Info("   Transaction Hash: %s", tx.Hash().Hex())
	logger.Info("   Recipient: %s", accFromGenMnemonic.GetEVMAddress().Hex())

	logger.Info("   Waiting for mint transaction...")
	_, err = client.WaitForEVMTransaction(tx.Hash())
	if err != nil {
		logger.Info("   ERROR: Error waiting for mint: %v", err)
		return
	}
	logger.Info("└─ SUCCESS: NFT #%d minted to new user", tokenID+2)

	// Step 8.1.3: User signs EIP-712 permit offline
	logger.Info("┌─ [Step ")
	offlineEVMClient := evm.NewEVMClient(*accFromGenMnemonic)

	logger.Info("   User signs permit for NFT #%d transfer", tokenID+2)
	logger.Info("   From (User): %s", accFromGenMnemonic.GetEVMAddress().Hex())
	logger.Info("   To (Charlie): %s", ChalieEVMAddress)
	logger.Info("   NOTE: This is just a signature, NO transaction, NO gas needed!")

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
		logger.Info("   ERROR: Failed to sign permit: %v", err)
		return
	}
	logger.Info("└─ SUCCESS")

	// Step 8.1.4: Admin broadcasts transfer with permit
	logger.Info("┌─ [Step ")
	logger.Info("   NOTE: Admin pays ALL gas, user pays NOTHING")
	logger.Info("   Transferring NFT #%d with permit", tokenID+2)

	transferTx, err := evmClient.TransferWithPermit(
		contractAddress,
		accFromGenMnemonic.GetEVMAddress(),    // From
		common.HexToAddress(ChalieEVMAddress), // To
		big.NewInt(int64(tokenID+2)),
		permitSig,
	)
	if err != nil {
		logger.Info("   ERROR: Failed to execute transfer with permit: %v", err)
		return
	}

	logger.Info("   Transaction Hash: %s", transferTx.Hash().Hex())
	logger.Info("   Waiting for confirmation...")
	_, err = client.WaitForEVMTransaction(transferTx.Hash())
	if err != nil {
		logger.Info("   ERROR: Error waiting for transfer: %v", err)
		return
	}

	logger.Info("└─ SUCCESS")
	logger.Info("   Transfer Summary:")
	logger.Info("      NFT #%d transferred", tokenID+2)
	logger.Info("      Owner (User): %s (paid 0 gas! 🎉)", accFromGenMnemonic.GetEVMAddress().Hex())
	logger.Info("      Destination: %s", ChalieEVMAddress)
	logger.Info("      Gas paid by: Admin (%s)", acc.GetEVMAddress().Hex())
	logger.Info("      Transaction Hash: %s", transferTx.Hash().Hex())
	logger.Info("      Method: EIP-2612 Permit (transferWithPermit)")
	logger.Info("========================================")

	// Step 9: Verify ownership
	logger.Info("========================================")
	logger.Info("[Step ")
	logger.Info("========================================")

	logger.Info("Checking owner of NFT #%d...", tokenID+2)
	currentOwner := evmClient.TokenOwner(contractAddress, tokenID+2)
	logger.Info("Current owner: %s", currentOwner.Hex())
	logger.Info("  (Should be Charlie: %s)", ChalieEVMAddress)
	logger.Info("========================================")

	// Step 10: Burn NFT directly
	logger.Info("========================================")
	logger.Info("[Step ")
	logger.Info("========================================")

	// First mint a new NFT to burn
	burnTokenID := tokenID + 3
	logger.Info("Minting NFT #%d to burn...", burnTokenID)
	tx, err = evmClient.MintCertificateNFT(contractAddress, burnTokenID)
	if err != nil {
		logger.Error("Failed to mint NFT for burning: %v", err)
		return
	}
	logger.Info("  Transaction Hash: %s", tx.Hash().Hex())

	logger.Info("Waiting for mint transaction...")
	_, err = client.WaitForEVMTransaction(tx.Hash())
	if err != nil {
		logger.Error("Error waiting for mint: %v", err)
		return
	}
	logger.Info("SUCCESS: NFT #%d minted", burnTokenID)

	// Now burn it
	logger.Info("Burning NFT #%d...", burnTokenID)
	logger.Info("  Owner: %s", acc.GetEVMAddress().Hex())

	burnTx, err := evmClient.BurnCertificateNFT(contractAddress, burnTokenID)
	if err != nil {
		logger.Error("Failed to burn NFT: %v", err)
		return
	}
	logger.Info("  Transaction Hash: %s", burnTx.Hash().Hex())
	logger.Info("  Nonce: %v", burnTx.Nonce())

	logger.Info("Waiting for burn transaction to be mined...")
	_, err = client.WaitForEVMTransaction(burnTx.Hash())
	if err != nil {
		logger.Error("Error waiting for burn: %v", err)
		return
	}

	logger.Info("SUCCESS: NFT #%d burned successfully", burnTokenID)

	// Verify the token was burned by checking owner
	logger.Info("Verifying burn...")
	logger.Info("Querying owner of burned NFT #%d...", burnTokenID)
	burnedOwner := evmClient.TokenOwner(contractAddress, burnTokenID)
	logger.Info("  Owner address: %s", burnedOwner.Hex())

	zeroAddress := common.HexToAddress("0x0000000000000000000000000000000000000000")
	if burnedOwner == zeroAddress {
		logger.Info("  ✓")
	} else {
		logger.Info("  ⚠ WARNING: Token still has owner: %s", burnedOwner.Hex())
	}
	logger.Info("========================================")

	// Step 11: Gasless Burn with Permit
	logger.Info("========================================")
	logger.Info("[Step ")
	logger.Info("Admin pays gas for user's NFT burn")
	logger.Info("========================================")

	logger.Info("Process Overview:")
	logger.Info("  1. Mint NFT to user account (the one we created earlier)")
	logger.Info("  2. User signs EIP-712 permit message offline (no gas)")
	logger.Info("  3. Admin broadcasts burn with permit (admin pays all gas)")

	// Step 11.1: Mint NFT to user account
	gaslessBurnTokenID := burnTokenID + 1
	logger.Info("┌─ [Step 11.1] Minting NFT #%d to user account for gasless burn", gaslessBurnTokenID)
	tx, err = evmClient.MintCertificateNFTToDestination(contractAddress, gaslessBurnTokenID, accFromGenMnemonic.GetEVMAddress())
	if err != nil {
		logger.Info("   ERROR: Failed to mint NFT: %v", err)
		return
	}
	logger.Info("   Transaction Hash: %s", tx.Hash().Hex())
	logger.Info("   Recipient: %s", accFromGenMnemonic.GetEVMAddress().Hex())

	logger.Info("   Waiting for mint transaction...")
	_, err = client.WaitForEVMTransaction(tx.Hash())
	if err != nil {
		logger.Info("   ERROR: Error waiting for mint: %v", err)
		return
	}
	logger.Info("└─ SUCCESS: NFT #%d minted to user", gaslessBurnTokenID)

	// Step 11.2: User signs EIP-712 permit for burn
	logger.Info("┌─ [Step ")
	offlineEVMClientForBurn := evm.NewEVMClient(*accFromGenMnemonic)

	logger.Info("   User signs permit for NFT #%d burn", gaslessBurnTokenID)
	logger.Info("   Owner (User): %s", accFromGenMnemonic.GetEVMAddress().Hex())
	logger.Info("   NOTE: This is just a signature, NO transaction, NO gas needed!")

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
		logger.Info("   ERROR: Failed to sign burn permit: %v", err)
		return
	}
	logger.Info("└─ SUCCESS")

	// Step 11.3: Admin broadcasts burn with permit
	logger.Info("┌─ [Step ")
	logger.Info("   NOTE: Admin pays ALL gas, user pays NOTHING")
	logger.Info("   Burning NFT #%d with permit", gaslessBurnTokenID)

	burnWithPermitTx, err := evmClient.BurnWithPermit(
		contractAddress,
		accFromGenMnemonic.GetEVMAddress(), // From (owner)
		big.NewInt(int64(gaslessBurnTokenID)),
		burnPermitSig,
	)
	if err != nil {
		logger.Info("   ERROR: Failed to execute burn with permit: %v", err)
		return
	}

	logger.Info("   Transaction Hash: %s", burnWithPermitTx.Hash().Hex())
	logger.Info("   Waiting for confirmation...")
	_, err = client.WaitForEVMTransaction(burnWithPermitTx.Hash())
	if err != nil {
		logger.Info("   ERROR: Error waiting for burn: %v", err)
		return
	}

	logger.Info("└─ SUCCESS")
	logger.Info("   Burn Summary:")
	logger.Info("      NFT #%d burned", gaslessBurnTokenID)
	logger.Info("      Owner (User): %s (paid 0 gas! 🎉)", accFromGenMnemonic.GetEVMAddress().Hex())
	logger.Info("      Gas paid by: Admin (%s)", acc.GetEVMAddress().Hex())
	logger.Info("      Transaction Hash: %s", burnWithPermitTx.Hash().Hex())
	logger.Info("      Method: EIP-2612 Permit (burnWithPermit)")

	// Verify the token was burned by checking owner
	logger.Info("   Verifying burn...")
	logger.Info("   Querying owner of burned NFT #%d...", gaslessBurnTokenID)
	gaslessBurnedOwner := evmClient.TokenOwner(contractAddress, gaslessBurnTokenID)
	logger.Info("   Owner address: %s", gaslessBurnedOwner.Hex())

	zeroAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")
	if gaslessBurnedOwner == zeroAddress {
		logger.Info("   ✓")
	} else {
		logger.Info("   ⚠ WARNING: Token still has owner: %s", gaslessBurnedOwner.Hex())
	}
	logger.Info("========================================")

	// Summary
	logger.Info("╔════════════════════════════════════════╗")
	logger.Info("║          EXECUTION SUMMARY             ║")
	logger.Info("╚════════════════════════════════════════╝")
	logger.Info("Schema Information:")
	logger.Info("  Schema Code: %s\n", nftSchemaName)

	logger.Info("Contract Information:")
	logger.Info("  Address: %s", contractAddress.Hex())
	logger.Info("  Name: %s", contractName)
	logger.Info("  Symbol: %s\n", contractSymbol)

	logger.Info("NFTs Created:")
	logger.Info("  NFT #1: Transferred to Bob")
	logger.Info("  NFT #2: Minted to Bob")
	logger.Info("  NFT #3: Gasless transfer to Charlie")
	logger.Info("  NFT #4: Burned directly by owner")
	logger.Info("  NFT #5: Gasless burn via permit")

	logger.Info("Final Owner of NFT #%d:", tokenID+2)
	logger.Info("  Address: %s", currentOwner.Hex())
	logger.Info("Quick start completed successfully!")
	logger.Info("╔════════════════════════════════════════╗")
}
