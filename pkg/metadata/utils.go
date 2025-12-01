package metadata

import (
	"encoding/json"

	"github.com/thesixnetwork/lbb-sdk-go/pkg/metadata/assets"
	nftmngrtypes "github.com/thesixnetwork/six-protocol/v4/x/nftmngr/types"
)

func GetSchemaInput() (*nftmngrtypes.NFTSchemaINPUT, error) {
	var schemaInput *nftmngrtypes.NFTSchemaINPUT

	schemaBytes, err := assets.GetJSONSchema()
	if err != nil {
		return schemaInput, err
	}

	// valiate schema input is valid
	if err := json.Unmarshal(schemaBytes, &schemaInput); err != nil {
		return schemaInput, err
	}


	return schemaInput, nil
}
