package metadata

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/thesixnetwork/lbb-sdk-go/account"
)

type MetadataClient struct {
	account.Account
}

func NewMetadataClient(a account.Account) *MetadataClient {
	return &MetadataClient{
		Account: a,
	}
}

func (m MetadataClient) GetCodec() codec.BinaryCodec {
	return m.Codec
}
