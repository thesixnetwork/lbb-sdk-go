# Audit Fixes - Code Examples

This document provides specific code changes to fix all HIGH and MEDIUM severity issues found in the audit.

---

## 🔴 HIGH SEVERITY FIXES

### H-1: Fix Nonce Increment Timing in Permit Functions

**Affected Files:**
- `contracts/src/CertAutoID.sol` (lines 127-143)
- `contracts/src/Cert.sol` (lines 118-134)

**Current Code (VULNERABLE):**

```solidity
function permit(
    address owner,
    address spender,
    uint256 tokenId,
    uint256 deadline,
    uint8 v,
    bytes32 r,
    bytes32 s
) public {
    if (block.timestamp > deadline) {
        revert SignatureExpired();
    }

    bytes32 structHash = keccak256(
        abi.encode(
            PERMIT_TYPEHASH,
            owner,
            spender,
            tokenId,
            _nonces[owner]++,  // ❌ PROBLEM: Incremented before validation
            deadline
        )
    );

    bytes32 hash = _hashTypedDataV4(structHash);
    address signer = ECDSA.recover(hash, v, r, s);

    if (signer != owner) {
        revert InvalidSigner();  // Nonce already consumed!
    }

    if (ownerOf(tokenId) != owner) {
        revert InvalidSigner();
    }

    _approve(spender, tokenId, owner);
    emit PermitUsed(owner, spender, tokenId);
}
```

**Fixed Code:**

```solidity
function permit(
    address owner,
    address spender,
    uint256 tokenId,
    uint256 deadline,
    uint8 v,
    bytes32 r,
    bytes32 s
) public {
    if (block.timestamp > deadline) {
        revert SignatureExpired();
    }

    uint256 currentNonce = _nonces[owner];  // ✅ Read without incrementing

    bytes32 structHash = keccak256(
        abi.encode(
            PERMIT_TYPEHASH,
            owner,
            spender,
            tokenId,
            currentNonce,  // ✅ Use cached value
            deadline
        )
    );

    bytes32 hash = _hashTypedDataV4(structHash);
    address signer = ECDSA.recover(hash, v, r, s);

    if (signer != owner) {
        revert InvalidSigner();  // Nonce not consumed if invalid
    }

    if (ownerOf(tokenId) != owner) {
        revert InvalidSigner();
    }

    _nonces[owner]++;  // ✅ Only increment after all validations pass

    _approve(spender, tokenId, owner);
    emit PermitUsed(owner, spender, tokenId);
}
```

**Apply Same Fix to `permitForAll`:**

```solidity
function permitForAll(
    address owner,
    address operator,
    bool approved,
    uint256 deadline,
    uint8 v,
    bytes32 r,
    bytes32 s
) public {
    if (block.timestamp > deadline) {
        revert SignatureExpired();
    }

    uint256 currentNonce = _nonces[owner];  // ✅ Read without incrementing

    bytes32 structHash = keccak256(
        abi.encode(
            PERMIT_FOR_ALL_TYPEHASH,
            owner,
            operator,
            approved,
            currentNonce,  // ✅ Use cached value
            deadline
        )
    );

    bytes32 hash = _hashTypedDataV4(structHash);
    address signer = ECDSA.recover(hash, v, r, s);

    if (signer != owner) {
        revert InvalidSigner();
    }

    _nonces[owner]++;  // ✅ Only increment after validation passes

    _setApprovalForAll(owner, operator, approved);
    emit PermitForAllUsed(owner, operator, approved);
}
```

---

### H-2: Add Token Existence Check in Cert.sol

**Affected File:** `contracts/src/Cert.sol` (line 44)

**Current Code:**

```solidity
function safeMint(address to, uint256 tokenId) public onlyOwner {
    _safeMint(to, tokenId);
}
```

**Fixed Code:**

```solidity
function safeMint(address to, uint256 tokenId) public onlyOwner {
    if (_ownerOf(tokenId) != address(0)) {
        revert TokenAlreadyMinted();
    }
    _safeMint(to, tokenId);
    emit safeMintEvent(to, tokenId);  // Also add missing event emission
}
```

**Add Custom Error at Top of Contract:**

```solidity
error NonExistentTokenURI();
error InvalidSignature();
error SignatureExpired();
error InvalidSigner();
error TokenAlreadyMinted();  // ✅ Add this new error
```

---

## 🟡 MEDIUM SEVERITY FIXES

### M-1: Validate Spender/Operator is Not Zero Address

**Affected Files:** Both contracts, `permit` and `permitForAll` functions

**Add to `permit` function (after deadline check):**

```solidity
function permit(
    address owner,
    address spender,
    uint256 tokenId,
    uint256 deadline,
    uint8 v,
    bytes32 r,
    bytes32 s
) public {
    if (block.timestamp > deadline) {
        revert SignatureExpired();
    }
    
    // ✅ Add this validation
    if (spender == address(0)) {
        revert InvalidSpender();
    }

    // ... rest of function
}
```

**Add to `permitForAll` function:**

```solidity
function permitForAll(
    address owner,
    address operator,
    bool approved,
    uint256 deadline,
    uint8 v,
    bytes32 r,
    bytes32 s
) public {
    if (block.timestamp > deadline) {
        revert SignatureExpired();
    }
    
    // ✅ Add this validation
    if (operator == address(0)) {
        revert InvalidOperator();
    }

    // ... rest of function
}
```

**Add Custom Errors:**

```solidity
error NonExistentTokenURI();
error InvalidSignature();
error SignatureExpired();
error InvalidSigner();
error TokenAlreadyMinted();
error InvalidSpender();    // ✅ Add this
error InvalidOperator();   // ✅ Add this
```

---

### M-2: Add Missing Event Emission in Cert.sol

**Affected File:** `contracts/src/Cert.sol` (line 44)

**Current Code:**

```solidity
function safeMint(address to, uint256 tokenId) public onlyOwner {
    _safeMint(to, tokenId);
}
```

**Fixed Code:**

```solidity
function safeMint(address to, uint256 tokenId) public onlyOwner {
    if (_ownerOf(tokenId) != address(0)) {
        revert TokenAlreadyMinted();
    }
    _safeMint(to, tokenId);
    emit safeMintEvent(to, tokenId);  // ✅ Add event emission
}
```

---

### M-3: Improve TransferWithPermit and BurnWithPermit Flow

**Current Issue:** If transfer/burn fails, permit is still consumed.

**Option 1: Validate Before Permit (Recommended)**

```solidity
function transferWithPermit(
    address from,
    address to,
    uint256 tokenId,
    uint256 deadline,
    uint8 v,
    bytes32 r,
    bytes32 s
) public {
    // ✅ Validate the transfer can happen BEFORE consuming permit
    if (ownerOf(tokenId) != from) {
        revert InvalidOwner();
    }
    if (to == address(0)) {
        revert InvalidRecipient();
    }
    
    // Now consume the permit
    permit(from, msg.sender, tokenId, deadline, v, r, s);
    
    // Execute transfer
    safeTransferFrom(from, to, tokenId);
}
```

**Apply Same Pattern to BurnWithPermit:**

```solidity
function burnWithPermit(
    address from,
    uint256 tokenId,
    uint256 deadline,
    uint8 v,
    bytes32 r,
    bytes32 s
) public {
    // ✅ Validate burn can happen BEFORE consuming permit
    if (ownerOf(tokenId) != from) {
        revert InvalidOwner();
    }

    // Now consume the permit
    permit(from, msg.sender, tokenId, deadline, v, r, s);

    // Execute burn
    burn(tokenId);
}
```

**Add Custom Error:**

```solidity
error InvalidOwner();      // ✅ Add this
error InvalidRecipient();  // ✅ Add this
```

**Option 2: Inline Permit Logic (Alternative)**

If you want more control, you can inline the permit validation:

```solidity
function transferWithPermit(
    address from,
    address to,
    uint256 tokenId,
    uint256 deadline,
    uint8 v,
    bytes32 r,
    bytes32 s
) public {
    // Validate deadline
    if (block.timestamp > deadline) {
        revert SignatureExpired();
    }
    
    // Validate addresses
    if (to == address(0)) {
        revert InvalidRecipient();
    }
    if (ownerOf(tokenId) != from) {
        revert InvalidOwner();
    }
    
    // Validate signature
    uint256 currentNonce = _nonces[from];
    bytes32 structHash = keccak256(
        abi.encode(
            PERMIT_TYPEHASH,
            from,
            msg.sender,
            tokenId,
            currentNonce,
            deadline
        )
    );
    bytes32 hash = _hashTypedDataV4(structHash);
    address signer = ECDSA.recover(hash, v, r, s);
    
    if (signer != from) {
        revert InvalidSigner();
    }
    
    // Only increment nonce after all validations
    _nonces[from]++;
    
    // Approve and transfer
    _approve(msg.sender, tokenId, from);
    emit PermitUsed(from, msg.sender, tokenId);
    
    safeTransferFrom(from, to, tokenId);
}
```

---

### M-4: Standardize Event Naming

**Option 1: Remove Custom Event (Use Standard Transfer)**

Remove the custom `safeMintEvent` and rely on OpenZeppelin's standard `Transfer` event:

```solidity
// Remove this event declaration
// event safeMintEvent(address to, uint256 tokenId);

function safeMint(address to, uint256 tokenId) public onlyOwner {
    if (_ownerOf(tokenId) != address(0)) {
        revert TokenAlreadyMinted();
    }
    _safeMint(to, tokenId);
    // No custom event needed - Transfer event is automatically emitted
}
```

**Option 2: Rename to More Descriptive Name**

```solidity
// Replace safeMintEvent with more descriptive name
event CertificateMinted(address indexed to, uint256 indexed tokenId, uint256 timestamp);

function safeMint(address to, uint256 tokenId) public onlyOwner {
    if (_ownerOf(tokenId) != address(0)) {
        revert TokenAlreadyMinted();
    }
    _safeMint(to, tokenId);
    emit CertificateMinted(to, tokenId, block.timestamp);
}
```

---

### M-5: Add Max Supply Limit to CertAutoID

**Affected File:** `contracts/src/CertAutoID.sol`

**Add Constant:**

```solidity
contract LBBCertAutoID is ERC721, ERC721Enumerable, ERC721Burnable, Ownable, EIP712 {
    using Strings for uint256;
    string private _baseTokenURI;
    uint256 private _nextTokenId;
    
    uint256 private constant MAX_SUPPLY = 1_000_000;  // ✅ Add max supply

    // ... rest of contract
}
```

**Update safeMint Function:**

```solidity
function safeMint(address to) public onlyOwner returns (uint256) {
    if (_nextTokenId >= MAX_SUPPLY) {
        revert MaxSupplyReached();
    }
    uint256 tokenId = _nextTokenId++;
    _safeMint(to, tokenId);
    emit safeMintEvent(to, tokenId);
    return tokenId;
}
```

**Add Custom Error:**

```solidity
error NonExistentTokenURI();
error InvalidSignature();
error SignatureExpired();
error InvalidSigner();
error MaxSupplyReached();  // ✅ Add this
```

**Make Max Supply Configurable (Optional):**

```solidity
uint256 private immutable MAX_SUPPLY;

constructor(
    string memory name,
    string memory symbol,
    string memory baseURI,
    address initialOwner,
    uint256 maxSupply  // ✅ Add parameter
) ERC721(name, symbol) EIP712(name, "1") Ownable(initialOwner) {
    _baseTokenURI = baseURI;
    _nextTokenId = 1;
    MAX_SUPPLY = maxSupply > 0 ? maxSupply : type(uint256).max;  // 0 = unlimited
}
```

---

## 🟢 LOW SEVERITY FIXES

### L-1: Rename CertAutoID Contract

**Affected File:** `contracts/src/CertAutoID.sol`

**Current:**

```solidity
contract LBBCert is ERC721, ERC721Enumerable, ERC721Burnable, Ownable, EIP712 {
```

**Fixed:**

```solidity
contract LBBCertAutoID is ERC721, ERC721Enumerable, ERC721Burnable, Ownable, EIP712 {
```

---

### L-2: Add NatSpec Documentation

**Example for `permit` function:**

```solidity
/**
 * @notice Allows token owner to approve spender via EIP-712 signature
 * @dev Implements EIP-2612 style permit for ERC-721
 * @param owner The current owner of the token
 * @param spender The address being approved
 * @param tokenId The token ID to approve
 * @param deadline Unix timestamp after which signature expires
 * @param v Recovery byte of the signature
 * @param r First 32 bytes of signature
 * @param s Second 32 bytes of signature
 * @custom:security Validates signature and increments nonce only after validation
 * @custom:emits PermitUsed
 */
function permit(
    address owner,
    address spender,
    uint256 tokenId,
    uint256 deadline,
    uint8 v,
    bytes32 r,
    bytes32 s
) public {
    // ... implementation
}
```

---

### L-3: Add BaseURI Change Event

**Add Event Declaration:**

```solidity
event BaseURIUpdated(string indexed oldBaseURI, string indexed newBaseURI);
```

**Update setBaseURI Function:**

```solidity
function setBaseURI(string calldata baseURI) external onlyOwner {
    string memory oldURI = _baseTokenURI;
    _baseTokenURI = baseURI;
    emit BaseURIUpdated(oldURI, baseURI);
}
```

---

### L-4: Lock OpenZeppelin Version

**In Contract Files:**

**Current:**

```solidity
import {ERC721} from "openzeppelin-contracts/token/ERC721/ERC721.sol";
```

**Fixed:**

Update your package manager to lock to exact version, or use git submodule with specific commit.

In `foundry.toml`:

```toml
[dependencies]
openzeppelin-contracts = { version = "5.5.0" }  # Exact version, no ^
```

---

### L-5: Standardize Custom Errors

**Create Error Section at Top:**

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {ERC721} from "openzeppelin-contracts/token/ERC721/ERC721.sol";
// ... other imports

/**
 * @dev Custom errors for gas efficiency
 */
error NonExistentTokenURI();
error InvalidSignature();
error SignatureExpired();
error InvalidSigner();
error TokenAlreadyMinted();
error InvalidSpender();
error InvalidOperator();
error InvalidOwner();
error InvalidRecipient();
error MaxSupplyReached();

contract LBBCert is ERC721, ERC721Enumerable, ERC721Burnable, Ownable, EIP712 {
    // ... contract code
}
```

---

## 🔧 SDK FIXES

### SDK-M-1: Add Private Key Security Warnings

**File:** `account/account.go`

**Update GetMnemonic Function:**

```go
// GetMnemonic returns the mnemonic phrase (use with extreme caution)
// WARNING: This exposes sensitive cryptographic material. Only use this for:
// - Backup purposes in secure environments
// - Migration to other wallets
// - Testing and development
// Never log, transmit, or store the result in plain text.
func (a *Account) GetMnemonic() string {
    // Log warning when accessed
    log.Println("WARNING: Mnemonic accessed. Ensure proper security measures.")
    return a.mnemonic
}
```

**Better Approach - Require Confirmation:**

```go
// GetMnemonic returns the mnemonic phrase with explicit confirmation required
// The confirm parameter must be the string "I understand the security risks"
func (a *Account) GetMnemonic(confirm string) (string, error) {
    if confirm != "I understand the security risks" {
        return "", fmt.Errorf("must explicitly confirm security risks to access mnemonic")
    }
    
    log.Printf("WARNING: Mnemonic accessed for account %s at %s", 
        a.accountName, time.Now().Format(time.RFC3339))
    
    return a.mnemonic, nil
}
```

---

### SDK-M-2: Add Context Timeouts

**File:** `client/client.go`

**Update WaitForTransaction:**

```go
// WaitForTransaction waits for a Cosmos transaction to be mined
func (c *Client) WaitForTransaction(txHash string) error {
    if txHash == "" {
        return fmt.Errorf("transaction hash cannot be empty")
    }

    fmt.Printf("Waiting for transaction %s to be mined...\n", txHash)

    // ✅ Create context with timeout
    ctx, cancel := context.WithTimeout(c.ctx, transactionTimeout)
    defer cancel()

    ticker := time.NewTicker(transactionPollInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():  // ✅ Use context timeout
            return fmt.Errorf("timeout waiting for transaction %s: %w", txHash, ctx.Err())
        case <-ticker.C:
            output, err := authtx.QueryTx(c.cosmosClientCTX, txHash)
            if err != nil {
                continue
            }

            if output.Empty() {
                return fmt.Errorf("no transaction found with hash %s", txHash)
            }

            if output.Code != 0 {
                return fmt.Errorf("transaction %s failed with code %d: %s", 
                    txHash, output.Code, output.RawLog)
            }

            fmt.Printf("Transaction %s successfully mined in block %d\n", 
                txHash, output.Height)
            return nil
        }
    }
}
```

---

### SDK-M-3: Make Gas Buffer Configurable

**File:** `pkg/evm/ethclient.go`

**Add Configuration:**

```go
const (
    DefaultGasBuffer = 120  // 20% buffer
    MinGasBuffer     = 100  // No buffer
    MaxGasBuffer     = 200  // 100% buffer
)

type EVMClient struct {
    account.Account
    gasBuffer uint64  // ✅ Add configurable gas buffer
}

func NewEVMClient(a account.Account) *EVMClient {
    return &EVMClient{
        Account:   a,
        gasBuffer: DefaultGasBuffer,  // ✅ Initialize with default
    }
}

// ✅ Add setter for gas buffer
func (e *EVMClient) SetGasBuffer(buffer uint64) error {
    if buffer < MinGasBuffer || buffer > MaxGasBuffer {
        return fmt.Errorf("gas buffer must be between %d and %d", MinGasBuffer, MaxGasBuffer)
    }
    e.gasBuffer = buffer
    return nil
}
```

**Update GasLimit Function:**

```go
func (e *EVMClient) GasLimit(callMsg ethereum.CallMsg) (uint64, error) {
    goCtx := e.GetClient().GetContext()
    ethClient := e.GetClient().GetETHClient()

    gasLimit, err := ethClient.EstimateGas(goCtx, callMsg)
    if err != nil {
        fmt.Printf("ERROR EstimateGas : %v \n", err)
        return gasLimit, err
    }
    
    // ✅ Use configurable gas buffer
    gasLimit = gasLimit * e.gasBuffer / 100
    return gasLimit, nil
}
```

---

### SDK-L-3: Add Input Validation

**Example for MintCertificateNFT:**

```go
func (e *EVMClient) MintCertificateNFT(
    contractAddress common.Address,
    tokenID uint64,
) (tx *types.Transaction, err error) {
    // ✅ Add validation
    if contractAddress == (common.Address{}) {
        return nil, fmt.Errorf("contract address cannot be zero address")
    }
    
    if tokenID == 0 {
        return nil, fmt.Errorf("token ID must be greater than 0")
    }

    goCtx := e.GetClient().GetContext()
    ethClient := e.GetClient().GetETHClient()

    // ... rest of implementation
}
```

---

## 📋 Testing Examples

### Test for Fixed Nonce Issue

```solidity
// Test file: contracts/test/Cert.t.sol

function testPermitWithInvalidSignatureDoesNotConsumeNonce() public {
    // Setup
    address owner = address(1);
    address spender = address(2);
    uint256 tokenId = 1;
    
    // Mint token to owner
    vm.prank(owner);
    myNft.safeMint(owner, tokenId);
    
    // Get initial nonce
    uint256 nonceBefore = myNft.nonces(owner);
    
    // Try permit with invalid signature
    vm.expectRevert(LBBCert.InvalidSigner.selector);
    myNft.permit(
        owner,
        spender,
        tokenId,
        block.timestamp + 1 hours,
        27,  // Invalid v
        bytes32(0),  // Invalid r
        bytes32(0)   // Invalid s
    );
    
    // Verify nonce was NOT incremented
    uint256 nonceAfter = myNft.nonces(owner);
    assertEq(nonceBefore, nonceAfter, "Nonce should not change on invalid signature");
}
```

### Test for Token Existence Check

```solidity
function testCannotMintDuplicateTokenId() public {
    address owner = address(1);
    uint256 tokenId = 1;
    
    // First mint should succeed
    vm.prank(owner);
    myNft.safeMint(address(2), tokenId);
    
    // Second mint with same ID should fail
    vm.prank(owner);
    vm.expectRevert(LBBCert.TokenAlreadyMinted.selector);
    myNft.safeMint(address(3), tokenId);
}
```

---

## 🚀 Deployment Script Updates

After applying fixes, update deployment scripts:

```javascript
// script/Deploy.s.sol

contract DeployScript is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        
        vm.startBroadcast(deployerPrivateKey);
        
        // Deploy with max supply parameter
        LBBCertAutoID cert = new LBBCertAutoID(
            "My Certificate",
            "CERT",
            "ipfs://base-uri/",
            msg.sender,
            1_000_000  // Max supply
        );
        
        console.log("Deployed LBBCertAutoID at:", address(cert));
        
        vm.stopBroadcast();
    }
}
```

---

## ✅ Verification Checklist

After applying all fixes:

- [ ] All HIGH severity issues fixed
- [ ] All MEDIUM severity issues fixed
- [ ] Tests added for all fixes
- [ ] Test coverage > 90%
- [ ] All tests passing
- [ ] Gas benchmarks run
- [ ] Documentation updated
- [ ] Deployment scripts updated
- [ ] Testnet deployment successful
- [ ] Code review completed

---

**Note:** Always test thoroughly on testnet before mainnet deployment!