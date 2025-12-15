# LBB SDK Go - Examples

Welcome to the LBB SDK Go examples directory! This collection of examples will help you understand how to use the SDK to create and manage certificate NFTs with metadata on the blockchain.

## What's in This Directory

### Step-by-step Examples

Step-by-step examples that teach you each function:

1. **[01_generate_wallet.go](./01_generate_wallet.go)** - Generate a new wallet with mnemonic
2. **[02_create_account.go](./02_create_account.go)** - Connect to network and create account
3. **[03_deploy_schema.go](./03_deploy_schema.go)** - Deploy certificate schema
4. **[04_mint_metadata.go](./04_mint_metadata.go)** - Mint certificate metadata
5. **[05_deploy_contract.go](./05_deploy_contract.go)** - Deploy EVM NFT contract
6. **[06_mint_nft.go](./06_mint_nft.go)** - Mint certificate NFT
7. **[07_transfer_nft.go](./07_transfer_nft.go)** - Transfer NFT to another address
8. **[08_freeze_metadata.go](./08_freeze_metadata.go)** - Freeze and unfreeze metadata
9. **[09_query_nft.go](./09_query_nft.go)** - Query NFT information
10. **[10_balance_operations.go](./10_balance_operations.go)** - Query and transfer balances
11. **[11_query_metadata.go](./11_query_metadata.go)** - Query schema and certificate metadata
12. **[12_query_evm.go](./12_query_evm.go)** - Query EVM information (gas, nonce, ownership)

### Full Example

- **[main.go](./main.go)** - Complete workflow in a single file

### Documentation
- **This README** - Overview and quick start

## Quick Start

#### Learn Step by Step (Recommended)

Follow the numbered examples in order:

```bash
go run 01_generate_wallet.go
go run 02_create_account.go
go run 03_deploy_schema.go
go run 04_mint_metadata.go
go run 05_deploy_contract.go
# ... and so on
```

#### Path B: Run Complete Example

Run the full workflow:

```bash
go run main.go
```

### 3. Update Configuration

After running deployment scripts, update these values in subsequent scripts:

- **Contract Address**: From `05_deploy_contract.go` output
- **Schema Name**: Your chosen schema name from `03_deploy_schema.go`
- **Mnemonic**: Use test mnemonic or your own

## Documentation

- **[Tutorial Guide](../docs/TUTORIAL.md)** - Comprehensive step-by-step tutorial
- **[Workflow Architecture](../docs/WORKFLOW_ARCHITECTURE.md)** - Understanding dual-layer architecture and parallel execution
- **[Root README](../readme.md)** - Main SDK documentation and quick reference

## This example covers with

1. **Account Management**
   - Generate wallets and mnemonics
   - Create accounts with Cosmos and EVM addresses
   - Connect to testnet and mainnet

2. **Schema Operations**
   - Deploy certificate schemas
   - Manage metadata structure
   - Mint metadata instances
   - Freeze and unfreeze certificates

3. **EVM Contract Operations**
   - Deploy NFT contracts
   - Link contracts to schemas
   - Interact with smart contracts

4. **NFT Management**
   - Mint certificate NFTs
   - Transfer ownership
   - Query token information

5. **Balance Operations**
   - Query all balances
   - Query Cosmos layer balance (usix)
   - Query EVM layer balance (asix)
   - Transfer tokens between addresses

6. **Query Operations**
   - Query NFT schemas and metadata
   - Check executor permissions
   - Query EVM information (gas price, chain ID, nonce)
   - Verify NFT ownership
   - Check transaction receipts

## Network Information

- **Testnet (fivenet)**: For development and testing
  - Use `client.NewClient(ctx, false)`
  - Free test tokens available

- **Mainnet (sixnet)**: For production
  - Use `client.NewClient(ctx, true)`
  - Real tokens required

## Common Operations

### Balance Operations

The SDK provides comprehensive balance query and transfer operations:

#### Query Balances

```go
// Create balance client for queries
bal:= balance.NewBalance(*acc)

// Get all balances
allBalances, err := bal.GetBalance()
if err != nil {
    panic(fmt.Sprintf("Failed to get balance: %v", err))
}
fmt.Printf("All balances: %v\n", allBalances)

// Get Cosmos layer balance (usix)
cosmosBalance, err := bal.GetCosmosBalance()
if err != nil {
    panic(fmt.Sprintf("Failed to get cosmos balance: %v", err))
}
fmt.Printf("Cosmos balance: %v\n", cosmosBalance)

// Get EVM layer balance (asix)
evmBalance, err := bal.GetEVMBalance()
if err != nil {
    panic(fmt.Sprintf("Failed to get EVM balance: %v", err))
}
fmt.Printf("EVM balance: %v\n", evmBalance)
```

#### Transfer Tokens

```go
// Create balance message client for transactions
balMsg, err := balance.NewBalanceMsg(*acc)
if err != nil {
    panic(fmt.Sprintf("Failed to create balance msg: %v", err))
}

// Define amount to send (1 SIX = 1,000,000 usix)
amount := sdk.NewCoins(sdk.NewInt64Coin("usix", 1000000))

// Send balance (returns immediately)
res, err := balMsg.SendBalance("6x1recipient_address_here", amount)
if err != nil {
    panic(fmt.Sprintf("Failed to send balance: %v", err))
}
fmt.Printf("Transfer tx: %s\n", res.TxHash)

// OR send and wait for confirmation
res, err = balMsg.SendBalanceAndWait("6x1recipient_address_here", amount)
if err != nil {
    panic(fmt.Sprintf("Failed to send balance: %v", err))
}
fmt.Printf("Transfer confirmed: %s\n", res.TxHash)
```

#### Balance Denominations

- **usix**: Cosmos layer token (1 SIX = 1,000,000 usix)
  - Used for Cosmos transactions (metadata, schemas, etc.)
- **asix**: EVM layer token (1 SIX = 1,000,000,000,000,000,000 asix)
  - Used for EVM transactions (contract deployment, NFT minting)

### Metadata Query Operations

Query certificate schemas and metadata on the Cosmos layer:

#### Query Schema Information

```go
metaClient:= metadata.NewMetadata(*acc)
// Get NFT schema details
schema, err := metaClient.GetNFTSchema(schemaName)
if err != nil {
    panic(fmt.Sprintf("Failed to get schema: %v", err))
}
fmt.Printf("Schema: %v\n", schema)
fmt.Printf("Owner: %s\n", schema.NftSchemaBase.Owner)
fmt.Printf("Contract Address: %s\n", schema.NftSchemaBase.OriginData.OriginContractAddress)
```

#### Query Certificate Metadata

```go
// Get certificate metadata for a specific token
nftData, err := metaClient.GetNFTMetadata(schemaName, "1")
if err != nil {
    panic(fmt.Sprintf("Failed to get metadata: %v", err))
}
fmt.Printf("Token ID: %s\n", nftData.TokenId)
fmt.Printf("Owner: %s\n", nftData.OwnerAddressType)

// Access certificate attributes
if nftData.OnchainAttributes != nil {
    fmt.Printf("Status: %s\n", nftData.OnchainAttributes.Status.Value)
    fmt.Printf("Weight: %s\n", nftData.OnchainAttributes.Weight.Value)
}
```

### EVM Query Operations

Query NFT ownership and EVM layer information:

#### Query Gas Price and Chain Info

```go
evmClient := evm.NewEVMClient(*acc)

// Get current gas price
gasPrice, err := evmClient.GasPrice()
if err != nil {
    panic(fmt.Sprintf("Failed to get gas price: %v", err))
}
fmt.Printf("Gas price: %v wei\n", gasPrice)

// Get chain ID
chainID, err := evmClient.ChainID()
if err != nil {
    panic(fmt.Sprintf("Failed to get chain ID: %v", err))
}
fmt.Printf("Chain ID: %v\n", chainID)

// Get nonce (transaction count)
nonce, err := evmClient.GetNonce()
if err != nil {
    panic(fmt.Sprintf("Failed to get nonce: %v", err))
}
fmt.Printf("Nonce: %d\n", nonce)
```

#### Query NFT Ownership

```go
// Get token owner
owner := evmClient.TokenOwner(contractAddress, tokenId)
fmt.Printf("Token owner: %s\n", owner.Hex())

// Verify ownership
if owner.Hex() == acc.GetEVMAddress().Hex() {
    fmt.Println("You own this NFT")
}
```

#### Check Transaction Receipt

```go
// Check transaction receipt
err = evmClient.CheckTransactionReceipt(tx.Hash())
if err != nil {
    panic(fmt.Sprintf("Transaction check failed: %v", err))
}
// This will print transaction details, gas used, and status
```

### Certificate Freeze/Unfreeze

Manage certificate state on the Cosmos layer:

```go
meta, err := metadata.NewMetadataMsg(*acc, schemaName)

// Freeze a certificate
res, err := meta.FreezeCertificate("1")
if err != nil {
    panic(fmt.Sprintf("Failed to freeze certificate: %v", err))
}
fmt.Printf("Certificate frozen, tx: %s\n", res.TxHash)

// Wait for confirmation
err = client.WaitForTransaction(res.TxHash)
if err != nil {
    panic(fmt.Sprintf("Transaction failed: %v", err))
}

// Unfreeze a certificate
res, err = meta.UnfreezeCertificate("1")
if err != nil {
    panic(fmt.Sprintf("Failed to unfreeze certificate: %v", err))
}
fmt.Printf("Certificate unfrozen, tx: %s\n", res.TxHash)
```

## Procedures

1. **Start with numbered examples** - They build on each other
2. **Use the test mnemonic** - It has testnet tokens for testing
3. **Save important values** - Contract addresses, schema names, tx hashes
4. **Read the comments** - Each example has detailed explanations
5. **EVM and Cosmos Layer** - User able to send transactions to each layer at the same time.

## Quick Reference

### Generate Wallet
```bash
go run 01_generate_wallet.go
```

### Create Account
```bash
go run 02_create_account.go
```

### Deploy Schema
```bash
go run 03_deploy_schema.go
# Save the schema name!
```

### Mint Metadata
```bash
go run 04_mint_metadata.go
# Uses the schema from step 3
```

### Deploy Contract
```bash
go run 05_deploy_contract.go
# Save the contract address!
```

### Mint NFT
```bash
# Update contractAddress in the file first
go run 06_mint_nft.go
```

### Transfer NFT
```bash
# Update contractAddress and recipientAddress
go run 07_transfer_nft.go
```

### Freeze/Unfreeze
```bash
go run 08_freeze_metadata.go
```

### Query NFT
```bash
# Update contractAddress
go run 09_query_nft.go
```

### Balance Operations
```bash
go run 10_balance_operations.go
# Uncomment code to send tokens
```

### Query Metadata
```bash
# Update schemaName
go run 11_query_metadata.go
```

### Query EVM Information
```bash
go run 12_query_evm.go
# Update contractAddress for NFT ownership queries
```

## Security Notes

- ⚠️ **Never commit mnemonics** to version control
- ⚠️ **Test on testnet first** before mainnet
- ⚠️ **Secure production keys** properly
- ⚠️ **Use environment variables** for sensitive data

## Best Practices

### Error Handling

Always check for errors and wait for transaction confirmation:

```go
res, err := meta.BroadcastTx(msg)
if err != nil {
    panic(fmt.Sprintf("Failed to broadcast: %v", err))
}

err = client.WaitForTransaction(res.TxHash)
if err != nil {
    panic(fmt.Sprintf("Transaction failed: %v", err))
}
```

### Account Creation

```go
// Generate new mnemonic
mnemonic, err := account.GenerateMnemonic()

// Or use existing mnemonic
mnemonic := "your twelve word mnemonic phrase here..."

// Create account
acc, err := account.NewAccount(client, "my-account", mnemonic, "")
if err != nil {
    panic(fmt.Sprintf("Failed to create account: %v", err))
}
```

### Schema Naming

Schema codes follow the format: `{ORGNAME}.{SCHEMACODE}`

```go
const schemaName = "myorg.lbbv01"
```

Choose a unique organization name to avoid conflicts.

### Transaction Best Practices

- **Check balances** before sending transactions to ensure sufficient funds
- **Use `SendBalanceAndWait()`** for confirmed transfers
- **Wait for confirmations** using `client.WaitForTransaction()` or `client.WaitForEVMTransaction()`
- **Handle errors gracefully** - don't assume transactions will succeed
- **Query first** - verify state before making changes

### Dual-Layer Transactions

The SDK supports simultaneous transactions on both Cosmos and EVM layers:

> **Important:** You can execute Deploy Schema (Cosmos) and Deploy Contract (EVM) at the same time without nonce/sequence conflicts because they operate on separate layers.

**Examples of parallel execution:**
- Deploy Schema + Deploy NFT Contract (different layers)
- Mint Metadata + Mint NFT (different layers)
- Freeze Metadata + Transfer NFT (different layers)

**Sequential execution required:**
- Multiple Cosmos transactions (must wait for confirmation)
- Multiple EVM transactions (nonce must increment)

> **Best Practice:** For reliability on the Cosmos layer, deploy the certificate schema and wait for transaction confirmation before creating certificate metadata (e.g., use `meta.DeployCertificateSchema()`, wait for confirmation, then call `meta.CreateCertificateMetadata()`). However, you can batch multiple messages in one Cosmos transaction using `meta.BroadcastTx(msg1, msg2, ...)`.

> **Note:** CreateCertificateMetadata and MintCertificateNFT can be executed at the same time without invalid nonce or sequence conflicts, because they operate on separate layers (Cosmos and EVM).

## Troubleshooting

### "Failed to create account"
- Verify your mnemonic is valid (12 or 24 words)
- Ensure the network is accessible
- Check that the client is properly initialized

### "Insufficient funds" error
- Ensure your account has enough tokens for gas fees
- For testnet, request tokens from the faucet
- Check your balance before transactions using balance query operations

### Transaction timeout
- Network might be congested
- Increase timeout in client configuration
- Verify node connectivity

### Schema already exists
- Schema codes must be unique
- Use a different organization name or schema code
- Check existing schemas before deploying using `GetNFTSchema()`

### "Failed to get metadata" or "NFT not found"
- Ensure the schema exists and is deployed
- Verify the token ID is correct
- Make sure metadata was created for that token
- Check that you're querying the correct schema name

### Nonce issues (EVM transactions)
- If nonce errors occur, query current nonce using `GetNonce()`
- Wait for pending transactions to confirm
- Each transaction increments the nonce automatically

## 📈 Workflow Overview

```
01. Generate Wallet
    ↓
02. Create Account
    ↓
    ┌─────────────────────────────┬─────────────────────────────┐
    ↓                             ↓                             ↓
03. Deploy Schema        05. Deploy NFT Contract    10. Balance Operations
    (Cosmos layer)           (EVM layer)                (Query & Transfer)
    ↓                             ↓
04. Mint Metadata         06. Mint NFT
    (Cosmos layer)           (EVM layer)
    ↓                             ↓
08. Freeze/Unfreeze       07. Transfer NFT
    Metadata (Optional)      (Optional)
    ↓                             ↓
11. Query Metadata        09. Query NFT Information
    & Schema                  ↓
    ↓                     12. Query EVM Information
    └─────────────────────────────┘

Note: Steps 03 and 05 can be executed simultaneously
      (separate Cosmos and EVM layers)
```

### Example Projects

Use these examples to build:
- **Certificate Management System**: Issue, verify, and manage certificates
- **Supply Chain Tracking**: Track products with certificate NFTs
- **Credential Verification**: Issue and verify educational or professional credentials
- **Digital Asset Registry**: Maintain a registry of authenticated assets
- **Compliance Systems**: Automate compliance certificate management

### Development Tips

- Start with the step-by-step examples to understand each component
- Use the complete example (main.go) as a reference for full workflows
- Query operations are free - use them liberally for verification
- Test all operations on testnet before mainnet deployment
- Keep transaction hashes for audit trails
- Implement proper logging for production applications
