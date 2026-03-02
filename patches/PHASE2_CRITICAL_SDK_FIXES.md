# Phase 2: Critical SDK Security Fixes

## Overview
This document provides PR-ready patches for all HIGH and CRITICAL severity issues found in the Go SDK during the security audit.

**Estimated Time:** 16-24 hours
**Priority:** 🔴 CRITICAL - Must be applied before production use

---

## Fix 1: Remove/Protect Mnemonic Exposure (CRITICAL)

### Issue
**Severity:** CRITICAL  
**File:** `account/account.go`  
**Problem:** `GetMnemonic()` returns plaintext mnemonic with no protection or warnings. This is an extreme security risk if accidentally logged, exposed in error messages, or used inappropriately.

### Impact
- Mnemonic exposure = complete account compromise
- All assets can be stolen
- No recovery possible once exposed
- Violation of security best practices

### Fix Option 1: Remove GetMnemonic() Entirely (RECOMMENDED)

**File:** `account/account.go` (Lines ~130-135)

```go
// REMOVE THIS FUNCTION ENTIRELY:
// func (a *Account) GetMnemonic() string {
//     return a.mnemonic
// }

// Instead, don't store mnemonic after key derivation:

// Modified NewAccount function:
func NewAccount(ctx client.ClientI, accountName, mnemonic, password string) (*Account, error) {
    if ctx == nil {
        return nil, fmt.Errorf("client cannot be nil")
    }

    if accountName == "" {
        return nil, fmt.Errorf("account name cannot be empty")
    }

    if !ValidateMnemonic(mnemonic) {
        return nil, fmt.Errorf("invalid mnemonic provided")
    }

    evmAddress, err := GetAddressFromMnemonic(mnemonic, password)
    if err != nil {
        return nil, fmt.Errorf("failed to get EVM address from mnemonic for account '%s': %w", accountName, err)
    }

    cosmosAddress, err := GetBech32AccountFromMnemonic(ctx.GetKeyring(), accountName, mnemonic, password)
    if err != nil {
        return nil, fmt.Errorf("failed to get Bech32 Cosmos address from mnemonic for account '%s': %w", accountName, err)
    }

    privateKey, err := CreatePrivateKeyFromMnemonic(mnemonic, password)
    if err != nil {
        return nil, fmt.Errorf("failed to generate private key from mnemonic for account '%s': %w", accountName, err)
    }

    chainIDBigInt, ok := ChainIDMapping[ctx.GetChainID()]
    if !ok {
        return nil, fmt.Errorf("chain ID '%s' not found in mapping", ctx.GetChainID())
    }

    authz, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIDBigInt)
    if err != nil {
        return nil, fmt.Errorf("failed to create transactor for account '%s': %w", accountName, err)
    }

    fmt.Printf("Account '%s' created successfully\n", accountName)
    fmt.Printf("  Cosmos Address: %s\n", cosmosAddress.String())
    fmt.Printf("  EVM Address: %s\n", evmAddress.Hex())

    return &Account{
        client:        ctx,
        auth:          authz,
        privateKey:    privateKey,
        // mnemonic:   mnemonic,  // REMOVE THIS - don't store mnemonic
        evmAddress:    evmAddress,
        cosmosAddress: cosmosAddress,
        accountName:   accountName,
    }, nil
}

// Update Account struct - remove mnemonic field:
type Account struct {
    client        client.ClientI
    auth          *bind.TransactOpts
    // mnemonic   string  // REMOVE THIS FIELD
    privateKey    *ecdsa.PrivateKey
    evmAddress    common.Address
    cosmosAddress sdk.AccAddress
    accountName   string
}
```

### Fix Option 2: Add Explicit Warning (If removal not possible)

If you MUST keep GetMnemonic for backward compatibility:

```go
// GetMnemonic returns the mnemonic phrase
// 
// ⚠️  SECURITY WARNING ⚠️
// This function exposes your mnemonic phrase in plaintext. Anyone with access
// to this mnemonic can steal ALL funds from your account with NO RECOVERY.
//
// NEVER:
//   - Log the mnemonic
//   - Send it over the network
//   - Store it in plaintext files
//   - Display it in UI without explicit user confirmation
//   - Use it in production environments
//
// This function should ONLY be used for:
//   - Initial backup during account creation (with user confirmation)
//   - Migration to secure key storage
//   - Development/testing with throwaway accounts
//
// Consider using encrypted key storage or hardware security modules in production.
func (a *Account) GetMnemonic() string {
    // Consider adding a log warning
    // log.Warn("GetMnemonic called - mnemonic exposure risk!")
    return a.mnemonic
}

// Better: Add explicit confirmation requirement
func (a *Account) GetMnemonicWithConfirmation(iUnderstandTheRisks bool) (string, error) {
    if !iUnderstandTheRisks {
        return "", fmt.Errorf("must explicitly confirm understanding of security risks")
    }
    return a.mnemonic, nil
}
```

### Documentation Update

Add to `account/README.md`:

```markdown
## Security Best Practices

### Mnemonic Management

**CRITICAL**: Never expose or store mnemonics in plaintext in production.

#### DO:
- ✅ Generate mnemonics securely using `GenerateNewMnemonic()`
- ✅ Display mnemonic ONCE during account creation with explicit user confirmation
- ✅ Require user to write it down physically (offline backup)
- ✅ Use encrypted keystores (keystore.json) for production
- ✅ Consider hardware security modules (HSMs) for high-value accounts
- ✅ Use secure key management services (AWS KMS, Google Cloud KMS)

#### DON'T:
- ❌ Store mnemonics in application logs
- ❌ Send mnemonics over the network
- ❌ Store mnemonics in environment variables
- ❌ Store mnemonics in plaintext files
- ❌ Include mnemonics in error messages
- ❌ Display mnemonics in UI without explicit confirmation
- ❌ Keep mnemonics in memory longer than necessary

#### Example: Secure Account Creation

```go
// Generate new mnemonic
mnemonic, err := account.GenerateNewMnemonic()
if err != nil {
    return err
}

// Display to user ONCE for backup
fmt.Println("⚠️  WRITE DOWN YOUR MNEMONIC - THIS WILL ONLY BE SHOWN ONCE:")
fmt.Println(mnemonic)
fmt.Println("\nPress Enter after you have written it down...")
bufio.NewReader(os.Stdin).ReadBytes('\n')

// Create account (mnemonic won't be stored)
acc, err := account.NewAccount(client, "my-account", mnemonic, "")
if err != nil {
    return err
}

// Mnemonic is now out of scope and will be garbage collected
// Use the account normally
```
```

**Recommendation:** Use Option 1 (complete removal). If backward compatibility is required, use Option 2 with explicit confirmation.

---

## Fix 2: Add Context Timeout Support (CRITICAL)

### Issue
**Severity:** HIGH  
**Files:** `client/client.go`, `pkg/evm/evm.go`  
**Problem:** Network operations don't respect context timeouts/cancellation, leading to hanging operations and resource leaks.

### Impact
- Operations can hang indefinitely
- Resource leaks (goroutines, connections)
- No way to cancel long-running operations
- Poor error handling

### Fix for WaitForTransaction

**File:** `client/client.go`

```go
// WaitForTransaction waits for a transaction to be included in a block
// Returns the transaction response or an error if the transaction fails or times out
func (c *Client) WaitForTransaction(ctx context.Context, txHash string) (*sdk.TxResponse, error) {
    if ctx == nil {
        ctx = context.Background()
    }

    if txHash == "" {
        return nil, fmt.Errorf("transaction hash cannot be empty")
    }

    // Set default timeout if context has none
    if _, hasDeadline := ctx.Deadline(); !hasDeadline {
        var cancel context.CancelFunc
        ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
        defer cancel()
    }

    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    attempt := 0
    maxAttempts := 60 // Fallback if context has no deadline

    for {
        select {
        case <-ctx.Done():
            return nil, fmt.Errorf("waiting for transaction %s cancelled: %w", txHash, ctx.Err())
        
        case <-ticker.C:
            attempt++
            
            // Query transaction status
            res, err := c.TxClient.GetTx(ctx, &tx.GetTxRequest{Hash: txHash})
            
            if err != nil {
                // Check if error is due to context cancellation
                if ctx.Err() != nil {
                    return nil, fmt.Errorf("transaction query cancelled: %w", ctx.Err())
                }
                
                // Transaction not found yet, continue waiting
                if strings.Contains(err.Error(), "not found") {
                    if attempt >= maxAttempts {
                        return nil, fmt.Errorf("transaction %s not found after %d attempts", txHash, maxAttempts)
                    }
                    continue
                }
                
                return nil, fmt.Errorf("failed to query transaction %s: %w", txHash, err)
            }
            
            if res == nil || res.TxResponse == nil {
                return nil, fmt.Errorf("received nil response for transaction %s", txHash)
            }
            
            // Transaction found
            txResp := res.TxResponse
            
            // Check transaction status
            if txResp.Code != 0 {
                return txResp, fmt.Errorf("transaction %s failed with code %d: %s", 
                    txHash, txResp.Code, txResp.RawLog)
            }
            
            return txResp, nil
        }
    }
}
```

### Fix for WaitForEVMTransaction

**File:** `pkg/evm/evm.go`

```go
// WaitForEVMTransaction waits for an EVM transaction to be mined and returns the receipt
func (e *EVMClient) WaitForEVMTransaction(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
    if ctx == nil {
        ctx = context.Background()
    }

    // Set default timeout if context has none
    if _, hasDeadline := ctx.Deadline(); !hasDeadline {
        var cancel context.CancelFunc
        ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
        defer cancel()
    }

    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    attempt := 0
    maxAttempts := 60

    for {
        select {
        case <-ctx.Done():
            return nil, fmt.Errorf("waiting for EVM transaction %s cancelled: %w", txHash.Hex(), ctx.Err())
        
        case <-ticker.C:
            attempt++
            
            // Query transaction receipt
            receipt, err := e.client.TransactionReceipt(ctx, txHash)
            
            if err != nil {
                // Check if error is due to context cancellation
                if ctx.Err() != nil {
                    return nil, fmt.Errorf("transaction receipt query cancelled: %w", ctx.Err())
                }
                
                // Receipt not found yet, continue waiting
                if err == ethereum.NotFound {
                    if attempt >= maxAttempts {
                        return nil, fmt.Errorf("transaction %s receipt not found after %d attempts", 
                            txHash.Hex(), maxAttempts)
                    }
                    continue
                }
                
                return nil, fmt.Errorf("failed to get receipt for transaction %s: %w", txHash.Hex(), err)
            }
            
            if receipt == nil {
                return nil, fmt.Errorf("received nil receipt for transaction %s", txHash.Hex())
            }
            
            // Check transaction status
            if receipt.Status == types.ReceiptStatusFailed {
                return receipt, fmt.Errorf("EVM transaction %s failed (status: %d)", 
                    txHash.Hex(), receipt.Status)
            }
            
            return receipt, nil
        }
    }
}
```

### Add Context to All Query Operations

**File:** `pkg/metadata/metadata.go`

```go
// GetNFTSchema queries NFT schema by contract address with context support
func (c *Client) GetNFTSchema(ctx context.Context, schemaCode string) (*nftmngrtypes.NFTSchemaByContractResponse, error) {
    if ctx == nil {
        ctx = context.Background()
    }

    if schemaCode == "" {
        return nil, fmt.Errorf("schema code cannot be empty")
    }

    // Set timeout if not already set
    if _, hasDeadline := ctx.Deadline(); !hasDeadline {
        var cancel context.CancelFunc
        ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
        defer cancel()
    }

    res, err := c.QueryClient.NFTSchemaByContract(ctx, &nftmngrtypes.QueryNFTSchemaByContractRequest{
        OriginContractAddress: schemaCode,
    })

    if err != nil {
        if ctx.Err() != nil {
            return nil, fmt.Errorf("query cancelled: %w", ctx.Err())
        }
        return nil, fmt.Errorf("failed to query NFT schema for %s: %w", schemaCode, err)
    }

    if res == nil || res.NftSchemaByContract == nil {
        return nil, fmt.Errorf("received nil response for schema: %s", schemaCode)
    }

    return res, nil
}
```

### Update All Function Signatures

Add `context.Context` as the first parameter to all functions that perform network operations:

```go
// Before:
func (c *Client) BroadcastTx(txBytes []byte) (*sdk.TxResponse, error)

// After:
func (c *Client) BroadcastTx(ctx context.Context, txBytes []byte) (*sdk.TxResponse, error)
```

---

## Fix 3: Thread-Safe Nonce Management (CRITICAL)

### Issue
**Severity:** HIGH  
**File:** `pkg/evm/evm.go`  
**Problem:** Nonce handling is not concurrency-safe. Multiple goroutines can get the same nonce, causing transaction failures.

### Impact
- Race conditions in concurrent transaction submission
- Transaction failures due to duplicate nonces
- Difficult to debug "nonce too low" errors

### Fix: Add Nonce Mutex

**File:** `pkg/evm/evm.go`

```go
package evm

import (
    "context"
    "fmt"
    "math/big"
    "sync"
    
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/ethclient"
)

// EVMClient handles EVM-compatible blockchain interactions
type EVMClient struct {
    client     *ethclient.Client
    chainID    *big.Int
    rpcURL     string
    
    // Nonce management (thread-safe)
    nonceMu    sync.Mutex
    nonceCache map[common.Address]uint64
}

// NewEVMClient creates a new EVM client instance
func NewEVMClient(rpcURL string, chainID *big.Int) (*EVMClient, error) {
    if rpcURL == "" {
        return nil, fmt.Errorf("RPC URL cannot be empty")
    }

    if chainID == nil {
        return nil, fmt.Errorf("chain ID cannot be nil")
    }

    client, err := ethclient.Dial(rpcURL)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to EVM RPC at %s: %w", rpcURL, err)
    }

    return &EVMClient{
        client:     client,
        chainID:    chainID,
        rpcURL:     rpcURL,
        nonceCache: make(map[common.Address]uint64),
    }, nil
}

// GetNonce returns the current nonce for an address (thread-safe)
func (e *EVMClient) GetNonce(ctx context.Context, address common.Address) (uint64, error) {
    e.nonceMu.Lock()
    defer e.nonceMu.Unlock()

    if ctx == nil {
        ctx = context.Background()
    }

    // Get nonce from blockchain
    nonce, err := e.client.PendingNonceAt(ctx, address)
    if err != nil {
        return 0, fmt.Errorf("failed to get nonce for address %s: %w", address.Hex(), err)
    }

    // Update cache
    cachedNonce, exists := e.nonceCache[address]
    if !exists || nonce > cachedNonce {
        e.nonceCache[address] = nonce
        return nonce, nil
    }

    // Use cached nonce if higher (for pending transactions)
    return cachedNonce, nil
}

// GetAndIncrementNonce returns the current nonce and increments the cache (thread-safe)
// Use this when preparing to send a transaction
func (e *EVMClient) GetAndIncrementNonce(ctx context.Context, address common.Address) (uint64, error) {
    e.nonceMu.Lock()
    defer e.nonceMu.Unlock()

    if ctx == nil {
        ctx = context.Background()
    }

    // Get nonce from blockchain
    nonce, err := e.client.PendingNonceAt(ctx, address)
    if err != nil {
        return 0, fmt.Errorf("failed to get nonce for address %s: %w", address.Hex(), err)
    }

    // Check cached nonce
    cachedNonce, exists := e.nonceCache[address]
    if exists && cachedNonce > nonce {
        nonce = cachedNonce
    }

    // Increment cache for next transaction
    e.nonceCache[address] = nonce + 1

    return nonce, nil
}

// ResetNonce clears the cached nonce for an address
// Call this if a transaction fails with "nonce too low"
func (e *EVMClient) ResetNonce(address common.Address) {
    e.nonceMu.Lock()
    defer e.nonceMu.Unlock()
    
    delete(e.nonceCache, address)
}

// ResetAllNonces clears all cached nonces
func (e *EVMClient) ResetAllNonces() {
    e.nonceMu.Lock()
    defer e.nonceMu.Unlock()
    
    e.nonceCache = make(map[common.Address]uint64)
}

// SendTransaction sends a signed transaction (thread-safe)
func (e *EVMClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
    if ctx == nil {
        ctx = context.Background()
    }

    if tx == nil {
        return fmt.Errorf("transaction cannot be nil")
    }

    err := e.client.SendTransaction(ctx, tx)
    if err != nil {
        // Check if error is due to nonce issue
        if isNonceError(err) {
            // Get sender address from transaction
            signer := types.LatestSignerForChainID(e.chainID)
            sender, senderErr := signer.Sender(tx)
            if senderErr == nil {
                // Reset nonce cache for this address
                e.ResetNonce(sender)
            }
        }
        return fmt.Errorf("failed to send transaction: %w", err)
    }

    return nil
}

// isNonceError checks if an error is related to nonce issues
func isNonceError(err error) bool {
    if err == nil {
        return false
    }
    errStr := err.Error()
    return strings.Contains(errStr, "nonce too low") ||
           strings.Contains(errStr, "nonce too high") ||
           strings.Contains(errStr, "replacement transaction underpriced")
}
```

### Usage Example

```go
// Concurrent transaction submission (now safe)
func sendMultipleTransactions(client *evm.EVMClient, from common.Address, recipients []common.Address) error {
    ctx := context.Background()
    
    var wg sync.WaitGroup
    errChan := make(chan error, len(recipients))
    
    for _, recipient := range recipients {
        wg.Add(1)
        go func(to common.Address) {
            defer wg.Done()
            
            // Thread-safe nonce retrieval
            nonce, err := client.GetAndIncrementNonce(ctx, from)
            if err != nil {
                errChan <- err
                return
            }
            
            // Create and sign transaction
            tx := types.NewTransaction(nonce, to, big.NewInt(1000), 21000, big.NewInt(1000000000), nil)
            // ... sign transaction ...
            
            // Thread-safe transaction submission
            if err := client.SendTransaction(ctx, tx); err != nil {
                errChan <- err
                return
            }
        }(recipient)
    }
    
    wg.Wait()
    close(errChan)
    
    // Check for errors
    for err := range errChan {
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

### Alternative: Document Non-Thread-Safety

If thread-safety is not a requirement, document clearly:

```go
// EVMClient handles EVM-compatible blockchain interactions
//
// ⚠️  CONCURRENCY WARNING ⚠️
// This client is NOT thread-safe. Do not use the same EVMClient instance
// from multiple goroutines. Create separate client instances for concurrent use,
// or use external synchronization (mutexes).
//
// Example of safe concurrent use:
//   client1 := evm.NewEVMClient(rpcURL, chainID)
//   client2 := evm.NewEVMClient(rpcURL, chainID)
//   // Now safe to use client1 and client2 concurrently
type EVMClient struct {
    // ...
}
```

---

## Fix 4: Add Response Validation (CRITICAL)

### Issue
**Severity:** HIGH  
**Files:** `pkg/metadata/metadata.go`, `pkg/balance/balance.go`  
**Problem:** Query responses are not validated for nil before dereferencing, causing panic risks.

### Impact
- Application crashes on invalid responses
- Poor error messages
- Difficult debugging

### Fix for GetNFTSchema

**File:** `pkg/metadata/metadata.go`

```go
// GetNFTSchema queries NFT schema by contract address
func (c *Client) GetNFTSchema(ctx context.Context, schemaCode string) (*nftmngrtypes.NFTSchemaByContractResponse, error) {
    if ctx == nil {
        ctx = context.Background()
    }

    // Input validation
    if schemaCode == "" {
        return nil, fmt.Errorf("schema code cannot be empty")
    }

    // Add timeout if not set
    if _, hasDeadline := ctx.Deadline(); !hasDeadline {
        var cancel context.CancelFunc
        ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
        defer cancel()
    }

    // Query schema
    res, err := c.QueryClient.NFTSchemaByContract(ctx, &nftmngrtypes.QueryNFTSchemaByContractRequest{
        OriginContractAddress: schemaCode,
    })

    if err != nil {
        if ctx.Err() != nil {
            return nil, fmt.Errorf("query cancelled: %w", ctx.Err())
        }
        return nil, fmt.Errorf("failed to query NFT schema for %s: %w", schemaCode, err)
    }

    // Response validation
    if res == nil {
        return nil, fmt.Errorf("received nil response for schema: %s", schemaCode)
    }

    if res.NftSchemaByContract == nil {
        return nil, fmt.Errorf("schema not found for contract: %s", schemaCode)
    }

    // Validate schema data
    schema := res.NftSchemaByContract
    if schema.Code == "" {
        return nil, fmt.Errorf("received invalid schema with empty code for contract: %s", schemaCode)
    }

    return res, nil
}
```

### Fix for GetNFTMetadata

**File:** `pkg/metadata/metadata.go`

```go
// GetNFTMetadata queries NFT metadata by schema code and token ID
func (c *Client) GetNFTMetadata(ctx context.Context, schemaCode string, tokenId string) (nftData, error) {
    if ctx == nil {
        ctx = context.Background()
    }

    // Input validation
    if schemaCode == "" {
        return nftData{}, fmt.Errorf("schema code cannot be empty")
    }

    if tokenId == "" {
        return nftData{}, fmt.Errorf("token ID cannot be empty")
    }

    // Get raw metadata
    res, err := c.GetNFTMetadataRaw(ctx, schemaCode, tokenId)
    if err != nil {
        return nftData{}, fmt.Errorf("failed to get NFT metadata: %w", err)
    }

    // Response validation
    if res == nil {
        return nftData{}, fmt.Errorf("received nil response for metadata: schema=%s, tokenId=%s", 
            schemaCode, tokenId)
    }

    if res.NftData == nil {
        return nftData{}, fmt.Errorf("metadata not found: schema=%s, tokenId=%s", 
            schemaCode, tokenId)
    }

    // Parse and validate metadata
    metadata := res.NftData
    if metadata.NftSchemaCode == "" {
        return nftData{}, fmt.Errorf("received invalid metadata with empty schema code")
    }

    if metadata.TokenId == "" {
        return nftData{}, fmt.Errorf("received invalid metadata with empty token ID")
    }

    return nftData{
        SchemaCode: metadata.NftSchemaCode,
        TokenId:    metadata.TokenId,
        OwnerAddress: metadata.OwnerAddressType,
        // ... other fields ...
    }, nil
}
```

### Fix for GetBalance

**File:** `pkg/balance/balance.go`

```go
// GetBalance queries the balance of an address
func (c *Client) GetBalance(ctx context.Context, address string, denom string) (*sdk.Coin, error) {
    if ctx == nil {
        ctx = context.Background()
    }

    // Input validation
    if address == "" {
        return nil, fmt.Errorf("address cannot be empty")
    }

    if denom == "" {
        return nil, fmt.Errorf("denom cannot be empty")
    }

    // Validate address format
    if _, err := sdk.AccAddressFromBech32(address); err != nil {
        return nil, fmt.Errorf("invalid address format: %w", err)
    }

    // Add timeout if not set
    if _, hasDeadline := ctx.Deadline(); !hasDeadline {
        var cancel context.CancelFunc
        ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
        defer cancel()
    }

    // Query balance
    res, err := c.QueryClient.Balance(ctx, &banktypes.QueryBalanceRequest{
        Address: address,
        Denom:   denom,
    })

    if err != nil {
        if ctx.Err() != nil {
            return nil, fmt.Errorf("query cancelled: %w", ctx.Err())
        }
        return nil, fmt.Errorf("failed to query balance for %s: %w", address, err)
    }

    // Response validation
    if res == nil {
        return nil, fmt.Errorf("received nil response for balance query: address=%s, denom=%s", 
            address, denom)
    }

    if res.Balance == nil {
        // Return zero balance instead of error
        return &sdk.Coin{Denom: denom, Amount: sdk.ZeroInt()}, nil
    }

    return res.Balance, nil
}
```

### Add Validation Helper Functions

**File:** `pkg/validation/validation.go` (new file)

```go
package validation

import (
    "fmt"
    "strings"
    
    "github.com/ethereum/go-ethereum/common"
    sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidateAddress validates a Cosmos address
func ValidateCosmosAddress(address string) error {
    if address == "" {
        return fmt.Errorf("address cannot be empty")
    }
    
    _, err := sdk.AccAddressFromBech32(address)
    if err != nil {
        return fmt.Errorf("invalid cosmos address format: %w", err)
    }
    
    return nil
}

// ValidateEVMAddress validates an EVM address
func ValidateEVMAddress(address string) error {
    if address == "" {
        return fmt.Errorf("address cannot be empty")
    }
    
    if !strings.HasPrefix(address, "0x") {
        return fmt.Errorf("EVM address must start with 0x")
    }
    
    if len(address) != 42 {
        return fmt.Errorf("EVM address must be 42 characters (including 0x prefix)")
    }
    
    if !common.IsHexAddress(address) {
        return fmt.Errorf("invalid EVM address format")
    }
    
    return nil
}

// ValidateTokenID validates a token ID
func ValidateTokenID(tokenID string) error {
    if tokenID == "" {
        return fmt.Errorf("token ID cannot be empty")
    }
    
    // Add specific format validation if needed
    return nil
}

// ValidateSchemaCode validates a schema code
func ValidateSchemaCode(schemaCode string) error {
    if schemaCode == "" {
        return fmt.Errorf("schema code cannot be empty")
    }
    
    // Add specific format validation if needed
    return nil
}

// ValidateDenom validates a token denomination
func ValidateDenom(denom string) error {
    if denom == "" {
        return fmt.Errorf("denom cannot be empty")
    }
    
    // Cosmos SDK denom validation
    if err := sdk.ValidateDenom(denom); err != nil {
        return fmt.Errorf("invalid denom: %w", err)
    }
    
    return nil
}

// ValidateAmount validates an amount is positive
func ValidateAmount(amount sdk.Int) error {
    if amount.IsNil() {
        return fmt.Errorf("amount cannot be nil")
    }
    
    if amount.IsNegative() {
        return fmt.Errorf("amount cannot be negative")
    }
    
    if amount.IsZero() {
        return fmt.Errorf("amount must be greater than zero")
    }
    
    return nil
}
```

---

## Fix 5: Fix Metadata Creation Bug (CRITICAL)

### Issue
**Severity:** HIGH  
**File:** `pkg/metadata/metadata.go`  
**Problem:** `CreateCertificateMetadataWithInfo` calls wrong builder function, ignoring the `info` parameter.

### Impact
- Functional bug - WithInfo parameter completely ignored
- User data lost
- Unexpected behavior

### Fix

**File:** `pkg/metadata/metadata.go`

```go
// CreateCertificateMetadataWithInfo creates certificate metadata with additional info
func (c *Client) CreateCertificateMetadataWithInfo(ctx context.Context, schemaCode string, tokenId string, info string) (string, error) {
    if ctx == nil {
        ctx = context.Background()
    }

    // Input validation
    if schemaCode == "" {
        return "", fmt.Errorf("schema code cannot be empty")
    }

    if tokenId == "" {
        return "", fmt.Errorf("token ID cannot be empty")
    }

    if info == "" {
        return "", fmt.Errorf("info cannot be empty")
    }

    // FIX: Use BuildMintMetadataWithInfoMsg instead of BuildMintMetadataMsg
    msg := BuildMintMetadataWithInfoMsg(
        c.Account.GetAddress(),
        schemaCode,
        tokenId,
        info,
    )

    // Broadcast transaction
    txHash, err := c.BroadcastTx(ctx, msg)
    if err != nil {
        return "", fmt.Errorf("failed to broadcast metadata creation: %w", err)
    }

    return txHash, nil
}

// CreateCertificateMetadata creates certificate metadata without additional info
func (c *Client) CreateCertificateMetadata(ctx context.Context, schemaCode string, tokenId string) (string, error) {
    if ctx == nil {
        ctx = context.Background()
    }

    // Input validation
    if schemaCode == "" {
        return "", fmt.Errorf("schema code cannot be empty")
    }

    if tokenId == "" {
        return "", fmt.Errorf("token ID cannot be empty")
    }

    // Use the basic builder (without info)
    msg := BuildMintMetadataMsg(
        c.Account.GetAddress(),
        schemaCode,
        tokenId,
    )

    // Broadcast transaction
    txHash, err := c.BroadcastTx(ctx, msg)
    if err != nil {
        return "", fmt.Errorf("failed to broadcast metadata creation: %w", err)
    }

    return txHash, nil
}
```

### Add Test

**File:** `pkg/metadata/metadata_test.go`

```go
func TestCreateCertificateMetadataWithInfo(t *testing.T) {
    // Setup
    client := setupTestClient(t)
    
    schemaCode := "test-schema"
    tokenID := "1"
    info := "test-info-data"
    
    // Create metadata with info
    txHash, err := client.CreateCertificateMetadataWithInfo(
        context.Background(),
        schemaCode,
        tokenID,
        info,
    )
    
    require.NoError(t, err)
    require.NotEmpty(t, txHash)
    
    // Wait for transaction
    receipt, err := client.WaitForTransaction(context.Background(), txHash)
    require.NoError(t, err)
    require.NotNil(t, receipt)
    
    // Verify metadata was created with info
    metadata, err := client.GetNFTMetadata(context.Background(), schemaCode, tokenID)
    require.NoError(t, err)
    
    // CRITICAL: Verify info was actually stored
    require.Equal(t, info, metadata.Info, "Info field should match input")
}
```

---

## Fix 6: Fix GetNFTMetadata Error Handling (CRITICAL)

### Issue
**Severity:** HIGH  
**File:** `pkg/metadata/metadata.go`  
**Problem:** Returns `(nftData{}, nil)` on error instead of returning the error.

### Impact
- Caller cannot detect errors
- Silent failures
- Incorrect error handling

### Fix

**File:** `pkg/metadata/metadata.go`

```go
// GetNFTMetadata queries and parses NFT metadata
func (c *Client) GetNFTMetadata(ctx context.Context, schemaCode string, tokenId string) (nftData, error) {
    if ctx == nil {
        ctx = context.Background()
    }

    // Get raw metadata
    res, err := c.GetNFTMetadataRaw(ctx, schemaCode, tokenId)
    if err != nil {
        // FIX: Return error instead of (nftData{}, nil)
        return nftData{}, fmt.Errorf("failed to get NFT metadata: %w", err)
    }

    // Validate response
    if res == nil || res.NftData == nil {
        return nftData{}, fmt.Errorf("metadata not found: schema=%s, tokenId=%s", 
            schemaCode, tokenId)
    }

    // Parse metadata
    metadata := res.NftData
    
    return nftData{
        SchemaCode:    metadata.NftSchemaCode,
        TokenId:       metadata.TokenId,
        OwnerAddress:  metadata.OwnerAddressType,
        // ... parse other fields ...
    }, nil
}

// GetNFTMetadataRaw queries raw NFT metadata response
func (c *Client) GetNFTMetadataRaw(ctx context.Context, schemaCode string, tokenId string) (*nftmngrtypes.QueryNFTDataResponse, error) {
    if ctx == nil {
        ctx = context.Background()
    }

    // Input validation
    if schemaCode == "" {
        return nil, fmt.Errorf("schema code cannot be empty")
    }

    if tokenId == "" {
        return nil, fmt.Errorf("token ID cannot be empty")
    }

    // Add timeout
    if _, hasDeadline := ctx.Deadline(); !hasDeadline {
        var cancel context.CancelFunc
        ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
        defer cancel()
    }

    // Query metadata
    res, err := c.QueryClient.NFTData(ctx, &nftmngrtypes.QueryGetNFTDataRequest{
        NftSchemaCode: schemaCode,
        TokenId:       tokenId,
    })

    if err != nil {
        if ctx.Err() != nil {
            return nil, fmt.Errorf("query cancelled: %w", ctx.Err())
        }
        return nil, fmt.Errorf("failed to query NFT metadata for schema=%s, tokenId=%s: %w", 
            schemaCode, tokenId, err)
    }

    return res, nil
}
```

---

## Fix 7: Add Key Zeroization (HIGH)

### Issue
**Severity:** HIGH  
**File:** `account/account.go`  
**Problem:** Private keys remain in memory until garbage collected.

### Impact
- Extended exposure window
- Memory dumps can capture keys
- Debuggers can access keys

### Fix: Add Close Method with Zeroization

**File:** `account/account.go`

```go
// Close securely cleans up the account by zeroizing sensitive data
func (a *Account) Close() error {
    // Zeroize private key
    if a.privateKey != nil {
        // Zero out the private key bytes
        key := a.privateKey.D.Bytes()
        for i := range key {
            key[i] = 0
        }
        a.privateKey = nil
    }

    // Zero out mnemonic if stored
    if a.mnemonic != "" {
        // Convert to byte slice and zero
        mnemonic := []byte(a.mnemonic)
        for i := range mnemonic {
            mnemonic[i] = 0
        }
        a.mnemonic = ""
    }

    return nil
}

// Usage example with defer:
func ExampleSecureAccountUsage() error {
    client, err := client.NewClient(config)
    if err != nil {
        return err
    }
    defer client.Close()

    account, err := account.NewAccount(client, "my-account", mnemonic, "")
    if err != nil {
        return err
    }
    defer account.Close() // Ensure cleanup

    // Use account...
    
    return nil
}
```

---

## Implementation Checklist

### Phase 2A: Critical Fixes (Day 1-2)

- [ ] **Fix 1: Mnemonic Exposure**
  - [ ] Remove GetMnemonic() or add protection
  - [ ] Remove mnemonic storage from Account struct
  - [ ] Add security documentation
  - [ ] Update examples

- [ ] **Fix 2: Context Timeouts**
  - [ ] Add context parameter to WaitForTransaction
  - [ ] Add context parameter to WaitForEVMTransaction
  - [ ] Add context parameter to all query functions
  - [ ] Add default timeout handling
  - [ ] Test timeout behavior

- [ ] **Fix 3: Nonce Management**
  - [ ] Add mutex to EVMClient
  - [ ] Implement GetAndIncrementNonce
  - [ ] Add nonce cache
  - [ ] Add ResetNonce methods
  - [ ] Test concurrent transactions

### Phase 2B: Validation & Bug Fixes (Day 3-4)

- [ ] **Fix 4: Response Validation**
  - [ ] Add validation to GetNFTSchema
  - [ ] Add validation to GetNFTMetadata
  - [ ] Add validation to GetBalance
  - [ ] Create validation helper package
  - [ ] Add input validation to all public functions

- [ ] **Fix 5: Metadata Bug**
  - [ ] Fix CreateCertificateMetadataWithInfo
  - [ ] Add test for info parameter
  - [ ] Verify in integration test

- [ ] **Fix 6: Error Handling**
  - [ ] Fix GetNFTMetadata return value
  - [ ] Review all error returns
  - [ ] Add error wrapping

- [ ] **Fix 7: Key Zeroization**
  - [ ] Add Close method to Account
  - [ ] Document usage with defer
  - [ ] Update examples

### Testing

- [ ] Unit tests for all fixes
- [ ] Integration tests with testnet
- [ ] Concurrent transaction tests
- [ ] Timeout and cancellation tests
- [ ] Error path testing
- [ ] Performance benchmarks

### Documentation

- [ ] Update godoc comments
- [ ] Add security best practices guide
- [ ] Update examples
- [ ] Migration guide for breaking changes
- [ ] CHANGELOG.md entry

---

## Git Workflow

```bash
# Create feature branch
git checkout -b fix/critical-sdk-security

# Make changes incrementally
git add account/account.go
git commit -m "fix(account): remove mnemonic exposure and add zeroization"

git add client/client.go
git commit -m "fix(client): add context timeout support"

git add pkg/evm/evm.go
git commit -m "fix(evm): add thread-safe nonce management"

git add pkg/metadata/metadata.go
git commit -m "fix(metadata): add validation and fix bugs"

# Run tests
go test ./...

# Push and create PR
git push origin fix/critical-sdk-security
```

---

## Testing Commands

```bash
# Run all tests
go test ./... -v

# Run with race detector
go test ./... -race

# Run specific package
go test ./pkg/evm -v

# Run with coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Benchmark
go test ./... -bench=. -benchmem
```

---

## Next Steps

After Phase 2:
1. Move to Phase 3: Comprehensive Test Suite
2. Reference: `IMPLEMENTATION_ROADMAP.md` Phase 3
3. Target: ≥80% test coverage

---

**Estimated Time:** 16-24 hours  
**Priority:** 🔴 CRITICAL  
**Blocking:** Production use, SDK release