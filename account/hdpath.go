package account

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethaccounts "github.com/ethereum/go-ethereum/accounts"
)

var (
	CoinType    uint32 = 60
	BIP44HDPath string = ethaccounts.DefaultBaseDerivationPath.String()
)

type (
	HDPathIterator func() ethaccounts.DerivationPath
)

// HDPathIterator receives a base path as a string and a boolean for the desired iterator type and
// returns a function that iterates over the base HD path, returning the string.
func NewHDPathIterator(basePath string) (HDPathIterator, error) {
	hdPath, err := ethaccounts.ParseDerivationPath(basePath)
	if err != nil {
		return nil, err
	}

	return ethaccounts.DefaultIterator(hdPath), nil
}

func GetFullBIP44Path() string {
	return fmt.Sprintf("m/%d'/%d'/0'/0/0", sdk.Purpose, sdk.CoinType)
}
