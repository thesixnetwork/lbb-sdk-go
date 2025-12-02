// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title GoldCertificationNFT
 * @dev ERC721 contract for Gold Certification NFTs with Cosmos metadata integration
 *
 * This contract represents the EVM side of a hybrid NFT system where:
 * - NFT tokens exist on EVM chains (for trading and market access)
 * - Rich metadata and updates are managed on Cosmos (for cost efficiency)
 * - Token URIs point to a bridge API that fetches live data from Cosmos
 */
contract GoldCertificateNFT is ERC721, ERC721URIStorage, Ownable {
    uint256 private _nextTokenId = 1;

    // Base URI for metadata (points to Cosmos bridge API)
    string private _baseTokenURI;

    // Contract metadata for OpenSea
    string public contractURI;

    // Mapping to track custom token IDs (e.g., "GOLD001")
    mapping(string => uint256) public stringToTokenId;
    mapping(uint256 => string) public tokenIdToString;

    // Events
    event BaseURIUpdated(string newBaseURI);
    event ContractURIUpdated(string newContractURI);
    event TokenMinted(
        address indexed to,
        uint256 indexed tokenId,
        string tokenIdString
    );

    /**
     * @dev Constructor
     * @param name The name of the NFT collection
     * @param symbol The symbol of the NFT collection
     * @param baseTokenURI Base URI pointing to Cosmos metadata bridge
     * @param _contractURI URI for contract metadata (OpenSea collection info)
     * @param initialOwner Address that will own the contract
     */
    constructor(
        string memory name,
        string memory symbol,
        string memory baseTokenURI,
        string memory _contractURI,
        address initialOwner
    ) ERC721(name, symbol) Ownable(initialOwner) {
        _baseTokenURI = baseTokenURI;
        contractURI = _contractURI;
    }

    /**
     * @dev Mint NFT with custom string token ID (e.g., "GOLD001")
     * @param to Address to mint the NFT to
     * @param tokenIdString Custom token ID string
     * @return The numeric token ID
     */
    function safeMint(
        address to,
        string memory tokenIdString
    ) public onlyOwner returns (uint256) {
        require(
            bytes(tokenIdString).length > 0,
            "Token ID string cannot be empty"
        );
        require(
            stringToTokenId[tokenIdString] == 0,
            "Token ID string already exists"
        );

        uint256 tokenId = _nextTokenId++;

        // Map string ID to numeric ID
        stringToTokenId[tokenIdString] = tokenId;
        tokenIdToString[tokenId] = tokenIdString;

        // Mint the NFT
        _safeMint(to, tokenId);

        // Set token URI pointing to Cosmos metadata
        string memory metadataURI = string(
            abi.encodePacked(_baseTokenURI, tokenIdString)
        );
        _setTokenURI(tokenId, metadataURI);

        emit TokenMinted(to, tokenId, tokenIdString);

        return tokenId;
    }

    /**
     * @dev Batch mint multiple Gold NFTs
     * @param recipients Array of recipient addresses
     * @param tokenIdStrings Array of token ID strings
     * @return Array of numeric token IDs
     */
    function batchMint(
        address[] memory recipients,
        string[] memory tokenIdStrings
    ) public onlyOwner returns (uint256[] memory) {
        require(
            recipients.length == tokenIdStrings.length,
            "Arrays length mismatch"
        );
        require(recipients.length > 0, "Arrays cannot be empty");

        uint256[] memory tokenIds = new uint256[](recipients.length);

        for (uint256 i = 0; i < recipients.length; i++) {
            tokenIds[i] = safeMint(recipients[i], tokenIdStrings[i]);
        }

        return tokenIds;
    }

    /**
     * @dev Update base URI (points to Cosmos metadata bridge)
     * @param newBaseURI New base URI for metadata
     */
    function setBaseURI(string memory newBaseURI) public onlyOwner {
        _baseTokenURI = newBaseURI;
        emit BaseURIUpdated(newBaseURI);
    }

    /**
     * @dev Update contract metadata URI for OpenSea
     * @param newContractURI New contract metadata URI
     */
    function setContractURI(string memory newContractURI) public onlyOwner {
        contractURI = newContractURI;
        emit ContractURIUpdated(newContractURI);
    }

    /**
     * @dev Get the base URI
     */
    function _baseURI() internal view override returns (string memory) {
        return _baseTokenURI;
    }

    /**
     * @dev Override tokenURI to use stored URI
     */
    function tokenURI(
        uint256 tokenId
    ) public view override(ERC721, ERC721URIStorage) returns (string memory) {
        return super.tokenURI(tokenId);
    }

    /**
     * @dev Get token ID from string
     * @param tokenIdString The string token ID
     * @return The numeric token ID (0 if not found)
     */
    function getTokenIdFromString(
        string memory tokenIdString
    ) public view returns (uint256) {
        return stringToTokenId[tokenIdString];
    }

    /**
     * @dev Get string token ID from numeric ID
     * @param tokenId The numeric token ID
     * @return The string token ID
     */
    function getStringFromTokenId(
        uint256 tokenId
    ) public view returns (string memory) {
        require(_ownerOf(tokenId) != address(0), "Token does not exist");
        return tokenIdToString[tokenId];
    }

    /**
     * @dev Check if a string token ID exists
     * @param tokenIdString The string token ID to check
     * @return True if the token exists
     */
    function stringTokenExists(
        string memory tokenIdString
    ) public view returns (bool) {
        return stringToTokenId[tokenIdString] != 0;
    }

    /**
     * @dev Get total supply
     * @return Total number of tokens minted
     */
    function totalSupply() public view returns (uint256) {
        return _nextTokenId - 1;
    }

    /**
     * @dev Override supportsInterface for multiple inheritance
     */
    function supportsInterface(
        bytes4 interfaceId
    ) public view override(ERC721, ERC721URIStorage) returns (bool) {
        return super.supportsInterface(interfaceId);
    }
}
