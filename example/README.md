> [!IMPORTANT]
> If the example in `example/main.go` fails, please refer to `cmd/main.go` in the root directory—this is the most up-to-date usage pattern for the LBB SDK.
>
> The instructions, patterns, and methods here are for reference and may be slightly out of sync with the latest SDK features.  
> Always cross-check with `cmd/main.go` for the recommended workflow and error handling.

# LBB SDK Go - Examples

This directory contains example code demonstrating how to use the LBB SDK Go for certificate management on the SIX Protocol.

## Examples

### 1. Complete Example (`main.go`)

A comprehensive example showcasing all SDK features including:

- Account creation and management
- Mnemonic generation
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

If this fails or does not act as expected, **switch to**:
```bash
cd cmd
go run main.go
```

and check the code in `cmd/main.go` for the latest SDK usage.

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
if err != nil {
    panic(fmt.Sprintf("Failed to create client: %v", err))
}
```

**Mainnet (Sixnet):**
```go
client, err := client.NewClient(context.Background(), true)
if err != nil {
    panic(fmt.Sprintf("Failed to create client: %v", err))
}
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
if err != nil {
    panic(fmt.Sprintf("Failed to create client: %v", err))
}
```

## Understanding the Flow

### Certificate Issuance Flow

1. **Create Wallet** → Generate mnemonic and derive addresses
2. **Deploy Schema** → Define certificate structure on Cosmos layer
3. **Mint Metadata** → Create certificate metadata on Cosmos layer
4. **Deploy Contract** → Deploy NFT contract on EVM layer
5. **Mint NFT** → Create NFT token on EVM layer
6. **Transfer** _(optional)_ → Transfer ownership to recipient
7. **Manage** _(optional)_ → Freeze/unfreeze certificates as needed

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

_(All example snippets below are for reference. See [`cmd/main.go`](../cmd/main.go) if unsure about compatibility or error handling.)_

### Account Creation

```go
// Generate new mnemonic
mnemonic, err := account.GenerateMnemonic()
if err != nil {
    panic(fmt.Sprintf("Failed to generate mnemonic: %v", err))
}

fmt.Println("Mnemonic generated")
fmt.Println("*Important** write this mnemonic phrase in a safe place.")
fmt.Printf("\nMnemonic: %s\n\n", mnemonic)

// Create account from mnemonic
acc, err := account.NewAccount(client, "alice", mnemonic, "")
if err != nil {
    panic("ERROR CREATE ACCOUNT: NewAccount returned nil - check mnemonic and keyring initialization")
}

fmt.Printf("Account created\n")
fmt.Printf("  EVM Address: %s\n", acc.GetEVMAddress().Hex())
fmt.Printf("  Cosmos Address: %s\n", acc.GetCosmosAddress().String())
```

### Schema Naming

Schema codes follow the format: `{ORGNAME}.{SCHEMACODE}`

```go
const schemaName = "myorg.lbbv01"
```

Choose a unique organization name to avoid conflicts.

### Deploy Schema and Mint Metadata (Recommended Pattern)

The recommended pattern is to build multiple messages and broadcast them together:

```go
meta, err := metadata.NewMetadataMsg(*acc, schemaName)
if err != nil {
    fmt.Printf("NewMetadataMsg error: %v\n", err)
    return
}

// Build deploy schema message
msgDeploySchema, err := meta.BuildDeployMsg()
if err != nil {
    fmt.Printf("Failed to build deploy message: %v\n", err)
    return
}

// Build mint metadata message
msgCreateMetadata, err := meta.BuildMintMetadataMsg("1")
if err != nil {
    fmt.Printf("Failed to build create metadata: %v\n", err)
    return
}

// Combine messages and broadcast
var msgs []sdk.Msg
msgs = append(msgs, msgDeploySchema, msgCreateMetadata)

res, err := meta.BroadcastTxAndWait(msgs...)
if err != nil {
    fmt.Printf("Broadcast Tx error: %v\n", err)
    return
}

fmt.Printf("Schema deployed and metadata minted\n")
fmt.Printf("  Schema Code: %s\n", schemaName)
fmt.Printf("  Transaction: %s\n", res.TxHash)
```

> **Note:** `BroadcastTxAndWait()` will automatically wait for transaction confirmation. For more control, use `BroadcastTx()` followed by `client.WaitForTransaction()`.

### Certificate Freeze/Unfreeze

Manage certificate state on the Cosmos layer:

```go
meta, err := metadata.NewMetadataMsg(*acc, schemaName)
if err != nil {
    fmt.Printf("NewMetadataMsg error: %v\n", err)
    return
}

// Freeze a certificate
res, err := meta.FreezeCertificate("1")
if err != nil {
    fmt.Printf("Freeze error: %v\n", err)
    return
}

// Wait for confirmation
err = client.WaitForTransaction(res.TxHash)
if err != nil {
    fmt.Printf("Error waiting for freeze: %v\n", err)
    return
}

fmt.Printf("Certificate frozen, tx: %s\n", res.TxHash)

// Unfreeze a certificate
res, err = meta.UnfreezeCertificate("1")
if err != nil {
    fmt.Printf("Unfreeze error: %v\n", err)
    return
}

fmt.Printf("Certificate unfrozen, tx: %s\n", res.TxHash)
```

### Deploy EVM NFT Contract

```go
evmClient := evm.NewEVMClient(*acc)

// Deploy certificate contract
contractAddress, tx, err := evmClient.DeployCertificateContract(
    "MyCertificate",
    "CERT",
    schemaName,
)
if err != nil {
    fmt.Printf("EVM deploy error: %v\n", err)
    return
}

fmt.Printf("Deploy Tx: %s\n", tx.Hash().Hex())
fmt.Printf("Deploy at Nonce: %v\n", tx.Nonce())

// Wait for deployment transaction to be mined
_, err = client.WaitForEVMTransaction(tx.Hash())
if err != nil {
    fmt.Printf("Error waiting for deployment: %v\n", err)
    return
}

fmt.Printf("Contract deployed at: %s\n", contractAddress.Hex())
```

### Mint NFT

```go
evmClient := evm.NewEVMClient(*acc)

tokenId := uint64(1)
tx, err := evmClient.MintCertificateNFT(contractAddress, tokenId)
if err != nil {
    fmt.Printf("EVM mint error: %v\n", err)
    return
}

fmt.Printf("Mint Tx: %s\n", tx.Hash().Hex())
fmt.Printf("Mint at Nonce: %v\n", tx.Nonce())

// Wait for mint transaction
_, err = client.WaitForEVMTransaction(tx.Hash())
if err != nil {
    fmt.Printf("Error waiting for mint: %v\n", err)
    return
}

fmt.Printf("NFT minted with token ID: %d\n", tokenId)
```

### Transfer NFT

```go
import "github.com/ethereum/go-ethereum/common"

evmClient := evm.NewEVMClient(*acc)

recipientEVM := "0x8a28fb81A084Ac7A276800957a19a6054BF86E4D"
tokenId := uint64(1)

tx, err := evmClient.TransferCertificateNFT(
    contractAddress,
    common.HexToAddress(recipientEVM),
    tokenId,
)
if err != nil {
    fmt.Printf("EVM transfer error: %v\n", err)
    return
}

fmt.Printf("Transfer Tx: %s\n", tx.Hash().Hex())
fmt.Printf("Transfer at Nonce: %v\n", tx.Nonce())

// Wait for transfer transaction
_, err = client.WaitForEVMTransaction(tx.Hash())
if err != nil {
    fmt.Printf("Error waiting for transfer: %v\n", err)
    return
}

fmt.Printf("NFT transferred to: %s\n", recipientEVM)
```

### Verify NFT Ownership

```go
evmClient := evm.NewEVMClient(*acc)

currentOwner := evmClient.TokenOwner(contractAddress, tokenId)
fmt.Printf("Current owner: %s\n", currentOwner.Hex())
```

### Balance Operations

Send tokens on the Cosmos layer:

```go
import (
    "cosmossdk.io/math"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/thesixnetwork/lbb-sdk-go/pkg/balance"
)

balanceMsg, err := balance.NewBalanceMsg(*acc)
if err != nil {
    fmt.Printf("NewBalanceMsg error: %v\n", err)
    return
}

sendAmount := sdk.Coin{
    Amount: math.NewInt(1000000),
    Denom:  "usix",
}

res, err := balanceMsg.SendBalanceAndWait(
    "6x1recipient_address_here",
    sdk.NewCoins(sendAmount),
)
if err != nil {
    fmt.Printf("Send error: %v\n", err)
    return
}

fmt.Printf("Transfer tx: %s\n", res.TxHash)
```

### Complete Flow Example

Here's the typical flow from `example/main.go`:

```go
const (
    contractName   = "MyCertificate"
    contractSymbol = "CERT"
    schemaName     = "myorg.lbbv01"
    recipientEVM   = "0x8a28fb81A084Ac7A276800957a19a6054BF86E4D"
)

func main() {
    ctx := context.Background()
    
    // 1. Initialize client
    client, err := client.NewClient(ctx, false)
    if err != nil {
        panic(fmt.Sprintf("Failed to create client: %v", err))
    }
    
    // 2. Generate mnemonic and create account
    mnemonic, err := account.GenerateMnemonic()
    if err != nil {
        panic(fmt.Sprintf("Failed to generate mnemonic: %v", err))
    }
    
    acc, err := account.NewAccount(client, "alice", mnemonic, "")
    if err != nil {
        panic("Failed to create account")
    }
    
    // 3. Deploy schema and mint metadata
    meta, err := metadata.NewMetadataMsg(*acc, schemaName)
    if err != nil {
        panic(fmt.Sprintf("NewMetadataMsg error: %v", err))
    }
    
    msgDeploySchema, err := meta.BuildDeployMsg()
    if err != nil {
        panic(fmt.Sprintf("Failed to build deploy message: %v", err))
    }
    
    msgCreateMetadata, err := meta.BuildMintMetadataMsg("1")
    if err != nil {
        panic(fmt.Sprintf("Failed to build metadata: %v", err))
    }
    
    var msgs []sdk.Msg
    msgs = append(msgs, msgDeploySchema, msgCreateMetadata)
    
    res, err := meta.BroadcastTxAndWait(msgs...)
    if err != nil {
        panic(fmt.Sprintf("Broadcast error: %v", err))
    }
    
    // 4. Deploy EVM contract
    evmClient := evm.NewEVMClient(*acc)
    contractAddress, tx, err := evmClient.DeployCertificateContract(
        contractName,
        contractSymbol,
        schemaName,
    )
    if err != nil {
        panic(fmt.Sprintf("Deploy error: %v", err))
    }
    
    _, err = client.WaitForEVMTransaction(tx.Hash())
    if err != nil {
        panic(fmt.Sprintf("Wait error: %v", err))
    }
    
    // 5. Mint NFT
    tokenId := uint64(1)
    tx, err = evmClient.MintCertificateNFT(contractAddress, tokenId)
    if err != nil {
        panic(fmt.Sprintf("Mint error: %v", err))
    }
    
    _, err = client.WaitForEVMTransaction(tx.Hash())
    if err != nil {
        panic(fmt.Sprintf("Wait error: %v", err))
    }
    
    // 6. Freeze/Unfreeze certificate
    res, err = meta.FreezeCertificate("1")
    if err != nil {
        panic(fmt.Sprintf("Freeze error: %v", err))
    }
    
    err = client.WaitForTransaction(res.TxHash)
    if err != nil {
        panic(fmt.Sprintf("Wait error: %v", err))
    }
    
    res, err = meta.UnfreezeCertificate("1")
    if err != nil {
        panic(fmt.Sprintf("Unfreeze error: %v", err))
    }
    
    // 7. Transfer NFT
    tx, err = evmClient.TransferCertificateNFT(
        contractAddress,
        common.HexToAddress(recipientEVM),
        tokenId,
    )
    if err != nil {
        panic(fmt.Sprintf("Transfer error: %v", err))
    }
    
    _, err = client.WaitForEVMTransaction(tx.Hash())
    if err != nil {
        panic(fmt.Sprintf("Wait error: %v", err))
    }
    
    // 8. Verify ownership
    currentOwner := evmClient.TokenOwner(contractAddress, tokenId)
    fmt.Printf("Current owner: %s\n", currentOwner.Hex())
}
```

## Key Differences Between Layers

### Cosmos Layer Operations
- Use `meta.BroadcastTx()` or `meta.BroadcastTxAndWait()`
- Wait with `client.WaitForTransaction(res.TxHash)`
- Handle multiple messages in one transaction
- Operations: Deploy schema, mint metadata, freeze/unfreeze

### EVM Layer Operations
- Direct function calls return `*types.Transaction`
- Wait with `client.WaitForEVMTransaction(tx.Hash())`
- Each transaction is separate
- Operations: Deploy contract, mint NFT, transfer NFT

> **Important:** Cosmos and EVM layers operate independently. You can execute operations on both layers simultaneously without nonce/sequence conflicts.

## Error Handling Best Practices

Always check errors and wait for confirmation:

```go
// For Cosmos transactions
res, err := meta.BroadcastTx(msg)
if err != nil {
    fmt.Printf("Broadcast error: %v\n", err)
    return
}

err = client.WaitForTransaction(res.TxHash)
if err != nil {
    fmt.Printf("Transaction failed: %v\n", err)
    return
}

// For EVM transactions
tx, err := evmClient.MintCertificateNFT(contractAddress, tokenId)
if err != nil {
    fmt.Printf("Mint error: %v\n", err)
    return
}

_, err = client.WaitForEVMTransaction(tx.Hash())
if err != nil {
    fmt.Printf("Transaction failed: %v\n", err)
    return
}
```

## Modifying the Examples

### Change Network

Edit the client initialization:
```go
// false = testnet (fivenet), true = mainnet (sixnet)
client, err := client.NewClient(context.Background(), false)
```

### Use Your Own Mnemonic

Replace the test mnemonic with your own:
```go
const myMnemonic = "your twelve word mnemonic phrase..."
acc, err := account.NewAccount(client, "my-wallet", myMnemonic, "")
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

### Invalid nonce errors
- Ensure you wait for each EVM transaction to complete
- Don't send multiple EVM transactions without waiting
- Use `client.WaitForEVMTransaction()` between transactions

## Next Steps

1. Review the [main USAGE.md](../USAGE.md) for detailed documentation
2. Explore the SDK source code to understand available methods
3. Build your own application using these examples as a template
4. Check the test files for additional usage patterns

## Reference Files

- **Latest Example:** [`cmd/main.go`](../cmd/main.go) - Most up-to-date patterns
- **This Example:** [`example/main.go`](main.go) - User-friendly walkthrough
- **Documentation:** [`../USAGE.md`](../USAGE.md) - Comprehensive API reference

---

**Always refer to [`cmd/main.go`](../cmd/main.go) for the latest, most robust example and error handling patterns.**