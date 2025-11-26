package metadata

import (
	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/thesixnetwork/lbb-sdk-go/client"
)

type MetadataClient struct {
	client.ClientI
}

func NewMetadataClient(c client.ClientI) *MetadataClient {
	return &MetadataClient{
		ClientI: c,
	}
}

func (mc *MetadataClient) GetContext() client.Context {
	return mc.ClientI.GetContext()
}

func (mc *MetadataClient) GetCosmosClientCTX() cosmosclient.Context {
	return mc.ClientI.GetCosmosClientCTX()
}
