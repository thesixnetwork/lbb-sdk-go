package client

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func CheckTransactionReceipt(goCtx context.Context, client *ethclient.Client, txhash common.Hash) (*types.Receipt, error) {
	receipt, err := client.TransactionReceipt(goCtx, txhash)
	if err != nil {
		return nil, fmt.Errorf("failed to get receipt: %w", err)
	}

	return receipt, nil
}
