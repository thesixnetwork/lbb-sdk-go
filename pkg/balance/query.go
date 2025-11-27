package balance

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/thesixnetwork/lbb-sdk-go/account"
)

const (
	SIX_BASE_DENOM = "usix"
	SIX_EVM_DENOM = "asix"
)

type BalanceClient struct {
	account.Account
}

type BalanceClientI interface {
	GetBalance() (sdk.Coins, error)
	GetCosmosBalance() (sdk.Coin, error)
	GetEVMBalance() (sdk.Coin, error)
}

func (a *BalanceClient) GetBalance() (sdk.Coins, error) {
	ctx := a.GetClientCTX()
	queryClient := banktypes.NewQueryClient(ctx)

	res, err := queryClient.AllBalances(a.Context, &banktypes.QueryAllBalancesRequest{
		Address: a.GetCosmosAddress().String(),
	})
	if err != nil {
		return sdk.Coins{}, err
	}

	return res.Balances, nil
}

func (a *BalanceClient) GetCosmosBalance() (sdk.Coin, error) {
	ctx := a.GetClientCTX()
	queryClient := banktypes.NewQueryClient(ctx)

	res, err := queryClient.Balance(a.Context, &banktypes.QueryBalanceRequest{
		Address: a.GetCosmosAddress().String(),
		Denom:   SIX_BASE_DENOM,
	})
	if err != nil {
		return sdk.Coin{}, err
	}

	return *res.Balance, nil
}

func (a *BalanceClient) GetEVMBalance() (sdk.Coin, error) {
	ctx := a.GetClientCTX()
	queryClient := banktypes.NewQueryClient(ctx)

	addr := a.GetEVMAddress().Bytes()
	bech32AccAddress := sdk.AccAddress(addr)

	res, err := queryClient.Balance(a.Context, &banktypes.QueryBalanceRequest{
		Address: bech32AccAddress.String(),
		Denom:   SIX_EVM_DENOM,
	})
	if err != nil {
		return sdk.Coin{}, err
	}

	return *res.Balance, nil
}
