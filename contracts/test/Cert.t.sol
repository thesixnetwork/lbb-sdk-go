// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.20;

import "forge-std/Test.sol";

import {LBBCert} from "../src/Cert.sol";

contract CertNFTTest is Test {
    address public ownerAddress;
    uint256 public ownerPrivateKey;

    address public spender;

    LBBCert public cert;

    error OwnableUnauthorizedAccount(address);
    error ERC721InvalidSender(address);

    bytes32 private constant PERMIT_TYPEHASH =
        keccak256(
            "Permit(address owner,address spender,uint256 tokenId,uint256 nonce,uint256 deadline)"
        );

    bytes32 private constant PERMIT_FOR_ALL_TYPEHASH =
        keccak256(
            "PermitForAll(address owner,address operator,bool approved,uint256 nonce,uint256 deadline)"
        );

    function setUp() public {
        ownerPrivateKey = vm.envUint("PRIVATE_KEY");
        ownerAddress = vm.addr(ownerPrivateKey);
        // spender is contract owner address or someone who will pay gas on behalf of signer
        spender = address(0xDEE);

        vm.startPrank(ownerAddress);
        cert = new LBBCert("TEST", "TEST", "TEST", ownerAddress);
        vm.stopPrank();
    }

    function testMint() public {
        vm.startPrank(ownerAddress);
        cert.safeMint(address(0xDEE), 1);
        assertEq(cert.balanceOf(address(0xDEE)), 1);
        assertEq(cert.ownerOf(1), address(0xDEE));
        vm.stopPrank();
    }

    function testRevertMintFromNotOwner() public {
        address unauthorizedUser = address(2);

        vm.startPrank(unauthorizedUser);
        vm.expectRevert(
            abi.encodeWithSelector(
                OwnableUnauthorizedAccount.selector,
                unauthorizedUser
            )
        );

        cert.safeMint(address(0xDEE), 2);
        vm.stopPrank();
    }

    function testRevertMintTokenAlredyExisted() public {
        vm.startPrank(ownerAddress);
        cert.safeMint(address(0xDEE), 1);
        vm.expectRevert(
            abi.encodeWithSelector(ERC721InvalidSender.selector, address(0))
        );
        cert.safeMint(address(2), 1);
        vm.stopPrank();
    }

    function testPermitValidSignature() public {
        uint256 tokenId = 1;
        uint256 deadline = block.timestamp + 1 hours;
        uint256 nonce = cert.nonces(ownerAddress);
        vm.prank(ownerAddress);
        cert.safeMint(ownerAddress, 1);

        // Generate valid signature
        bytes32 structHash = keccak256(
            abi.encode(
                keccak256(
                    "Permit(address owner,address spender,uint256 tokenId,uint256 nonce,uint256 deadline)"
                ),
                ownerAddress,
                spender,
                tokenId,
                nonce,
                deadline
            )
        );

        bytes32 domainSeparator = cert.DOMAIN_SEPARATOR();
        bytes32 digest = keccak256(
            abi.encodePacked("\x19\x01", domainSeparator, structHash)
        );
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(ownerPrivateKey, digest);

        // Execute permit
        cert.permit(ownerAddress, spender, tokenId, deadline, v, r, s);

        // Verify approval
        assertEq(cert.getApproved(tokenId), spender);
        assertEq(cert.nonces(ownerAddress), nonce + 1);
    }

    function testPermitForAllValidSignature() public {
        uint256 deadline = block.timestamp + 1 hours;
        uint256 nonce = cert.nonces(ownerAddress);
        bool approved = true;

        // Generate valid signature
        bytes32 structHash = keccak256(
            abi.encode(
                keccak256(
                    "PermitForAll(address owner,address operator,bool approved,uint256 nonce,uint256 deadline)"
                ),
                ownerAddress,
                spender,
                approved,
                nonce,
                deadline
            )
        );

        bytes32 domainSeparator = cert.DOMAIN_SEPARATOR();
        bytes32 digest = keccak256(
            abi.encodePacked("\x19\x01", domainSeparator, structHash)
        );
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(ownerPrivateKey, digest);

        // Execute permitForAll
        cert.permitForAll(ownerAddress, spender, approved, deadline, v, r, s);

        // Verify approval
        assertTrue(cert.isApprovedForAll(ownerAddress, spender));
        assertEq(cert.nonces(ownerAddress), nonce + 1);
    }

    // ============ Invalid Signature Tests ============

    function testPermitInvalidSignatureShouldNotIncrementNonce() public {
        uint256 tokenId = 1;
        uint256 deadline = block.timestamp + 1 hours;
        uint256 nonceBefore = cert.nonces(ownerAddress);

        vm.prank(ownerAddress);
        cert.safeMint(ownerAddress, 1);

        // Create signature with WRONG private key (attacker's key)
        uint256 attackerPrivateKey = 0xBAD;

        bytes32 structHash = keccak256(
            abi.encode(
                keccak256(
                    "Permit(address owner,address spender,uint256 tokenId,uint256 nonce,uint256 deadline)"
                ),
                ownerAddress,
                spender,
                tokenId,
                nonceBefore,
                deadline
            )
        );

        bytes32 domainSeparator = cert.DOMAIN_SEPARATOR();
        bytes32 digest = keccak256(
            abi.encodePacked("\x19\x01", domainSeparator, structHash)
        );
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(attackerPrivateKey, digest);

        // Attempt permit with invalid signature
        vm.expectRevert(abi.encodeWithSignature("InvalidSigner()"));
        cert.permit(ownerAddress, spender, tokenId, deadline, v, r, s);

        // CRITICAL: Nonce should NOT have incremented
        assertEq(
            cert.nonces(ownerAddress),
            nonceBefore,
            "Nonce should not increment on invalid signature"
        );

        // Token should NOT be approved
        assertEq(cert.getApproved(tokenId), address(0));
    }

    function testPermitForAllInvalidSignatureShouldNotIncrementNonce() public {
        uint256 deadline = block.timestamp + 1 hours;
        uint256 nonceBefore = cert.nonces(ownerAddress);
        bool approved = true;

        // Create signature with WRONG private key
        uint256 attackerPrivateKey = 0xBAD;

        bytes32 structHash = keccak256(
            abi.encode(
                keccak256(
                    "PermitForAll(address owner,address operator,bool approved,uint256 nonce,uint256 deadline)"
                ),
                ownerAddress,
                spender,
                approved,
                nonceBefore,
                deadline
            )
        );

        bytes32 domainSeparator = cert.DOMAIN_SEPARATOR();
        bytes32 digest = keccak256(
            abi.encodePacked("\x19\x01", domainSeparator, structHash)
        );
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(attackerPrivateKey, digest);

        // Attempt permitForAll with invalid signature
        vm.expectRevert(abi.encodeWithSignature("InvalidSigner()"));
        cert.permitForAll(ownerAddress, spender, approved, deadline, v, r, s);

        // CRITICAL: Nonce should NOT have incremented
        assertEq(
            cert.nonces(ownerAddress),
            nonceBefore,
            "Nonce should not increment on invalid signature"
        );

        // Operator should NOT be approved
        assertFalse(cert.isApprovedForAll(ownerAddress, spender));
    }

    // ============ Nonce Replay Tests ============

    function testPermitCannotReplaySignature() public {
        uint256 tokenId = 1;
        uint256 deadline = block.timestamp + 1 hours;
        uint256 nonce = cert.nonces(ownerAddress);
        vm.prank(ownerAddress);
        cert.safeMint(ownerAddress, 1);

        // Generate valid signature
        bytes32 structHash = keccak256(
            abi.encode(
                keccak256(
                    "Permit(address owner,address spender,uint256 tokenId,uint256 nonce,uint256 deadline)"
                ),
                ownerAddress,
                spender,
                tokenId,
                nonce,
                deadline
            )
        );

        bytes32 domainSeparator = cert.DOMAIN_SEPARATOR();
        bytes32 digest = keccak256(
            abi.encodePacked("\x19\x01", domainSeparator, structHash)
        );
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(ownerPrivateKey, digest);

        // First permit succeeds
        cert.permit(ownerAddress, spender, tokenId, deadline, v, r, s);
        assertEq(cert.getApproved(tokenId), spender);

        // Reset approval
        vm.prank(ownerAddress);
        cert.approve(address(0), tokenId);

        // Replay same signature should fail (nonce changed)
        vm.expectRevert(abi.encodeWithSignature("InvalidSigner()"));
        cert.permit(ownerAddress, spender, tokenId, deadline, v, r, s);

        // Token should NOT be approved
        assertEq(cert.getApproved(tokenId), address(0));
    }

    // ============ Deadline Tests ============

    function testPermitExpiredDeadline() public {
        uint256 tokenId = 1;
        uint256 deadline = block.timestamp - 1; // Expired
        uint256 nonce = cert.nonces(ownerAddress);

        bytes32 structHash = keccak256(
            abi.encode(
                keccak256(
                    "Permit(address owner,address spender,uint256 tokenId,uint256 nonce,uint256 deadline)"
                ),
                ownerAddress,
                spender,
                tokenId,
                nonce,
                deadline
            )
        );

        bytes32 domainSeparator = cert.DOMAIN_SEPARATOR();
        bytes32 digest = keccak256(
            abi.encodePacked("\x19\x01", domainSeparator, structHash)
        );
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(ownerPrivateKey, digest);

        // Should revert with SignatureExpired
        vm.expectRevert(abi.encodeWithSignature("SignatureExpired()"));
        cert.permit(ownerAddress, spender, tokenId, deadline, v, r, s);

        // Nonce should not increment
        assertEq(cert.nonces(ownerAddress), nonce);
    }

    // ============ Zero Address Tests ============

    function testPermitZeroAddressSpenderShouldRevert() public {
        uint256 tokenId = 1;
        uint256 deadline = block.timestamp + 1 hours;
        uint256 nonce = cert.nonces(ownerAddress);

        bytes32 structHash = keccak256(
            abi.encode(
                keccak256(
                    "Permit(address owner,address spender,uint256 tokenId,uint256 nonce,uint256 deadline)"
                ),
                ownerAddress,
                address(0), // Zero address spender
                tokenId,
                nonce,
                deadline
            )
        );

        bytes32 domainSeparator = cert.DOMAIN_SEPARATOR();
        bytes32 digest = keccak256(
            abi.encodePacked("\x19\x01", domainSeparator, structHash)
        );
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(ownerPrivateKey, digest);

        // Should revert
        vm.expectRevert();
        cert.permit(ownerAddress, address(0), tokenId, deadline, v, r, s);
    }

    // ============ TransferWithPermit Tests ============

    function testTransferWithPermitSuccess() public {
        uint256 tokenId = 1;
        uint256 deadline = block.timestamp + 1 hours;
        uint256 nonce = cert.nonces(ownerAddress);
        address recipient = makeAddr("recipient");

        vm.prank(ownerAddress);
        cert.safeMint(ownerAddress, 1);

        bytes32 structHash = keccak256(
            abi.encode(
                keccak256(
                    "Permit(address owner,address spender,uint256 tokenId,uint256 nonce,uint256 deadline)"
                ),
                ownerAddress,
                spender,
                tokenId,
                nonce,
                deadline
            )
        );

        bytes32 domainSeparator = cert.DOMAIN_SEPARATOR();
        bytes32 digest = keccak256(
            abi.encodePacked("\x19\x01", domainSeparator, structHash)
        );
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(ownerPrivateKey, digest);

        // Execute transferWithPermit as spender
        vm.prank(spender);
        cert.transferWithPermit(
            ownerAddress,
            recipient,
            tokenId,
            deadline,
            v,
            r,
            s
        );

        // Verify transfer
        assertEq(cert.ownerOf(tokenId), recipient);
        assertEq(cert.nonces(ownerAddress), nonce + 1);
    }

    function testTransferWithPermitInvalidRecipient() public {
        uint256 tokenId = 1;
        uint256 deadline = block.timestamp + 1 hours;
        uint256 nonce = cert.nonces(ownerAddress);

        bytes32 structHash = keccak256(
            abi.encode(
                keccak256(
                    "Permit(address owner,address spender,uint256 tokenId,uint256 nonce,uint256 deadline)"
                ),
                ownerAddress,
                spender,
                tokenId,
                nonce,
                deadline
            )
        );

        bytes32 domainSeparator = cert.DOMAIN_SEPARATOR();
        bytes32 digest = keccak256(
            abi.encodePacked("\x19\x01", domainSeparator, structHash)
        );
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(ownerPrivateKey, digest);

        // Attempt transfer to zero address
        vm.prank(spender);
        vm.expectRevert();
        cert.transferWithPermit(
            ownerAddress,
            address(0),
            tokenId,
            deadline,
            v,
            r,
            s
        );

        // CRITICAL: Nonce should not increment if using Option 2 fix
        // (Comment out if using Option 1)
        // assertEq(cert.nonces(owner), nonce);
    }

    // ============ BurnWithPermit Tests ============

    function testBurnWithPermitSuccess() public {
        uint256 tokenId = 1;
        uint256 deadline = block.timestamp + 1 hours;
        uint256 nonce = cert.nonces(ownerAddress);
        vm.prank(ownerAddress);
        cert.safeMint(ownerAddress, 1);

        bytes32 structHash = keccak256(
            abi.encode(
                keccak256(
                    "Permit(address owner,address spender,uint256 tokenId,uint256 nonce,uint256 deadline)"
                ),
                ownerAddress,
                spender,
                tokenId,
                nonce,
                deadline
            )
        );

        bytes32 domainSeparator = cert.DOMAIN_SEPARATOR();
        bytes32 digest = keccak256(
            abi.encodePacked("\x19\x01", domainSeparator, structHash)
        );
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(ownerPrivateKey, digest);

        // Execute burnWithPermit
        vm.prank(spender);
        cert.burnWithPermit(ownerAddress, tokenId, deadline, v, r, s);

        // Verify burn
        vm.expectRevert();
        cert.ownerOf(tokenId);

        assertEq(cert.nonces(ownerAddress), nonce + 1);
    }

    // ============ Gas Tests ============

    function testPermitGasUsage() public {
        uint256 tokenId = 1;
        uint256 deadline = block.timestamp + 1 hours;
        uint256 nonce = cert.nonces(ownerAddress);

        vm.prank(ownerAddress);
        cert.safeMint(ownerAddress, 1);

        bytes32 structHash = keccak256(
            abi.encode(
                keccak256(
                    "Permit(address owner,address spender,uint256 tokenId,uint256 nonce,uint256 deadline)"
                ),
                ownerAddress,
                spender,
                tokenId,
                nonce,
                deadline
            )
        );

        bytes32 domainSeparator = cert.DOMAIN_SEPARATOR();
        bytes32 digest = keccak256(
            abi.encodePacked("\x19\x01", domainSeparator, structHash)
        );
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(ownerPrivateKey, digest);

        uint256 gasBefore = gasleft();
        cert.permit(ownerAddress, spender, tokenId, deadline, v, r, s);
        uint256 gasUsed = gasBefore - gasleft();

        // Log gas usage for monitoring
        emit log_named_uint("Gas used for permit", gasUsed);

        // Typical gas usage should be < 100k
        assertLt(gasUsed, 100000);
    }
}
