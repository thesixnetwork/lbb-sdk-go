# LBB SDK Go - Quick Start Guide

Welcome to the LBB SDK Go! This guide will get you up and running in minutes.

## 🚀 Installation

```bash
# Install the SDK
go get github.com/thesixnetwork/lbb-sdk-go

# Navigate to examples
cd example

# Install dependencies
go mod download
```

## 📖 Documentation

- **[Tutorial Guide](./docs/TUTORIAL.md)** - Step-by-step beginner guide
- **[Workflow Architecture](./docs/WORKFLOW_ARCHITECTURE.md)** - Understanding dual-layer architecture
- **[Examples Guide](./example/README.md)** - Detailed documentation for all 12 examples
- **[Root README](./readme.md)** - Complete API reference

## 🎯 5-Minute Quick Start

### Option 1: Run Complete Example

```bash
cd example
go run main.go
```

This runs the complete workflow in one script, demonstrating:
- Wallet generation
- Account creation
- Schema deployment
- Metadata minting
- Contract deployment
- NFT minting
- NFT transfer
- Metadata freezing/unfreezing

### Option 2: Step-by-Step Examples

Run examples in sequence to learn each feature:

```bash
# 1. Generate a wallet
go run 01_generate_wallet.go

# 2. Create an account
go run 02_create_account.go

# 3. Deploy certificate schema
go run 03_deploy_schema.go

# 4. Mint certificate metadata
go run 04_1_mint_metadata.go

# 5. Deploy NFT contract
go run 05_deploy_contract.go
# ⚠️ Save the contract address from output!

# 6. Mint NFT (update contractAddress in file first)
go run 06_mint_nft.go

# 7. Transfer NFT (update contractAddress and recipientAddress)
go run 07_transfer_nft.go

# 8. Freeze/unfreeze metadata
go run 08_freeze_metadata.go

# 9. Query NFT information
go run 09_query_nft.go

# 10. Balance operations
go run 10_balance_operations.go

# 11. Query metadata and schema
go run 11_query_metadata.go

# 12. Query EVM information
go run 12_query_evm.go
```

## 🔑 Key Concepts

### Dual-Layer Architecture

The SDK operates on two independent layers:

```
┌──────────────────────────┐  ┌──────────────────────────┐
│   Cosmos Layer (usix)    │  │    EVM Layer (asix)      │
│  ────────────────────    │  │  ──────────────────      │
│  • Schema Deployment     │  │  • Contract Deployment   │
│  • Metadata Management   │  │  • NFT Minting           │
│  • Certificate Freezing  │  │  • NFT Transfers         │
│  • Address: 6x...        │  │  • Address: 0x...        │
└──────────────────────────┘  └──────────────────────────┘
         ↓                              ↓
   Can execute simultaneously without conflicts!
```

**Key Advantage:** Deploy Schema and Deploy Contract at the same time for faster execution!

### Network Selection

```go
// Testnet (fivenet) - for development
client, err := client.NewClient(ctx, false)

// Mainnet (sixnet) - for production
client, err := client.NewClient(ctx, true)
```

## 💡 Common Operations

### Balance Operations

```go
import "github.com/thesixnetwork/lbb-sdk-go/pkg/balance"

// Query balances
bal := balance.NewBalance(*acc)
cosmosBalance, _ := bal.GetCosmosBalance()  // usix
evmBalance, _ := bal.GetEVMBalance()        // asix

// Transfer tokens
balMsg, _ := balance.NewBalanceMsg(*acc)
amount := sdk.NewCoins(sdk.NewInt64Coin("usix", 1000000))
res, _ := balMsg.SendBalanceAndWait(recipientAddr, amount)
```

### Query Operations

```go
import "github.com/thesixnetwork/lbb-sdk-go/pkg/metadata"

// Query schema
metaClient, _ := metadata.NewMetadata(*acc)
schema, _ := metaClient.GetNFTSchema(schemaName)

// Query metadata
nftData, _ := metaClient.GetNFTMetadata(schemaName, tokenId)

// Check permissions
isExecutor, _ := metaClient.GetIsExecutor(schemaName, address)
```

### EVM Queries

```go
import "github.com/thesixnetwork/lbb-sdk-go/pkg/evm"

evmClient := evm.NewEVMClient(*acc)

// Get gas price
gasPrice, _ := evmClient.GasPrice()

// Get nonce
nonce, _ := evmClient.GetNonce()

// Query NFT owner
owner := evmClient.TokenOwner(contractAddress, tokenId)
```

## 🎓 Learning Path

### Beginner

1. ✅ Run `main.go` to see the complete workflow
2. ✅ Read the [Tutorial Guide](./docs/TUTORIAL.md)
3. ✅ Run examples 01-09 in sequence

### Intermediate

4. ✅ Learn about [balance operations](./example/10_balance_operations.go)
5. ✅ Explore [metadata queries](./example/11_query_metadata.go)
6. ✅ Understand [EVM queries](./example/12_query_evm.go)
7. ✅ Study [Workflow Architecture](./docs/WORKFLOW_ARCHITECTURE.md)

### Advanced

8. ✅ Implement parallel execution for performance
9. ✅ Build custom schemas for your use case
10. ✅ Integrate with your application
11. ✅ Deploy to production on mainnet

## 🔥 Parallel Execution (Advanced)

Execute operations on different layers simultaneously:

```go
var wg sync.WaitGroup
wg.Add(2)

// Deploy Schema (Cosmos) - runs in parallel
go func() {
    defer wg.Done()
    msg, _ := meta.BuildDeployMsg()
    meta.BroadcastTxAndWait(msg)
}()

// Deploy Contract (EVM) - runs in parallel
go func() {
    defer wg.Done()
    evmClient.DeployCertificateContract(name, symbol, schemaName)
}()

wg.Wait()
```

**Result:** 42% faster than sequential execution! See [Workflow Architecture](./docs/WORKFLOW_ARCHITECTURE.md) for details.

## 🧪 Testing

### Get Test Tokens

For testnet (fivenet):

1. Generate a wallet: `go run 01_generate_wallet.go`
2. Request faucet tokens from [SIX Protocol Discord](https://discord.com/channels/940155834426613811/1352152627290443858)
3. Use your wallet address to receive tokens

### Using Test Mnemonic

All examples use a pre-funded test mnemonic:

```go
const exampleMnemonic = account.TestMnemonic
```

This works out of the box on testnet!

## 📊 Complete Workflow

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

✨ Deploy Schema and Deploy NFT Contract can run simultaneously!
```

## ⚡ Performance Tips

1. **Use Parallel Execution**: Run Cosmos and EVM operations simultaneously
2. **Batch Cosmos Messages**: Combine multiple messages in one transaction
3. **Wait for Confirmations**: Always confirm before dependent operations
4. **Query Before Actions**: Verify state before making changes

## 🛠️ Troubleshooting

### Common Issues

**"Insufficient funds"**
```bash
# Check balance
go run 10_balance_operations.go
# Request faucet tokens if needed
```

**"Schema already exists"**
```bash
# Use a unique schema name
const schemaName = "myorg.unique123"
```

**"Contract address not found"**
```bash
# Make sure to update contractAddress constant
# Copy from step 05_deploy_contract.go output
```

**"Nonce issues"**
```bash
# Query current nonce
go run 12_query_evm.go
# Wait for pending transactions to complete
```

## 📚 Example Projects

Build these with the SDK:

- **Certificate Management Portal** - Issue and verify certificates
- **Supply Chain Tracking** - Track products with NFT certificates
- **Credential Verification** - Educational or professional credentials
- **Digital Asset Registry** - Maintain authenticated asset records
- **Compliance Systems** - Automate certificate lifecycle management

## 🎯 Next Steps

1. **Explore Examples**: Check out all 12 examples in the [example/](./example/) directory
2. **Read Architecture**: Understand the [Workflow Architecture](./docs/WORKFLOW_ARCHITECTURE.md)
3. **Build Your App**: Use examples as templates for your use case
4. **Deploy Production**: Test on testnet, then deploy to mainnet

## 🔐 Security Notes

⚠️ **Important Security Practices:**

- Never commit mnemonics to version control
- Use environment variables for sensitive data
- Test thoroughly on testnet before mainnet
- Secure production keys with proper key management
- Always validate transaction results via [testnet-scham](https://sixchain.io/fivenet/), [mainnet-scham](https://sixchain.io/sixnet/)

## 📖 Additional Resources

| Resource | Description |
|----------|-------------|
| [Tutorial Guide](./docs/TUTORIAL.md) | Complete beginner tutorial |
| [Workflow Architecture](./docs/WORKFLOW_ARCHITECTURE.md) | Dual-layer architecture guide |
| [Examples README](./example/README.md) | Detailed examples documentation |
| [Root README](./readme.md) | API reference and features |
| [SIX Protocol Docs](https://docs.sixprotocol.net) | Official protocol documentation |

---

**Ready to build?** Start with `go run main.go` and explore the examples! 🚀

For questions, check the [Tutorial Guide](./docs/TUTORIAL.md) or reach out to the community.
