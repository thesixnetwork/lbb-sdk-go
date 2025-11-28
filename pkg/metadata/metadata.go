package metadata

import (
	"github.com/thesixnetwork/lbb-sdk-go/account"
)

type MetadataClient struct {
	account.Account
}

func NewMetadataClient(a account.Account) *MetadataClient {
	return &MetadataClient{
		a,
	}
}
