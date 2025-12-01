package metadata

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/thesixnetwork/lbb-sdk-go/account"
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

func (m *MetadataMsg) DeployCertificateSchema() (res *sdk.TxResponse, err error) {
	schemaInput, err := GetSchemaInput()
	if err != nil {
		return res, err
	}

	schemaName := strings.ReplaceAll(m.nftSchemaCode, ".", "_")
	schemaInput.Code = m.nftSchemaCode
	schemaInput.Owner = m.GetCosmosAddress().String()
	schemaInput.Name = schemaName
	schemaInput.Description = schemaName

	schemaBytes, err := json.Marshal(&schemaInput)
	if err != nil {
		return res, err
	}

	base64Schema := base64.StdEncoding.EncodeToString(schemaBytes)

	msg := &nftmngrtypes.MsgCreateNFTSchema{
		Creator:         m.GetCosmosAddress().String(),
		NftSchemaBase64: base64Schema,
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

func (m *MetadataMsg) CreateCertificateMetadata(tokenID string) (res *sdk.TxResponse, err error) {
	_ = tokenID
	return m.BroadcastTx(nil)
}
