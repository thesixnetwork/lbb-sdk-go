# LBB SDK Go - Tutorial Examples

Welcome! This tutorial will guide you through the complete process of working with the LBB SDK for Go. Each example script focuses on a specific function to help you understand the workflow step by step.

## 📚 Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Setup](#setup)
- [Tutorial Flow](#tutorial-flow)
- [Example Scripts](#example-scripts)
- [Common Issues](#common-issues)
- [Additional Resources](#additional-resources)

## Overview

The LBB SDK allows you to create and manage certificate NFTs with metadata on the blockchain. This involves both Cosmos SDK operations (for metadata/schema) and EVM operations (for NFT contracts).

### What You'll Learn

1. How to generate wallets and create accounts
2. How to deploy certificate schemas (metadata structure)
3. How to deploy EVM NFT contracts
4. How to mint, transfer, and manage NFTs
5. How to freeze/unfreeze certificate metadata
6. How to query NFT information

## Prerequisites

### Required Knowledge

- Basic understanding of Go programming
- Familiarity with blockchain concepts
- Understanding of NFTs and smart contracts (helpful but not required)

### Software Requirements

- Go 1.21 or higher
- Internet connection (to connect to the testnet)

### Account Requirements

- For testing: Use the provided test mnemonic (already included in examples)
- For production: You'll need tokens in your account for transaction fees

## Setup

### 1. Install Dependencies

Navigate to the example directory and install dependencies:

```bash
cd example
go mod download
```

### 2. Verify Setup

Run the first example to verify everything works:

```bash
go run 01_generate_wallet.go
```

You should see a mnemonic phrase generated successfully.

## Tutorial Flow

The examples are numbered in the recommended order of execution:

```
01_generate_wallet.go
    ↓
02_create_account.go
    ↓
03_deploy_schema.go
    ↓
04_deploy_contract.go
    ↓
05_mint_nft.go
    ↓
06_transfer_nft.go
    ↓
07_freeze_metadata.go
    ↓
08_query_nft.go
```

### Quick Start Path

If you want to run everything quickly using the test mnemonic:

```bash
# Run all examples in sequence
go run 02_create_account.go
go run 03_deploy_schema.go
go run 04_deploy_contract.go
# Update contract address in 05_mint_nft.go
go run 05_mint_nft.go
# Update contract address in 06_transfer_nft.go
go run 06_transfer_nft.go
go run 07_freeze_metadata.go
# Update contract address in 08_query_nft.go
go run 08_query_nft.go
```

## Example Scripts

### 01_generate_wallet.go

**Purpose:** Generate a new wallet with a mnemonic phrase

**What it does:**
- Creates a random 24-word mnemonic (BIP-39 standard)
- Displays the mnemonic for backup

**Usage:**
```bash
go run 01_generate_wallet.go
```

**Key Points:**
- ⚠️ **Save your mnemonic securely!** This is the only way to recover your wallet
- Never share your mnemonic with anyone
- For testing, you can use the built-in test mnemonic

**Output Example:**
```
Mnemonic: word1 word2 word3 ... word24
```

---

### 02_create_account.go

**Purpose:** Connect to the network and create an account from a mnemonic

**What it does:**
- Connects to fivenet (testnet)
- Creates an account from mnemonic
- Displays both Cosmos and EVM addresses

**Usage:**
```bash
go run 02_create_account.go
```

**Key Concepts:**
- **Cosmos Address (6x...)**: Used for Cosmos SDK operations
- **EVM Address (0x...)**: Used for Ethereum-compatible operations
- Both addresses are derived from the same private key

**Configuration:**
```go
const (
    exampleMnemonic = account.TestMnemonic  // Replace with your mnemonic
    accountName = "my-account"               // Any name you prefer
    accountPassword = ""                     // Optional keyring password
)
```

**Output Example:**
```
Account Name:    my-account
EVM Address:     0x1234...
Cosmos Address:  6x1234...
```

---

### 03_deploy_schema.go

**Purpose:** Deploy a certificate schema (metadata structure)

**What it does:**
- Connects to the network
- Creates a metadata schema
- Deploys the schema to the blockchain
- Mints the first metadata instance

**Usage:**
```bash
go run 03_deploy_schema.go
```

**Key Concepts:**
- **Schema**: Defines the structure for your certificates
- **Schema Name Format**: `{ORGNAME}.{SchemaCode}` (e.g., "mycompany.cert01")
- Schemas are required before deploying EVM contracts

**Configuration:**
```go
const (
    schemaName = "myorg.lbbv01"  // Change to your organization name
    initialTokenId = "1"          // First metadata token ID
)
```

**Important:**
- You need tokens in your account for transaction fees
- Schema name must be unique
- Save the schema name for use in contract deployment

**Output Example:**
```
Schema Code:       myorg.lbbv01
Transaction Hash:  ABC123...
```

---

### 04_deploy_contract.go

**Purpose:** Deploy an EVM NFT contract linked to your schema

**What it does:**
- Connects to the network
- Creates an EVM client
- Deploys a Certificate NFT contract
- Links the contract to your schema

**Usage:**
```bash
go run 04_deploy_contract.go
```

**Key Concepts:**
- **Contract**: The smart contract that manages NFTs
- **ERC-721**: Standard for NFTs
- Contract is linked to the schema deployed in step 03

**Configuration:**
```go
const (
    contractName = "MyCertificate"    // Human-readable name
    contractSymbol = "CERT"            // Short symbol (3-5 chars)
    schemaName = "myorg.lbbv01"       // Must match step 03
)
```

**Important:**
- ⚠️ **Save the contract address!** You'll need it for all future operations
- Schema must be deployed first
- Contract deployment costs gas fees

**Output Example:**
```
Contract Address:  0xABCD1234...
Transaction Hash:  0x5678...
```

---

### 05_mint_nft.go

**Purpose:** Mint a new NFT from your deployed contract

**What it does:**
- Connects to the network
- Mints a certificate NFT with a specific token ID
- Verifies ownership

**Usage:**
```bash
# Update contractAddress in the file first!
go run 05_mint_nft.go
```

**Configuration:**
```go
const (
    contractAddress = "0x..."  // From step 04 output
    tokenId = uint64(1)        // Unique token ID
)
```

**Important:**
- Contract must be deployed first
- Token ID must be unique (can't mint the same ID twice)
- NFT is minted to your address

**Output Example:**
```
Token ID:     1
Owner:        0x1234...
Transaction:  0xABCD...
```

---

### 06_transfer_nft.go

**Purpose:** Transfer an NFT to another address

**What it does:**
- Verifies you own the NFT
- Transfers the NFT to a recipient
- Verifies the new owner

**Usage:**
```bash
# Update contractAddress and recipientAddress first!
go run 06_transfer_nft.go
```

**Configuration:**
```go
const (
    contractAddress = "0x..."      // From step 04
    tokenId = uint64(1)            // Token to transfer
    recipientAddress = "0x..."     // Recipient's address
)
```

**Important:**
- You must own the NFT to transfer it
- Recipient address must be valid
- After transfer, you can't transfer it again (new owner can)

**Output Example:**
```
Previous Owner:  0x1111...
New Owner:       0x2222...
Transaction:     0xABCD...
```

---

### 07_freeze_metadata.go

**Purpose:** Freeze and unfreeze certificate metadata

**What it does:**
- Freezes the metadata (prevents modifications)
- Unfreezes the metadata (allows modifications)

**Usage:**
```bash
go run 07_freeze_metadata.go
```

**Key Concepts:**
- **Freeze**: Locks metadata to prevent changes (useful for finalized certificates)
- **Unfreeze**: Unlocks metadata to allow updates

**Configuration:**
```go
const (
    schemaName = "myorg.lbbv01"  // From step 03
    tokenId = "1"                 // Metadata token ID
)
```

**Use Cases:**
- Freeze: Lock academic certificates after issuance
- Freeze: Prevent tampering with credentials
- Unfreeze: Allow corrections or updates

**Output Example:**
```
Freeze Tx:    0x1111...
Unfreeze Tx:  0x2222...
Final Status: UNFROZEN
```

---

### 08_query_nft.go

**Purpose:** Query NFT information and verify ownership

**What it does:**
- Queries NFT ownership
- Displays token information
- Verifies if you own the NFT

**Usage:**
```bash
# Update contractAddress first!
go run 08_query_nft.go
```

**Configuration:**
```go
const (
    contractAddress = "0x..."  // From step 04
    tokenId = uint64(1)        // Token to query
)
```

**Key Points:**
- Queries are read-only (no gas fees)
- Anyone can query NFT ownership
- Useful for verification and auditing

**Output Example:**
```
Token ID:      1
Current Owner: 0x1234...
```

---

## Common Issues

### Issue: "Failed to create client"

**Cause:** Network connection issue or wrong network configuration

**Solution:**
- Check your internet connection
- Verify you're using the correct network (testnet vs mainnet)
- In `client.NewClient(ctx, false)`, `false` = testnet, `true` = mainnet

### Issue: "Insufficient funds"

**Cause:** Your account doesn't have enough tokens for transaction fees

**Solution:**
- For testnet: Request tokens from the faucet
- Use the test mnemonic which has testnet tokens
- Check your balance before transactions

### Issue: "Schema not found"

**Cause:** Schema wasn't deployed or wrong schema name used

**Solution:**
- Ensure you ran `03_deploy_schema.go` successfully
- Verify the schema name matches exactly (case-sensitive)
- Check the transaction was confirmed

### Issue: "Contract address not found"

**Cause:** Contract wasn't deployed or wrong address used

**Solution:**
- Ensure you ran `04_deploy_contract.go` successfully
- Copy the exact contract address from the deployment output
- Don't use the placeholder `0x0000...` address

### Issue: "Token ID already exists"

**Cause:** Trying to mint an NFT with a token ID that's already minted

**Solution:**
- Use a different token ID
- Query existing tokens to find available IDs
- Increment the token ID (e.g., use 2, 3, 4...)

### Issue: "You don't own this NFT"

**Cause:** Trying to transfer an NFT you don't own

**Solution:**
- Verify ownership with `08_query_nft.go`
- Make sure you're using the correct mnemonic/account
- Check if the NFT was already transferred

## Additional Resources

### Network Information

- **Testnet (fivenet)**: For development and testing
- **Mainnet (sixnet)**: For production use

### Useful Links

- SDK Documentation: Check the main README
- API Reference: See the pkg documentation
- Community Support: [Add your support channels]

### Best Practices

1. **Always use testnet first** before deploying to mainnet
2. **Save important values**: Contract addresses, schema names, transaction hashes
3. **Secure your mnemonic**: Never commit it to version control
4. **Test transfers** with small amounts first
5. **Verify transactions** before proceeding to next steps

### Configuration Tips

Create a `config.go` file to store common values:

```go
package main

const (
    // Network
    IsMainnet = false  // false = testnet, true = mainnet
    
    // Account
    MyMnemonic = "your mnemonic here"
    
    // Deployed Resources
    MySchemaName = "myorg.lbbv01"
    MyContractAddress = "0x..."
)
```

Then import these values in your scripts instead of hardcoding them.

## Next Steps

After completing these tutorials, you can:

1. **Build an application** using the SDK
2. **Create custom metadata schemas** for your use case
3. **Integrate with a frontend** to display NFTs
4. **Deploy to mainnet** when ready for production

## Need Help?

If you encounter issues not covered here:

1. Check the main SDK documentation
2. Review the error messages carefully
3. Ensure all prerequisites are completed
4. Verify your configuration values
5. Reach out to the community for support

---

Happy building! 🚀