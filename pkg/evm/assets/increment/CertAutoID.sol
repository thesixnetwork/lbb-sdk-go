// SPDX-License-Identifier: MIT
// Compatible with OpenZeppelin Contracts ^5.5.0
pragma solidity ^0.8.4;

import {ERC721} from "openzeppelin-contracts/token/ERC721/ERC721.sol";
import {ERC721Enumerable} from "openzeppelin-contracts/token/ERC721/extensions/ERC721Enumerable.sol";
import {ERC721Burnable} from "openzeppelin-contracts/token/ERC721/extensions/ERC721Burnable.sol";
import {Strings} from "openzeppelin-contracts/utils/Strings.sol";
import {Ownable} from "openzeppelin-contracts/access/Ownable.sol";
import {ECDSA} from "openzeppelin-contracts/utils/cryptography/ECDSA.sol";
import {EIP712} from "openzeppelin-contracts/utils/cryptography/EIP712.sol";

error NonExistentTokenURI();
error InvalidSignature();
error SignatureExpired();
error InvalidSigner();

contract LBBCert is ERC721, ERC721Enumerable, ERC721Burnable, Ownable, EIP712 {
    using Strings for uint256;
    string private _baseTokenURI;
    uint256 private _nextTokenId;

    // Nonces for permit
    mapping(address => uint256) private _nonces;

    // EIP-712 Type Hashes
    bytes32 private constant PERMIT_TYPEHASH =
        keccak256(
            "Permit(address owner,address spender,uint256 tokenId,uint256 nonce,uint256 deadline)"
        );

    bytes32 private constant PERMIT_FOR_ALL_TYPEHASH =
        keccak256(
            "PermitForAll(address owner,address operator,bool approved,uint256 nonce,uint256 deadline)"
        );

    // EVENT
    event safeMintEvent(address to, uint256 tokenId);
    event PermitUsed(
        address indexed owner,
        address indexed spender,
        uint256 tokenId
    );
    event PermitForAllUsed(
        address indexed owner,
        address indexed operator,
        bool approved
    );

    constructor(
        string memory name,
        string memory symbol,
        string memory baseURI,
        address initialOwner
    ) ERC721(name, symbol) EIP712(name, "1") Ownable(initialOwner) {
        _baseTokenURI = baseURI;
        _nextTokenId = 1;
    }

    function safeMint(address to) public onlyOwner returns (uint256) {
        uint256 tokenId = _nextTokenId++;
        _safeMint(to, tokenId);
        return tokenId;
    }

    function nextTokenId() public view virtual returns (uint256) {
        return _nextTokenId;
    }

    // BASE URI
    function _baseURI() internal view virtual override returns (string memory) {
        return _baseTokenURI;
    }

    function setBaseURI(string calldata baseURI) external onlyOwner {
        _baseTokenURI = baseURI;
    }

    function tokenURI(
        uint256 tokenId
    ) public view virtual override returns (string memory) {
        if (ownerOf(tokenId) == address(0)) {
            revert NonExistentTokenURI();
        }
        return
            bytes(_baseTokenURI).length > 0
                ? string(abi.encodePacked(_baseTokenURI, tokenId.toString()))
                : "";
    }

    // ============ EIP-2612 Style Permit Functions ============

    /**
     * @dev Returns the current nonce for `owner`. This value must be included in the signature.
     */
    function nonces(address owner) public view returns (uint256) {
        return _nonces[owner];
    }

    /**
     * @dev Returns the domain separator for the current chain.
     */
    function DOMAIN_SEPARATOR() external view returns (bytes32) {
        return _domainSeparatorV4();
    }

    /**
     * @dev Permit approval for a specific token using EIP-712 signature
     * @param owner The owner of the token
     * @param spender The address to approve
     * @param tokenId The token ID to approve
     * @param deadline The deadline timestamp for the signature
     * @param v The recovery byte of the signature
     * @param r Half of the ECDSA signature
     * @param s Half of the ECDSA signature
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
        if (block.timestamp > deadline) {
            revert SignatureExpired();
        }

        bytes32 structHash = keccak256(
            abi.encode(
                PERMIT_TYPEHASH,
                owner,
                spender,
                tokenId,
                _nonces[owner]++,
                deadline
            )
        );

        bytes32 hash = _hashTypedDataV4(structHash);
        address signer = ECDSA.recover(hash, v, r, s);

        if (signer != owner) {
            revert InvalidSigner();
        }

        if (ownerOf(tokenId) != owner) {
            revert InvalidSigner();
        }

        _approve(spender, tokenId, owner);
        emit PermitUsed(owner, spender, tokenId);
    }

    /**
     * @dev Permit approval for all tokens using EIP-712 signature (setApprovalForAll)
     * @param owner The owner granting approval
     * @param operator The operator to approve/revoke
     * @param approved Whether to approve or revoke
     * @param deadline The deadline timestamp for the signature
     * @param v The recovery byte of the signature
     * @param r Half of the ECDSA signature
     * @param s Half of the ECDSA signature
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
        if (block.timestamp > deadline) {
            revert SignatureExpired();
        }

        bytes32 structHash = keccak256(
            abi.encode(
                PERMIT_FOR_ALL_TYPEHASH,
                owner,
                operator,
                approved,
                _nonces[owner]++,
                deadline
            )
        );

        bytes32 hash = _hashTypedDataV4(structHash);
        address signer = ECDSA.recover(hash, v, r, s);

        if (signer != owner) {
            revert InvalidSigner();
        }

        _setApprovalForAll(owner, operator, approved);
        emit PermitForAllUsed(owner, operator, approved);
    }

    /**
     * @dev Transfer token using permit signature (gasless transfer)
     * @param from The current owner
     * @param to The recipient
     * @param tokenId The token ID to transfer
     * @param deadline The deadline timestamp for the signature
     * @param v The recovery byte of the signature
     * @param r Half of the ECDSA signature
     * @param s Half of the ECDSA signature
     */
    function transferWithPermit(
        address from,
        address to,
        uint256 tokenId,
        uint256 deadline,
        uint8 v,
        bytes32 r,
        bytes32 s
    ) public {
        // First, validate the permit and approve msg.sender
        permit(from, msg.sender, tokenId, deadline, v, r, s);

        // Then transfer the token
        safeTransferFrom(from, to, tokenId);
    }

    function burnWithPermit(
        address from,
        uint256 tokenId,
        uint256 deadline,
        uint8 v,
        bytes32 r,
        bytes32 s
    ) public {
        permit(from, msg.sender, tokenId, deadline, v, r, s);

        burn(tokenId);
    }

    // ============ End of Permit Functions ============

    // The following functions are overrides required by Solidity.

    // The following functions are overrides required by Solidity.

    function _update(
        address to,
        uint256 tokenId,
        address auth
    ) internal override(ERC721, ERC721Enumerable) returns (address) {
        return super._update(to, tokenId, auth);
    }

    function _increaseBalance(
        address account,
        uint128 value
    ) internal override(ERC721, ERC721Enumerable) {
        super._increaseBalance(account, value);
    }

    function supportsInterface(
        bytes4 interfaceId
    ) public view override(ERC721, ERC721Enumerable) returns (bool) {
        return super.supportsInterface(interfaceId);
    }
}
