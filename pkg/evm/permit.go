package evm

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"

	"github.com/thesixnetwork/lbb-sdk-go/pkg/evm/assets"
)

// PermitSignature contains the signature components for a permit
type PermitSignature struct {
	V        uint8
	R        [32]byte
	S        [32]byte
	Deadline *big.Int
}

// EIP712Domain represents the domain separator parameters
type EIP712Domain struct {
	Name              string
	Version           string
	ChainID           *big.Int
	VerifyingContract common.Address
}

// SignPermitForAll creates an EIP-712 signature for permitForAll (gasless setApprovalForAll)
// This allows a user to sign offline and have anyone broadcast the approval transaction
func (e *EVMClient) SignPermitForAll(
	contractName string,
	contractAddress common.Address,
	operator common.Address,
	approved bool,
	deadline *big.Int,
) (*PermitSignature, error) {
	chainID, err := e.ChainID()
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Get current nonce
	nonce, err := e.GetPermitNonce(contractAddress, e.GetEVMAddress())
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	// Create EIP-712 typed data
	typedData := apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"PermitForAll": []apitypes.Type{
				{Name: "owner", Type: "address"},
				{Name: "operator", Type: "address"},
				{Name: "approved", Type: "bool"},
				{Name: "nonce", Type: "uint256"},
				{Name: "deadline", Type: "uint256"},
			},
		},
		PrimaryType: "PermitForAll",
		Domain: apitypes.TypedDataDomain{
			Name:              contractName,
			Version:           "1",
			ChainId:           (*math.HexOrDecimal256)(chainID),
			VerifyingContract: contractAddress.Hex(),
		},
		Message: apitypes.TypedDataMessage{
			"owner":    e.GetEVMAddress().Hex(),
			"operator": operator.Hex(),
			"approved": approved,
			"nonce":    (*math.HexOrDecimal256)(nonce),
			"deadline": (*math.HexOrDecimal256)(deadline),
		},
	}

	// Hash the typed data
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return nil, fmt.Errorf("failed to hash domain: %w", err)
	}

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, fmt.Errorf("failed to hash message: %w", err)
	}

	// Create the final hash
	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	hash := crypto.Keccak256Hash(rawData)

	// Sign the hash
	signature, err := crypto.Sign(hash.Bytes(), e.GetPrivateKey())
	if err != nil {
		return nil, fmt.Errorf("failed to sign: %w", err)
	}

	// Extract v, r, s
	var r, s [32]byte
	copy(r[:], signature[:32])
	copy(s[:], signature[32:64])
	v := signature[64] + 27 // Ethereum uses 27/28 instead of 0/1

	return &PermitSignature{
		V:        v,
		R:        r,
		S:        s,
		Deadline: deadline,
	}, nil
}

// SignPermit creates an EIP-712 signature for permit (gasless approval for specific token)
func (e *EVMClient) SignPermit(
	contractName string,
	contractAddress common.Address,
	spender common.Address,
	tokenID *big.Int,
	deadline *big.Int,
) (*PermitSignature, error) {
	chainID, err := e.ChainID()
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	nonce, err := e.GetPermitNonce(contractAddress, e.GetEVMAddress())
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	typedData := apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"Permit": []apitypes.Type{
				{Name: "owner", Type: "address"},
				{Name: "spender", Type: "address"},
				{Name: "tokenId", Type: "uint256"},
				{Name: "nonce", Type: "uint256"},
				{Name: "deadline", Type: "uint256"},
			},
		},
		PrimaryType: "Permit",
		Domain: apitypes.TypedDataDomain{
			Name:              contractName,
			Version:           "1",
			ChainId:           (*math.HexOrDecimal256)(chainID),
			VerifyingContract: contractAddress.Hex(),
		},
		Message: apitypes.TypedDataMessage{
			"owner":    e.GetEVMAddress().Hex(),
			"spender":  spender.Hex(),
			"tokenId":  (*math.HexOrDecimal256)(tokenID),
			"nonce":    (*math.HexOrDecimal256)(nonce),
			"deadline": (*math.HexOrDecimal256)(deadline),
		},
	}

	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return nil, fmt.Errorf("failed to hash domain: %w", err)
	}

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, fmt.Errorf("failed to hash message: %w", err)
	}

	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	hash := crypto.Keccak256Hash(rawData)

	signature, err := crypto.Sign(hash.Bytes(), e.GetPrivateKey())
	if err != nil {
		return nil, fmt.Errorf("failed to sign: %w", err)
	}

	var r, s [32]byte
	copy(r[:], signature[:32])
	copy(s[:], signature[32:64])
	v := signature[64] + 27

	return &PermitSignature{
		V:        v,
		R:        r,
		S:        s,
		Deadline: deadline,
	}, nil
}

// GetPermitNonce gets the current nonce for an address from the contract
func (e *EVMClient) GetPermitNonce(contractAddress common.Address, owner common.Address) (*big.Int, error) {
	goCtx := e.GetClient().GetContext()
	ethClient := e.GetClient().GetETHClient()

	// Call the nonces(address) function
	data := crypto.Keccak256([]byte("nonces(address)"))[:4]
	paddedAddress := common.LeftPadBytes(owner.Bytes(), 32)
	data = append(data, paddedAddress...)

	msg := map[string]interface{}{
		"to":   contractAddress.Hex(),
		"data": "0x" + common.Bytes2Hex(data),
	}

	var result string
	err := ethClient.Client().CallContext(goCtx, &result, "eth_call", msg, "latest")
	if err != nil {
		return nil, fmt.Errorf("failed to call nonces: %w", err)
	}

	nonce := new(big.Int)
	nonce.SetString(result[2:], 16) // Remove "0x" prefix
	return nonce, nil
}

// VerifyPermitSignature verifies that a permit signature is valid (useful for testing)
func VerifyPermitSignature(
	owner common.Address,
	signature *PermitSignature,
	hash []byte,
) (bool, error) {
	// Convert signature back to bytes
	sig := make([]byte, 65)
	copy(sig[:32], signature.R[:])
	copy(sig[32:64], signature.S[:])
	sig[64] = signature.V - 27

	// Recover the public key
	pubKey, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return false, fmt.Errorf("failed to recover public key: %w", err)
	}

	// Get the address from public key
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	return recoveredAddr == owner, nil
}

// SignPermitMessage signs a raw message with the private key (general purpose)
func SignPermitMessage(privateKey *ecdsa.PrivateKey, hash []byte) (*PermitSignature, error) {
	signature, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %w", err)
	}

	var r, s [32]byte
	copy(r[:], signature[:32])
	copy(s[:], signature[32:64])
	v := signature[64] + 27

	return &PermitSignature{
		V: v,
		R: r,
		S: s,
	}, nil
}

// ExecutePermitForAll broadcasts a permitForAll transaction using a pre-signed signature
// Anyone can call this and pay for gas, but the approval is for the owner who signed
func (e *EVMClient) ExecutePermitForAll(
	contractAddress common.Address,
	owner common.Address,
	operator common.Address,
	approved bool,
	signature *PermitSignature,
) (*types.Transaction, error) {
	goCtx := e.GetClient().GetContext()
	ethClient := e.GetClient().GetETHClient()

	// Get contract ABI
	stringABI, err := assets.GetContractABIString()
	if err != nil {
		return nil, fmt.Errorf("failed to get contract ABI: %w", err)
	}

	contractABI, err := abi.JSON(strings.NewReader(stringABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	// Pack the permitForAll function call
	data, err := contractABI.Pack("permitForAll", owner, operator, approved, signature.Deadline, signature.V, signature.R, signature.S)
	if err != nil {
		return nil, fmt.Errorf("failed to pack permitForAll: %w", err)
	}

	// Estimate gas
	gasLimit, err := e.GasLimit(ethereum.CallMsg{
		From: e.GetEVMAddress(),
		To:   &contractAddress,
		Data: data,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas: %w", err)
	}

	nonce, err := e.GetNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	gasPrice, err := e.GasPrice()
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	// Create transaction
	tx := types.NewTransaction(nonce, contractAddress, big.NewInt(0), gasLimit, gasPrice, data)

	chainID, err := e.ChainID()
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Sign with broadcaster's key (they pay for gas)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), e.GetPrivateKey())
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Broadcast
	err = ethClient.SendTransaction(goCtx, signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	return signedTx, nil
}

// TransferWithPermit transfers an NFT using a permit signature (completely gasless for owner)
// The broadcaster pays for gas to execute both the permit and transfer in one transaction
func (e *EVMClient) TransferWithPermit(
	contractAddress common.Address,
	from common.Address,
	to common.Address,
	tokenID *big.Int,
	signature *PermitSignature,
) (*types.Transaction, error) {
	goCtx := e.GetClient().GetContext()
	ethClient := e.GetClient().GetETHClient()

	// Get contract ABI
	stringABI, err := assets.GetContractABIString()
	if err != nil {
		return nil, fmt.Errorf("failed to get contract ABI: %w", err)
	}

	contractABI, err := abi.JSON(strings.NewReader(stringABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	// Pack the transferWithPermit function call
	data, err := contractABI.Pack("transferWithPermit", from, to, tokenID, signature.Deadline, signature.V, signature.R, signature.S)
	if err != nil {
		return nil, fmt.Errorf("failed to pack transferWithPermit: %w", err)
	}

	// Estimate gas
	gasLimit, err := e.GasLimit(ethereum.CallMsg{
		From: e.GetEVMAddress(),
		To:   &contractAddress,
		Data: data,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas: %w", err)
	}

	nonce, err := e.GetNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	gasPrice, err := e.GasPrice()
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	// Create transaction
	tx := types.NewTransaction(nonce, contractAddress, big.NewInt(0), gasLimit, gasPrice, data)

	chainID, err := e.ChainID()
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Sign with broadcaster's key (they pay for gas)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), e.GetPrivateKey())
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Broadcast
	err = ethClient.SendTransaction(goCtx, signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	return signedTx, nil
}
