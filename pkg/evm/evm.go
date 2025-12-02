package evm

import (
	"github.com/thesixnetwork/lbb-sdk-go/account"
)

type EVMClient struct {
	account.Account
}

func NewEVMClient(a account.Account) *EVMClient {
	return &EVMClient{
		Account: a,
	}
}

