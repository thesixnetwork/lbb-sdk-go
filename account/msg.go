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
	ctx := a.GetClientCTX().WithFromName(a.accountName)

	factory := clienttx.Factory{}.
		WithTxConfig(ctx.TxConfig).
		WithAccountRetriever(ctx.AccountRetriever).
		WithChainID(a.ChainID).
		WithGas(GasLimit).
		// WithFees(Fee)
		WithGasPrices(GasPrice).
		WithKeybase(a.GetKeyring()).
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
	txf, err := a.Prepare(a.CosmosClientCTX)
	if err != nil {
		return &sdk.TxResponse{}, err
	}

	if txf.SimulateAndExecute() || ctx.Simulate {
		if ctx.Offline {
			return &sdk.TxResponse{}, errors.New("cannot estimate gas in offline mode")
		}

		_, adjusted, err := clienttx.CalculateGas(ctx, txf, msgs...)
		if err != nil {
			return &sdk.TxResponse{}, err
		}

		txf = txf.WithGas(adjusted)
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", clienttx.GasEstimateResponse{GasEstimate: txf.Gas()})
	}

	if ctx.Simulate {
		return &sdk.TxResponse{}, nil
	}

	tx, err := txf.BuildUnsignedTx(msgs...)
	if err != nil {
		return &sdk.TxResponse{}, err
	}

	if err = clienttx.Sign(ctx.CmdContext, txf, ctx.FromName, tx, true); err != nil {
		return &sdk.TxResponse{}, err
	}

	txBytes, err := ctx.TxConfig.TxEncoder()(tx.GetTx())
	if err != nil {
		return &sdk.TxResponse{}, err
	}

	return ctx.BroadcastTx(txBytes)
}
