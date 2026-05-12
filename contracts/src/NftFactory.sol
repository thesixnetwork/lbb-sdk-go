// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {ERC721} from "openzeppelin-contracts/token/ERC721/ERC721.sol";
import {ERC721Enumerable} from "openzeppelin-contracts/token/ERC721/extensions/ERC721Enumerable.sol";
import {ERC721Burnable} from "openzeppelin-contracts/token/ERC721/extensions/ERC721Burnable.sol";
import {Strings} from "openzeppelin-contracts/utils/Strings.sol";
import {Base64} from "openzeppelin-contracts/utils/Base64.sol";
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
 * @dev ERC721 NFT contract with two-tier role management and backend-signature-authorised minting.
 *
 * Roles:
 *  - SuperAdmin: set at deployment; manages admin list and the mint signer address.
 *  - Admin:      relayer wallets that submit mint transactions and pay gas.
 *
 * Minting model (EIP-712):
 *  - The backend wallet (_mintSigner) signs a Mint / MintBatch struct off-chain.
 *  - An admin submits the transaction on-chain with that signature.
 *  - Both conditions must hold: caller is an admin AND signature comes from _mintSigner.
 *
 * Burn / Transfer:
 *  - Token owners initiate burns and transfers via standard ERC721 or EIP-712 permit flows.
 *  - No admin-forced overrides exist; all actions require the token owner's authorisation.
 *
 * Security notes:
 *  - Admin list uses swap-and-pop for O(1) removal without gaps.
 *  - ReentrancyGuard is applied to all minting and transfer paths.
 *  - Nonces are shared across Mint and Permit operations (per recipient / owner address).
 */
contract NftFactory is ERC721, ERC721Enumerable, ERC721Burnable, EIP712, ReentrancyGuard {
    using Strings for uint256;

    // ============ State Variables ============

    address private _superAdmin;
    address private _mintSigner;   // backend wallet that authorises every mint

    mapping(address => bool) private _isAdminMap;
    address[] private _adminList;
    mapping(address => uint256) private _adminListIndex; // 0-based index into _adminList

    mapping(address => uint256) private _nonces;

    // ============ On-chain Metadata ============

    struct TokenData {
        string name;
        string certID;
        uint256 createdAt;
    }

    mapping(uint256 => TokenData) private _tokenMetadata;

    // ============ EIP-712 Type Hashes ============

    bytes32 private constant PERMIT_TYPEHASH =
        keccak256("Permit(address owner,address spender,uint256 tokenId,uint256 nonce,uint256 deadline)");

    bytes32 private constant PERMIT_FOR_ALL_TYPEHASH =
        keccak256("PermitForAll(address owner,address operator,bool approved,uint256 nonce,uint256 deadline)");

    bytes32 private constant MINT_TYPEHASH =
        keccak256("Mint(address to,uint256 tokenId,string name,string certID,uint256 createdAt,uint256 nonce,uint256 deadline)");

    // ============ Events ============

    event AdminGranted(address indexed account, address indexed grantedBy);
    event AdminRevoked(address indexed account, address indexed revokedBy);
    event SuperAdminTransferred(address indexed previousSuperAdmin, address indexed newSuperAdmin);
    event MintSignerUpdated(address indexed previousSigner, address indexed newSigner);
    event TokenMinted(address indexed to, uint256 indexed tokenId, string certID, address indexed mintedBy);
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
     * @param superAdmin Address that receives the super admin role on deployment.
     * @param mintSigner Backend wallet whose EIP-712 signature authorises every mint.
     */
    constructor(
        string memory name,
        string memory symbol,
        address superAdmin,
        address mintSigner
    ) ERC721(name, symbol) EIP712(name, "1") {
        if (superAdmin == address(0)) revert ZeroAddressNotAllowed();
        if (mintSigner == address(0)) revert ZeroAddressNotAllowed();
        _superAdmin = superAdmin;
        _mintSigner = mintSigner;
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
     * @dev Returns the backend wallet currently authorised to sign mint requests.
     */
    function getMintSigner() external view returns (address) {
        return _mintSigner;
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

    // ============ Super Admin — Mint Signer Management ============

    /**
     * @dev Replace the backend mint signer with `newSigner`. Only super admin can call.
     *      Changing the signer immediately invalidates all previously issued off-chain
     *      mint signatures that have not yet been submitted.
     */
    function setMintSigner(address newSigner) external onlySuperAdmin {
        if (newSigner == address(0)) revert ZeroAddressNotAllowed();
        address previous = _mintSigner;
        _mintSigner = newSigner;
        emit MintSignerUpdated(previous, newSigner);
    }

    // ============ Admin — NFT Operations ============

    /**
     * @dev Mint a single token to `to` with on-chain certificate metadata.
     *
     * Dual-key security model:
     *  - `msg.sender` must be an admin or super admin (on-chain gas payer / relayer).
     *  - The signature (`v`, `r`, `s`) must be produced by `_mintSigner` over the
     *    EIP-712 Mint struct: {to, tokenId, name, certID, createdAt, nonce, deadline}.
     *
     * Nonce used is `_nonces[to]`, shared with the permit functions.
     * Nonce is incremented on success to prevent replay attacks.
     */
    function safeMintWithSignature(
        address to,
        uint256 tokenId,
        string calldata name,
        string calldata certID,
        uint256 createdAt,
        uint256 deadline,
        uint8 v,
        bytes32 r,
        bytes32 s
    ) external onlyAdmin nonReentrant {
        if (to == address(0)) revert ZeroAddressNotAllowed();
        if (block.timestamp > deadline) revert SignatureExpired();

        // Scoped block to keep stack depth under the EVM limit.
        {
            uint256 currentNonce = _nonces[to];
            bytes32 structHash = keccak256(
                abi.encode(
                    MINT_TYPEHASH,
                    to,
                    tokenId,
                    keccak256(bytes(name)),
                    keccak256(bytes(certID)),
                    createdAt,
                    currentNonce,
                    deadline
                )
            );
            if (ECDSA.recover(_hashTypedDataV4(structHash), v, r, s) != _mintSigner)
                revert InvalidSigner();
            _nonces[to] = currentNonce + 1;
        }

        _tokenMetadata[tokenId] = TokenData(name, certID, createdAt);
        _safeMint(to, tokenId);
        emit TokenMinted(to, tokenId, certID, msg.sender);
    }

    // ============ Read — Token Metadata ============

    /**
     * @dev Returns the raw on-chain metadata stored for `tokenId`.
     */
    function getTokenData(uint256 tokenId) external view returns (TokenData memory) {
        _requireOwned(tokenId);
        return _tokenMetadata[tokenId];
    }

    /**
     * @dev Returns a Base64-encoded JSON metadata URI for `tokenId`.
     *      Format: data:application/json;base64,<encoded JSON>
     */
    function tokenURI(uint256 tokenId) public view virtual override returns (string memory) {
        _requireOwned(tokenId);
        TokenData memory data = _tokenMetadata[tokenId];

        string memory json = string(abi.encodePacked(
            '{"name":"', data.name,
            '","description":"Certificate NFT","attributes":[',
            '{"trait_type":"Certificate ID","value":"', data.certID, '"},',
            '{"display_type":"date","trait_type":"Created At","value":', Strings.toString(data.createdAt), '}',
            ']}'
        ));

        return string(abi.encodePacked(
            "data:application/json;base64,",
            Base64.encode(bytes(json))
        ));
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
