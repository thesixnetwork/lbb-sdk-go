package balance

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/thesixnetwork/lbb-sdk-go/account"
)

type BalanceMsg struct {
	account.AccountMsg
}

func NewBalanceMsg(a account.Account) *BalanceMsg {
	return &BalanceMsg{
		AccountMsg: *account.NewAccountMsg(a),
	}
}

func (b *BalanceMsg) SendBalance(dest string, amount sdk.Coins) (res *sdk.TxResponse, err error) {
	sendMsg := &banktypes.MsgSend{
		FromAddress: b.GetCosmosAddress().String(),
		ToAddress:   dest,
		Amount:      amount,
	}

	return b.BroadcastTx(sendMsg)
}
