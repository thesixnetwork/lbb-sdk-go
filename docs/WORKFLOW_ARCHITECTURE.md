# LBB SDK Go - Workflow Architecture

## Overview

The LBB SDK Go operates on a dual-layer architecture, combining Cosmos SDK and EVM layers. This architecture allows for parallel execution of transactions on different layers while maintaining proper sequencing within each layer.

## Dual-Layer Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                        LBB SDK Go Architecture                      │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌───────────────────────────┐    ┌───────────────────────────┐     │
│  │     Cosmos Layer          │    │       EVM Layer           │     │
│  │   (Gen2 Data Layer)       │    │   (NFT Token Layer)       │     │
│  ├───────────────────────────┤    ├───────────────────────────┤     │
│  │                           │    │                           │     │
│  │  • Schema Deployment      │    │  • Contract Deployment    │     │
│  │  • Metadata Management    │    │  • NFT Minting            │     │
│  │  • Certificate Freezing   │    │  • NFT Transfers          │     │
│  │  • Data Mutations         │    │  • Ownership Tracking     │     │
│  │  • Executor Permissions   │    │  • ERC-721 Standard       │     │
│  │                           │    │                           │     │
│  │  Token: usix              │    │  Token: asix              │     │
│  │  Address: 6x...           │    │  Address: 0x...           │     │
│  │                           │    │                           │     │
│  └───────────────────────────┘    └───────────────────────────┘     │
│              ↓                              ↓                       │
│         Independent                    Independent                  │
│         Transaction                    Transaction                  │
│         Sequencing                     Sequencing                   │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

## Workflow Diagram

### Complete Certificate Lifecycle

```
                    Start: Generate Wallet
                              ↓
                      Create Account
                    (Cosmos: 6x... + EVM: 0x...)
                              ↓
            ┌─────────────────┴─────────────────┬─────────────────┐
            ↓                                   ↓                 ↓
    [COSMOS LAYER]                      [EVM LAYER]      [QUERY OPERATIONS]
            ↓                                   ↓                 ↓
  ┌─────────────────────┐           ┌─────────────────────┐  ┌──────────────┐
  │ 03. Deploy Schema   │           │ 05. Deploy Contract │  │ 10. Balance  │
  │   (myorg.cert01)    │ ← LINK →  │   (0xABC...)        │  │   Operations │
  └─────────────────────┘           └─────────────────────┘  └──────────────┘
            ↓                                   ↓                 ↓
  ┌─────────────────────┐           ┌─────────────────────┐  • Get Balance
  │ 04. Mint Metadata   │    CAN    │ 06. Mint NFT        │  • Transfer
  │   (Token ID: "1")   │ ← RUN  →  │   (Token ID: 1)     │  • Query All
  └─────────────────────┘  PARALLEL └─────────────────────┘
            ↓                                   ↓
  ┌─────────────────────┐           ┌─────────────────────┐
  │ 08. Freeze/Unfreeze │    CAN    │ 07. Transfer NFT    │
  │   (Lock Metadata)   │ ← RUN  →  │   (Change Owner)    │
  └─────────────────────┘  PARALLEL └─────────────────────┘
            ↓                                   ↓
  ┌─────────────────────┐           ┌─────────────────────┐
  │ 11. Query Metadata  │           │ 09. Query NFT       │
  │   • Schema Info     │           │   • Token Owner     │
  │   • Cert Data       │           │   • Ownership       │
  │   • Executors       │           └─────────────────────┘
  └─────────────────────┘                     ↓
                                    ┌─────────────────────┐
                                    │ 12. Query EVM       │
                                    │   • Gas Price       │
                                    │   • Chain ID        │
                                    │   • Nonce           │
                                    └─────────────────────┘
```

## Parallel Execution Capabilities

### ✅ Can Execute Simultaneously (Different Layers)

These operations can run in parallel without conflicts:

| Cosmos Layer Operation | EVM Layer Operation | Result |
|------------------------|---------------------|--------|
| Deploy Schema | Deploy NFT Contract | ✅ Both execute independently |
| Mint Metadata | Mint NFT | ✅ No nonce/sequence conflict |
| Freeze Metadata | Transfer NFT | ✅ Separate transaction pools |
| Update Metadata | Mint NFT | ✅ Independent operations |
| Query Schema | Query NFT Owner | ✅ Read-only, no conflicts |

**Example:**
```go
// These can be executed simultaneously
go func() {
    // Cosmos transaction
    meta.CreateCertificateMetadata("1")
}()

go func() {
    // EVM transaction
    evmClient.MintCertificateNFT(contractAddr, 1)
}()
```

### ⚠️ Must Execute Sequentially (Same Layer)

These operations must wait for confirmation:

#### Cosmos Layer (Sequential)
```
Deploy Schema → Wait → Mint Metadata → Wait → Freeze Metadata
```

**Why?** Cosmos uses sequence numbers. Each transaction must confirm before the next.

#### EVM Layer (Sequential)
```
Deploy Contract → Wait → Mint NFT → Wait → Transfer NFT
```

**Why?** EVM uses nonce. Each transaction increments the nonce sequentially.

### 📝 Batch Operations (Same Transaction)

#### Cosmos Layer Batching
You can batch multiple messages in ONE Cosmos transaction:

```go
meta := metadata.NewMetadataMsg(*acc, schemaName)

// Create multiple messages
msg1, _ := meta.BuildDeployMsg()
msg2, _ := meta.BuildMintMetadataMsg("1")
msg3, _ := meta.BuildMintMetadataMsg("2")

// Broadcast all in one transaction
res, err := meta.BroadcastTxAndWait(msg1, msg2, msg3)
```

**Advantage:** Single transaction fee, atomic execution (all succeed or all fail)

#### EVM Layer (No Batching)
EVM transactions must be sent individually, each with incrementing nonce.

## Transaction Flow Examples

### Example 1: Sequential Deployment (Safe)

```go
// Step 1: Deploy Schema (Cosmos)
msg := meta.BuildDeployMsg()
res, err := meta.BroadcastTxAndWait(msg)

// Step 2: Deploy Contract (EVM) - AFTER schema confirmed
contractAddr, tx, err := evmClient.DeployCertificateContract(name, symbol, schemaName)
client.WaitForEVMTransaction(tx.Hash())
```

**Use Case:** When contract deployment needs confirmed schema

### Example 2: Parallel Deployment (Faster)

```go
// Both can start simultaneously
var schemaErr, contractErr error
var wg sync.WaitGroup

wg.Add(2)

// Deploy Schema (Cosmos)
go func() {
    defer wg.Done()
    msg, _ := meta.BuildDeployMsg()
    _, schemaErr = meta.BroadcastTxAndWait(msg)
}()

// Deploy Contract (EVM)
go func() {
    defer wg.Done()
    _, tx, err := evmClient.DeployCertificateContract(name, symbol, schemaName)
    if err == nil {
        _, contractErr = client.WaitForEVMTransaction(tx.Hash())
    }
}()

wg.Wait()
```

**Use Case:** When operations are independent (saves time)

### Example 3: Minting Certificate (Parallel)

```go
var metadataErr, nftErr error
var wg sync.WaitGroup

wg.Add(2)

// Mint Metadata (Cosmos)
go func() {
    defer wg.Done()
    msg, _ := meta.BuildMintMetadataMsg("1")
    _, metadataErr = meta.BroadcastTxAndWait(msg)
}()

// Mint NFT (EVM)
go func() {
    defer wg.Done()
    tx, err := evmClient.MintCertificateNFT(contractAddr, 1)
    if err == nil {
        _, nftErr = client.WaitForEVMTransaction(tx.Hash())
    }
}()

wg.Wait()
```

**Use Case:** Fastest certificate issuance (both layers simultaneously)

## Best Practices

### 1. Understanding Layer Independence

```
Cosmos Layer: Manages DATA (schema, metadata, attributes)
EVM Layer:    Manages OWNERSHIP (NFT tokens, transfers)

Both layers are synchronized by Token ID but operate independently.
```

### 2. When to Use Sequential Execution

- When one operation depends on another's result
- When deploying for the first time (safer)
- When debugging or testing
- When gas/sequence management is critical

### 3. When to Use Parallel Execution

- Minting metadata + NFT for same token ID
- Deploying schema + contract (if schema name is predetermined)
- Freezing metadata while transferring NFT
- Any cross-layer operations that are independent

### 4. Transaction Confirmation

Always wait for confirmation before proceeding with dependent operations:

```go
// ✅ CORRECT
res, err := meta.BroadcastTx(msg1)
client.WaitForTransaction(res.TxHash)  // WAIT

res, err = meta.BroadcastTx(msg2)  // Now safe to send
```

```go
// ❌ INCORRECT
res1, _ := meta.BroadcastTx(msg1)
res2, _ := meta.BroadcastTx(msg2)  // May fail due to sequence number
```

### 5. Error Handling for Parallel Operations

```go
type Result struct {
    Layer string
    Error error
}

results := make(chan Result, 2)

// Cosmos operation
go func() {
    _, err := meta.BroadcastTxAndWait(cosmosMsg)
    results <- Result{"Cosmos", err}
}()

// EVM operation
go func() {
    tx, err := evmClient.MintCertificateNFT(addr, id)
    if err == nil {
        _, err = client.WaitForEVMTransaction(tx.Hash())
    }
    results <- Result{"EVM", err}
}()

// Collect results
for i := 0; i < 2; i++ {
    result := <-results
    if result.Error != nil {
        log.Printf("%s layer error: %v", result.Layer, result.Error)
    }
}
```

## Performance Comparison

### Sequential Execution
```
Deploy Schema (3s) → Deploy Contract (4s) = 7 seconds total
Mint Metadata (2s) → Mint NFT (3s) = 5 seconds total
                                      ─────────────
                                      12 seconds
```

### Parallel Execution
```
Deploy Schema (3s)  ┐
                    ├→ 4 seconds total
Deploy Contract (4s)┘

Mint Metadata (2s)  ┐
                    ├→ 3 seconds total
Mint NFT (3s)       ┘
                      ─────────────
                      7 seconds (42% faster!)
```

## Common Patterns

### Pattern 1: Certificate Issuance (Fast)
```go
// Parallel minting for speed
var wg sync.WaitGroup
wg.Add(2)

go func() {
    defer wg.Done()
    meta.CreateCertificateMetadata(tokenId)
}()

go func() {
    defer wg.Done()
    evmClient.MintCertificateNFT(contractAddr, tokenId)
}()

wg.Wait()
```

### Pattern 2: Certificate Issuance (Safe)
```go
// Sequential for reliability
res, err := meta.CreateCertificateMetadata(tokenId)
if err != nil {
    return err
}
client.WaitForTransaction(res.TxHash)

tx, err := evmClient.MintCertificateNFT(contractAddr, tokenId)
if err != nil {
    return err
}
client.WaitForEVMTransaction(tx.Hash())
```

### Pattern 3: Bulk Operations
```go
// Batch on Cosmos, iterate on EVM
var cosmosMessages []sdk.Msg
for i := 1; i <= 10; i++ {
    msg, _ := meta.BuildMintMetadataMsg(fmt.Sprintf("%d", i))
    cosmosMessages = append(cosmosMessages, msg)
}

// Single Cosmos transaction for all metadata
meta.BroadcastTxAndWait(cosmosMessages...)

// Individual EVM transactions for each NFT
for i := 1; i <= 10; i++ {
    tx, _ := evmClient.MintCertificateNFT(contractAddr, uint64(i))
    client.WaitForEVMTransaction(tx.Hash())
}
```

## Summary

### Key Takeaways

1. **Different Layers = Can Run Parallel**
   - Cosmos and EVM transactions are independent
   - No nonce/sequence conflicts across layers
   - Significant performance improvement

2. **Same Layer = Must Run Sequential**
   - Wait for confirmation before next transaction
   - Proper nonce/sequence management required
   - Prevents transaction failures

3. **Batching on Cosmos = Efficient**
   - Multiple messages in one transaction
   - Single fee, atomic execution
   - Recommended for bulk operations

4. **Choose Based on Use Case**
   - Parallel: Maximum speed, independent operations
   - Sequential: Maximum reliability, dependent operations
   - Batch: Efficient for multiple Cosmos operations

### Decision Matrix

| Scenario | Recommended Approach | Reason |
|----------|---------------------|---------|
| First-time deployment | Sequential | Safer, easier to debug |
| Production issuance | Parallel | Faster, proven stable |
| Bulk minting | Batch Cosmos + Sequential EVM | Most efficient |
| Testing | Sequential | Easier to track errors |
| High-volume production | Parallel with error handling | Maximum throughput |

---

**For more examples, see:**
- [Example Code](../example/) - 12+ working examples
- [Tutorial Guide](./TUTORIAL.md) - Step-by-step guide
- [Examples README](../example/README.md) - Detailed usage patterns
