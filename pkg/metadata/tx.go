package metadata

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/metadata/assets"
	nftmngrtypes "github.com/thesixnetwork/six-protocol/v4/x/nftmngr/types"
)

type MetadataMsg struct {
	account.AccountMsg
	nftSchemaCode string
}

func NewMetadataMsg(a account.Account, nftSchemaCode string) *MetadataMsg {
	return &MetadataMsg{
		AccountMsg:    *account.NewAccountMsg(a),
		nftSchemaCode: nftSchemaCode,
	}
}

func (m *MetadataMsg) BuildDeployMsg() (msg *nftmngrtypes.MsgCreateNFTSchema, err error) {
	var schemaInput nftmngrtypes.NFTSchemaINPUT
	schemaInputBytes, err := assets.GetJSONSchema()
	if err != nil {
		return msg, err
	}

	err = m.Codec.(*codec.ProtoCodec).UnmarshalJSON(schemaInputBytes, &schemaInput)
	if err != nil {
		return msg, err
	}

	schemaName := strings.ReplaceAll(m.nftSchemaCode, ".", "_")
	schemaInput.Code = m.nftSchemaCode
	schemaInput.Owner = m.GetCosmosAddress().String()
	schemaInput.Name = schemaName
	schemaInput.Description = schemaName

	schemaBytes, err := m.Codec.(*codec.ProtoCodec).MarshalJSON(&schemaInput)
	if err != nil {
		return msg, err
	}

	base64Schema := base64.StdEncoding.EncodeToString(schemaBytes)

	msg = &nftmngrtypes.MsgCreateNFTSchema{
		Creator:         m.GetCosmosAddress().String(),
		NftSchemaBase64: base64Schema,
	}

	return msg, nil
}

func (m *MetadataMsg) DeployCertificateSchema() (res *sdk.TxResponse, err error) {
	msg, err := m.BuildDeployMsg()
	if err != nil {
		return res, err
	}

	res, err = m.BroadcastTx(msg)
	if err != nil {
		return res, err
	}

	if res.Code != 0 {
		return res, fmt.Errorf("BroadcastTx error with reason: %v", res.Logs)
	}

	return res, nil
}

func (m *MetadataMsg) BuildMintMetadataMsg(tokenID string) (msg *nftmngrtypes.MsgCreateMetadata, err error) {
	var metadataInput nftmngrtypes.NftData

	metadataBytes, err := assets.GetJSONMetadata()
	if err != nil {
		return msg, err
	}
	err = m.Codec.(*codec.ProtoCodec).UnmarshalJSON(metadataBytes, &metadataInput)
	if err != nil {
		return msg, err
	}

	metadataInput.NftSchemaCode = m.nftSchemaCode
	metadataInput.OwnerAddressType = nftmngrtypes.OwnerAddressType_INTERNAL_ADDRESS
	metadataInput.TokenId = tokenID
	metadataInput.TokenOwner = m.GetCosmosAddress().String()

	metadataBytes, err = m.Codec.(*codec.ProtoCodec).MarshalJSON(&metadataInput)
	if err != nil {
		return msg, err
	}

	base64Metadata := base64.StdEncoding.EncodeToString(metadataBytes)

	msg = &nftmngrtypes.MsgCreateMetadata{
		Creator:       m.GetCosmosAddress().String(),
		NftSchemaCode: m.nftSchemaCode,
		TokenId:       tokenID,
		Base64NFTData: base64Metadata,
	}

	return msg, nil
}

func (m *MetadataMsg) CreateCertificateMetadata(tokenID string) (res *sdk.TxResponse, err error) {
	msg, err := m.BuildMintMetadataMsg(tokenID)
	if err != nil {
		return res, err
	}

	res, err = m.BroadcastTx(msg)
	if err != nil {
		return res, err
	}

	if res.Code != 0 {
		return res, fmt.Errorf("BroadcastTx error with reason: %v", res.Logs)
	}

	return res, nil
}

func (m MetadataMsg) FreezeCertificate(tokenID string) (res *sdk.TxResponse, err error) {
	msg := &nftmngrtypes.MsgPerformActionByAdmin{
		Creator:       m.GetCosmosAddress().String(),
		NftSchemaCode: m.nftSchemaCode,
		TokenId:       tokenID,
		Action:        "freeze_cert",
		RefId:         "",
		Parameters:    []*nftmngrtypes.ActionParameter{},
	}
	return m.BroadcastTx(msg)
}

func (m MetadataMsg) UnfreezeCertificate(tokenID string) (res *sdk.TxResponse, err error) {
	msg := &nftmngrtypes.MsgPerformActionByAdmin{
		Creator:       m.GetCosmosAddress().String(),
		NftSchemaCode: m.nftSchemaCode,
		TokenId:       tokenID,
		Action:        "unfreeze_cert",
		RefId:         "",
		Parameters:    []*nftmngrtypes.ActionParameter{},
	}
	return m.BroadcastTx(msg)
}
