# Gasless Operations Quick Reference

Quick reference for implementing gasless NFT operations using EIP-2612 permits.

## 🚀 Quick Start

### Gasless Transfer (3 Steps)

```go
// Step 1: User signs permit (FREE - no gas!)
userClient := evm.NewEVMClient(*userAcc)
permitSig, _ := userClient.SignPermit(
    "MyNFTCert",              // Contract name
    contractAddr,             // Contract address
    adminAddr,                // Spender
    big.NewInt(tokenId),      // Token ID
    big.NewInt(time.Now().Unix() + 3600), // Deadline (1 hour)
)

// Step 2: Admin broadcasts (pays all gas)
adminClient := evm.NewEVMClient(*adminAcc)
tx, _ := adminClient.TransferWithPermit(
    contractAddr,
    userAddr,                 // From
    recipientAddr,            // To
    big.NewInt(tokenId),
    permitSig,
)

// Step 3: Wait for confirmation
client.WaitForEVMTransaction(tx.Hash())
```

### Gasless Burn (3 Steps)

```go
// Step 1: User signs burn permit (FREE - no gas!)
userClient := evm.NewEVMClient(*userAcc)
permitSig, _ := userClient.SignPermit(
    "MyNFTCert",
    contractAddr,
    adminAddr,
    big.NewInt(tokenId),
    big.NewInt(time.Now().Unix() + 3600),
)

// Step 2: Admin broadcasts burn (pays all gas)
adminClient := evm.NewEVMClient(*adminAcc)
tx, _ := adminClient.BurnWithPermit(
    contractAddr,
    userAddr,                 // Owner
    big.NewInt(tokenId),
    permitSig,
)

// Step 3: Verify burn
owner := adminClient.TokenOwner(contractAddr, tokenId)
zeroAddr := common.HexToAddress("0x0000000000000000000000000000000000000000")
// owner == zeroAddr means successfully burned
```

## 📝 Function Signatures

### SignPermit

```go
func (e *EVMClient) SignPermit(
    contractName string,
    contractAddress common.Address,
    spender common.Address,
    tokenID *big.Int,
    deadline *big.Int,
) (*PermitSignature, error)
```

**Returns:** Signature with V, R, S components

### TransferWithPermit

```go
func (e *EVMClient) TransferWithPermit(
    contractAddress common.Address,
    from common.Address,
    to common.Address,
    tokenID *big.Int,
    signature *PermitSignature,
) (*types.Transaction, error)
```

**Returns:** Transaction object

### BurnWithPermit

```go
func (e *EVMClient) BurnWithPermit(
    contractAddress common.Address,
    from common.Address,
    tokenID *big.Int,
    signature *PermitSignature,
) (*types.Transaction, error)
```

**Returns:** Transaction object

## 🔑 Key Concepts

| Concept | Description |
|---------|-------------|
| **Permit** | Off-chain signature authorizing an action |
| **Spender** | Address allowed to execute the action (admin/relayer) |
| **Deadline** | Unix timestamp when permit expires |
| **Nonce** | Counter preventing signature reuse |
| **V, R, S** | ECDSA signature components |

## ⚡ Common Patterns

### Pattern 1: Set Deadline

```go
// 5 minutes from now
deadline := big.NewInt(time.Now().Unix() + 300)

// 1 hour from now
deadline := big.NewInt(time.Now().Unix() + 3600)

// 1 day from now
deadline := big.NewInt(time.Now().Unix() + 86400)
```

### Pattern 2: Create Accounts

```go
// Admin account (has funds)
adminAcc, _ := account.NewAccount(client, "admin", adminMnemonic, "")

// User account (NO funds needed!)
userAcc, _ := account.NewAccount(client, "user", userMnemonic, "")
```

### Pattern 3: Verify Ownership

```go
// Before signing permit
owner := evmClient.TokenOwner(contractAddr, tokenId)
if owner != userAddr {
    panic("User doesn't own this token")
}
```

### Pattern 4: Verify Burn

```go
// After burning
owner := evmClient.TokenOwner(contractAddr, tokenId)
zeroAddr := common.HexToAddress("0x0000000000000000000000000000000000000000")
if owner == zeroAddr {
    fmt.Println("✓ Successfully burned")
}
```

### Pattern 5: Error Handling

```go
permitSig, err := userClient.SignPermit(...)
if err != nil {
    return fmt.Errorf("failed to sign permit: %w", err)
}

tx, err := adminClient.TransferWithPermit(...)
if err != nil {
    return fmt.Errorf("failed to execute transfer: %w", err)
}

receipt, err := client.WaitForEVMTransaction(tx.Hash())
if err != nil {
    return fmt.Errorf("transaction failed: %w", err)
}
```

## ✅ Checklist

Before implementing gasless operations:

- [ ] Contract supports EIP-2612 permits
- [ ] Admin account has sufficient gas tokens
- [ ] Contract name matches deployment
- [ ] Contract address is correct
- [ ] Deadline is reasonable (not too long)
- [ ] User owns the NFT being operated on
- [ ] Error handling is implemented
- [ ] Transaction confirmation is awaited

## 🎯 Use Cases

### Onboarding
```go
// New user with ZERO balance
newUser, _ := account.NewAccount(client, "newuser", mnemonic, "")
// Balance: 0 - Can still transfer NFTs! 🎉
```

### Bulk Operations
```go
// Collect permits from multiple users (all free)
for _, user := range users {
    permit, _ := user.SignPermit(...)
    permits = append(permits, permit)
}
// Admin executes all (pays all gas)
```

### Certificate Revocation
```go
// User signs burn permit (free)
// Platform burns expired certificate (pays gas)
```

## 🔒 Security

### ✅ DO

- Set reasonable deadlines (minutes to hours)
- Validate signatures before execution
- Check ownership before signing
- Use trusted admin addresses
- Monitor gas costs
- Implement rate limiting

### ❌ DON'T

- Use very long deadlines (years)
- Store permits publicly
- Skip ownership checks
- Trust unvalidated signatures
- Ignore deadline expiration
- Reuse old permits (nonces prevent this anyway)

## 🐛 Troubleshooting

| Error | Solution |
|-------|----------|
| "Invalid signature" | Verify contract name, address, and chain ID |
| "Deadline expired" | Sign new permit with future deadline |
| "Not token owner" | Verify user owns the token |
| "Nonce mismatch" | Signature already used or nonce changed |
| "Transaction reverted" | Check contract state and permit validity |

## 📊 Gas Comparison

| Operation | Standard | Gasless | User Pays |
|-----------|----------|---------|-----------|
| Transfer  | ~50k gas | ~60k gas | 0 gas ✅ |
| Burn      | ~40k gas | ~50k gas | 0 gas ✅ |

**Note:** Gasless operations cost ~20% more gas but user pays nothing!

## 📚 Examples

### Complete Examples
- `07_1_gasless_transfer.go` - Full gasless transfer example
- `13_1_gasless_burn.go` - Full gasless burn example
- `GASLESS_OPERATIONS.md` - Detailed guide

### Run Examples
```bash
cd example

# Gasless transfer
go run 07_1_gasless_transfer.go

# Gasless burn
go run 13_1_gasless_burn.go
```

## 🔗 Standards

- **EIP-2612**: Permit extension for ERC-20/721
- **EIP-712**: Typed structured data hashing and signing
- **EIP-155**: Simple replay attack protection

## 💡 Tips

1. **Test on testnet first** - Always!
2. **Use short deadlines** - Reduce risk window
3. **Validate before executing** - Check ownership and deadline
4. **Monitor gas prices** - Adjust relayer strategy
5. **Log everything** - Track permits and executions
6. **Handle errors gracefully** - Don't panic in production
7. **Rate limit** - Prevent abuse

## 🎓 Learning Path

1. ✅ Run standard transfer (`07_0_transfer_nft.go`)
2. ✅ Run gasless transfer (`07_1_gasless_transfer.go`)
3. ✅ Compare the differences
4. ✅ Run standard burn (`13_0_burn_nft.go`)
5. ✅ Run gasless burn (`13_1_gasless_burn.go`)
6. ✅ Read full guide (`GASLESS_OPERATIONS.md`)
7. ✅ Implement in your project

## 🚀 Next Steps

- Build a relayer service
- Implement in your dApp
- Create meta-transaction system
- Add gasless marketplace
- Enable gasless certificate management

---

**For detailed explanations, see:** [GASLESS_OPERATIONS.md](./GASLESS_OPERATIONS.md)

**For working examples, see:**
- [07_1_gasless_transfer.go](./07_1_gasless_transfer.go)
- [13_1_gasless_burn.go](./13_1_gasless_burn.go)