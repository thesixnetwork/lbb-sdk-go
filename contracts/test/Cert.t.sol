// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.20;

import "forge-std/Test.sol";

import { GoldCertificateNFT } from "../src/GoldCertificate.sol";

contract CertNFTTest is Test {
    GoldCertificateNFT public myNft;

   //  function setUp() public {
   //      myNft = new GoldCertificateNFT("TEST","TEST", "", "");
   //      myNft.setPreMinteeAddress(address(0xDEE));
   //      myNft.setLimitedEditionSize(260);
   //      myNft.preMint(260);
   //  }

   //  function testMint() public {
   //      myNft.setLimitedEditionSize(myNft.limitedEditionSize() + 40);
   //      myNft.preMint(40);
   //      assertEq(myNft.balanceOf(address(0xDEE)), 300);
   //      assertEq(myNft.ownerOf(1), address(0xDEE));
   //  }

   //  function testRevertSetPreMintFromNotOwner() public {
   //      vm.startPrank(address(1));
   //      vm.expectRevert("Ownable: caller is not the owner");
   //      myNft.setPreMinteeAddress(address(0xDEE));
   //      vm.stopPrank();
   //  }

   //  function testRevertPreMintOverLimit() public {
   //      vm.expectRevert("Too many already minted");
   //      myNft.preMint(200);
   //  }
}
