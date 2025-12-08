# Changelog - v0.1.0

## Release Date: Mon Dec 8 2025

This is the first major refactoring release of the LBB SDK Go, transforming the codebase from prototype to production-ready with proper Go idioms and best practices.

---

## 🎉 Highlights

- **Proper Error Handling**: All constructors now return errors instead of nil
- **Interface-Based Design**: Complete interfaces for better testability and flexibility
- **Input Validation**: Comprehensive validation to prevent runtime panics
- **Better API**: Fluent builder pattern for configuration
- **Enhanced Documentation**: Complete examples and migration guides

---

## 🔥 Breaking Changes

### Client Package (`client/client.go`)

#### Constructor Changes
```go
// OLD
func NewClient(ctx context.Context, mainnet bool) (Client, error)

// NEW
func NewClient(ctx context.Context, mainnet bool) (*Client, error)
```

**Migration:**
```go
// Before
client, err := client.NewClient(ctx, false)

// After (no changes needed, but now returns pointer)
client, err := client.NewClient(ctx, false)
if err != nil {
    return fmt.Errorf("failed to create client: %w", err)
}
```

#### Field Access Changes
```go
// OLD
client.ChainID          // Direct field access
client.ETHClient        // Direct field access
client.EVMRPCCleint     // Typo!

// NEW
client.GetChainID()     // Method access
client.GetETHClient()   // Method access
client.GetEVMRPCClient() // Fixed typo, method access
```

### Account Package (`account/account.go`)

#### Constructor Changes
```go
// OLD
func NewAccount(ctx client.Client, accountName, mnemonic, password string) *Account

// NEW
func NewAccount(ctx client.ClientI, accountName, mnemonic, password string) (*Account, error)
```

**Migration:**
```go
// Before
account := account.NewAccount(client, "test", mnemonic, "")
if account == nil {
    log.Fatal("failed")
}

// After
account, err := account.NewAccount(client, "test", mnemonic, "")
if err != nil {
    return fmt.Errorf("failed to create account: %w", err)
}
```

#### Method Changes
```go
// OLD
func (a *Account) ValidateMnemonic(mnemonic string) bool

// NEW (Package-level function)
func ValidateMnemonic(mnemonic string) bool
```

**Migration:**
```go
// Before
if account.ValidateMnemonic(mnemonic) { ... }

// After
if account.ValidateMnemonic(mnemonic) { ... }
```

#### Field Access Changes
```go
// OLD
account.Client.GetChainID()  // Through embedded struct

// NEW
account.GetClient().GetChainID()  // Explicit method call
```

### Account Msg Package (`account/msg.go`)

#### Constructor Changes
```go
// OLD
func NewAccountMsg(a Account) *AccountMsg

// NEW
func NewAccountMsg(a AccountI) (*AccountMsg, error)
```

**Migration:**
```go
// Before
accountMsg := account.NewAccountMsg(*account)

// After
accountMsg, err := account.NewAccountMsg(account)
if err != nil {
    return fmt.Errorf("failed to create account msg: %w", err)
}
```

---

## ✨ New Features

### 1. New Account Creation Methods

```go
// Create account from existing private key
func NewAccountFromPrivateKey(ctx client.ClientI, accountName string, privateKey *ecdsa.PrivateKey) (*Account, error)
```

### 2. Mnemonic Helper Functions

```go
// Generate new mnemonic (renamed from GenerateMnemonic)
func GenerateNewMnemonic() (string, error)

// Validate mnemonic (now package-level)
func ValidateMnemonic(mnemonic string) bool
```

### 3. Convenience Transaction Methods

```go
// Broadcast and wait for confirmation in one call
func (a *AccountMsg) BroadcastTxAndWait(msgs ...sdk.Msg) (*sdk.TxResponse, error)
```

### 4. Builder Pattern for Configuration

```go
// Fluent API for AccountMsg
accountMsg.
    WithGas(2000000).
    WithMemo("test transaction").
    WithGasAdjustment(1.8).
    WithFees("1000usix").
    BroadcastTx(msgs...)
```

### 5. Client Configuration Methods

```go
// Fluent API for Client
client.
    WithFrom("cosmos1...").
    WithFromName("alice").
    WithBroadcastMode("sync")
```

### 6. Enhanced Wait Methods

```go
// Wait for Cosmos transaction with better error messages
func (c *Client) WaitForTransaction(txHash string) error

// Wait for EVM transaction with validation
func (c *Client) WaitForEVMTransaction(txHash common.Hash) (*types.Receipt, error)
```

### 7. String Representation

```go
// Debug-friendly string representation
func (a *Account) String() string
```

---

## 🔧 Improvements

### Error Handling
- ✅ All constructors return proper errors with context
- ✅ Error wrapping using `fmt.Errorf` with `%w`
- ✅ No more silent failures or nil returns
- ✅ Descriptive error messages with operation context

### Input Validation
- ✅ Nil pointer checks on all public functions
- ✅ Empty string validation
- ✅ Mnemonic validation before use
- ✅ Chain ID existence checks
- ✅ Transaction hash validation

### Code Quality
- ✅ Fixed typo: `EVMRPCCleint` → `EVMRPCClient`
- ✅ Consistent naming conventions
- ✅ Proper documentation comments on all exported items
- ✅ Encapsulation: private fields with public getters
- ✅ Interface-based design for testability

### Structure
- ✅ Composition over deep embedding
- ✅ Clear method ownership
- ✅ Complete interface definitions
- ✅ Immutable builder pattern

### Logging
- ✅ Structured logging output
- ✅ Clear success/failure messages
- ✅ Transaction details in logs

### Constants
- ✅ Better organization and documentation
- ✅ Named constants for timeouts
- ✅ Clear comments explaining values

---

## 📚 Documentation

### New Documentation Files
- **REVIEW_SUMMARY.md** - Quick overview of all changes
- **BEFORE_AFTER_COMPARISON.md** - Visual side-by-side comparisons
- **REFACTORING_NOTES.md** - Detailed explanations and rationale
- **VERIFICATION_CHECKLIST.md** - Testing checklist
- **REFACTORING_README.md** - Migration guide
- **example/EXAMPLES.md** - Comprehensive API examples

### Updated Examples
- **cmd/main.go** - Updated to use v0.1.0 API with proper error handling

---

## 🐛 Bug Fixes

- Fixed: Typo `EVMRPCCleint` → `EVMRPCClient`
- Fixed: Silent error swallowing in `NewClient` and `NewCustomClient`
- Fixed: Nil pointer panics from lack of validation
- Fixed: Context handling in wait functions
- Fixed: Error messages missing context

---

## 🔄 Deprecations

### Deprecated (Still Works, Use New API)
- None (this is the first major release)

### Removed
- None (clean break from prototype to v0.1.0)

---

## 📊 Statistics

- **Files Refactored:** 3 (client.go, account.go, msg.go)
- **New Methods Added:** 15+
- **Interfaces Completed:** 3 (ClientI, AccountI, AccountMsgI)
- **Documentation Pages:** 6
- **Breaking Changes:** 8
- **New Features:** 7
- **Bug Fixes:** 5

---

## 🧪 Testing

### Test Coverage Goals
- Unit tests for all public functions: >80%
- Integration tests for main flows: >60%
- Error path tests: >70%

### Test Updates Required
All existing tests need to be updated to:
1. Handle new error returns from constructors
2. Use pointer types for Client
3. Call package-level functions instead of methods where applicable

---

## 🚀 Migration Guide

### Step 1: Update Imports
No changes needed - package names remain the same.

### Step 2: Update Client Creation
```go
// Add error handling
client, err := client.NewClient(ctx, false)
if err != nil {
    return fmt.Errorf("failed to create client: %w", err)
}
```

### Step 3: Update Account Creation
```go
// Validate mnemonic first
if !account.ValidateMnemonic(mnemonic) {
    return fmt.Errorf("invalid mnemonic")
}

// Add error handling
account, err := account.NewAccount(client, "myaccount", mnemonic, "")
if err != nil {
    return fmt.Errorf("failed to create account: %w", err)
}
```

### Step 4: Update AccountMsg Creation
```go
// Add error handling
accountMsg, err := account.NewAccountMsg(account)
if err != nil {
    return fmt.Errorf("failed to create account msg: %w", err)
}
```

### Step 5: Update Field Access
```go
// Replace direct field access with getters
chainID := client.GetChainID()
ethClient := client.GetETHClient()
```

### Step 6: Test Thoroughly
Run your test suite and verify all error paths work correctly.

---

## 🎯 Upgrade Path

### From Prototype → v0.1.0
**Estimated Time:** 1-2 hours for small projects, 4-8 hours for large projects

**Steps:**
1. Update all constructor calls to handle errors
2. Replace direct field access with getter methods
3. Update ValidateMnemonic calls (now package-level)
4. Test all error paths
5. Review and update custom error handling

---

## 📝 Notes

### Why These Changes?
1. **Production Ready**: Proper error handling is critical for production use
2. **Go Idioms**: Following Go best practices and community standards
3. **Testability**: Interface-based design enables comprehensive testing
4. **Maintainability**: Clear structure and encapsulation ease future changes
5. **Safety**: Input validation prevents runtime panics

### Future Plans
- v0.2.0: Enhanced query capabilities
- v0.3.0: Connection pooling and retry logic
- v0.4.0: Observability and metrics

---

## 🙏 Acknowledgments

Special thanks to the Six Protocol team and all contributors who made this refactoring possible.

---

## 📞 Support

- **Documentation**: See REFACTORING_README.md
- **Examples**: See example/EXAMPLES.md
- **Issues**: Please report any issues on GitHub
- **Questions**: Contact the development team

---

**Full Changelog**: Initial Release → v0.1.0

**Download**: [GitHub Releases](https://github.com/thesixnetwork/lbb-sdk-go/releases/tag/v0.1.0)