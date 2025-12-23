package main

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/client"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/evm"
)

// This example demonstrates gasless NFT burning using EIP-2612 permit signatures.
// In a gasless burn, the NFT owner signs a permit message offline (no gas needed),
// and a relayer/admin broadcasts the burn transaction and pays for all gas fees.
//
// Usage:
//   go run 13_gasless_burn.go
//
// What this script does:
// 1. Creates two accounts: Admin (has funds) and User (no funds needed!)
// 2. Admin mints an NFT to the User
// 3. User signs an EIP-712 permit message offline (completely free, no blockchain interaction)
// 4. Admin broadcasts the burn using the permit (admin pays all gas)
// 5. Verifies the NFT was burned successfully (owner = zero address)
//
// Use Cases:
// - Allow users to burn NFTs without paying gas fees
// - Certificate revocation without user costs
// - Building gasless dApps where the platform pays for user transactions
// - Token cleanup and management services
// - Eco-friendly token destruction (users don't need tokens to burn tokens!)
//
// Prerequisites:
// - NFT contract must be deployed (see 05_deploy_contract.go)
// - Admin account must have tokens for gas fees
// - User account doesn't need any tokens!

const (
	// IMPORTANT: Replace with your deployed contract address
	contractAddress = "0x0000000000000000000000000000000000000000"

	// Contract name (must match the name used during deployment)
	contractName = "MyNFTCert"

	// Token ID to mint and burn
	tokenId = uint64(1)

	// Admin mnemonic (has funds to pay for gas)
	adminMnemonic = account.TestMnemonic

	// User mnemonic (no funds needed!)
	// In production, this would be the user's own wallet
	userMnemonic = "your user mnemonic here or generate a new one"
)

func main() {
	fmt.Println("=== Step 13: Gasless Burn (EIP-2612 Permit) ===")
	fmt.Println()
	fmt.Println("This example demonstrates how to enable gasless NFT burning")
	fmt.Println("where users can burn NFTs without paying any gas fees.")
	fmt.Println()

	// Validate configuration
	if contractAddress == "0x0000000000000000000000000000000000000000" {
		fmt.Println("ERROR: Please update the contractAddress constant")
		fmt.Println("   Use the contract address from step 05_deploy_contract.go")
		return
	}

	// Step 1: Setup client
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Step 1: Connecting to network")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	ctx := context.Background()
	client, err := client.NewClient(ctx, false) // false = testnet
	if err != nil {
		panic(fmt.Sprintf("Failed to create client: %v", err))
	}
	fmt.Println("✓ Connected to testnet")
	fmt.Println()

	// Step 2: Create Admin account (pays for gas)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Step 2: Creating Admin account (gas payer)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	adminAcc, err := account.NewAccount(client, "admin", adminMnemonic, "")
	if err != nil {
		panic(fmt.Sprintf("Failed to create admin account: %v", err))
	}
	fmt.Printf("✓ Admin account created\n")
	fmt.Printf("   Address: %s\n", adminAcc.GetEVMAddress().Hex())
	fmt.Println()

	// Step 3: Create User account (no funds needed!)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Step 3: Creating User account (NO FUNDS NEEDED!)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Generate a new mnemonic for the user if not provided
	actualUserMnemonic := userMnemonic
	if userMnemonic == "your user mnemonic here or generate a new one" {
		fmt.Println("Generating new user mnemonic...")
		actualUserMnemonic, err = account.GenerateMnemonic()
		if err != nil {
			panic(fmt.Sprintf("Failed to generate mnemonic: %v", err))
		}
		fmt.Printf("   Generated Mnemonic: %s\n", actualUserMnemonic)
		fmt.Println("   ⚠️  Save this mnemonic if you want to use it later")
		fmt.Println()
	}

	userAcc, err := account.NewAccount(client, "user", actualUserMnemonic, "")
	if err != nil {
		panic(fmt.Sprintf("Failed to create user account: %v", err))
	}
	fmt.Printf("✓ User account created\n")
	fmt.Printf("   Address: %s\n", userAcc.GetEVMAddress().Hex())
	fmt.Printf("   Balance: 0 (no funds needed! 🎉)\n")
	fmt.Println()

	// Step 4: Admin mints NFT to User
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Step 4: Admin mints NFT to User")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	adminEvmClient := evm.NewEVMClient(*adminAcc)
	contractAddr := common.HexToAddress(contractAddress)

	fmt.Printf("Minting NFT #%d to user...\n", tokenId)
	fmt.Printf("   Recipient: %s\n", userAcc.GetEVMAddress().Hex())

	mintTx, err := adminEvmClient.MintCertificateNFTToDestination(
		contractAddr,
		tokenId,
		userAcc.GetEVMAddress(),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to mint NFT: %v", err))
	}

	fmt.Printf("   Transaction Hash: %s\n", mintTx.Hash().Hex())
	fmt.Println("   Waiting for confirmation...")

	_, err = client.WaitForEVMTransaction(mintTx.Hash())
	if err != nil {
		panic(fmt.Sprintf("Error waiting for mint: %v", err))
	}

	fmt.Println("✓ NFT minted successfully")
	fmt.Println()

	// Step 5: Verify user owns the NFT
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Step 5: Verifying NFT ownership")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	currentOwner := adminEvmClient.TokenOwner(contractAddr, tokenId)
	fmt.Printf("Current owner: %s\n", currentOwner.Hex())
	fmt.Printf("Expected:      %s\n", userAcc.GetEVMAddress().Hex())

	if currentOwner.Hex() != userAcc.GetEVMAddress().Hex() {
		panic("Ownership verification failed!")
	}
	fmt.Println("✓ User owns the NFT")
	fmt.Println()

	// Step 6: User signs EIP-712 permit for burn offline (NO GAS!)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Step 6: User signs EIP-712 permit for burn (COMPLETELY OFFLINE)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	userEvmClient := evm.NewEVMClient(*userAcc)

	fmt.Println("💡 This is just a cryptographic signature:")
	fmt.Println("   • NO transaction sent to blockchain")
	fmt.Println("   • NO gas fees required")
	fmt.Println("   • NO need for tokens in wallet")
	fmt.Println("   • Can be done completely offline")
	fmt.Println()

	fmt.Printf("User signs permit for NFT #%d burn:\n", tokenId)
	fmt.Printf("   Owner (User):     %s\n", userAcc.GetEVMAddress().Hex())
	fmt.Printf("   Spender (Admin):  %s\n", adminAcc.GetEVMAddress().Hex())

	// Set deadline to 1 hour from now (Unix timestamp)
	deadline := big.NewInt(time.Now().Unix() + 3600)
	fmt.Printf("   Deadline:         %s\n", time.Unix(deadline.Int64(), 0).Format(time.RFC3339))
	fmt.Println()

	fmt.Println("Signing permit message for burn...")
	permitSig, err := userEvmClient.SignPermit(
		contractName,
		contractAddr,
		adminAcc.GetEVMAddress(), // Spender (admin/relay)
		big.NewInt(int64(tokenId)),
		deadline,
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to sign permit: %v", err))
	}

	fmt.Println("✓ Burn permit signed successfully!")
	fmt.Println()
	fmt.Println("🎉 User paid ZERO gas for this signature!")
	fmt.Println()

	// Step 7: Admin broadcasts burn with permit (ADMIN PAYS GAS)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Step 7: Admin broadcasts burn (ADMIN PAYS ALL GAS)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	fmt.Println("Admin executes burnWithPermit():")
	fmt.Printf("   Owner:            %s\n", userAcc.GetEVMAddress().Hex())
	fmt.Printf("   Token ID:         %d\n", tokenId)
	fmt.Printf("   Gas Payer:        %s (Admin)\n", adminAcc.GetEVMAddress().Hex())
	fmt.Println()

	fmt.Println("Broadcasting burn transaction...")
	burnTx, err := adminEvmClient.BurnWithPermit(
		contractAddr,
		userAcc.GetEVMAddress(), // From (owner)
		big.NewInt(int64(tokenId)),
		permitSig,
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to execute burn with permit: %v", err))
	}

	fmt.Printf("   Transaction Hash: %s\n", burnTx.Hash().Hex())
	fmt.Printf("   Nonce: %d\n", burnTx.Nonce())
	fmt.Println("   Waiting for confirmation...")

	receipt, err := client.WaitForEVMTransaction(burnTx.Hash())
	if err != nil {
		panic(fmt.Sprintf("Error waiting for burn: %v", err))
	}

	fmt.Println("✓ Burn completed!")
	fmt.Println()

	// Step 8: Verify NFT was burned (owner should be zero address)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Step 8: Verifying NFT was burned")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	burnedOwner := adminEvmClient.TokenOwner(contractAddr, tokenId)
	zeroAddress := common.HexToAddress("0x0000000000000000000000000000000000000000")

	fmt.Printf("Owner after burn: %s\n", burnedOwner.Hex())
	fmt.Printf("Zero address:     %s\n", zeroAddress.Hex())

	if burnedOwner == zeroAddress {
		fmt.Println("✓ NFT successfully burned (owner is zero address)")
	} else {
		panic(fmt.Sprintf("Burn verification failed! Owner is still: %s", burnedOwner.Hex()))
	}
	fmt.Println()

	// Display summary
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║                GASLESS BURN SUMMARY                        ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("Contract Address:      %s\n", contractAddress)
	fmt.Printf("Token ID:              %d\n", tokenId)
	fmt.Println()
	fmt.Printf("Original Owner (User): %s\n", userAcc.GetEVMAddress().Hex())
	fmt.Printf("   Gas Paid:           0 🎉 (completely free!)\n")
	fmt.Println()
	fmt.Printf("Final Owner:           %s\n", burnedOwner.Hex())
	fmt.Printf("   Status:             BURNED ♨️\n")
	fmt.Println()
	fmt.Printf("Gas Payer (Admin):     %s\n", adminAcc.GetEVMAddress().Hex())
	fmt.Printf("   Gas Used:           %d\n", receipt.GasUsed)
	fmt.Println()
	fmt.Printf("Transaction Hash:      %s\n", burnTx.Hash().Hex())
	fmt.Printf("Block Number:          %d\n", receipt.BlockNumber)
	fmt.Printf("Method:                burnWithPermit (EIP-2612)\n")
	fmt.Println()
	fmt.Println("────────────────────────────────────────────────────────────")
	fmt.Println()

	// Explanation
	fmt.Println("📚 What just happened:")
	fmt.Println()
	fmt.Println("1. User signed an EIP-712 permit message offline (no gas)")
	fmt.Println("2. Admin broadcasted the burn using the permit")
	fmt.Println("3. Admin paid ALL gas fees (user paid nothing!)")
	fmt.Println("4. NFT permanently burned (owner = zero address)")
	fmt.Println()
	fmt.Println("🎯 Use Cases:")
	fmt.Println()
	fmt.Println("• Certificate revocation without user paying gas")
	fmt.Println("• Token cleanup services (burn expired certificates)")
	fmt.Println("• Eco-friendly NFT destruction (no gas barrier)")
	fmt.Println("• Platform-managed token lifecycle")
	fmt.Println("• Gasless dApps where platform handles all costs")
	fmt.Println("• Compliance: revoke certificates without user friction")
	fmt.Println()
	fmt.Println("💡 Technical Details:")
	fmt.Println()
	fmt.Println("• Standard: EIP-2612 (Permit Extension for ERC-721)")
	fmt.Println("• Signature: EIP-712 structured data signing")
	fmt.Println("• Security: Includes deadline and nonce to prevent replay attacks")
	fmt.Println("• Flexibility: User signs offline, anyone can broadcast")
	fmt.Println("• Verification: Owner becomes zero address after burn")
	fmt.Println()
	fmt.Println("🔥 Burn vs Transfer:")
	fmt.Println()
	fmt.Println("• Burn: Permanently destroys the NFT (irreversible)")
	fmt.Println("• Transfer: Changes ownership to another address (reversible)")
	fmt.Println("• Both support gasless operations via EIP-2612 permits")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  • Compare with direct burn (user pays gas)")
	fmt.Println("  • Implement your own relayer service")
	fmt.Println("  • Build gasless certificate management system")
	fmt.Println("  • Try gasless transfer (07_1_gasless_transfer.go)")
	fmt.Println()
}
