package metadata

import (
	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/cosmos/cosmos-sdk/codec"
)

type MetadataClient struct {
	account.Account
}

func NewMetadataClient(a account.Account) *MetadataClient {
	return &MetadataClient{
		a,
	}
}

func (m MetadataClient) GetCodec() codec.BinaryCodec {
	return m.Codec 
}
