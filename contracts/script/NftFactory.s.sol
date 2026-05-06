// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.20;

import "forge-std/Script.sol";
import {NftFactory} from "../src/NftFactory.sol";

// ─────────────────────────────────────────────────────────────────────────────
//  Required environment variables
//  PRIVATE_KEY   — deployer's private key (must hold SUPER_ADMIN role intent)
//  SUPER_ADMIN   — address that will receive the super admin role
//
//  Optional
//  NFT_NAME      — collection name      (default: "LBB NFT")
//  NFT_SYMBOL    — collection symbol    (default: "LNFT")
//  NFT_BASE_URI  — base metadata URI   (default: "")
//
//  Run:
//    forge script contracts/script/NftFactory.s.sol:DeployNftFactory \
//      --rpc-url <RPC_URL> --broadcast --verify
// ─────────────────────────────────────────────────────────────────────────────
contract DeployNftFactory is Script {
    address superAdminAddress;

    function setUp() public {
        superAdminAddress = vm.envAddress("SUPER_ADMIN");
    }

    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");

        string memory name     = vm.envOr("NFT_NAME",     string("LBB NFT"));
        string memory symbol   = vm.envOr("NFT_SYMBOL",   string("LNFT"));
        string memory baseURI  = vm.envOr("NFT_BASE_URI", string(""));

        vm.startBroadcast(deployerPrivateKey);

        NftFactory nftFactory = new NftFactory(
            name,
            symbol,
            baseURI,
            superAdminAddress
        );

        vm.stopBroadcast();

        console.log("=== NftFactory Deployed ===");
        console.log("Contract address :", address(nftFactory));
        console.log("Super admin      :", nftFactory.getSuperAdmin());
        console.log("Name             :", nftFactory.name());
        console.log("Symbol           :", nftFactory.symbol());
        console.log("Total supply     :", nftFactory.totalSupply());
    }
}

// ─────────────────────────────────────────────────────────────────────────────
//  Grant admin role to a new address after deployment.
//
//  Required env vars:
//    PRIVATE_KEY        — super admin's private key
//    NFT_FACTORY        — deployed NftFactory address
//    ADMIN_ADDRESS      — address to be granted admin
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
//    PRIVATE_KEY        — super admin's private key
//    NFT_FACTORY        — deployed NftFactory address
//    ADMIN_ADDRESS      — admin address to revoke
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
//  Mint a single token.
//
//  Required env vars:
//    PRIVATE_KEY   — admin or super admin private key
//    NFT_FACTORY   — deployed NftFactory address
//    MINT_TO       — recipient address
//    TOKEN_ID      — token ID to mint
//
//  Run:
//    forge script contracts/script/NftFactory.s.sol:MintNftScript \
//      --rpc-url <RPC_URL> --broadcast
// ─────────────────────────────────────────────────────────────────────────────
contract MintNftScript is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address nftFactoryAddress  = vm.envAddress("NFT_FACTORY");
        address mintTo             = vm.envAddress("MINT_TO");
        uint256 tokenId            = vm.envUint("TOKEN_ID");

        vm.startBroadcast(deployerPrivateKey);

        NftFactory nftFactory = NftFactory(nftFactoryAddress);
        nftFactory.safeMint(mintTo, tokenId);

        vm.stopBroadcast();

        console.log("=== Token Minted ===");
        console.log("Token ID         :", tokenId);
        console.log("Owner            :", nftFactory.ownerOf(tokenId));
        console.log("Total supply     :", nftFactory.totalSupply());
    }
}
