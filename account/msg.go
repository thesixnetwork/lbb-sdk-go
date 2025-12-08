package account

import (
	"errors"
	"fmt"
	"os"

	clienttx "github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// Default gas settings for SixProtocol operations
	GasLimit      = uint64(1000000)
	GasPrice      = "1.25usix"
	GasAdjustment = 1.5
)

type AccountMsgI interface {
	BroadcastTx(msgs ...sdk.Msg) (*sdk.TxResponse, error)
	GenerateOrBroadcastTxWithFactory(msgs ...sdk.Msg) error
	GetFactory() clienttx.Factory
}

type AccountMsg struct {
	account AccountI
	factory clienttx.Factory
}

var _ AccountMsgI = (*AccountMsg)(nil)

// NewAccountMsg creates a new AccountMsg instance for transaction operations
func NewAccountMsg(acc AccountI) (*AccountMsg, error) {
	if acc == nil {
		return nil, fmt.Errorf("account cannot be nil")
	}

	// Get the underlying account to access client and account details
	account, ok := acc.(*Account)
	if !ok {
		return nil, fmt.Errorf("account must be of type *Account")
	}

	// Get client context
	client := account.GetClient()
	if client == nil {
		return nil, fmt.Errorf("account client cannot be nil")
	}

	ctx := account.client.GetClientCTX().
		WithFromName(account.GetAccountName()).
		WithFrom(account.GetAccountName()).
		WithFromAddress(account.GetCosmosAddress())

	ctx = client.SetClientCTX(ctx)

	// Create transaction factory with account settings
	factory := clienttx.Factory{}.
		WithTxConfig(ctx.TxConfig).
		WithAccountRetriever(ctx.AccountRetriever).
		WithChainID(client.GetChainID()).
		WithGas(GasLimit).
		WithGasPrices(GasPrice).
		WithKeybase(ctx.Keyring).
		WithFromName(ctx.FromName).
		WithGasAdjustment(GasAdjustment).
		WithFeeGranter(ctx.GetFeeGranterAddress()).
		WithFeePayer(ctx.FeePayer)

	return &AccountMsg{
		account: account,
		factory: factory,
	}, nil
}

// GetFactory returns the transaction factory
func (a *AccountMsg) GetFactory() clienttx.Factory {
	return a.factory
}

// GetAccount returns the underlying account
func (a *AccountMsg) GetAccount() AccountI {
	return a.account
}

// GenerateOrBroadcastTxWithFactory generates or broadcasts a transaction using the factory
func (a *AccountMsg) GenerateOrBroadcastTxWithFactory(msgs ...sdk.Msg) error {
	if len(msgs) == 0 {
		return errors.New("no messages provided to broadcast")
	}

	account, ok := a.account.(*Account)
	if !ok {
		return fmt.Errorf("account must be of type *Account")
	}

	ctx := account.client.GetClientCTX()
	return clienttx.GenerateOrBroadcastTxWithFactory(ctx, a.factory, msgs...)
}

// BroadcastTx builds, signs, and broadcasts a transaction with the provided messages
func (a *AccountMsg) BroadcastTx(msgs ...sdk.Msg) (*sdk.TxResponse, error) {
	if len(msgs) == 0 {
		return nil, errors.New("no messages provided to broadcast")
	}

	// Get the underlying account
	account, ok := a.account.(*Account)
	if !ok {
		return nil, fmt.Errorf("account must be of type *Account")
	}

	// Validate account address
	if account.cosmosAddress.Empty() {
		return nil, fmt.Errorf("account cosmos address is empty, account name: %s", account.GetAccountName())
	}

	// Get client context
	ctx := account.client.GetClientCTX()

	// Prepare the transaction factory with account and sequence numbers
	txf, err := a.factory.Prepare(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare transaction factory: %w", err)
	}

	// Handle gas estimation if needed
	if txf.SimulateAndExecute() || ctx.Simulate {
		if ctx.Offline {
			return nil, errors.New("cannot estimate gas in offline mode")
		}

		_, adjusted, err := clienttx.CalculateGas(ctx, txf, msgs...)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate gas for transaction (from: %s): %w",
				account.cosmosAddress.String(), err)
		}

		txf = txf.WithGas(adjusted)
		_, _ = fmt.Fprintf(os.Stderr, "estimated gas: %d\n", txf.Gas())
	}

	// If simulation mode, return early
	if ctx.Simulate {
		return &sdk.TxResponse{}, nil
	}

	// Build unsigned transaction
	tx, err := txf.BuildUnsignedTx(msgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to build unsigned transaction (from: %s, gas: %d): %w",
			account.cosmosAddress.String(), txf.Gas(), err)
	}

	// Sign the transaction
	if err = clienttx.Sign(ctx.CmdContext, txf, ctx.FromName, tx, true); err != nil {
		return nil, fmt.Errorf("failed to sign transaction (account: %s, from: %s): %w",
			ctx.FromName, account.cosmosAddress.String(), err)
	}

	// Encode the transaction
	txBytes, err := ctx.TxConfig.TxEncoder()(tx.GetTx())
	if err != nil {
		return nil, fmt.Errorf("failed to encode transaction: %w", err)
	}

	// Broadcast the transaction
	res, err := ctx.BroadcastTx(txBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to broadcast transaction (from: %s, chain: %s): %w",
			account.cosmosAddress.String(), account.client.GetChainID(), err)
	}

	// Check transaction result
	if res.Code != 0 {
		return res, fmt.Errorf("transaction failed with code %d: %s (from: %s, txhash: %s)",
			res.Code, res.RawLog, account.cosmosAddress.String(), res.TxHash)
	}

	fmt.Printf("Transaction broadcast successfully\n")
	fmt.Printf("  TxHash: %s\n", res.TxHash)
	fmt.Printf("  Code: %d\n", res.Code)
	fmt.Printf("  Gas Used: %d\n", res.GasUsed)

	return res, nil
}

// BroadcastTxAndWait broadcasts a transaction and waits for it to be mined
func (a *AccountMsg) BroadcastTxAndWait(msgs ...sdk.Msg) (*sdk.TxResponse, error) {
	res, err := a.BroadcastTx(msgs...)
	if err != nil {
		return res, err
	}

	account, ok := a.account.(*Account)
	if !ok {
		return res, fmt.Errorf("account must be of type *Account")
	}

	// Wait for transaction to be mined
	if err := account.client.WaitForTransaction(res.TxHash); err != nil {
		return res, fmt.Errorf("transaction broadcast succeeded but confirmation failed: %w", err)
	}

	return res, nil
}

// NOTE: THESE ARE UTILITIES METHOD (to modify default tx factory setting)
// WithGas returns a new AccountMsg with the specified gas limit
func (a *AccountMsg) WithGas(gas uint64) *AccountMsg {
	newAccountMsg := *a
	newAccountMsg.factory = newAccountMsg.factory.WithGas(gas)
	return &newAccountMsg
}

// WithGasAdjustment returns a new AccountMsg with the specified gas adjustment
func (a *AccountMsg) WithGasAdjustment(gasAdjustment float64) *AccountMsg {
	newAccountMsg := *a
	newAccountMsg.factory = newAccountMsg.factory.WithGasAdjustment(gasAdjustment)
	return &newAccountMsg
}

// WithGasPrices returns a new AccountMsg with the specified gas prices
func (a *AccountMsg) WithGasPrices(gasPrices string) *AccountMsg {
	newAccountMsg := *a
	newAccountMsg.factory = newAccountMsg.factory.WithGasPrices(gasPrices)
	return &newAccountMsg
}

// WithFees returns a new AccountMsg with the specified fees
func (a *AccountMsg) WithFees(fees string) *AccountMsg {
	newAccountMsg := *a
	newAccountMsg.factory = newAccountMsg.factory.WithFees(fees)
	return &newAccountMsg
}

// WithMemo returns a new AccountMsg with the specified memo
func (a *AccountMsg) WithMemo(memo string) *AccountMsg {
	newAccountMsg := *a
	newAccountMsg.factory = newAccountMsg.factory.WithMemo(memo)
	return &newAccountMsg
}

// WithTimeoutHeight returns a new AccountMsg with the specified timeout height
func (a *AccountMsg) WithTimeoutHeight(timeoutHeight uint64) *AccountMsg {
	newAccountMsg := *a
	newAccountMsg.factory = newAccountMsg.factory.WithTimeoutHeight(timeoutHeight)
	return &newAccountMsg
}
