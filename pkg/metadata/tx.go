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
	Metadata
	accountMsg *account.AccountMsg
	nftSchemaCode string
}

func NewMetadataMsg(a account.Account, nftSchemaCode string) (*MetadataMsg, error) {
	accountMsg, err := account.NewAccountMsg(&a)
	if err != nil {
		return nil, err
	}

	return &MetadataMsg{
		Metadata: Metadata{
			account: a,
		},
		accountMsg:    accountMsg,
		nftSchemaCode: nftSchemaCode,
	}, nil
}

func (m *MetadataMsg) BroadcastTx(msgs ...sdk.Msg) (*sdk.TxResponse, error) {
	return m.accountMsg.BroadcastTx(msgs...)
}

func (m *MetadataMsg) BuildDeployMsg() (msg *nftmngrtypes.MsgCreateNFTSchema, err error) {
	var schemaInput nftmngrtypes.NFTSchemaINPUT
	schemaInputBytes, err := assets.GetJSONSchema()
	if err != nil {
		return msg, err
	}

	err = m.GetCodec().(*codec.ProtoCodec).UnmarshalJSON(schemaInputBytes, &schemaInput)
	if err != nil {
		return msg, err
	}

	schemaName := strings.ReplaceAll(m.nftSchemaCode, ".", "_")
	schemaInput.Code = m.nftSchemaCode
	schemaInput.Owner = m.account.GetCosmosAddress().String()
	schemaInput.Name = schemaName
	schemaInput.Description = schemaName

	schemaBytes, err := m.GetCodec().(*codec.ProtoCodec).MarshalJSON(&schemaInput)
	if err != nil {
		return msg, err
	}

	base64Schema := base64.StdEncoding.EncodeToString(schemaBytes)

	msg = &nftmngrtypes.MsgCreateNFTSchema{
		Creator:         m.account.GetCosmosAddress().String(),
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
	err = m.GetCodec().(*codec.ProtoCodec).UnmarshalJSON(metadataBytes, &metadataInput)
	if err != nil {
		return msg, err
	}

	metadataInput.NftSchemaCode = m.nftSchemaCode
	metadataInput.OwnerAddressType = nftmngrtypes.OwnerAddressType_INTERNAL_ADDRESS
	metadataInput.TokenId = tokenID
	metadataInput.TokenOwner = m.account.GetCosmosAddress().String()

	metadataBytes, err = m.GetCodec().(*codec.ProtoCodec).MarshalJSON(&metadataInput)
	if err != nil {
		return msg, err
	}

	base64Metadata := base64.StdEncoding.EncodeToString(metadataBytes)

	msg = &nftmngrtypes.MsgCreateMetadata{
		Creator:       m.account.GetCosmosAddress().String(),
		NftSchemaCode: m.nftSchemaCode,
		TokenId:       tokenID,
		Base64NFTData: base64Metadata,
	}

	return msg, nil
}

func (m *MetadataMsg) BuildMintMetadataWithInfoMsg(tokenID string, info CertificateInfo) (msg *nftmngrtypes.MsgCreateMetadata, err error) {
	var metadataInput nftmngrtypes.NftData

	metadataBytes, err := assets.GetJSONMetadata()
	if err != nil {
		return msg, err
	}
	err = m.GetCodec().(*codec.ProtoCodec).UnmarshalJSON(metadataBytes, &metadataInput)
	if err != nil {
		return msg, err
	}

	metadataInput.NftSchemaCode = m.nftSchemaCode
	metadataInput.OwnerAddressType = nftmngrtypes.OwnerAddressType_INTERNAL_ADDRESS
	metadataInput.TokenId = tokenID
	metadataInput.TokenOwner = m.account.GetCosmosAddress().String()

	for _, attr := range metadataInput.OnchainAttributes {
		switch attr.Name {
		case "status":
			attr.Value = &nftmngrtypes.NftAttributeValue_StringAttributeValue{
				StringAttributeValue: &nftmngrtypes.StringAttributeValue{
					Value: info.Status,
				},
			}
		case "gold_standard":
			attr.Value = &nftmngrtypes.NftAttributeValue_StringAttributeValue{
				StringAttributeValue: &nftmngrtypes.StringAttributeValue{
					Value: info.GoldStandard,
				},
			}
		case "cert_number":
			attr.Value = &nftmngrtypes.NftAttributeValue_StringAttributeValue{
				StringAttributeValue: &nftmngrtypes.StringAttributeValue{
					Value: info.CertNumber,
				},
			}
		case "customer_id":
			attr.Value = &nftmngrtypes.NftAttributeValue_StringAttributeValue{
				StringAttributeValue: &nftmngrtypes.StringAttributeValue{
					Value: info.CustomerID,
				},
			}
		case "issue_date":
			attr.Value = &nftmngrtypes.NftAttributeValue_StringAttributeValue{
				StringAttributeValue: &nftmngrtypes.StringAttributeValue{
					Value: info.IssueDate,
				},
			}
		}
	}

	metadataBytes, err = m.GetCodec().(*codec.ProtoCodec).MarshalJSON(&metadataInput)
	if err != nil {
		return msg, err
	}

	base64Metadata := base64.StdEncoding.EncodeToString(metadataBytes)

	msg = &nftmngrtypes.MsgCreateMetadata{
		Creator:       m.account.GetCosmosAddress().String(),
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

func (m *MetadataMsg) CreateCertificateMetadataWithInfo(tokenID string, info CertificateInfo) (res *sdk.TxResponse, err error) {
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
		Creator:       m.account.GetCosmosAddress().String(),
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
		Creator:       m.account.GetCosmosAddress().String(),
		NftSchemaCode: m.nftSchemaCode,
		TokenId:       tokenID,
		Action:        "unfreeze_cert",
		RefId:         "",
		Parameters:    []*nftmngrtypes.ActionParameter{},
	}
	return m.BroadcastTx(msg)
}
