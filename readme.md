# LBB SDK Go

THE LBB SDK is a comprehensive software development kit designed to facilitate the creation, management, and deployment of applications using the SIX Protocol Dynamic NFT Metadata. It provides developers with a suite of tools, libraries, and documentation to streamline the development process.

**New to LBB SDK?** See the [Quick Start Guide](./QUICKSTART.md) for a 5-minute introduction!

## Quick Start

```bash
# Install the SDK
go get github.com/thesixnetwork/lbb-sdk-go

# Run the complete example
cd example
go run main.go

# Or run step-by-step examples
go run 01_generate_wallet.go
go run 02_create_account.go
go run 03_deploy_schema.go
# ... and so on
```

For detailed usage instructions, see [Tutorial](./docs/TUTORIAL.md) and [Examples Guide](./example/README.md).

## Features

### Core Features

- **🔐 Account Management**
  - Generate secure mnemonics (24 words)
  - Create accounts with Cosmos and EVM addresses
  - Keyring integration for secure key storage

- **📜 Schema Operations**
  - Deploy certificate schemas with custom attributes
  - Define metadata structure for certificates
  - Manage schema executors and permissions

- **🎨 NFT Contract Management**
  - Deploy ERC-721 compatible contracts
  - Link contracts to certificate schemas
  - Auto-increment or custom token IDs

- **💎 NFT Operations**
  - Mint certificate NFTs on EVM layer
  - Transfer NFT ownership
  - Query token ownership on-chain

- **📊 Metadata Management**
  - Create certificate metadata with custom attributes
  - Update metadata information
  - Freeze/unfreeze certificates status on chain

- **💰 Balance Operations**
  - Query balances (Cosmos & EVM layers)
  - Transfer tokens between addresses
  - Support for multiple denominations (usix, asix)

- **🔍 Query Capabilities**
  - Query NFT schemas and metadata
  - Retrieve gas prices and chain information
  - Verify NFT ownership
  - Transaction receipt verification

## Key Capabilities

- ✅ **Account Management** - Generate wallets, create accounts, manage keys
- ✅ **Schema Deployment** - Deploy and manage certificate schemas
- ✅ **NFT Contract Deployment** - Deploy EVM-compatible NFT contracts
- ✅ **Metadata Management** - Create, update, freeze/unfreeze certificate metadata
- ✅ **NFT Minting & Transfer** - Mint NFTs and transfer ownership
- ✅ **Balance Operations** - Query and transfer tokens (Cosmos & EVM layers)
- ✅ **Query Operations** - Query schemas, metadata, NFT ownership, and blockchain state
- ✅ **Dual-Layer Architecture** - Seamless integration between Cosmos and EVM layers

## Documentation

- **[Quick Start Guide](./QUICKSTART.md)** - 5-minute introduction for new users
- **[Tutorial Guide](./docs/TUTORIAL.md)** - Comprehensive tutorial with step-by-step examples
- **[Examples Guide](./example/README.md)** - Detailed guide for all example scripts with usage patterns
- **[Workflow Architecture](./docs/WORKFLOW_ARCHITECTURE.md)** - Dual-layer architecture and parallel execution guide
- **[Example Code](./example/)** - 12+ working examples demonstrating all SDK features
- For detailed API references and advanced tutorials, visit the official LBB SDK documentation site

## Example Scripts

The SDK includes comprehensive examples covering all features:

1. **[01_generate_wallet.go](./example/01_generate_wallet.go)** - Generate new wallet with mnemonic
2. **[02_create_account.go](./example/02_create_account.go)** - Connect to network and create account
3. **[03_deploy_schema.go](./example/03_deploy_schema.go)** - Deploy certificate schema
4. **[04_mint_metadata.go](./example/04_1_mint_metadata.go)** - Mint certificate metadata
5. **[05_deploy_contract.go](./example/05_deploy_contract.go)** - Deploy EVM NFT contract
6. **[06_mint_nft.go](./example/06_mint_nft.go)** - Mint certificate NFT
7. **[07_transfer_nft.go](./example/07_transfer_nft.go)** - Transfer NFT ownership
8. **[08_freeze_metadata.go](./example/08_freeze_metadata.go)** - Freeze/unfreeze metadata
9. **[09_query_nft.go](./example/09_query_nft.go)** - Query NFT information
10. **[10_balance_operations.go](./example/10_balance_operations.go)** - Query and transfer balances
11. **[11_query_metadata.go](./example/11_query_metadata.go)** - Query schema and certificate metadata
12. **[12_query_evm.go](./example/12_query_evm.go)** - Query EVM information

See the [Examples README](./example/README.md) for detailed usage instructions.

## Network Support

- **Fivenet** (Testnet) - For development and testing
  - Chain ID: fivenet
  - Free test tokens available via faucet
  - Use `client.NewClient(ctx, false)`

- **Sixnet** (Mainnet) - For production use
  - Chain ID: sixnet
  - Real tokens required
  - Use `client.NewClient(ctx, true)`

## Common Operations

### Balance Operations

Query and transfer tokens on both Cosmos and EVM layers:

```go
// Create balance client for queries
bal := balance.NewBalance(*acc)

// Get all balances
allBalances, err := bal.GetBalance()
fmt.Printf("All balances: %v\n", allBalances)

// Get Cosmos layer balance (usix)
cosmosBalance, err := bal.GetCosmosBalance()
fmt.Printf("Cosmos balance: %v\n", cosmosBalance)

// Get EVM layer balance (asix)
evmBalance, err := bal.GetEVMBalance()
fmt.Printf("EVM balance: %v\n", evmBalance)

// Transfer tokens
balMsg, err := balance.NewBalanceMsg(*acc)
amount := sdk.NewCoins(sdk.NewInt64Coin("usix", 1000000))
res, err := balMsg.SendBalanceAndWait("6x1recipient_address", amount)
fmt.Printf("Transfer tx: %s\n", res.TxHash)
```

**Balance Denominations:**
- **usix**: Cosmos layer token (1 SIX = 1,000,000 usix)
- **asix**: EVM layer token (1 SIX = 1,000,000,000,000,000,000 asix)

### Metadata Query Operations

Query certificate schemas and metadata:

```go
metaClient := metadata.NewMetadata(*acc)

// Get NFT schema details
schema, err := metaClient.GetNFTSchema(schemaName)
fmt.Printf("Schema Owner: %s\n", schema.NftSchemaBase.Owner)

// Get certificate metadata
nftData, err := metaClient.GetNFTMetadata(schemaName, "1")
fmt.Printf("Token Owner: %s\n", nftData.OwnerAddressType)

// Check executor permissions
isExecutor, err := metaClient.GetIsExecutor(schemaName, address)
fmt.Printf("Is executor: %v\n", isExecutor)

// Get all executors
executors, err := metaClient.GetExecutor(schemaName)
fmt.Printf("Executors: %v\n", executors)
```

### EVM Query Operations

Query NFT ownership and blockchain state:

```go
evmClient := evm.NewEVMClient(*acc)

// Get current gas price
gasPrice, err := evmClient.GasPrice()
fmt.Printf("Gas price: %v wei\n", gasPrice)

// Get chain ID
chainID, err := evmClient.ChainID()
fmt.Printf("Chain ID: %v\n", chainID)

// Get account nonce
nonce, err := evmClient.GetNonce()
fmt.Printf("Nonce: %d\n", nonce)

// Get token owner
owner := evmClient.TokenOwner(contractAddress, tokenId)
fmt.Printf("Token owner: %s\n", owner.Hex())

// Check transaction receipt
err = evmClient.CheckTransactionReceipt(tx.Hash())
```

### Certificate Management

Freeze and unfreeze certificates:

```go
meta, err := metadata.NewMetadataMsg(*acc, schemaName)

// Freeze a certificate
res, err := meta.FreezeCertificate("1")
fmt.Printf("Certificate frozen, tx: %s\n", res.TxHash)

// Wait for confirmation
err = client.WaitForTransaction(res.TxHash)

// Unfreeze a certificate
res, err = meta.UnfreezeCertificate("1")
fmt.Printf("Certificate unfrozen, tx: %s\n", res.TxHash)
```

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

// Create account
acc, err := account.NewAccount(client, "my-account", mnemonic, "")
if err != nil {
    panic(fmt.Sprintf("Failed to create account: %v", err))
}
```

### Schema Naming Convention

Schema codes follow the format: `{ORGNAME}.{SCHEMACODE}`

```go
const schemaName = "myorg.lbbv01"
```

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

## Troubleshooting

### Common Issues

**"Failed to create account"**
- Verify your mnemonic is valid (12 or 24 words)
- Ensure the network is accessible
- Check that the client is properly initialized

**"Insufficient funds"**
- Ensure your account has enough tokens for gas fees
- For testnet, request tokens from the faucet
- Check your balance using balance query operations

**"Schema already exists"**
- Schema codes must be unique
- Use a different organization name or schema code
- Check existing schemas using `GetNFTSchema()`

**Transaction timeout**
- Network might be congested
- Increase timeout in client configuration
- Verify node connectivity

**Nonce issues (EVM transactions)**
- Query current nonce using `GetNonce()`
- Wait for pending transactions to confirm
- Each transaction increments the nonce automatically

## Getting Test Tokens

For testing on fivenet (testnet), you'll need test tokens:

1. Generate a new wallet using the SDK
2. Request faucet tokens from the SIX Protocol [Discord](https://discord.com/channels/940155834426613811/1352152627290443858) or Telegram
3. Use the wallet address to receive test tokens

## Workflow Overview

```
Generate Wallet → Create Account
    ↓
    ┌─────────────────────────────────┬─────────────────────────────────┐
    ↓                                 ↓                                 ↓
Deploy Schema                  Deploy NFT Contract           Balance Operations
(Cosmos layer)                 (EVM layer)                   (Query & Transfer)
    ↓                                 ↓
Mint Metadata                  Mint NFT
(Cosmos layer)                 (EVM layer)
    ↓                                 ↓
Freeze/Unfreeze                Transfer NFT
Metadata (Optional)            (Optional)
    ↓                                 ↓
Query Metadata                 Query NFT Information
& Schema                              ↓
    ↓                          Query EVM Information
    └─────────────────────────────────┘

Note: Deploy Schema and Deploy NFT Contract can be executed simultaneously
      (they operate on separate Cosmos and EVM layers)
```

## Quick Reference

### Initialize Client

```go
// Testnet (fivenet)
client, err := client.NewClient(ctx, false)

// Mainnet (sixnet)
client, err := client.NewClient(ctx, true)

// Custom local node
client, err := client.NewCustomClient(
    ctx,
    "http://localhost:26657",  // Tendermint RPC
    "http://localhost:1317",   // Cosmos REST API
    "http://localhost:8545",   // EVM JSON-RPC
    "testnet",                 // Chain ID type
)
```

### Create Account

```go
// Generate new mnemonic
mnemonic, err := account.GenerateMnemonic()

// Create account from mnemonic
acc, err := account.NewAccount(client, "my-account", mnemonic, "")

// Get addresses
cosmosAddr := acc.GetCosmosAddress().String()  // 6x...
evmAddr := acc.GetEVMAddress().Hex()           // 0x...
```

### Deploy Schema

```go
meta, err := metadata.NewMetadataMsg(*acc, "myorg.schema01")
msg, err := meta.BuildDeployMsg()
res, err := meta.BroadcastTxAndWait(msg)
```

### Deploy Contract

```go
evmClient := evm.NewEVMClient(*acc)
contractAddr, tx, err := evmClient.DeployCertificateContract(
    "MyCertificate",
    "CERT",
    "myorg.schema01",
)
```

### Mint NFT

```go
tx, err := evmClient.MintCertificateNFT(contractAddress, tokenId)
_, err = client.WaitForEVMTransaction(tx.Hash())
```

### Query Balance

```go
bal := balance.NewBalance(*acc)
cosmosBalance, err := bal.GetCosmosBalance()
evmBalance, err := bal.GetEVMBalance()
```

### Transfer Tokens

```go
balMsg, err := balance.NewBalanceMsg(*acc)
amount := sdk.NewCoins(sdk.NewInt64Coin("usix", 1000000))
res, err := balMsg.SendBalanceAndWait(recipientAddress, amount)
```

### Query NFT Owner

```go
owner := evmClient.TokenOwner(contractAddress, tokenId)
```

## Additional Resources

### Learn More

- **[Quick Start Guide](./QUICKSTART.md)** - 5-minute introduction for new users
- **[Tutorial Guide](./docs/TUTORIAL.md)** - Step-by-step guide for beginners
- **[Examples Guide](./example/README.md)** - Detailed documentation for all examples
- **[Workflow Architecture](./docs/WORKFLOW_ARCHITECTURE.md)** - Understanding dual-layer architecture and parallel execution
- **[Example Code](./example/)** - 12+ working code examples
- **[SIX Protocol Docs](https://docs.sixprotocol.net)** - Official protocol documentation

### Use Cases

- **Certificate Management**: Issue, verify, and manage digital certificates
- **Supply Chain Tracking**: Track products with certificate NFTs
- **Credential Verification**: Issue educational or professional credentials
- **Digital Asset Registry**: Maintain authenticated asset registries
- **Compliance Systems**: Automate compliance certificate management
- **Identity Systems**: Build decentralized identity solutions

## Install Submodule

- **Forge Standard**

```bash
git submodule add https://github.com/foundry-rs/forge-std ./contracts/lib/forge-std
```

- **Openzeppelin**

```bash
git submodule add https://github.com/openzeppelin/openzeppelin-contracts ./contracts/lib/openzeppelin-contracts
```

## License

See [LICENSE](./LICENSE) file for details.

---

**Ready to get started?** Check out the [Quick Start Guide](./QUICKSTART.md) for a 5-minute intro, or dive into the [Tutorial Guide](./docs/TUTORIAL.md)! 🚀
