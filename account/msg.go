package account

import (
	"errors"
	"fmt"
	"os"

	clienttx "github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	GasLimit      = uint64(300000) // Higher gas limit for SixProtocol operations
	GasPrice      = "1.25usix"
	GasAdjustment = 1.5
)

type AccountMsg struct {
	Account
	clienttx.Factory
}

func NewAccountMsg(a Account) *AccountMsg {
	// use account info to set clintContext
	a.CosmosClientCTX = a.GetClientCTX().WithFromName(a.accountName).WithFrom(a.accountName).WithFromAddress(a.GetCosmosAddress())

	ctx := a.CosmosClientCTX

	factory := clienttx.Factory{}.
		WithTxConfig(ctx.TxConfig).
		WithAccountRetriever(ctx.AccountRetriever).
		WithChainID(a.ChainID).
		WithGas(GasLimit).
		// WithFees(Fee)
		WithGasPrices(GasPrice).
		WithKeybase(ctx.Keyring).
		WithFromName(ctx.FromName).
		// WithSequence()
		// WithAccountNumber(a.CosmosClientCTX.Account)
		// WithMemo()
		WithGasAdjustment(GasAdjustment).
		// WithSignMode()
		// WithTimeoutHeight()
		WithFeeGranter(ctx.GetFeeGranterAddress()).
		WithFeePayer(ctx.FeePayer)
	return &AccountMsg{
		Account: a,
		Factory: factory,
	}
}

func (a *AccountMsg) GenerateOrBroadcastTxWithFactory(msgs ...sdk.Msg) error {
	ctx := a.GetClientCTX()

	return clienttx.GenerateOrBroadcastTxWithFactory(ctx, a.Factory, msgs...)
}

func (a *AccountMsg) BroadcastTx(msgs ...sdk.Msg) (res *sdk.TxResponse, err error) {
	ctx := a.GetClientCTX()

	// Validate account and address
	if a.cosmosAddress.Empty() {
		return &sdk.TxResponse{}, fmt.Errorf("account cosmos address is empty, account name: %s", a.accountName)
	}

	if len(msgs) == 0 {
		return &sdk.TxResponse{}, errors.New("no messages provided to broadcast")
	}

	txf, err := a.Prepare(a.CosmosClientCTX)
	if err != nil {
		return &sdk.TxResponse{}, fmt.Errorf("failed to prepare transaction factory: %w", err)
	}

	if txf.SimulateAndExecute() || ctx.Simulate {
		if ctx.Offline {
			return &sdk.TxResponse{}, errors.New("cannot estimate gas in offline mode")
		}

		_, adjusted, err := clienttx.CalculateGas(ctx, txf, msgs...)
		if err != nil {
			return &sdk.TxResponse{}, fmt.Errorf("failed to calculate gas for transaction (from: %s): %w", a.cosmosAddress.String(), err)
		}

		txf = txf.WithGas(adjusted)
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", clienttx.GasEstimateResponse{GasEstimate: txf.Gas()})
	}

	if ctx.Simulate {
		return &sdk.TxResponse{}, nil
	}

	tx, err := txf.BuildUnsignedTx(msgs...)
	if err != nil {
		return &sdk.TxResponse{}, fmt.Errorf("failed to build unsigned transaction (from: %s, gas: %d): %w", a.cosmosAddress.String(), txf.Gas(), err)
	}

	if err = clienttx.Sign(ctx.CmdContext, txf, ctx.FromName, tx, true); err != nil {
		return &sdk.TxResponse{}, fmt.Errorf("failed to sign transaction (account: %s, from: %s): %w", ctx.FromName, a.cosmosAddress.String(), err)
	}

	txBytes, err := ctx.TxConfig.TxEncoder()(tx.GetTx())
	if err != nil {
		return &sdk.TxResponse{}, fmt.Errorf("failed to encode transaction: %w", err)
	}

	res, err = ctx.BroadcastTx(txBytes)
	if err != nil {
		return &sdk.TxResponse{}, fmt.Errorf("failed to broadcast transaction (from: %s, chain: %s): %w", a.cosmosAddress.String(), a.ChainID, err)
	}

	if res.Code != 0 {
		return res, fmt.Errorf("transaction failed with code %d: %s (from: %s, txhash: %s)", res.Code, res.RawLog, a.cosmosAddress.String(), res.TxHash)
	}

	return res, nil
}
