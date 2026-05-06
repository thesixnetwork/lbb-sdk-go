// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {ERC721} from "openzeppelin-contracts/token/ERC721/ERC721.sol";
import {ERC721Enumerable} from "openzeppelin-contracts/token/ERC721/extensions/ERC721Enumerable.sol";
import {ERC721Burnable} from "openzeppelin-contracts/token/ERC721/extensions/ERC721Burnable.sol";
import {Strings} from "openzeppelin-contracts/utils/Strings.sol";
import {ECDSA} from "openzeppelin-contracts/utils/cryptography/ECDSA.sol";
import {EIP712} from "openzeppelin-contracts/utils/cryptography/EIP712.sol";
import {ReentrancyGuard} from "openzeppelin-contracts/utils/ReentrancyGuard.sol";

// ============ Custom Errors ============
error NotSuperAdmin();
error NotAdmin();
error AdminAlreadyExists(address account);
error AdminNotFound(address account);
error ZeroAddressNotAllowed();
error SignatureExpired();
error InvalidSigner();
error InvalidTokenOwner();

/**
 * @title NftFactory
 * @dev ERC721 NFT contract with two-tier role management:
 *      - SuperAdmin: deploys the contract and can grant/revoke admin roles.
 *                   Has all admin privileges plus role management.
 *      - Admin:     can mint, burn, transfer tokens and update metadata.
 *
 * Also supports EIP-2612 style gasless permit operations (inherited from LBBCert design).
 *
 * Security notes:
 *  - Admin list uses swap-and-pop for O(1) removal without gaps.
 *  - adminBurn / adminTransfer bypass token-owner approval via internal _update(auth=0).
 *  - ReentrancyGuard is applied to all minting and transfer paths.
 */
contract NftFactory is ERC721, ERC721Enumerable, ERC721Burnable, EIP712, ReentrancyGuard {
    using Strings for uint256;

    // ============ State Variables ============

    address private _superAdmin;
    string private _baseTokenURI;

    mapping(address => bool) private _isAdminMap;
    address[] private _adminList;
    mapping(address => uint256) private _adminListIndex; // 0-based index into _adminList

    mapping(address => uint256) private _nonces;

    // ============ EIP-712 Type Hashes ============

    bytes32 private constant PERMIT_TYPEHASH =
        keccak256("Permit(address owner,address spender,uint256 tokenId,uint256 nonce,uint256 deadline)");

    bytes32 private constant PERMIT_FOR_ALL_TYPEHASH =
        keccak256("PermitForAll(address owner,address operator,bool approved,uint256 nonce,uint256 deadline)");

    // ============ Events ============

    event AdminGranted(address indexed account, address indexed grantedBy);
    event AdminRevoked(address indexed account, address indexed revokedBy);
    event SuperAdminTransferred(address indexed previousSuperAdmin, address indexed newSuperAdmin);
    event TokenMinted(address indexed to, uint256 indexed tokenId, address indexed mintedBy);
    event TokenBatchMinted(address indexed to, uint256[] tokenIds, address indexed mintedBy);
    event TokenBurned(uint256 indexed tokenId, address indexed burnedBy);
    event AdminTransfer(address indexed from, address indexed to, uint256 indexed tokenId, address transferredBy);
    event BaseURIUpdated(string newBaseURI, address indexed updatedBy);
    event PermitUsed(address indexed owner, address indexed spender, uint256 indexed tokenId);
    event PermitForAllUsed(address indexed owner, address indexed operator, bool approved);

    // ============ Modifiers ============

    /// @dev Only the super admin may call.
    modifier onlySuperAdmin() {
        if (msg.sender != _superAdmin) revert NotSuperAdmin();
        _;
    }

    /// @dev Super admin and any granted admin may call.
    modifier onlyAdmin() {
        if (msg.sender != _superAdmin && !_isAdminMap[msg.sender]) revert NotAdmin();
        _;
    }

    // ============ Constructor ============

    /**
     * @param name       Token collection name (also used in EIP-712 domain separator).
     * @param symbol     Token collection symbol.
     * @param baseURI    Initial base URI for token metadata.
     * @param superAdmin Address that receives the super admin role on deployment.
     */
    constructor(
        string memory name,
        string memory symbol,
        string memory baseURI,
        address superAdmin
    ) ERC721(name, symbol) EIP712(name, "1") {
        if (superAdmin == address(0)) revert ZeroAddressNotAllowed();
        _superAdmin = superAdmin;
        _baseTokenURI = baseURI;
    }

    // ============ Super Admin — Role Management ============

    /**
     * @dev Transfer the super admin role to a new address.
     *      The caller immediately loses the role; there is no two-step confirmation.
     *      Do NOT call this with an address you do not control.
     */
    function transferSuperAdmin(address newSuperAdmin) external onlySuperAdmin {
        if (newSuperAdmin == address(0)) revert ZeroAddressNotAllowed();
        address previous = _superAdmin;
        _superAdmin = newSuperAdmin;
        emit SuperAdminTransferred(previous, newSuperAdmin);
    }

    /**
     * @dev Grant admin role to `account`. Only super admin can call.
     *      Reverts if `account` is already an admin.
     */
    function grantAdmin(address account) external onlySuperAdmin {
        if (account == address(0)) revert ZeroAddressNotAllowed();
        if (_isAdminMap[account]) revert AdminAlreadyExists(account);
        _isAdminMap[account] = true;
        _adminListIndex[account] = _adminList.length;
        _adminList.push(account);
        emit AdminGranted(account, msg.sender);
    }

    /**
     * @dev Revoke admin role from `account`. Only super admin can call.
     *      Uses swap-and-pop to maintain a compact list with O(1) removal.
     */
    function revokeAdmin(address account) external onlySuperAdmin {
        if (!_isAdminMap[account]) revert AdminNotFound(account);
        _isAdminMap[account] = false;

        uint256 index = _adminListIndex[account];
        uint256 lastIndex = _adminList.length - 1;

        if (index != lastIndex) {
            address lastAdmin = _adminList[lastIndex];
            _adminList[index] = lastAdmin;
            _adminListIndex[lastAdmin] = index;
        }

        _adminList.pop();
        delete _adminListIndex[account];

        emit AdminRevoked(account, msg.sender);
    }

    // ============ Read — Role Queries ============

    /**
     * @dev Returns the current super admin address.
     */
    function getSuperAdmin() external view returns (address) {
        return _superAdmin;
    }

    /**
     * @dev Returns true if `account` holds admin or super admin privileges.
     */
    function isAdmin(address account) external view returns (bool) {
        return account == _superAdmin || _isAdminMap[account];
    }

    /**
     * @dev Returns the full list of granted admin addresses.
     *      Does NOT include the super admin address.
     */
    function listAdmins() external view returns (address[] memory) {
        return _adminList;
    }

    /**
     * @dev Returns how many admin addresses have been granted (excludes super admin).
     */
    function getAdminCount() external view returns (uint256) {
        return _adminList.length;
    }

    // ============ Admin — NFT Operations ============

    /**
     * @dev Mint a single token to `to`. Callable by any admin or super admin.
     */
    function safeMint(address to, uint256 tokenId) external onlyAdmin nonReentrant {
        if (to == address(0)) revert ZeroAddressNotAllowed();
        _safeMint(to, tokenId);
        emit TokenMinted(to, tokenId, msg.sender);
    }

    /**
     * @dev Mint multiple tokens to `to` in one transaction.
     *      Callable by any admin or super admin.
     */
    function safeMintBatch(address to, uint256[] calldata tokenIds) external onlyAdmin nonReentrant {
        if (to == address(0)) revert ZeroAddressNotAllowed();
        uint256 len = tokenIds.length;
        for (uint256 i = 0; i < len; ) {
            _safeMint(to, tokenIds[i]);
            unchecked { ++i; }
        }
        emit TokenBatchMinted(to, tokenIds, msg.sender);
    }

    /**
     * @dev Admin forced burn — destroys a token without needing owner approval.
     *      Bypasses the standard approval check by passing auth=address(0) to _update.
     *      Callable by any admin or super admin.
     */
    function adminBurn(uint256 tokenId) external onlyAdmin {
        _requireOwned(tokenId);
        // Passing address(0) as auth skips the ERC721 authorization check
        _update(address(0), tokenId, address(0));
        emit TokenBurned(tokenId, msg.sender);
    }

    /**
     * @dev Admin forced transfer — moves a token without owner approval.
     *      Bypasses the standard approval check by passing auth=address(0) to _update.
     *      Callable by any admin or super admin.
     */
    function adminTransfer(address from, address to, uint256 tokenId) external onlyAdmin nonReentrant {
        if (to == address(0)) revert ZeroAddressNotAllowed();
        address currentOwner = _requireOwned(tokenId);
        if (currentOwner != from) revert InvalidTokenOwner();
        // Passing address(0) as auth skips the ERC721 authorization check
        _update(to, tokenId, address(0));
        emit AdminTransfer(from, to, tokenId, msg.sender);
    }

    /**
     * @dev Update the base metadata URI for all tokens.
     *      Callable by any admin or super admin.
     */
    function setBaseURI(string calldata baseURI) external onlyAdmin {
        _baseTokenURI = baseURI;
        emit BaseURIUpdated(baseURI, msg.sender);
    }

    // ============ Read — Token Metadata ============

    function _baseURI() internal view virtual override returns (string memory) {
        return _baseTokenURI;
    }

    /**
     * @dev Returns the full metadata URI for `tokenId`.
     *      Reverts with ERC721NonexistentToken if the token does not exist.
     */
    function tokenURI(uint256 tokenId) public view virtual override returns (string memory) {
        _requireOwned(tokenId);
        return bytes(_baseTokenURI).length > 0
            ? string(abi.encodePacked(_baseTokenURI, tokenId.toString()))
            : "";
    }

    // ============ EIP-2612 Style Permit Functions ============

    /**
     * @dev Returns the current nonce for `owner`.
     *      Must be included when building an EIP-712 permit signature.
     */
    function nonces(address owner) public view returns (uint256) {
        return _nonces[owner];
    }

    /**
     * @dev Returns the EIP-712 domain separator for this contract and chain.
     */
    function DOMAIN_SEPARATOR() external view returns (bytes32) {
        return _domainSeparatorV4();
    }

    /**
     * @dev EIP-2612 style permit: approve `spender` for a specific `tokenId`
     *      using an off-chain EIP-712 signature from `owner`.
     *      Consumes one nonce from `owner`.
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
        if (block.timestamp > deadline) revert SignatureExpired();

        uint256 currentNonce = _nonces[owner];
        bytes32 structHash = keccak256(
            abi.encode(PERMIT_TYPEHASH, owner, spender, tokenId, currentNonce, deadline)
        );

        address signer = ECDSA.recover(_hashTypedDataV4(structHash), v, r, s);
        if (signer != owner) revert InvalidSigner();
        if (ownerOf(tokenId) != owner) revert InvalidSigner();

        _approve(spender, tokenId, owner);
        _nonces[owner] = currentNonce + 1;
        emit PermitUsed(owner, spender, tokenId);
    }

    /**
     * @dev EIP-712 style permitForAll: set or unset operator approval for all tokens
     *      owned by `owner`, using an off-chain signature.
     *      Consumes one nonce from `owner`.
     */
    function permitForAll(
        address owner,
        address operator,
        bool approved,
        uint256 deadline,
        uint8 v,
        bytes32 r,
        bytes32 s
    ) public {
        if (block.timestamp > deadline) revert SignatureExpired();

        uint256 currentNonce = _nonces[owner];
        bytes32 structHash = keccak256(
            abi.encode(PERMIT_FOR_ALL_TYPEHASH, owner, operator, approved, currentNonce, deadline)
        );

        address signer = ECDSA.recover(_hashTypedDataV4(structHash), v, r, s);
        if (signer != owner) revert InvalidSigner();

        _setApprovalForAll(owner, operator, approved);
        _nonces[owner] = currentNonce + 1;
        emit PermitForAllUsed(owner, operator, approved);
    }

    /**
     * @dev Gasless transfer: `owner` signs a permit off-chain; the caller (relayer or
     *      anyone) submits this transaction and pays the gas.
     */
    function transferWithPermit(
        address from,
        address to,
        uint256 tokenId,
        uint256 deadline,
        uint8 v,
        bytes32 r,
        bytes32 s
    ) public nonReentrant {
        if (to == address(0)) revert ZeroAddressNotAllowed();
        if (ownerOf(tokenId) != from) revert InvalidTokenOwner();
        permit(from, msg.sender, tokenId, deadline, v, r, s);
        safeTransferFrom(from, to, tokenId);
    }

    /**
     * @dev Gasless burn: `owner` signs a permit off-chain; the caller (relayer or
     *      anyone) submits this transaction and pays the gas.
     */
    function burnWithPermit(
        address from,
        uint256 tokenId,
        uint256 deadline,
        uint8 v,
        bytes32 r,
        bytes32 s
    ) public nonReentrant {
        if (ownerOf(tokenId) != from) revert InvalidTokenOwner();
        permit(from, msg.sender, tokenId, deadline, v, r, s);
        burn(tokenId);
    }

    // ============ Required Overrides ============

    function _update(address to, uint256 tokenId, address auth)
        internal override(ERC721, ERC721Enumerable) returns (address)
    {
        return super._update(to, tokenId, auth);
    }

    function _increaseBalance(address account, uint128 value)
        internal override(ERC721, ERC721Enumerable)
    {
        super._increaseBalance(account, value);
    }

    function supportsInterface(bytes4 interfaceId)
        public view override(ERC721, ERC721Enumerable) returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }
}
