// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.20;

import "forge-std/Script.sol";
import {MyToken} from "../src/token.sol";

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

        MyToken myToken = new MyToken(ownerAddress);
        nonceUp(ownerAddress);
        address myTokenAddress = address(myToken);
        console.log("NFT : ", myTokenAddress);

        myToken.safeMint(ownerAddress, 1);
        nonceUp(ownerAddress);

        vm.stopBroadcast();
    }

    function nonceUp(address signer) public {
        vm.setNonce(signer, currentNonce + uint64(1));
        currentNonce++;
    }
}
