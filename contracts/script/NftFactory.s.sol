// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.20;

import "forge-std/Script.sol";
import {NftFactory} from "../src/NftFactory.sol";

// ─────────────────────────────────────────────────────────────────────────────
//  Deploy NftFactory
//
//  Required env vars:
//    PRIVATE_KEY    — deployer private key
//    SUPER_ADMIN    — address that receives the super admin role
//    MINT_SIGNER    — backend wallet address authorised to sign mint requests
//
//  Optional:
//    NFT_NAME       — collection name     (default: "LBB NFT")
//    NFT_SYMBOL     — collection symbol   (default: "LNFT")
//
//  Run:
//    forge script contracts/script/NftFactory.s.sol:DeployNftFactory \
//      --rpc-url <RPC_URL> --broadcast --verify
// ─────────────────────────────────────────────────────────────────────────────
contract DeployNftFactory is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address superAdminAddress  = vm.envAddress("SUPER_ADMIN");
        address mintSignerAddress  = vm.envAddress("MINT_SIGNER");

        string memory name   = vm.envOr("NFT_NAME",   string("LBB NFT"));
        string memory symbol = vm.envOr("NFT_SYMBOL", string("LNFT"));

        vm.startBroadcast(deployerPrivateKey);

        NftFactory nftFactory = new NftFactory(
            name,
            symbol,
            superAdminAddress,
            mintSignerAddress
        );

        vm.stopBroadcast();

        console.log("=== NftFactory Deployed ===");
        console.log("Contract address :", address(nftFactory));
        console.log("Super admin      :", nftFactory.getSuperAdmin());
        console.log("Mint signer      :", nftFactory.getMintSigner());
        console.log("Name             :", nftFactory.name());
        console.log("Symbol           :", nftFactory.symbol());
    }
}

// ─────────────────────────────────────────────────────────────────────────────
//  Grant admin role to a new address.
//
//  Required env vars:
//    PRIVATE_KEY    — super admin private key
//    NFT_FACTORY    — deployed NftFactory address
//    ADMIN_ADDRESS  — address to grant admin role
//
//  Run:
//    forge script contracts/script/NftFactory.s.sol:GrantAdminScript \
//      --rpc-url <RPC_URL> --broadcast
// ─────────────────────────────────────────────────────────────────────────────
contract GrantAdminScript is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address nftFactoryAddress  = vm.envAddress("NFT_FACTORY");
        address adminToGrant       = vm.envAddress("ADMIN_ADDRESS");

        vm.startBroadcast(deployerPrivateKey);

        NftFactory nftFactory = NftFactory(nftFactoryAddress);
        nftFactory.grantAdmin(adminToGrant);

        vm.stopBroadcast();

        console.log("=== Admin Granted ===");
        console.log("New admin        :", adminToGrant);
        console.log("Total admins     :", nftFactory.getAdminCount());
    }
}

// ─────────────────────────────────────────────────────────────────────────────
//  Revoke admin role from an address.
//
//  Required env vars:
//    PRIVATE_KEY    — super admin private key
//    NFT_FACTORY    — deployed NftFactory address
//    ADMIN_ADDRESS  — admin address to revoke
//
//  Run:
//    forge script contracts/script/NftFactory.s.sol:RevokeAdminScript \
//      --rpc-url <RPC_URL> --broadcast
// ─────────────────────────────────────────────────────────────────────────────
contract RevokeAdminScript is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address nftFactoryAddress  = vm.envAddress("NFT_FACTORY");
        address adminToRevoke      = vm.envAddress("ADMIN_ADDRESS");

        vm.startBroadcast(deployerPrivateKey);

        NftFactory nftFactory = NftFactory(nftFactoryAddress);
        nftFactory.revokeAdmin(adminToRevoke);

        vm.stopBroadcast();

        console.log("=== Admin Revoked ===");
        console.log("Revoked admin    :", adminToRevoke);
        console.log("Remaining admins :", nftFactory.getAdminCount());
    }
}

// ─────────────────────────────────────────────────────────────────────────────
//  Update the authorised mint signer (backend wallet).
//
//  Required env vars:
//    PRIVATE_KEY    — super admin private key
//    NFT_FACTORY    — deployed NftFactory address
//    MINT_SIGNER    — new mint signer address
//
//  NOTE: Changing the signer immediately invalidates all previously issued
//        off-chain Mint signatures that have not yet been submitted on-chain.
//
//  Run:
//    forge script contracts/script/NftFactory.s.sol:SetMintSignerScript \
//      --rpc-url <RPC_URL> --broadcast
// ─────────────────────────────────────────────────────────────────────────────
contract SetMintSignerScript is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address nftFactoryAddress  = vm.envAddress("NFT_FACTORY");
        address newMintSigner      = vm.envAddress("MINT_SIGNER");

        vm.startBroadcast(deployerPrivateKey);

        NftFactory nftFactory = NftFactory(nftFactoryAddress);
        nftFactory.setMintSigner(newMintSigner);

        vm.stopBroadcast();

        console.log("=== Mint Signer Updated ===");
        console.log("New mint signer  :", nftFactory.getMintSigner());
    }
}

// ─────────────────────────────────────────────────────────────────────────────
//  Mint a single token using a backend EIP-712 signature.
//
//  Required env vars:
//    PRIVATE_KEY        — admin private key (pays gas, must hold admin role)
//    MINT_SIGNER_KEY    — mint signer private key (signs the authorisation)
//    NFT_FACTORY        — deployed NftFactory address
//    MINT_TO            — recipient address
//    TOKEN_ID           — token ID to mint
//    METADATA_BASE64    — Base64-encoded JSON metadata string (no data URI prefix)
//
//  How it works:
//    1. The script reads the current nonce for MINT_TO from the contract.
//    2. It builds the EIP-712 Mint digest and signs it with MINT_SIGNER_KEY.
//    3. It broadcasts safeMintWithSignature from the admin key with an empty
//       approvers array (populate for production multisig audit trails).
//
//  Run:
//    forge script contracts/script/NftFactory.s.sol:MintNftScript \
//      --rpc-url <RPC_URL> --broadcast
// ─────────────────────────────────────────────────────────────────────────────
contract MintNftScript is Script {
    // Must match NftFactory.MINT_TYPEHASH
    bytes32 private constant MINT_TYPEHASH =
        keccak256("Mint(address to,uint256 tokenId,string metadataBase64,uint256 nonce,uint256 deadline)");

    struct MintParams {
        address mintTo;
        uint256 tokenId;
        string  metadataBase64;
        uint256 deadline;
    }

    function _buildDigest(NftFactory nftFactory, MintParams memory p) internal view returns (bytes32) {
        bytes32 structHash = keccak256(
            abi.encode(
                MINT_TYPEHASH,
                p.mintTo,
                p.tokenId,
                keccak256(bytes(p.metadataBase64)),
                nftFactory.nonces(p.mintTo),
                p.deadline
            )
        );
        return keccak256(abi.encodePacked("\x19\x01", nftFactory.DOMAIN_SEPARATOR(), structHash));
    }

    function run() external {
        uint256 adminPrivateKey      = vm.envUint("PRIVATE_KEY");
        uint256 mintSignerPrivateKey = vm.envUint("MINT_SIGNER_KEY");
        NftFactory nftFactory        = NftFactory(vm.envAddress("NFT_FACTORY"));

        MintParams memory p = MintParams({
            mintTo:         vm.envAddress("MINT_TO"),
            tokenId:        vm.envUint("TOKEN_ID"),
            metadataBase64: vm.envString("METADATA_BASE64"),
            deadline:       block.timestamp + 1 hours
        });

        (uint8 v, bytes32 r, bytes32 s) = vm.sign(mintSignerPrivateKey, _buildDigest(nftFactory, p));

        // approvers is the off-chain multisig audit trail; populate as needed.
        address[] memory approvers = new address[](0);

        vm.startBroadcast(adminPrivateKey);
        nftFactory.safeMintWithSignature(p.mintTo, p.tokenId, p.metadataBase64, approvers, p.deadline, v, r, s);
        vm.stopBroadcast();

        console.log("=== Token Minted ===");
        console.log("Token ID         :", p.tokenId);
        console.log("Owner            :", nftFactory.ownerOf(p.tokenId));
        console.log("Total supply     :", nftFactory.totalSupply());
    }
}

