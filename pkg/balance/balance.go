package balance

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/thesixnetwork/lbb-sdk-go/account"
)

const (
	BaseDenom = "usix"
	EVMDenom  = "asix"
)

type Balance struct {
	account account.Account
}

type BalanceI interface {
	GetBalance() (sdk.Coins, error)
	GetCosmosBalance() (sdk.Coin, error)
	GetEVMBalance() (sdk.Coin, error)
	GetAccount() account.Account
}

var _ BalanceI = (*Balance)(nil)

func NewBalance(acc account.Account) *Balance {
	return &Balance{
		account: acc,
	}
}

func (b *Balance) GetAccount() account.Account {
	return b.account
}

// GetBalance retrieves all balances for the account
func (b *Balance) GetBalance() (sdk.Coins, error) {
	goCtx := b.account.GetClient().GetContext()
	clientCtx := b.account.GetClient().GetClientCTX()
	queryClient := banktypes.NewQueryClient(clientCtx)

	res, err := queryClient.AllBalances(goCtx, &banktypes.QueryAllBalancesRequest{
		Address: b.account.GetCosmosAddress().String(),
	})
	if err != nil {
		return sdk.Coins{}, err
	}

	return res.Balances, nil
}

// GetCosmosBalance retrieves the Cosmos native token balance (usix)
func (b *Balance) GetCosmosBalance() (sdk.Coin, error) {
	goCtx := b.account.GetClient().GetContext()
	clientCtx := b.account.GetClient().GetClientCTX()
	queryClient := banktypes.NewQueryClient(clientCtx)

	res, err := queryClient.Balance(goCtx, &banktypes.QueryBalanceRequest{
		Address: b.account.GetCosmosAddress().String(),
		Denom:   BaseDenom,
	})
	if err != nil {
		return sdk.Coin{}, err
	}

	return *res.Balance, nil
}

// GetEVMBalance retrieves the EVM token balance (asix)
func (b *Balance) GetEVMBalance() (sdk.Coin, error) {
	goCtx := b.account.GetClient().GetContext()
	clientCtx := b.account.GetClient().GetClientCTX()
	queryClient := banktypes.NewQueryClient(clientCtx)

	addr := b.account.GetEVMAddress().Bytes()
	bech32AccAddress := sdk.AccAddress(addr)

	res, err := queryClient.Balance(goCtx, &banktypes.QueryBalanceRequest{
		Address: bech32AccAddress.String(),
		Denom:   EVMDenom,
	})
	if err != nil {
		return sdk.Coin{}, err
	}

	return *res.Balance, nil
}

// GetBalanceByDenom retrieves the balance for a specific denomination
func (b *Balance) GetBalanceByDenom(denom string) (sdk.Coin, error) {
	goCtx := b.account.GetClient().GetContext()
	clientCtx := b.account.GetClient().GetClientCTX()
	queryClient := banktypes.NewQueryClient(clientCtx)

	res, err := queryClient.Balance(goCtx, &banktypes.QueryBalanceRequest{
		Address: b.account.GetCosmosAddress().String(),
		Denom:   denom,
	})
	if err != nil {
		return sdk.Coin{}, err
	}

	return *res.Balance, nil
}
