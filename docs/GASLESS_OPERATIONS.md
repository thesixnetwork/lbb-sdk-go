# Gasless Operations Guide

This guide explains how to implement gasless NFT operations using EIP-2612 permit signatures, allowing users to interact with blockchain without paying gas fees.

## Table of Contents

- [Overview](#overview)
- [How It Works](#how-it-works)
- [EIP-2612 Permit Standard](#eip-2612-permit-standard)
- [Use Cases](#use-cases)
- [Gasless Transfer](#gasless-transfer)
- [Gasless Burn](#gasless-burn)
- [Security Considerations](#security-considerations)
- [Implementation Guide](#implementation-guide)
- [FAQ](#faq)

## Overview

Gasless operations enable users to perform blockchain transactions without having tokens for gas fees. This is achieved through **EIP-2612 permit signatures**, where:

1. **User** signs a message offline (free, no gas)
2. **Relayer/Admin** broadcasts the transaction (pays all gas)
3. **Result**: User's action is executed, but they paid nothing!

### Benefits

- 🎉 **No gas fees for users** - Remove the biggest barrier to blockchain adoption
- 🚀 **Better UX** - Users don't need tokens to interact with your dApp
- 💰 **Platform control** - You choose who pays for what
- 🔒 **Secure** - User maintains full control through signatures
- ♻️ **Eco-friendly** - No need for users to acquire and hold gas tokens

## How It Works

### Traditional Flow (User Pays Gas)

```
User → Signs Transaction → Pays Gas → Blockchain
```

**Problem**: User needs tokens for gas fees

### Gasless Flow (Admin Pays Gas)

```
User → Signs Permit Offline (FREE) → Admin Broadcasts → Pays Gas → Blockchain
```

**Solution**: User signs for free, admin handles gas fees

### Technical Flow

1. **User Signs EIP-712 Message** (Offline, No Gas)
   - Message includes: tokenId, spender, deadline
   - Produces signature: (v, r, s)
   - Completely offline, no blockchain interaction

2. **Signature Sent to Relayer** (Off-chain)
   - Via API, WebSocket, or other communication
   - No gas cost, just data transfer

3. **Relayer Executes Transaction** (On-chain)
   - Calls `transferWithPermit()` or `burnWithPermit()`
   - Includes user's signature
   - Relayer pays all gas fees

4. **Smart Contract Validates** (On-chain)
   - Recovers signer from signature
   - Verifies signer owns the NFT
   - Verifies deadline hasn't passed
   - Executes the action

## EIP-2612 Permit Standard

EIP-2612 extends ERC-20 and ERC-721 with permit functionality, enabling gasless approvals and operations.

### Permit Structure

```go
type Permit struct {
    owner    address    // NFT owner
    spender  address    // Who can execute the action
    tokenId  uint256    // Which token
    deadline uint256    // Unix timestamp
    v        uint8      // Signature component
    r        bytes32    // Signature component
    s        bytes32    // Signature component
}
```

### EIP-712 Domain

```go
{
    name: "MyNFTCert",                           // Contract name
    version: "1",                                // Version
    chainId: 26,                                 // Chain ID
    verifyingContract: "0x..."                   // Contract address
}
```

### Typed Data Structure

```solidity
struct Permit {
    address owner;
    address spender;
    uint256 tokenId;
    uint256 nonce;
    uint256 deadline;
}
```

## Use Cases

### 1. Onboarding New Users

**Problem**: New users don't have gas tokens
**Solution**: Let them sign permits, platform pays gas

```go
// New user with ZERO balance can still transfer NFTs!
newUser, _ := account.NewAccount(client, "newuser", mnemonic, "")
// Balance: 0 (no tokens!)

// User signs permit (completely free)
permit, _ := userClient.SignPermit(...)

// Platform broadcasts and pays gas
admin.TransferWithPermit(..., permit)
```

### 2. Certificate Revocation

**Problem**: Users won't pay to revoke expired certificates
**Solution**: Gasless burn where platform pays

```go
// User signs burn permit (free)
burnPermit, _ := userClient.SignPermit(...)

// Platform burns expired certificate (pays gas)
admin.BurnWithPermit(contractAddr, owner, tokenId, burnPermit)
```

### 3. Bulk Operations

**Problem**: Expensive to perform many operations
**Solution**: Collect signatures, batch execute

```go
// Collect permits from multiple users (all free)
permits := []Permit{}
for _, user := range users {
    permit, _ := user.SignPermit(...)
    permits = append(permits, permit)
}

// Admin executes all at once (pays all gas)
for _, permit := range permits {
    admin.TransferWithPermit(..., permit)
}
```

### 4. Marketplace Without Gas

**Problem**: Users need gas to list/buy NFTs
**Solution**: Platform handles all transactions

```go
// Seller signs listing (free)
listingPermit, _ := seller.SignPermit(...)

// Buyer signs purchase (free)  
purchasePermit, _ := buyer.SignPermit(...)

// Marketplace executes trade (pays all gas)
marketplace.ExecuteTrade(listingPermit, purchasePermit)
```

## Gasless Transfer

Transfer NFT ownership without the owner paying gas fees.

### Example Code

```go
package main

import (
    "math/big"
    "time"
    "github.com/thesixnetwork/lbb-sdk-go/account"
    "github.com/thesixnetwork/lbb-sdk-go/client"
    "github.com/thesixnetwork/lbb-sdk-go/pkg/evm"
)

func gaslessTransfer() {
    // 1. Setup accounts
    adminAcc, _ := account.NewAccount(client, "admin", adminMnemonic, "")
    userAcc, _ := account.NewAccount(client, "user", userMnemonic, "")
    
    // 2. User signs permit OFFLINE (NO GAS!)
    userClient := evm.NewEVMClient(*userAcc)
    deadline := big.NewInt(time.Now().Unix() + 3600) // 1 hour
    
    permitSig, err := userClient.SignPermit(
        "MyNFTCert",              // Contract name
        contractAddress,          // Contract address
        adminAcc.GetEVMAddress(), // Spender (who executes)
        big.NewInt(1),            // Token ID
        deadline,                 // Deadline
    )
    
    // 3. Admin broadcasts transfer (PAYS ALL GAS)
    adminClient := evm.NewEVMClient(*adminAcc)
    tx, err := adminClient.TransferWithPermit(
        contractAddress,
        userAcc.GetEVMAddress(),  // From
        recipientAddress,         // To
        big.NewInt(1),            // Token ID
        permitSig,                // User's signature
    )
    
    // 4. Wait for confirmation
    client.WaitForEVMTransaction(tx.Hash())
    
    fmt.Println("✓ Transfer complete! User paid ZERO gas!")
}
```

### Complete Example

See: `07_1_gasless_transfer.go`

```bash
go run example/07_1_gasless_transfer.go
```

## Gasless Burn

Permanently destroy an NFT without the owner paying gas fees.

### Example Code

```go
func gaslessBurn() {
    // 1. Setup accounts
    adminAcc, _ := account.NewAccount(client, "admin", adminMnemonic, "")
    userAcc, _ := account.NewAccount(client, "user", userMnemonic, "")
    
    // 2. User signs burn permit OFFLINE (NO GAS!)
    userClient := evm.NewEVMClient(*userAcc)
    deadline := big.NewInt(time.Now().Unix() + 3600)
    
    permitSig, err := userClient.SignPermit(
        "MyNFTCert",              // Contract name
        contractAddress,          // Contract address
        adminAcc.GetEVMAddress(), // Spender (who executes)
        big.NewInt(1),            // Token ID
        deadline,                 // Deadline
    )
    
    // 3. Admin broadcasts burn (PAYS ALL GAS)
    adminClient := evm.NewEVMClient(*adminAcc)
    tx, err := adminClient.BurnWithPermit(
        contractAddress,
        userAcc.GetEVMAddress(),  // Owner
        big.NewInt(1),            // Token ID
        permitSig,                // User's signature
    )
    
    // 4. Wait for confirmation
    client.WaitForEVMTransaction(tx.Hash())
    
    // 5. Verify burn (owner = zero address)
    owner := adminClient.TokenOwner(contractAddress, 1)
    zeroAddr := common.HexToAddress("0x0000000000000000000000000000000000000000")
    
    if owner == zeroAddr {
        fmt.Println("✓ Burn complete! User paid ZERO gas!")
    }
}
```

### Complete Example

See: `13_1_gasless_burn.go`

```bash
go run example/13_1_gasless_burn.go
```

## Security Considerations

### 1. Deadline Protection

Always set a deadline to prevent replay attacks:

```go
// Good: 1 hour deadline
deadline := big.NewInt(time.Now().Unix() + 3600)

// Bad: No deadline or too far in future
deadline := big.NewInt(999999999999) // DON'T DO THIS
```

### 2. Nonce Management

The smart contract tracks nonces to prevent signature reuse:

```solidity
mapping(address => uint256) public nonces;

function _useNonce(address owner) internal returns (uint256 current) {
    current = nonces[owner];
    nonces[owner]++;
}
```

### 3. Signature Validation

Never trust signatures blindly:

```go
// Contract validates:
// 1. Signature is valid (correct signer)
// 2. Signer owns the token
// 3. Deadline hasn't passed
// 4. Nonce is correct
```

### 4. Front-Running Protection

Use deadlines and nonces to prevent front-running:

```go
// Short deadline reduces front-running window
deadline := big.NewInt(time.Now().Unix() + 300) // 5 minutes
```

### 5. Permit Storage

**NEVER** store permits in public places:

```go
// ❌ DON'T: Store in database visible to others
db.Save(permit)

// ✅ DO: Send directly to trusted relayer
relayer.Submit(permit)
```

### 6. Access Control

Only trusted addresses should be spenders:

```go
// ✅ DO: Use trusted admin/relayer address
spender := trustedRelayerAddress

// ❌ DON'T: Let anyone be spender
spender := untrustedAddress
```

## Implementation Guide

### Step 1: Setup Accounts

```go
// Admin account (has gas tokens)
adminAcc, err := account.NewAccount(
    client, 
    "admin", 
    adminMnemonic, 
    "",
)

// User account (NO gas tokens needed!)
userAcc, err := account.NewAccount(
    client, 
    "user", 
    userMnemonic, 
    "",
)
```

### Step 2: User Signs Permit

```go
userClient := evm.NewEVMClient(*userAcc)

// Set deadline (e.g., 1 hour from now)
deadline := big.NewInt(time.Now().Unix() + 3600)

// Sign permit message (completely offline, no gas)
permitSig, err := userClient.SignPermit(
    contractName,             // e.g., "MyNFTCert"
    contractAddress,          // Deployed contract address
    adminAcc.GetEVMAddress(), // Who can execute
    big.NewInt(tokenId),      // Token ID
    deadline,                 // Unix timestamp
)
if err != nil {
    panic(fmt.Sprintf("Failed to sign permit: %v", err))
}

// permitSig contains: V, R, S signature components
```

### Step 3: Send Signature to Relayer

```go
// Option A: API call
response, err := http.Post(
    "https://relayer.example.com/api/submit",
    "application/json",
    permitSignatureJSON,
)

// Option B: WebSocket
ws.Send(permitSignature)

// Option C: Direct (same application)
executeGaslessOperation(permitSig)
```

### Step 4: Relayer Executes Transaction

```go
adminClient := evm.NewEVMClient(*adminAcc)

// For transfer
tx, err := adminClient.TransferWithPermit(
    contractAddress,
    fromAddress,              // Token owner
    toAddress,                // Recipient
    big.NewInt(tokenId),
    permitSig,                // User's signature
)

// For burn
tx, err := adminClient.BurnWithPermit(
    contractAddress,
    ownerAddress,             // Token owner
    big.NewInt(tokenId),
    permitSig,                // User's signature
)

// Wait for confirmation
receipt, err := client.WaitForEVMTransaction(tx.Hash())
```

### Step 5: Verify Operation

```go
// Verify transfer
newOwner := adminClient.TokenOwner(contractAddress, tokenId)
if newOwner == toAddress {
    fmt.Println("✓ Transfer successful!")
}

// Verify burn
owner := adminClient.TokenOwner(contractAddress, tokenId)
zeroAddress := common.HexToAddress("0x0000000000000000000000000000000000000000")
if owner == zeroAddress {
    fmt.Println("✓ Burn successful!")
}
```

## FAQ

### Q: Does the user need any tokens at all?

**A:** No! The user account can have zero balance. They only need to sign a message, which is free.

### Q: Who pays for the gas?

**A:** The relayer/admin who broadcasts the transaction pays all gas fees.

### Q: Can the permit be reused?

**A:** No. Each permit can only be used once due to the nonce system.

### Q: What if the deadline expires?

**A:** The transaction will fail. The user needs to sign a new permit with a new deadline.

### Q: Is this secure?

**A:** Yes, when implemented correctly. The signature proves the user authorized the action, and the smart contract validates everything.

### Q: Can the relayer steal my NFT?

**A:** No. The relayer can only execute what the user explicitly signed. They cannot change the recipient or token ID.

### Q: What happens if I sign a permit but don't send it?

**A:** Nothing. The permit is just a signature - it doesn't do anything until someone broadcasts a transaction using it.

### Q: Can I cancel a permit?

**A:** Not directly, but you can:
1. Transfer the NFT to another address (makes permit invalid)
2. Wait for the deadline to expire
3. Execute a different action that changes the nonce

### Q: What's the difference between permit and approve?

**A:** 
- **approve()**: Requires gas, gives ongoing permission
- **permit**: Free signature, one-time permission with deadline

### Q: Can I use this on any NFT contract?

**A:** Only contracts that implement EIP-2612 permit functionality. Standard ERC-721 contracts don't support this.

### Q: How do I build a relayer service?

**A:** 
1. Create API endpoint to receive permits
2. Validate permit structure and deadline
3. Queue permits for execution
4. Execute using admin account
5. Return transaction hash to user

### Q: What are the gas costs?

**A:** Similar to regular operations, but slightly higher due to signature verification:
- Regular transfer: ~50,000 gas
- Permit transfer: ~60,000 gas
- Regular burn: ~40,000 gas  
- Permit burn: ~50,000 gas

### Q: Can I batch multiple permits?

**A:** Yes! Collect multiple permits and execute them in sequence or use a batch transaction contract.

### Q: What if the user changes their mind?

**A:** If the permit hasn't been executed yet, you can:
1. Simply not execute it
2. Wait for deadline to expire
3. Transfer the NFT (invalidates permit)

## Examples

### Complete Examples

1. **Gasless Transfer**: `example/07_1_gasless_transfer.go`
2. **Gasless Burn**: `example/13_1_gasless_burn.go`

### Run Examples

```bash
# Update contract address in the example file first
cd example

# Gasless transfer
go run 07_1_gasless_transfer.go

# Gasless burn
go run 13_1_gasless_burn.go
```

### Compare with Standard Operations

```bash
# Standard transfer (user pays gas)
go run 07_0_transfer_nft.go

# Gasless transfer (admin pays gas)
go run 07_1_gasless_transfer.go

# Standard burn (user pays gas)
go run 13_0_burn_nft.go

# Gasless burn (admin pays gas)
go run 13_1_gasless_burn.go
```

## Best Practices

### 1. Set Reasonable Deadlines

```go
// ✅ Good: Short deadline for immediate actions
deadline := big.NewInt(time.Now().Unix() + 300) // 5 minutes

// ✅ Good: Longer deadline for async operations
deadline := big.NewInt(time.Now().Unix() + 3600) // 1 hour

// ❌ Bad: Too far in future
deadline := big.NewInt(time.Now().Unix() + 86400*365) // 1 year
```

### 2. Validate Before Execution

```go
// Check ownership before executing
currentOwner := client.TokenOwner(contractAddr, tokenId)
if currentOwner != permitSigner {
    return errors.New("signer doesn't own token")
}

// Check deadline
if time.Now().Unix() > deadline.Int64() {
    return errors.New("permit expired")
}
```

### 3. Handle Errors Gracefully

```go
tx, err := adminClient.TransferWithPermit(...)
if err != nil {
    // Log error
    log.Printf("Transfer failed: %v", err)
    
    // Notify user
    notifyUser(user, "Transaction failed: " + err.Error())
    
    // Don't panic in production!
    return err
}
```

### 4. Monitor Gas Costs

```go
receipt, _ := client.WaitForEVMTransaction(tx.Hash())

// Track gas usage
totalGas := receipt.GasUsed
gasPrice := tx.GasPrice()
totalCost := new(big.Int).Mul(totalGas, gasPrice)

// Alert if costs are too high
if totalCost.Cmp(maxAcceptableCost) > 0 {
    log.Printf("WARNING: High gas cost: %v", totalCost)
}
```

### 5. Implement Rate Limiting

```go
// Limit permits per user
if userPermitCount[user] > maxPermitsPerHour {
    return errors.New("rate limit exceeded")
}

// Limit total executions
if totalExecutions > maxExecutionsPerBlock {
    return errors.New("system capacity reached")
}
```

## Conclusion

Gasless operations using EIP-2612 permits enable a significantly better user experience by removing the gas fee barrier. Users can interact with your blockchain application without ever needing to acquire gas tokens, making blockchain technology more accessible to everyone.

### Key Takeaways

- ✅ Users sign messages offline (completely free)
- ✅ Relayer/admin broadcasts and pays gas
- ✅ Secure through cryptographic signatures
- ✅ Protected by deadlines and nonces
- ✅ Works for transfers, burns, and more

### Next Steps

1. Run the example scripts
2. Implement gasless operations in your dApp
3. Build a relayer service for your users
4. Test thoroughly on testnet
5. Monitor gas costs in production

For more information, see:
- [EIP-2612 Specification](https://eips.ethereum.org/EIPS/eip-2612)
- [EIP-712 Typed Data](https://eips.ethereum.org/EIPS/eip-712)
- [Example Code](./07_1_gasless_transfer.go)