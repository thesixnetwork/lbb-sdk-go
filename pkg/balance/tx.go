package balance

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/thesixnetwork/lbb-sdk-go/account"
)

type BalanceMsg struct {
	Balance
	accountMsg *account.AccountMsg
}

type BalanceMsgI interface {
	BalanceI
	BuildSendMsg(dest string, amount sdk.Coins) (*banktypes.MsgSend, error)
	SendBalance(dest string, amount sdk.Coins) (*sdk.TxResponse, error)
	SendBalanceAndWait(dest string, amount sdk.Coins) (*sdk.TxResponse, error)
	BroadcastTx(msgs ...sdk.Msg) (*sdk.TxResponse, error)
	WithGas(gas uint64) *BalanceMsg
	WithGasAdjustment(gasAdjustment float64) *BalanceMsg
	WithGasPrices(gasPrices string) *BalanceMsg
	WithFees(fees string) *BalanceMsg
	WithMemo(memo string) *BalanceMsg
	WithTimeoutHeight(timeoutHeight uint64) *BalanceMsg
}

var _ BalanceMsgI = (*BalanceMsg)(nil)

func NewBalanceMsg(acc account.Account) (*BalanceMsg, error) {
	accountMsg, err := account.NewAccountMsg(&acc)
	if err != nil {
		return nil, err
	}

	return &BalanceMsg{
		Balance: Balance{
			account: acc,
		},
		accountMsg: accountMsg,
	}, nil
}

// BuildSendMsg builds a MsgSend message without broadcasting
func (b *BalanceMsg) BuildSendMsg(dest string, amount sdk.Coins) (*banktypes.MsgSend, error) {
	msg := &banktypes.MsgSend{
		FromAddress: b.account.GetCosmosAddress().String(),
		ToAddress:   dest,
		Amount:      amount,
	}
	return msg, nil
}

// SendBalance sends tokens to a destination address
func (b *BalanceMsg) SendBalance(dest string, amount sdk.Coins) (*sdk.TxResponse, error) {
	sendMsg, err := b.BuildSendMsg(dest, amount)
	if err != nil {
		return nil, err
	}

	return b.accountMsg.BroadcastTx(sendMsg)
}

// SendBalanceAndWait sends tokens and waits for the transaction to be confirmed
func (b *BalanceMsg) SendBalanceAndWait(dest string, amount sdk.Coins) (*sdk.TxResponse, error) {
	sendMsg, err := b.BuildSendMsg(dest, amount)
	if err != nil {
		return nil, err
	}

	return b.accountMsg.BroadcastTxAndWait(sendMsg)
}

// BroadcastTx broadcasts one or more messages
// This allows for batch operations or custom message types
func (b *BalanceMsg) BroadcastTx(msgs ...sdk.Msg) (*sdk.TxResponse, error) {
	return b.accountMsg.BroadcastTx(msgs...)
}

// NOTE: THESE ARE UTILITIES METHOD ALLOW use to modify tx factory setting of the package.

// WithGas returns a new BalanceMsg with the specified gas limit
func (b *BalanceMsg) WithGas(gas uint64) *BalanceMsg {
	newBalanceMsg := *b
	newBalanceMsg.accountMsg = newBalanceMsg.accountMsg.WithGas(gas)
	return &newBalanceMsg
}

// WithGasAdjustment returns a new BalanceMsg with the specified gas adjustment
func (b *BalanceMsg) WithGasAdjustment(gasAdjustment float64) *BalanceMsg {
	newBalanceMsg := *b
	newBalanceMsg.accountMsg = newBalanceMsg.accountMsg.WithGasAdjustment(gasAdjustment)
	return &newBalanceMsg
}

// WithGasPrices returns a new BalanceMsg with the specified gas prices
func (b *BalanceMsg) WithGasPrices(gasPrices string) *BalanceMsg {
	newBalanceMsg := *b
	newBalanceMsg.accountMsg = newBalanceMsg.accountMsg.WithGasPrices(gasPrices)
	return &newBalanceMsg
}

// WithFees returns a new BalanceMsg with the specified fees
func (b *BalanceMsg) WithFees(fees string) *BalanceMsg {
	newBalanceMsg := *b
	newBalanceMsg.accountMsg = newBalanceMsg.accountMsg.WithFees(fees)
	return &newBalanceMsg
}

// WithMemo returns a new BalanceMsg with the specified memo
func (b *BalanceMsg) WithMemo(memo string) *BalanceMsg {
	newBalanceMsg := *b
	newBalanceMsg.accountMsg = newBalanceMsg.accountMsg.WithMemo(memo)
	return &newBalanceMsg
}

// WithTimeoutHeight returns a new BalanceMsg with the specified timeout height
func (b *BalanceMsg) WithTimeoutHeight(timeoutHeight uint64) *BalanceMsg {
	newBalanceMsg := *b
	newBalanceMsg.accountMsg = newBalanceMsg.accountMsg.WithTimeoutHeight(timeoutHeight)
	return &newBalanceMsg
}
