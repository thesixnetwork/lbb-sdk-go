package metadata

import (
	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	account "github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/client"
	nftmngrtypes "github.com/thesixnetwork/six-protocol/v4/x/nftmngr/types"
)

type MetadataClient struct {
	account.Account
}

func NewMetadataClient(a account.Account) *MetadataClient {
	return &MetadataClient{
		a,
	}
}

func (mc *MetadataClient) GetClient() client.Client {
	return mc.Client
}

func (mc *MetadataClient) GetClientCTX() cosmosclient.Context {
	return mc.GetClient().CosmosClientCTX
}

func (mc *MetadataClient) GetNFTSchema(nftSchemaCode string) (nftmngrtypes.NFTSchemaQueryResult, error) {
	ctx := mc.GetClientCTX()
	queryClient := nftmngrtypes.NewQueryClient(ctx)

	res, err := queryClient.NFTSchema(
		mc.Context,
		&nftmngrtypes.QueryGetNFTSchemaRequest{Code: nftSchemaCode},
	)
	if err != nil {
		return nftmngrtypes.NFTSchemaQueryResult{}, err
	}

	return nftmngrtypes.NFTSchemaQueryResult{
		Code:              res.NFTSchema.Code,
		Name:              res.NFTSchema.Name,
		Owner:             res.NFTSchema.Owner,
		Description:       res.NFTSchema.Description,
		OriginData:        res.NFTSchema.OriginData,
		OnchainData:       res.NFTSchema.OnchainData,
		IsVerified:        res.NFTSchema.IsVerified,
		MintAuthorization: res.NFTSchema.MintAuthorization,
	}, nil
}

func (mc *MetadataClient) GetNFTMetadata(nftSchemaCode, tokenID string) (nftmngrtypes.NftData, error) {
	ctx := mc.GetClientCTX()
	queryClient := nftmngrtypes.NewQueryClient(ctx)

	res, err := queryClient.NftData(mc.Context, &nftmngrtypes.QueryGetNftDataRequest{
		NftSchemaCode: nftSchemaCode,
		TokenId:       tokenID,
	})
	if err != nil {
		return nftmngrtypes.NftData{}, nil
	}

	return nftmngrtypes.NftData{
		NftSchemaCode:     res.NftData.NftSchemaCode,
		TokenId:           res.NftData.TokenId,
		TokenOwner:        res.NftData.TokenOwner,
		OwnerAddressType:  res.NftData.OwnerAddressType,
		OriginImage:       res.NftData.OriginImage,
		OnchainImage:      res.NftData.OnchainImage,
		TokenUri:          res.NftData.TokenUri,
		OriginAttributes:  res.NftData.OriginAttributes,
		OnchainAttributes: res.NftData.OnchainAttributes,
	}, nil
}

func (mc *MetadataClient) GetExecutor(nftSchemaCode string) ([]string, error) {
	ctx := mc.GetClientCTX()
	queryClient := nftmngrtypes.NewQueryClient(ctx)

	res, err := queryClient.ExecutorOfSchema(
		mc.Context,
		&nftmngrtypes.QueryGetExecutorOfSchemaRequest{NftSchemaCode: nftSchemaCode},
	)
	if err != nil {
		return []string{}, err
	}

	var executor []string

	executor = append(executor, res.ExecutorOfSchema.ExecutorAddress...)

	return executor, nil
}

func (mc *MetadataClient) GetIsExecutor(nftSchemaCode, executorAddress string) (bool, error) {
	ctx := mc.GetClientCTX()
	queryClient := nftmngrtypes.NewQueryClient(ctx)

	res, err := queryClient.ActionExecutor(mc.Context, &nftmngrtypes.QueryGetActionExecutorRequest{
		NftSchemaCode:   nftSchemaCode,
		ExecutorAddress: executorAddress,
	})
	if err != nil {
		return false, err
	}

	if executorAddress == res.ActionExecutor.ExecutorAddress {
		return true, nil
	}

	return false, nil
}
