# Implementation Roadmap - LBB SDK Go Security Fixes

This document outlines the prioritized implementation plan for addressing all audit findings in both the smart contracts and Go SDK.

---

## Executive Summary

**Total Estimated Effort:** ~100-120 hours (2-3 weeks with 1-2 engineers)

**Critical Issues to Fix:** 8 HIGH severity issues
**Important Issues:** 12 MEDIUM severity issues
**Enhancement Issues:** 15+ LOW severity issues

---

## Phase 1: Critical Smart Contract Fixes (IMMEDIATE - 4-8 hours)

### Priority: 🔴 CRITICAL - Must fix before any deployment

#### 1.1 Fix Permit Nonce Vulnerability (HIGH)
**Files:** `contracts/src/CertAutoID.sol`, `contracts/src/Cert.sol`
**Issue:** Nonce incremented before signature validation allows attacker to exhaust nonces
**Impact:** DOS attack vector, permit system can be broken
**Effort:** 30 minutes

**Fix:**
- Move `_nonces[owner]++` to AFTER signature validation in `permit()` and `permitForAll()`
- Add unit tests for invalid signature scenarios

#### 1.2 Add Token Existence Check (HIGH)
**Files:** `contracts/src/Cert.sol`
**Issue:** Manual `safeMint` lacks existence check for tokenId
**Impact:** Unclear error messages, potential duplicate mint issues
**Effort:** 15 minutes

**Fix:**
- Add `require(!_exists(tokenId), "Token already minted")` before minting

#### 1.3 Add Zero-Address Validation (MEDIUM-HIGH)
**Files:** `contracts/src/CertAutoID.sol`, `contracts/src/Cert.sol`
**Issue:** No validation for spender/operator addresses in permit functions
**Impact:** Invalid permits can be created, wasting gas
**Effort:** 30 minutes

**Fix:**
- Add `require(spender != address(0), "Invalid spender")` in `permit()`
- Add `require(operator != address(0), "Invalid operator")` in `permitForAll()`

#### 1.4 Fix Permit Consumption in Transfer/Burn (MEDIUM)
**Files:** `contracts/src/CertAutoID.sol`, `contracts/src/Cert.sol`
**Issue:** `transferWithPermit` and `burnWithPermit` consume permit even if operation fails
**Impact:** Wasted signatures, poor UX
**Effort:** 1-2 hours

**Fix:**
- Refactor to validate operation before consuming permit
- Or inline permit logic and only increment nonce on success

#### 1.5 Add Comprehensive Permit Tests (HIGH)
**Files:** `contracts/test/Cert.t.sol` (new test file recommended)
**Issue:** Zero permit test coverage
**Impact:** Cannot verify security fixes work correctly
**Effort:** 4-6 hours

**Tests to Add:**
- Valid permit flow (approve, transfer)
- Invalid signature rejection
- Nonce replay protection
- Expired deadline rejection
- Zero-address validation
- permitForAll scenarios
- transferWithPermit edge cases
- burnWithPermit edge cases

**Deliverable:** All contract tests pass, 100% coverage on permit functions

---

## Phase 2: Critical SDK Fixes (IMMEDIATE - 16-24 hours)

### Priority: 🔴 CRITICAL - Must fix before production use

#### 2.1 Remove/Protect Mnemonic Exposure (HIGH)
**Files:** `account/account.go`
**Issue:** `GetMnemonic()` returns plaintext mnemonic with no protection
**Impact:** Massive security risk if used in production
**Effort:** 1 hour

**Options:**
1. **Recommended:** Remove `GetMnemonic()` entirely
2. **Alternative:** Add explicit confirmation parameter: `GetMnemonicUnsafe(confirm bool)`
3. **Best:** Don't store mnemonic after key derivation

**Fix:**
```go
// Remove GetMnemonic() method entirely
// Add godoc warning on NewAccountFromMnemonic about secure mnemonic handling
```

#### 2.2 Fix Private Key Memory Security (HIGH)
**Files:** `account/account.go`
**Issue:** Private keys stored in memory unencrypted
**Impact:** Memory dumps, debuggers can expose keys
**Effort:** 2-4 hours

**Fix:**
- Add godoc warnings about key storage limitations
- Document best practices (use HSM, encrypted keystores in production)
- Optional: Implement key zeroization on `Close()` method
- Future enhancement: Support encrypted key storage

#### 2.3 Add Context Timeout Support (HIGH)
**Files:** `client/client.go`, `pkg/evm/evm.go`
**Issue:** Network operations don't respect context timeouts/cancellation
**Impact:** Resource leaks, hanging operations, poor error handling
**Effort:** 3-4 hours

**Functions to Fix:**
- `WaitForTransaction()`
- `WaitForEVMTransaction()`
- `QueryNFTSchemaByContract()`
- `QueryNFTMetadata()`
- All RPC calls

**Fix Pattern:**
```go
func (c *Client) WaitForTransaction(ctx context.Context, txHash string) (*sdk.TxResponse, error) {
    if ctx == nil {
        ctx = context.Background()
    }
    
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return nil, fmt.Errorf("waiting for transaction cancelled: %w", ctx.Err())
        case <-ticker.C:
            // Check transaction status
        }
    }
}
```

#### 2.4 Implement Thread-Safe Nonce Management (HIGH)
**Files:** `pkg/evm/evm.go`
**Issue:** Nonce handling not concurrency-safe, causes duplicate nonces
**Impact:** Transaction failures, race conditions
**Effort:** 4-6 hours

**Fix:**
```go
type EVMClient struct {
    // ... existing fields
    nonceMu    sync.Mutex
    nonceCache map[string]uint64 // cache per address
}

func (e *EVMClient) getNonceWithLock(ctx context.Context, from common.Address) (uint64, error) {
    e.nonceMu.Lock()
    defer e.nonceMu.Unlock()
    
    // Get nonce from chain or cache
    // Increment cached nonce for next use
}
```

**Alternative:** Document that SDK is not concurrency-safe and must be used serially

#### 2.5 Add Response Validation (HIGH)
**Files:** `pkg/metadata/metadata.go`, `pkg/balance/balance.go`
**Issue:** Query responses not validated (nil dereference risk)
**Impact:** Panic/crash on invalid responses
**Effort:** 2-3 hours

**Functions to Fix:**
- `GetNFTSchema()` - check nil before dereferencing
- `GetNFTMetadata()` - check nil before dereferencing
- `GetBalance()` - validate response structure
- All query functions

**Fix Pattern:**
```go
func (c *Client) GetNFTSchema(schemaCode string) (*nftmngrtypes.NFTSchemaByContractResponse, error) {
    if schemaCode == "" {
        return nil, fmt.Errorf("schema code cannot be empty")
    }
    
    res, err := c.QueryClient.NFTSchemaByContract(c.Context, &nftmngrtypes.QueryNFTSchemaByContractRequest{
        OriginContractAddress: schemaCode,
    })
    
    if err != nil {
        return nil, fmt.Errorf("failed to query NFT schema: %w", err)
    }
    
    if res == nil || res.NftSchemaByContract == nil {
        return nil, fmt.Errorf("received nil response for schema: %s", schemaCode)
    }
    
    return res, nil
}
```

#### 2.6 Fix Metadata Creation Bug (HIGH)
**Files:** `pkg/metadata/metadata.go`
**Issue:** `CreateCertificateMetadataWithInfo` calls wrong builder function
**Impact:** Functional bug - WithInfo parameter ignored
**Effort:** 15 minutes

**Fix:**
```go
func (c *Client) CreateCertificateMetadataWithInfo(schemaCode string, tokenId string, info string) (string, error) {
    // Change from BuildMintMetadataMsg to BuildMintMetadataWithInfoMsg
    msg := BuildMintMetadataWithInfoMsg(c.Account.GetAddress(), schemaCode, tokenId, info)
    // ... rest of function
}
```

#### 2.7 Fix GetNFTMetadata Error Handling (HIGH)
**Files:** `pkg/metadata/metadata.go`
**Issue:** Returns `(nftData{}, nil)` on error instead of returning error
**Impact:** Caller cannot detect errors
**Effort:** 10 minutes

**Fix:**
```go
func (c *Client) GetNFTMetadata(schemaCode string, tokenId string) (nftData, error) {
    res, err := c.GetNFTMetadataRaw(schemaCode, tokenId)
    if err != nil {
        return nftData{}, err // Return the error, not nil
    }
    // ... rest
}
```

**Deliverable:** All HIGH severity SDK issues resolved, basic validation in place

---

## Phase 3: Comprehensive Test Suite (Week 1-2, 32-40 hours)

### Priority: 🟡 HIGH - Required before production

#### 3.1 Smart Contract Test Suite
**Effort:** 12-16 hours

**Tests to Implement:**
- ✅ Basic minting (already exists)
- ❌ Permit workflow tests
  - Valid permit generation and verification
  - EIP-712 signature validation
  - Nonce management and replay protection
  - Deadline enforcement
  - Invalid signature rejection
- ❌ Permit edge cases
  - Expired permits
  - Invalid signers
  - Nonce exhaustion scenarios
  - Zero-address validation
- ❌ transferWithPermit tests
  - Successful transfer
  - Failed transfer (invalid recipient)
  - Permit consumption behavior
- ❌ burnWithPermit tests
  - Successful burn
  - Failed burn scenarios
- ❌ Gas estimation tests
- ❌ Integration tests (full workflows)

**Target Coverage:** ≥ 90% line coverage, 100% on critical paths

#### 3.2 SDK Unit Tests
**Effort:** 12-16 hours

**Packages to Test:**

**account/**
- ✅ NewAccountFromMnemonic (basic exists)
- ❌ NewAccountFromPrivateKey
- ❌ Key derivation paths
- ❌ Address generation
- ❌ Message signing
- ❌ EIP-712 permit signing

**client/**
- ❌ Client initialization
- ❌ Connection handling
- ❌ Transaction broadcast
- ❌ Transaction waiting with timeout
- ❌ Query operations
- ❌ Error handling

**pkg/evm/**
- ❌ EVM client initialization
- ❌ Contract deployment
- ❌ Transaction signing
- ❌ Nonce management (including concurrency)
- ❌ Gas estimation
- ❌ Permit signing and verification
- ❌ Transaction execution

**pkg/metadata/**
- ❌ Metadata message builders
- ❌ Metadata creation
- ❌ Metadata queries
- ❌ Schema operations

**pkg/balance/**
- ❌ Balance queries
- ❌ Send operations
- ❌ Multi-send

**Target Coverage:** ≥ 80% line coverage

#### 3.3 Integration Tests
**Effort:** 8-12 hours

**Test Scenarios:**
- End-to-end NFT minting workflow
- Certificate creation with metadata
- Gasless transfer using permits
- Concurrent transaction handling
- Network failure scenarios
- Timeout and retry behavior
- Full permit workflow (sign → verify → execute)

**Setup:**
- Local test chain (Ganache/Hardhat or Six Protocol testnet)
- Automated test data generation
- Cleanup procedures

**Deliverable:** Full test suite with ≥80% coverage, CI/CD integration

---

## Phase 4: Medium Priority Improvements (Week 2-3, 24-32 hours)

### Priority: 🟠 MEDIUM - Important for production quality

#### 4.1 Make Gas Buffer Configurable (MEDIUM)
**Files:** `pkg/evm/evm.go`
**Issue:** Hardcoded 20% gas buffer
**Effort:** 1 hour

**Fix:**
```go
type EVMClientConfig struct {
    RPCURL          string
    ChainID         *big.Int
    GasBufferFactor float64 // Default: 1.2 (20% buffer)
}

func (e *EVMClient) estimateGas(ctx context.Context, tx *types.Transaction) (uint64, error) {
    estimated, err := e.client.EstimateGas(ctx, ethereum.CallMsg{...})
    if err != nil {
        return 0, err
    }
    
    bufferFactor := e.config.GasBufferFactor
    if bufferFactor <= 1.0 {
        bufferFactor = 1.2 // Default 20%
    }
    
    return uint64(float64(estimated) * bufferFactor), nil
}
```

#### 4.2 Add Retry Logic (MEDIUM)
**Files:** `client/client.go`, `pkg/evm/evm.go`
**Issue:** No retry for transient network failures
**Effort:** 3-4 hours

**Fix:**
```go
type RetryConfig struct {
    MaxRetries  int
    InitialWait time.Duration
    MaxWait     time.Duration
    Multiplier  float64
}

func (c *Client) broadcastWithRetry(ctx context.Context, tx []byte) (*sdk.TxResponse, error) {
    var lastErr error
    wait := c.retryConfig.InitialWait
    
    for i := 0; i <= c.retryConfig.MaxRetries; i++ {
        if i > 0 {
            select {
            case <-ctx.Done():
                return nil, ctx.Err()
            case <-time.After(wait):
                wait = time.Duration(float64(wait) * c.retryConfig.Multiplier)
                if wait > c.retryConfig.MaxWait {
                    wait = c.retryConfig.MaxWait
                }
            }
        }
        
        resp, err := c.BroadcastTx(tx)
        if err == nil {
            return resp, nil
        }
        
        // Only retry on transient errors
        if !isRetryableError(err) {
            return nil, err
        }
        lastErr = err
    }
    
    return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}
```

#### 4.3 Add Input Validation (MEDIUM)
**Files:** All packages
**Issue:** Missing validation in many functions
**Effort:** 4-6 hours

**Validation to Add:**
- Zero-address checks for all address parameters
- Empty string checks for all string parameters
- Token ID format validation
- Schema code format validation
- Amount validation (>0)
- Gas limit bounds
- Chain ID validation

**Pattern:**
```go
func validateAddress(addr string) error {
    if addr == "" {
        return fmt.Errorf("address cannot be empty")
    }
    if !strings.HasPrefix(addr, "0x") || len(addr) != 42 {
        return fmt.Errorf("invalid address format: %s", addr)
    }
    return nil
}

func validateTokenID(tokenID string) error {
    if tokenID == "" {
        return fmt.Errorf("token ID cannot be empty")
    }
    // Add format-specific validation
    return nil
}
```

#### 4.4 Add Rate Limiting (MEDIUM)
**Files:** `client/client.go`
**Issue:** No rate limiting for API calls
**Effort:** 2-3 hours

**Fix:**
```go
import "golang.org/x/time/rate"

type Client struct {
    // ... existing fields
    rateLimiter *rate.Limiter
}

func NewClient(config Config) (*Client, error) {
    // ... existing code
    
    client := &Client{
        // ... existing fields
        rateLimiter: rate.NewLimiter(rate.Limit(config.RateLimit), config.RateBurst),
    }
    
    return client, nil
}

func (c *Client) rateLimit(ctx context.Context) error {
    if c.rateLimiter != nil {
        return c.rateLimiter.Wait(ctx)
    }
    return nil
}
```

#### 4.5 Improve Error Messages (MEDIUM)
**Files:** All packages
**Issue:** Generic error messages lack context
**Effort:** 3-4 hours

**Improvements:**
- Add relevant context to all errors (addresses, amounts, chain IDs)
- Use error wrapping consistently (`fmt.Errorf("context: %w", err)`)
- Define custom error types for common cases
- Add error codes for client error handling

**Example:**
```go
var (
    ErrInvalidAddress     = errors.New("invalid address")
    ErrInsufficientFunds  = errors.New("insufficient funds")
    ErrTransactionFailed  = errors.New("transaction failed")
    ErrContextCancelled   = errors.New("context cancelled")
)

type TransactionError struct {
    TxHash  string
    Reason  string
    Code    uint32
    Wrapped error
}

func (e *TransactionError) Error() string {
    return fmt.Sprintf("transaction %s failed (code %d): %s", e.TxHash, e.Code, e.Reason)
}

func (e *TransactionError) Unwrap() error {
    return e.Wrapped
}
```

#### 4.6 Add Client Close Method (MEDIUM)
**Files:** `client/client.go`, `pkg/evm/evm.go`
**Issue:** No graceful shutdown, resource leaks
**Effort:** 1-2 hours

**Fix:**
```go
func (c *Client) Close() error {
    c.cancel() // Cancel context
    
    // Close gRPC connections
    if c.grpcConn != nil {
        if err := c.grpcConn.Close(); err != nil {
            return fmt.Errorf("failed to close gRPC connection: %w", err)
        }
    }
    
    // Zeroize sensitive data (optional)
    if c.Account != nil {
        c.Account.Zeroize()
    }
    
    return nil
}

func (e *EVMClient) Close() error {
    if e.client != nil {
        e.client.Close()
    }
    return nil
}
```

#### 4.7 Remove Test Constants from Production Code (MEDIUM)
**Files:** Various
**Issue:** Test mnemonics/keys in production code
**Effort:** 1 hour

**Fix:**
- Move all test constants to `*_test.go` files
- Add build tags if needed: `// +build !production`
- Document that example keys should never be used in production

**Deliverable:** Production-quality SDK with robust error handling

---

## Phase 5: Documentation & Polish (Week 3, 12-16 hours)

### Priority: 🟢 LOW-MEDIUM - Important for adoption

#### 5.1 Add Godoc Comments (LOW)
**Files:** All packages
**Effort:** 4-6 hours

**Requirements:**
- Package-level documentation for each package
- Function-level documentation for all exported functions
- Example code in godoc where appropriate
- Security warnings where relevant

#### 5.2 Structured Logging (LOW)
**Files:** All packages
**Effort:** 3-4 hours

**Recommendation:**
```go
import "log/slog"

type Client struct {
    // ... existing fields
    logger *slog.Logger
}

func (c *Client) BroadcastTx(tx []byte) (*sdk.TxResponse, error) {
    c.logger.Debug("broadcasting transaction",
        slog.Int("size", len(tx)),
        slog.String("chain_id", c.ChainID),
    )
    
    // ... implementation
    
    c.logger.Info("transaction broadcast successful",
        slog.String("tx_hash", resp.TxHash),
        slog.Int64("height", resp.Height),
    )
}
```

#### 5.3 Add Metrics/Observability (LOW)
**Files:** Core packages
**Effort:** 4-6 hours

**Metrics to Track:**
- Transaction success/failure rate
- Transaction latency
- RPC call latency
- Error rates by type
- Active connections
- Nonce cache hit rate

**Recommendation:** Use Prometheus client library

#### 5.4 Security Documentation (MEDIUM)
**Effort:** 2-3 hours

**Documents to Create:**
- SECURITY.md - Security best practices
- Key management guidelines
- Production deployment checklist
- Threat model documentation
- Incident response guide

**Deliverable:** Well-documented, production-ready SDK

---

## Phase 6: Optional Enhancements (Future)

### Priority: 🔵 OPTIONAL - Nice to have

#### 6.1 Encrypted Key Storage
**Effort:** 8-12 hours
- Implement encrypted keystore (keystore.json format)
- Add password-based key derivation (PBKDF2 or scrypt)
- Support key import/export

#### 6.2 HSM Integration
**Effort:** 16-24 hours
- Add interface for external signers
- Support hardware security modules
- Integrate with cloud KMS (AWS KMS, GCP KMS)

#### 6.3 Connection Pooling
**Effort:** 4-6 hours
- Pool gRPC connections
- Implement connection reuse
- Add health checks

#### 6.4 Response Caching
**Effort:** 3-4 hours
- Cache transaction receipts
- Cache query results (with TTL)
- Implement cache invalidation

#### 6.5 EIP-4494 Support
**Effort:** 8-12 hours
- Add EIP-4494 permit support to contracts
- Update SDK to support new permit format
- Add comprehensive tests

#### 6.6 Contract Refactoring
**Effort:** 6-8 hours
- Create shared base contract
- Eliminate code duplication
- Standardize event emissions

---

## Implementation Timeline

### Week 1: Critical Fixes
**Days 1-2:** Phase 1 - Smart Contract Fixes
- [ ] Fix permit nonce vulnerability
- [ ] Add token existence check
- [ ] Add zero-address validation
- [ ] Fix permit consumption
- [ ] Add permit tests
- [ ] Run full test suite
- [ ] Deploy to testnet for validation

**Days 3-5:** Phase 2 - Critical SDK Fixes
- [ ] Remove/protect mnemonic exposure
- [ ] Add context timeout support
- [ ] Implement thread-safe nonce management
- [ ] Add response validation
- [ ] Fix metadata bugs
- [ ] Manual testing of critical paths

### Week 2: Testing & Medium Priority
**Days 1-3:** Phase 3 - Test Suite
- [ ] Smart contract test completion
- [ ] SDK unit tests
- [ ] Integration tests
- [ ] CI/CD setup
- [ ] Coverage reports

**Days 4-5:** Phase 4 - Medium Priority (Part 1)
- [ ] Configurable gas buffer
- [ ] Retry logic
- [ ] Input validation
- [ ] Rate limiting

### Week 3: Polish & Documentation
**Days 1-2:** Phase 4 - Medium Priority (Part 2)
- [ ] Improve error messages
- [ ] Add Close methods
- [ ] Remove test constants
- [ ] Code review

**Days 3-5:** Phase 5 - Documentation & Polish
- [ ] Godoc comments
- [ ] Structured logging
- [ ] Metrics/observability
- [ ] Security documentation
- [ ] Final review and release prep

---

## Success Criteria

### Smart Contracts
- ✅ All HIGH severity issues fixed
- ✅ ≥90% test coverage
- ✅ All tests passing
- ✅ Deployed and validated on testnet
- ✅ No known security vulnerabilities

### SDK
- ✅ All HIGH severity issues fixed
- ✅ ≥80% test coverage
- ✅ All tests passing
- ✅ Integration tests validated on testnet
- ✅ Complete godoc documentation
- ✅ Security best practices documented
- ✅ No plaintext secret exposure

### Overall
- ✅ External security audit (recommended after fixes)
- ✅ Performance benchmarks established
- ✅ Production deployment checklist complete
- ✅ Monitoring and alerting configured
- ✅ Incident response procedures documented

---

## Risk Mitigation

### Technical Risks
1. **Breaking Changes**: Changes to nonce management may break existing code
   - **Mitigation**: Version bump (v2.0.0), migration guide, deprecation period

2. **Test Coverage Gaps**: May not catch all edge cases
   - **Mitigation**: Fuzzing, external audit, bug bounty program

3. **Performance Impact**: Added validation and locking may slow operations
   - **Mitigation**: Benchmark before/after, optimize hot paths

### Process Risks
1. **Timeline Slippage**: Fixes may take longer than estimated
   - **Mitigation**: Focus on HIGH severity first, defer optional items

2. **Resource Constraints**: May need more engineering time
   - **Mitigation**: Prioritize ruthlessly, consider external contractors

3. **Regression Risk**: Fixes may introduce new bugs
   - **Mitigation**: Comprehensive test suite, staged rollout, monitoring

---

## Post-Implementation

### Deployment Strategy
1. Deploy fixed contracts to testnet
2. Run SDK against testnet for 1 week
3. Bug bounty program on testnet
4. External security audit
5. Address audit findings
6. Mainnet deployment (phased if possible)
7. Monitor closely for 2 weeks

### Monitoring Plan
- Transaction success rate ≥99%
- RPC call latency <100ms p95
- Error rate <1%
- Zero security incidents
- Alert on anomalies

### Maintenance
- Weekly dependency updates
- Monthly security reviews
- Quarterly external audits
- Continuous test coverage improvement
- Community bug reports triaged within 24h

---

## Get Started

To begin implementation, start with Phase 1:

```bash
# 1. Create a feature branch
git checkout -b fix/security-audit-phase1

# 2. Fix contract nonce vulnerability
cd contracts/src
# Edit CertAutoID.sol and Cert.sol

# 3. Run tests
forge test

# 4. Commit and create PR
git commit -am "fix: move nonce increment after signature validation"
git push origin fix/security-audit-phase1
```

Then follow the phase-by-phase checklist in each section above.

---

## Questions or Concerns?

If you have questions about any item in this roadmap:
1. Review the detailed audit reports (SECURITY_AUDIT.md, SDK_DEEP_AUDIT.md)
2. Check AUDIT_FIXES.md for specific code examples
3. Consult with the security team
4. Escalate to engineering leadership if timeline/resources are insufficient

---

**Last Updated:** 2024
**Document Owner:** Engineering Team
**Review Cadence:** Weekly during implementation, monthly post-deployment