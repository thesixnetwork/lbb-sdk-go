package evm

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/evm/assets"
)

type EVMClient struct {
	account.Account
}

func NewEVMClient(a account.Account) *EVMClient {
	return &EVMClient{
		Account: a,
	}
}

func (e *EVMClient) GasPrice() (*big.Int, error) {
	gasPrice, err := e.ETHClient.SuggestGasPrice(e.GetContext())
	if err != nil {
		return gasPrice, err
	}
	return gasPrice, nil
}

func (e *EVMClient) GasLimit(callMsg ethereum.CallMsg) (uint64, error) {
	gasLimit, err := e.ETHClient.EstimateGas(e.GetContext(), callMsg)
	if err != nil {
		fmt.Printf("ERROR EstimateGas : %v \n", err)
		return gasLimit, err
	}
	gasLimit = gasLimit * 120 / 100
	return gasLimit, nil
}

func (e *EVMClient) EstimateDeployGas(contractABI abi.ABI, bytecode []byte, constructorArgs ...interface{}) (uint64, error) {
	var data []byte
	if len(constructorArgs) > 0 {
		packedArgs, err := contractABI.Pack("", constructorArgs...)
		if err != nil {
			return 0, fmt.Errorf("failed to pack constructor args: %w", err)
		}
		data = append(bytecode, packedArgs...)
	} else {
		data = bytecode
	}

	// Estimate gas for deployment
	gasLimit, err := e.ETHClient.EstimateGas(e.GetContext(), ethereum.CallMsg{
		From: e.GetEVMAddress(),
		Data: data,
	})
	if err != nil {
		return 0, fmt.Errorf("gas estimation failed: %w", err)
	}

	// Add 20% buffer for safety
	gasLimit = gasLimit * 120 / 100
	return gasLimit, nil
}

func (e *EVMClient) ChainID() (*big.Int, error) {
	chainID, err := e.ETHClient.NetworkID(e.GetContext())
	if err != nil {
		return chainID, err
	}
	return chainID, err
}

func (e *EVMClient) GetNonce() (nonce uint64, err error) {
	nonce, err = e.ETHClient.PendingNonceAt(e.GetContext(), e.GetEVMAddress())
	if err != nil {
		return nonce, err
	}

	return nonce, nil
}


// WaitForTransaction add this line to remove annoying lint for the love god
/*
* NOTE:: ON Production both blocktime on mainnet and testnet are the same, which is 6.3 at maximux
* So I time out must be more than that so I will use 3 blocks at most
*/
func (e *EVMClient) WaitForTransaction(txHash common.Hash) (*types.Receipt, error) {
	fmt.Printf("Waiting for transaction %s to be mined...\n", txHash.Hex())

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timeout := time.After(20 * time.Second)

	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("timeout waiting for transaction to be mined")
		case <-ticker.C:
			receipt, err := e.ETHClient.TransactionReceipt(e.GetContext(), txHash)
			if err == nil {
				if receipt.Status == 0 {
					return receipt, fmt.Errorf("transaction failed")
				}
				fmt.Printf("Transaction mined in block %d\n", receipt.BlockNumber.Uint64())
				return receipt, nil
			}
			// Transaction not yet mined, continue waiting
		}
	}
}

func (e *EVMClient) CheckTransactionReceipt(txHash common.Hash) error {
	receipt, err := e.ETHClient.TransactionReceipt(e.GetContext(), txHash)
	if err != nil {
		return fmt.Errorf("failed to get receipt: %w", err)
	}

	fmt.Printf("Transaction: %s\n", txHash.Hex())
	fmt.Printf("  Block Number: %d\n", receipt.BlockNumber.Uint64())
	fmt.Printf("  Status: %d (1=success, 0=failed)\n", receipt.Status)
	fmt.Printf("  Gas Used: %d\n", receipt.GasUsed)
	fmt.Printf("  Contract Address: %s\n", receipt.ContractAddress.Hex())

	if receipt.Status == 0 {
		return fmt.Errorf("transaction failed")
	}

	return nil
}

func (e *EVMClient) DynamicABI(contractAddress common.Address, functionName string, args interface{}) (tx *types.Transaction, err error) {
	stringABI, err := assets.GetContractABIString()
	if err != nil {
		return &types.Transaction{}, err
	}

	contractABI, err := abi.JSON(strings.NewReader(stringABI))
	if err != nil {
		return &types.Transaction{}, err
	}

	// Pack the function call
	data, err := contractABI.Pack(functionName, args)
	if err != nil {
		return nil, fmt.Errorf("failed to pack data: %w", err)
	}

	gasLimit, err := e.GasLimit(ethereum.CallMsg{
		From: e.GetEVMAddress(),
		To:   &contractAddress,
		Data: data,
	})
	if err != nil {
		return &types.Transaction{}, err
	}

	nonce, err := e.GetNonce()
	if err != nil {
		return &types.Transaction{}, err
	}

	gasPrice, err := e.GasPrice()
	if err != nil {
		return &types.Transaction{}, err
	}

	tx = types.NewTransaction(nonce, contractAddress, big.NewInt(0), gasLimit, gasPrice, data)

	chainID, err := e.ChainID()
	if err != nil {
		return &types.Transaction{}, err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), e.GetPrivateKey())
	if err != nil {
		return &types.Transaction{}, err
	}

	err = e.ETHClient.SendTransaction(e.GetContext(), signedTx)
	if err != nil {
		return &types.Transaction{}, err
	}

	return signedTx, nil
}
