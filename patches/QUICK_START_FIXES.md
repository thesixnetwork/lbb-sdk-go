# Quick Start: Immediate Action Checklist

This is your fast-track guide to implementing the most critical security fixes identified in the audit. Follow this step-by-step to secure your codebase.

---

## 🔴 STOP: Read This First

**Current Status:** Your codebase has **8 CRITICAL/HIGH severity security issues** that must be fixed before any production deployment.

**Time Required:** 
- Critical Contract Fixes: 4-8 hours
- Critical SDK Fixes: 16-24 hours
- Total: ~3-4 days with testing

**Team Assignment:**
- Engineer 1: Smart Contract fixes (Phase 1)
- Engineer 2: SDK fixes (Phase 2)
- Both: Integration testing (Phase 3)

---

## Week 1, Day 1-2: Smart Contract Critical Fixes

### Morning: Setup & Planning (1 hour)

```bash
# 1. Create feature branch
cd /path/to/lbb-sdk-go
git checkout -b fix/security-audit-critical
git pull origin main

# 2. Read audit reports
open SECURITY_AUDIT.md
open SDK_DEEP_AUDIT.md
open AUDIT_ACTION_ITEMS.md

# 3. Review implementation roadmap
open IMPLEMENTATION_ROADMAP.md
```

### Afternoon: Contract Fix #1 - Nonce Vulnerability (2 hours)

**Priority:** 🔴 CRITICAL - DOS vulnerability

#### Step 1: Fix CertAutoID.sol

```bash
cd contracts/src
code CertAutoID.sol
```

**Find:** Lines ~116-154 (permit function)

**Change from:**
```solidity
_nonces[owner]++  // ❌ WRONG - increments before validation
```

**Change to:**
```solidity
uint256 currentNonce = _nonces[owner];  // ✅ Store current nonce
// ... validation code ...
_nonces[owner] = currentNonce + 1;      // ✅ Only increment after validation passes
```

#### Step 2: Apply same fix to Cert.sol

```bash
code Cert.sol
```

Repeat the same changes for both `permit()` and `permitForAll()` functions.

#### Step 3: Add zero-address checks

In both contracts, add at start of `permit()`:
```solidity
if (spender == address(0)) {
    revert InvalidSigner(); // or create InvalidSpender error
}
```

In both contracts, add at start of `permitForAll()`:
```solidity
if (operator == address(0)) {
    revert InvalidSigner(); // or create InvalidOperator error
}
```

#### Step 4: Test

```bash
cd contracts
forge build
forge test -vv

# Should see all tests passing
# ✅ testSafeMint (gas: 123456)
# ✅ testSafeMintMultiple (gas: 234567)
```

**Commit:**
```bash
git add contracts/src/CertAutoID.sol contracts/src/Cert.sol
git commit -m "fix(contracts): prevent nonce increment before signature validation

- Move nonce increment to after signature validation
- Add zero-address validation for spender/operator
- Prevents DOS attack via invalid signatures

Security: HIGH-CRITICAL"
```

---

### Day 1 Evening: Contract Fix #2 - Token Existence (30 min)

**Priority:** 🔴 HIGH

```bash
code contracts/src/Cert.sol
```

**Find:** safeMint function (~line 58)

**Change from:**
```solidity
function safeMint(address to, uint256 tokenId) public onlyOwner {
    _safeMint(to, tokenId);
    emit safeMintEvent(to, tokenId);
}
```

**Change to:**
```solidity
function safeMint(address to, uint256 tokenId) public onlyOwner {
    if (_ownerOf(tokenId) != address(0)) {
        revert("Token already minted");
    }
    _safeMint(to, tokenId);
    emit safeMintEvent(to, tokenId);
}
```

**Test & Commit:**
```bash
forge test -vv
git add contracts/src/Cert.sol
git commit -m "fix(contracts): add token existence check to prevent duplicate mints"
```

---

## Day 2: Contract Tests & Validation

### Morning: Create Comprehensive Test Suite (4 hours)

```bash
cd contracts/test
touch PermitTest.t.sol
code PermitTest.t.sol
```

**Copy test template from:** `patches/PHASE1_CRITICAL_CONTRACT_FIXES.md` (lines 465-850)

**Run tests:**
```bash
forge test --match-path test/PermitTest.t.sol -vvv

# Expected output:
# ✅ testPermitValidSignature
# ✅ testPermitInvalidSignatureShouldNotIncrementNonce  ⬅️ CRITICAL TEST
# ✅ testPermitCannotReplaySignature
# ✅ testPermitExpiredDeadline
# ✅ testTransferWithPermitSuccess
# ✅ testBurnWithPermitSuccess
```

**Check coverage:**
```bash
forge coverage

# Target: ≥90% coverage on permit functions
```

**Commit:**
```bash
git add contracts/test/PermitTest.t.sol
git commit -m "test(contracts): add comprehensive permit test suite

- Tests for valid/invalid signatures
- Nonce replay protection tests
- Zero-address validation tests
- Transfer/burn with permit tests

Coverage: 95% on permit functions"
```

### Afternoon: Deploy & Validate on Testnet (2 hours)

```bash
# Deploy to testnet
forge script script/Deploy.s.sol --rpc-url $TESTNET_RPC --broadcast

# Test with SDK
cd ..
go run example/permit/main.go

# Manual testing checklist:
# ✅ Mint token
# ✅ Generate permit signature
# ✅ Execute gasless transfer
# ✅ Verify nonce increments only on success
# ✅ Test invalid signature rejection
```

**Push contracts:**
```bash
git push origin fix/security-audit-critical
```

---

## Day 3-4: SDK Critical Fixes

### Day 3 Morning: Fix #1 - Remove Mnemonic Exposure (1 hour)

**Priority:** 🔴 CRITICAL - Complete account compromise risk

```bash
cd account
code account.go
```

#### Step 1: Remove GetMnemonic()

**Find and DELETE:** Lines ~130-135
```go
// DELETE THIS ENTIRE FUNCTION:
// func (a *Account) GetMnemonic() string {
//     return a.mnemonic
// }
```

#### Step 2: Remove mnemonic field from struct

**Find:** Account struct definition (~line 30)

**Change from:**
```go
type Account struct {
    client        client.ClientI
    auth          *bind.TransactOpts
    mnemonic      string              // ❌ REMOVE THIS
    privateKey    *ecdsa.PrivateKey
    evmAddress    common.Address
    cosmosAddress sdk.AccAddress
    accountName   string
}
```

**Change to:**
```go
type Account struct {
    client        client.ClientI
    auth          *bind.TransactOpts
    // mnemonic removed for security
    privateKey    *ecdsa.PrivateKey
    evmAddress    common.Address
    cosmosAddress sdk.AccAddress
    accountName   string
}
```

#### Step 3: Remove mnemonic from NewAccount

**Find:** NewAccount function return (~line 95)

**Change from:**
```go
return &Account{
    client:        ctx,
    auth:          authz,
    privateKey:    privateKey,
    mnemonic:      mnemonic,  // ❌ REMOVE THIS
    evmAddress:    evmAddress,
    cosmosAddress: cosmosAddress,
    accountName:   accountName,
}, nil
```

**Change to:**
```go
return &Account{
    client:        ctx,
    auth:          authz,
    privateKey:    privateKey,
    evmAddress:    evmAddress,
    cosmosAddress: cosmosAddress,
    accountName:   accountName,
}, nil
```

**Test & Commit:**
```bash
cd ..
go test ./account/... -v
git add account/account.go
git commit -m "fix(account): remove mnemonic storage and exposure

- Remove GetMnemonic() function
- Remove mnemonic field from Account struct
- Mnemonic only used during key derivation, not stored

Security: CRITICAL - prevents mnemonic exposure"
```

---

### Day 3 Afternoon: Fix #2 - Context Timeouts (3 hours)

**Priority:** 🔴 HIGH - Resource leaks, hanging operations

#### Step 1: Fix WaitForTransaction

```bash
code client/client.go
```

**Find:** WaitForTransaction function

**Replace with version from:** `patches/PHASE2_CRITICAL_SDK_FIXES.md` (lines 248-320)

Key changes:
- Add `context.Context` as first parameter
- Add timeout handling
- Add context cancellation support
- Add better error messages

#### Step 2: Fix WaitForEVMTransaction

```bash
code pkg/evm/evm.go
```

**Replace with version from:** `patches/PHASE2_CRITICAL_SDK_FIXES.md` (lines 325-385)

#### Step 3: Add context to query functions

Update all query functions in:
- `pkg/metadata/metadata.go`
- `pkg/balance/balance.go`

Pattern:
```go
// Before:
func (c *Client) GetNFTSchema(schemaCode string) (*Response, error)

// After:
func (c *Client) GetNFTSchema(ctx context.Context, schemaCode string) (*Response, error) {
    if ctx == nil {
        ctx = context.Background()
    }
    // ... add timeout if not set ...
}
```

**Test:**
```bash
go test ./client/... -v
go test ./pkg/evm/... -v
go test ./pkg/metadata/... -v
```

**Commit:**
```bash
git add client/client.go pkg/evm/evm.go pkg/metadata/metadata.go pkg/balance/balance.go
git commit -m "fix(client): add context timeout support to network operations

- Add context.Context parameter to all network operations
- Implement timeout handling
- Support context cancellation
- Prevent resource leaks

Security: HIGH - prevents hanging operations"
```

---

### Day 4 Morning: Fix #3 - Thread-Safe Nonces (3 hours)

**Priority:** 🔴 HIGH - Race conditions, transaction failures

```bash
code pkg/evm/evm.go
```

**Replace EVMClient struct and nonce methods with version from:**
`patches/PHASE2_CRITICAL_SDK_FIXES.md` (lines 470-650)

Key additions:
- `nonceMu sync.Mutex`
- `nonceCache map[common.Address]uint64`
- `GetNonce()` with locking
- `GetAndIncrementNonce()` with locking
- `ResetNonce()` methods

**Test concurrency:**
```bash
# Create test file
cat > pkg/evm/nonce_test.go << 'EOF'
package evm

import (
    "context"
    "sync"
    "testing"
    "github.com/ethereum/go-ethereum/common"
)

func TestConcurrentNonceAccess(t *testing.T) {
    client := setupTestClient(t)
    address := common.HexToAddress("0x123...")
    
    var wg sync.WaitGroup
    nonces := make([]uint64, 10)
    
    // Get 10 nonces concurrently
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(index int) {
            defer wg.Done()
            nonce, err := client.GetAndIncrementNonce(context.Background(), address)
            if err != nil {
                t.Error(err)
            }
            nonces[index] = nonce
        }(i)
    }
    
    wg.Wait()
    
    // Verify all nonces are unique
    seen := make(map[uint64]bool)
    for _, nonce := range nonces {
        if seen[nonce] {
            t.Errorf("Duplicate nonce: %d", nonce)
        }
        seen[nonce] = true
    }
}
EOF

go test ./pkg/evm/... -race -v
```

**Commit:**
```bash
git add pkg/evm/evm.go pkg/evm/nonce_test.go
git commit -m "fix(evm): add thread-safe nonce management

- Add mutex for nonce operations
- Implement nonce caching
- Add GetAndIncrementNonce for atomic operations
- Add nonce reset methods

Security: HIGH - prevents race conditions in concurrent transactions"
```

---

### Day 4 Afternoon: Fixes #4-7 - Validation & Bug Fixes (3 hours)

#### Fix #4: Response Validation

**Files:** `pkg/metadata/metadata.go`, `pkg/balance/balance.go`

Add nil checks before dereferencing all responses:
```go
if res == nil {
    return nil, fmt.Errorf("received nil response")
}
if res.Data == nil {
    return nil, fmt.Errorf("data not found")
}
```

#### Fix #5: Metadata Bug

**File:** `pkg/metadata/metadata.go`

**Find:** CreateCertificateMetadataWithInfo

**Change:**
```go
// WRONG:
msg := BuildMintMetadataMsg(...)

// CORRECT:
msg := BuildMintMetadataWithInfoMsg(...)
```

#### Fix #6: Error Handling

**File:** `pkg/metadata/metadata.go`

**Find:** GetNFTMetadata return statements

**Change all:**
```go
// WRONG:
return nftData{}, nil

// CORRECT:
return nftData{}, err
```

#### Fix #7: Key Zeroization

**File:** `account/account.go`

**Add Close method:**
```go
func (a *Account) Close() error {
    if a.privateKey != nil {
        key := a.privateKey.D.Bytes()
        for i := range key {
            key[i] = 0
        }
        a.privateKey = nil
    }
    return nil
}
```

**Test all:**
```bash
go test ./... -v
go test ./... -race
```

**Commit:**
```bash
git add pkg/metadata/metadata.go pkg/balance/balance.go account/account.go
git commit -m "fix(sdk): add validation and fix critical bugs

- Add response validation to prevent nil dereference
- Fix CreateCertificateMetadataWithInfo to use correct builder
- Fix GetNFTMetadata error return value
- Add Close() method with key zeroization

Security: HIGH - multiple critical bugs fixed"
```

---

## Day 5: Integration Testing & Validation

### Morning: Integration Tests (3 hours)

```bash
# Test full workflow
cd example
go run main.go

# Checklist:
# ✅ Account creation without mnemonic exposure
# ✅ Context timeouts work correctly
# ✅ Concurrent transactions succeed
# ✅ Metadata with info is created correctly
# ✅ Errors are properly returned
# ✅ No panics on nil responses
```

### Afternoon: Documentation & PR (2 hours)

#### Update Documentation

```bash
# Create security guide
cat > SECURITY_BEST_PRACTICES.md << 'EOF'
# Security Best Practices

## Mnemonic Management
- Never log mnemonics
- Never store in plaintext
- Display only once during backup
- Use hardware wallets for production

## Key Management
- Call account.Close() when done
- Use defer for cleanup
- Consider encrypted keystores

## Context Usage
- Always pass context with timeout
- Handle cancellation errors
- Use context.WithTimeout for long operations
EOF

git add SECURITY_BEST_PRACTICES.md
```

#### Create PR

```bash
# Final commit
git add .
git commit -m "docs: add security best practices guide"

# Push
git push origin fix/security-audit-critical

# Create PR with template from patches/PHASE1_CRITICAL_CONTRACT_FIXES.md
```

---

## Validation Checklist

Before merging, verify:

### Smart Contracts
- [ ] All permit functions increment nonce AFTER validation
- [ ] Zero-address checks in place
- [ ] Token existence check in manual safeMint
- [ ] All tests pass: `forge test -vv`
- [ ] Coverage ≥90% on permit functions: `forge coverage`
- [ ] Deployed and tested on testnet
- [ ] No compiler warnings

### SDK
- [ ] GetMnemonic() removed or protected
- [ ] Mnemonic not stored in Account struct
- [ ] All network functions accept context.Context
- [ ] Context timeouts implemented
- [ ] Nonce management is thread-safe
- [ ] All responses validated before dereferencing
- [ ] CreateCertificateMetadataWithInfo uses correct builder
- [ ] GetNFTMetadata returns errors correctly
- [ ] Close() method implemented with key zeroization
- [ ] All tests pass: `go test ./... -v`
- [ ] Race detector clean: `go test ./... -race`
- [ ] Integration tests pass

### Documentation
- [ ] Security best practices documented
- [ ] Breaking changes documented
- [ ] Examples updated
- [ ] CHANGELOG.md updated
- [ ] PR description complete

---

## Emergency Rollback Plan

If issues are discovered after deployment:

```bash
# 1. Immediately revert
git revert <commit-hash>
git push origin main --force-with-lease

# 2. Notify team
# Post in #security channel

# 3. Assess damage
# Check logs for exploitation attempts

# 4. Fix and redeploy
# Apply hotfix
# Test thoroughly
# Deploy with monitoring
```

---

## Post-Deployment Monitoring

First 48 hours after deployment:

### Metrics to Watch
- Transaction success rate (should be ≥99%)
- Permit function gas usage (should be <100k)
- Nonce-related errors (should be near 0)
- Context timeout errors (track patterns)
- Response validation errors (should be rare)

### Alerts to Set
- ⚠️  Transaction failure rate >1%
- 🔴 Security-related errors (nonce manipulation attempts)
- 🟡 High context timeout rate
- 🟡 Unusual error patterns

### Log Monitoring
```bash
# Watch for security issues
grep -i "nonce" logs/*.log
grep -i "invalid signature" logs/*.log
grep -i "context" logs/*.log
grep -i "panic" logs/*.log
```

---

## Success Criteria

### You're done when:
1. ✅ All HIGH/CRITICAL issues fixed
2. ✅ All tests passing (contracts and SDK)
3. ✅ Deployed to testnet successfully
4. ✅ Integration tests pass
5. ✅ PR approved and merged
6. ✅ Documentation complete
7. ✅ Team trained on new security practices
8. ✅ Monitoring in place

---

## Need Help?

### Audit Reports
- Detailed smart contract audit: `SECURITY_AUDIT.md`
- Detailed SDK audit: `SDK_DEEP_AUDIT.md`
- Complete action items: `AUDIT_ACTION_ITEMS.md`

### Code Examples
- Contract fixes: `patches/PHASE1_CRITICAL_CONTRACT_FIXES.md`
- SDK fixes: `patches/PHASE2_CRITICAL_SDK_FIXES.md`

### Implementation Plan
- Full roadmap: `IMPLEMENTATION_ROADMAP.md`

### Get Assistance
1. Review audit reports for context
2. Check code examples in patches/
3. Ask in team security channel
4. Escalate to security team if blocked

---

**Remember:** Security is not optional. These fixes must be completed before any production deployment.

Good luck! 🚀