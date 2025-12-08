package metadata

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/thesixnetwork/lbb-sdk-go/account"
	nftmngrtypes "github.com/thesixnetwork/six-protocol/v4/x/nftmngr/types"
)

type Metadata struct {
	account account.Account
}

type MetadataI interface {
	GetNFTSchema(string) (nftmngrtypes.NFTSchemaQueryResult, error)
	GetNFTMetadata(string, string) (nftmngrtypes.NftData, error)
	GetExecutor(string) ([]string, error)
	GetIsExecutor(string, string) (bool, error)
	GetAccount() account.Account
}

var _ MetadataI = (*Metadata)(nil)

func NewMetadata(a account.Account) *Metadata {
	return &Metadata{
		account: a,
	}
}

func (m Metadata) GetCodec() codec.BinaryCodec {
	return m.account.GetClient().GetClientCTX().Codec
}

func (m *Metadata) GetNFTSchema(nftSchemaCode string) (nftmngrtypes.NFTSchemaQueryResult, error) {
	goCtx := m.account.GetClient().GetContext()
	clientCtx := m.account.GetClient().GetClientCTX()

	queryClient := nftmngrtypes.NewQueryClient(clientCtx)

	res, err := queryClient.NFTSchema(
		goCtx,
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

func (m *Metadata) GetNFTMetadata(nftSchemaCode, tokenID string) (nftmngrtypes.NftData, error) {
	goCtx := m.account.GetClient().GetContext()
	clientCtx := m.account.GetClient().GetClientCTX()

	queryClient := nftmngrtypes.NewQueryClient(clientCtx)

	res, err := queryClient.NftData(goCtx, &nftmngrtypes.QueryGetNftDataRequest{
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

func (m *Metadata) GetExecutor(nftSchemaCode string) ([]string, error) {
	goCtx := m.account.GetClient().GetContext()
	clientCtx := m.account.GetClient().GetClientCTX()

	queryClient := nftmngrtypes.NewQueryClient(clientCtx)

	res, err := queryClient.ExecutorOfSchema(
		goCtx,
		&nftmngrtypes.QueryGetExecutorOfSchemaRequest{NftSchemaCode: nftSchemaCode},
	)
	if err != nil {
		return []string{}, err
	}

	var executor []string

	executor = append(executor, res.ExecutorOfSchema.ExecutorAddress...)

	return executor, nil
}

func (m *Metadata) GetIsExecutor(nftSchemaCode, executorAddress string) (bool, error) {
	goCtx := m.account.GetClient().GetContext()
	clientCtx := m.account.GetClient().GetClientCTX()
	queryClient := nftmngrtypes.NewQueryClient(clientCtx)

	res, err := queryClient.ActionExecutor(goCtx, &nftmngrtypes.QueryGetActionExecutorRequest{
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

func (m *Metadata) GetAccount() account.Account {
	return m.account
}
