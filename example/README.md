# LBB SDK Go - Examples

This directory contains example code demonstrating how to use the LBB SDK Go for certificate management on the SIX Protocol.

## Examples

### 1. Complete Example (`main.go`)

A comprehensive example showcasing all SDK features including:

- ccount creation and management
- Balance transfers
- Schema deployment
- Metadata minting
- Certificate freezing/unfreezing
- EVM contract deployment
- NFT minting and transfers
- Ownership verification

**Run the complete example:**
```bash
cd example
go run main.go
```

## Prerequisites

1. **Go 1.21+** installed
2. **Network access** to testnet (fivenet) or mainnet (sixnet)
3. **Funded account** with tokens for gas fees

## Getting Test Tokens

For testing on fivenet (testnet), you'll need test tokens:

1. Generate a new wallet using the SDK
2. Request faucet tokens from the SIX Protocol Discord [Faucet](https://discord.com/channels/940155834426613811/1352152627290443858) or Telegram
3. Use the wallet address to receive test tokens


## Configuration

### Network Selection

**Testnet (Fivenet):**
```go
client, err := client.NewClient(context.Background(), false)
```

**Mainnet (Sixnet):**
```go
client, err := client.NewClient(context.Background(), true)
```

### Custom Local Node

For local development with your own node:
```go
client, err := client.NewCustomClient(
    context.Background(),
    "http://localhost:26657",  // Tendermint RPC
    "http://localhost:1317",   // Cosmos REST API
    "http://localhost:8545",   // EVM JSON-RPC
    "testnet",                 // Chain ID type
)
```

## Understanding the Flow

### Certificate Issuance Flow

1. **Create Wallet** → Generate mnemonic and derive addresses
2. **Deploy Schema** → Define certificate structure on Cosmos layer
3. **Deploy Contract** → Deploy NFT contract on EVM layer
4. **Mint NFT** → Create NFT token on EVM layer
5. **Create Metadata** → Attach certificate data on Cosmos layer
6. **Transfer** _(optional)_ → Transfer ownership to recipient
7. **Manage** _(optional)_ → Freeze/unfreeze, update metadata as needed


### Dual-Layer Architecture

The LBB SDK operates on two layers:

**Cosmos Layer (Gen2 Data Layer):**
- Schema definitions
- Certificate metadata
- Data mutations and transforms
- Certificate freezing/unfreezing

**EVM Layer:**
- NFT smart contracts
- Token minting
- Token transfers
- Ownership tracking

Both layers work together to provide a complete certificate management solution.

## Common Patterns

### Method Usage
For existing NFT Contract and Schema. After init client and account if using the same mnemonic of these owner (contract and schema)
```go
client, err := client.NewClient(ctx, false)
...
acc := account.NewAccount(client, "quickstart", mnemonic, "")
meta := metadata.NewMetadataMsg(*acc, schemaName)
evmClient := evm.NewEVMClient(*acc)

createMetadata, err := meta.CreateCertificateMetadata(tokenIdStr)
if err != nil {
    // error handling
}

mintTx, err := evmClient.MintCertificateNFT(contractAddress, tokenId)
if err != nil {
    // error handling
}

```
> **Note:** On the Cosmos layer, you can append multiple messages and broadcast them in one transaction using `meta.BroadcastTx(msg)`. However, for reliability, it's recommended to deploy the certificate schema and wait for transaction confirmation before creating certificate metadata (e.g., use `meta.DeployCertificateSchema()`, wait for confirmation, then call `meta.CreateCertificateMetadata()`).

> **NOTE** We can CreateCertificateMetadata and MintCertificateNFT at the same time without invalid nonce or suquence of transaction, because both of them are on the separete layer

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
acc := account.NewAccount(client, "my-account", mnemonic, "")
if acc == nil {
    panic("Failed to create account")
}
```

### Schema Naming

Schema codes follow the format: `{ORGNAME}.{SCHEMACODE}`

```go
const schemaName = "myorg.lbbv01"
```

Choose a unique organization name to avoid conflicts.

### Balance Operations

The SDK provides balance query and transfer operations:

```go
// Query balances
balanceClient := balance.BalanceClient{Account: *acc}

// Get all balances
allBalances, err := balanceClient.GetBalance()
if err != nil {
    panic(fmt.Sprintf("Failed to get balance: %v", err))
}
fmt.Printf("All balances: %v\n", allBalances)

// Get Cosmos layer balance (usix)
cosmosBalance, err := balanceClient.GetCosmosBalance()
if err != nil {
    panic(fmt.Sprintf("Failed to get cosmos balance: %v", err))
}
fmt.Printf("Cosmos balance: %v\n", cosmosBalance)

// Get EVM layer balance (asix)
evmBalance, err := balanceClient.GetEVMBalance()
if err != nil {
    panic(fmt.Sprintf("Failed to get EVM balance: %v", err))
}
fmt.Printf("EVM balance: %v\n", evmBalance)

// Transfer tokens
balanceMsg := balance.NewBalanceMsg(*acc)
amount := sdk.NewCoins(sdk.NewInt64Coin("usix", 1000000))
res, err := balanceMsg.SendBalance("6x1recipient_address_here", amount)
if err != nil {
    panic(fmt.Sprintf("Failed to send balance: %v", err))
}
fmt.Printf("Transfer tx: %s\n", res.TxHash)
```

### Metadata Query Operations

Query certificate schemas and metadata:

```go
metaClient := metadata.NewMetadataClient(*acc)

// Get NFT schema details
schema, err := metaClient.GetNFTSchema(schemaName)
if err != nil {
    panic(fmt.Sprintf("Failed to get schema: %v", err))
}
fmt.Printf("Schema: %v\n", schema)

// Get certificate metadata
nftData, err := metaClient.GetNFTMetadata(schemaName, "1")
if err != nil {
    panic(fmt.Sprintf("Failed to get metadata: %v", err))
}
fmt.Printf("NFT Data: %v\n", nftData)

// Check if address is executor
isExecutor, err := metaClient.GetIsExecutor(schemaName, acc.GetCosmosAddress().String())
if err != nil {
    panic(fmt.Sprintf("Failed to check executor: %v", err))
}
fmt.Printf("Is executor: %v\n", isExecutor)

// Get all executors for schema
executors, err := metaClient.GetExecutor(schemaName)
if err != nil {
    panic(fmt.Sprintf("Failed to get executors: %v", err))
}
fmt.Printf("Executors: %v\n", executors)
```

### Certificate Freeze/Unfreeze

Manage certificate state on the Cosmos layer:

```go
meta := metadata.NewMetadataMsg(*acc, schemaName)

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

### EVM Query Operations

Query NFT ownership and contract details:

```go
evmClient := evm.NewEVMClient(*acc)

// Get token owner
owner := evmClient.TokenOwner(contractAddress, tokenId)
fmt.Printf("Token owner: %s\n", owner.Hex())

// Get current gas price
gasPrice, err := evmClient.GasPrice()
if err != nil {
    panic(fmt.Sprintf("Failed to get gas price: %v", err))
}
fmt.Printf("Gas price: %v\n", gasPrice)

// Get chain ID
chainID, err := evmClient.ChainID()
if err != nil {
    panic(fmt.Sprintf("Failed to get chain ID: %v", err))
}
fmt.Printf("Chain ID: %v\n", chainID)

// Get nonce
nonce, err := evmClient.GetNonce()
if err != nil {
    panic(fmt.Sprintf("Failed to get nonce: %v", err))
}
fmt.Printf("Nonce: %d\n", nonce)

// Check transaction receipt
err = evmClient.CheckTransactionReceipt(tx.Hash())
if err != nil {
    panic(fmt.Sprintf("Transaction check failed: %v", err))
}
```

### Auto-Increment Token ID Contract

Deploy and use a contract with auto-incrementing token IDs:

```go
// Deploy auto-increment contract
contractAddress, tx, err := evmClient.DeployCertIDIncrementContract(
    "AutoCert",
    "ACERT",
    schemaName,
)
if err != nil {
    panic(fmt.Sprintf("Failed to deploy contract: %v", err))
}

_, err = client.WaitForEVMTransaction(tx.Hash())
if err != nil {
    panic(fmt.Sprintf("Error waiting for contract deployment: %v", err))
}

// Mint with auto-increment (no need to specify token ID)
mintTx, err := evmClient.MintCertificateNFTAutoIncrement(contractAddress)
if err != nil {
    panic(fmt.Sprintf("Failed to mint NFT: %v", err))
}

_, err = client.WaitForEVMTransaction(mintTx.Hash())
if err != nil {
    panic(fmt.Sprintf("Error waiting for mint: %v", err))
}

// Get the current token ID counter
currentTokenId, err := evmClient.GetCurrentTokenID(contractAddress)
if err != nil {
    panic(fmt.Sprintf("Failed to get current token ID: %v", err))
}
fmt.Printf("Current token ID: %d\n", currentTokenId)
```

## Modifying the Examples

### Change Network

Edit the client initialization:
```go
// Change true to false for mainnet
client, err := client.NewClient(context.Background(), true)
```

### Use Your Own Mnemonic

Replace the test mnemonic with your own:
```go
const myMnemonic = "your twelve word mnemonic phrase..."
acc := account.NewAccount(client, "my-wallet", myMnemonic, "")
```

### Customize Certificate Details

Modify the schema name and contract details:
```go
const (
    contractName   = "YourCertificateName"
    contractSymbol = "YOURCERT"
    schemaName     = "yourorg.yourschema"
)
```

## Troubleshooting

### "Failed to create account"
- Verify your mnemonic is valid (12 or 24 words)
- Ensure the network is accessible
- Check that the client is properly initialized

### "Insufficient funds" error
- Ensure your account has enough tokens for gas fees
- For testnet, request tokens from the faucet
- Check your balance before transactions

### Transaction timeout
- Network might be congested
- Increase timeout in client configuration
- Verify node connectivity

### Schema already exists
- Schema codes must be unique
- Use a different organization name or schema code
- Check existing schemas before deploying

## Next Steps

1. Review the [main USAGE.md](../USAGE.md) for detailed documentation
2. Explore the SDK source code to understand available methods
3. Build your own application using these examples as a template