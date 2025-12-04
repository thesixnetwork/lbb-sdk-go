// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.20;

import "forge-std/Script.sol";
import { LBBCert } from "../src/CertAutoID.sol";

contract DeployScript is Script {
  address ownerAddress;
  uint64 currentNonce;

  function setUp() public {
    ownerAddress = vm.envAddress("OWNER");
    currentNonce = vm.getNonce(ownerAddress);
  }

  function run() external {
    uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
    vm.startBroadcast(deployerPrivateKey);

    LBBCert certNFT = new LBBCert("CERT", "CERTIFICATE", "https://picsum.photos/", ownerAddress);
    nonceUp(ownerAddress);
    address certNFTAddress = address(certNFT);
    console.log("cert : ", certNFTAddress);

    certNFT.safeMint(ownerAddress);
    nonceUp(ownerAddress);

    vm.stopBroadcast();
  }

  function nonceUp(address signer) public {
    vm.setNonce(signer, currentNonce + uint64(1));
    currentNonce++;
  }
}

contract MintScript is Script {
  address ownerAddress;
  address certAddress;
  uint64 currentNonce;

  function setUp() public {
    ownerAddress = vm.envAddress("OWNER");
    currentNonce = vm.getNonce(ownerAddress);
    string memory nftContractInfoPath = "./broadcast/CertAutoID.s.sol/150/run-latest.json";
    string memory nftContractInfo = vm.readFile(nftContractInfoPath);
    bytes memory certNFTJsonParsed = vm.parseJson(nftContractInfo, ".transactions[0].contractAddress");

    certAddress = abi.decode(certNFTJsonParsed, (address));
  }

  function run() external {
    uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
    vm.startBroadcast(deployerPrivateKey);

    LBBCert certNFT = LBBCert(payable(certAddress));
    certNFT.safeMint(ownerAddress);
    nonceUp(ownerAddress);

    vm.stopBroadcast();
  }

  function nonceUp(address signer) public {
    vm.setNonce(signer, currentNonce + uint64(1));
    currentNonce++;
  }
}

contract TransferToken is Script {
  address contractAdrress;
  address ownerAddress;
  address certNFTContractAddress;

  function setUp() public {
    ownerAddress = vm.envAddress("OWNER");
    string memory nftContractInfoPath = "./broadcast/CertAutoID.s.sol/666/run-latest.json";
    string memory nftContractInfo = vm.readFile(nftContractInfoPath);
    bytes memory certNFTJsonParsed = vm.parseJson(nftContractInfo, ".transactions[0].contractAddress");

    certNFTContractAddress = abi.decode(certNFTJsonParsed, (address));
  }

  function run() external {
    uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
    LBBCert certNFT = LBBCert(payable(certNFTContractAddress));
    vm.startBroadcast(deployerPrivateKey);

    certNFT.transferFrom(ownerAddress, 0x3753C81072A56072840990D3D02f354Efb7425A3, 1);

    vm.stopBroadcast();
  }
}
