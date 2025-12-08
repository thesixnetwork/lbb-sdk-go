package metadata

import (
	nftmngrtypes "github.com/thesixnetwork/six-protocol/v4/x/nftmngr/types"

	"github.com/thesixnetwork/lbb-sdk-go/pkg/metadata/assets"
)

func GetSchemaByteFromJSON() ([]byte, error) {
	schemaBytes, err := assets.GetJSONSchema()
	if err != nil {
		return []byte{}, err
	}

	return schemaBytes, nil
}

func GetMetadataByteFromJSON() ([]byte, error) {
	metadataByte, err := assets.GetJSONMetadata()
	if err != nil {
		return []byte{}, err
	}

	return metadataByte, nil
}

func NewMockSchema() nftmngrtypes.NFTSchemaINPUT {
	newSchema := nftmngrtypes.NFTSchemaINPUT{
		Code:        "sixprotocol.divine_elite",
		Name:        "Main Collection",
		Description: "Description of this schema",
		Owner:       "6x1myrlxmmasv6yq4axrxmdswj9kv5gc0ppx95rmq",
		OriginData: &nftmngrtypes.OriginData{
			OriginBaseUri:         "https://google.com/",
			UriRetrievalMethod:    nftmngrtypes.URIRetrievalMethod_TOKEN,
			OriginChain:           "SIXNET",
			OriginContractAddress: "0x00000000000000000000000000000000000",
			AttributeOverriding:   nftmngrtypes.AttributeOverriding_CHAIN,
			MetadataFormat:        "opensea",
			OriginAttributes: []*nftmngrtypes.AttributeDefinition{
				{
					Name: "something",
					DefaultMintValue: &nftmngrtypes.DefaultMintValue{
						Value: &nftmngrtypes.DefaultMintValue_StringAttributeValue{
							StringAttributeValue: &nftmngrtypes.StringAttributeValue{
								Value: "Something You see if nothing set",
							},
						},
					},
					DataType:          "string",
					Required:          true,
					DisplayValueField: "Somthing You See",
					DisplayOption: &nftmngrtypes.DisplayOption{
						Opensea: &nftmngrtypes.OpenseaDisplayOption{
							DisplayType: "string",
							TraitType:   "Name of Something You See",
							MaxValue:    0,
						},
					},
					HiddenOveride:       false,
					HiddenToMarketplace: false,
				},
			},
		},
		OnchainData: &nftmngrtypes.OnChainData{
			NftAttributes: []*nftmngrtypes.AttributeDefinition{
				{
					Name: "something global",
					DefaultMintValue: &nftmngrtypes.DefaultMintValue{
						Value: &nftmngrtypes.DefaultMintValue_StringAttributeValue{
							StringAttributeValue: &nftmngrtypes.StringAttributeValue{
								Value: "Something Global You see if nothing set",
							},
						},
					},
					DataType:          "string",
					Required:          true,
					DisplayValueField: "Somthing Global You See",
					DisplayOption: &nftmngrtypes.DisplayOption{
						Opensea: &nftmngrtypes.OpenseaDisplayOption{
							DisplayType: "string",
							TraitType:   "Name of Something Global You See",
							MaxValue:    0,
						},
					},
					HiddenOveride:       false,
					HiddenToMarketplace: false,
				},
			},
			TokenAttributes: []*nftmngrtypes.AttributeDefinition{
				{
					Name: "something per token",
					DefaultMintValue: &nftmngrtypes.DefaultMintValue{
						Value: &nftmngrtypes.DefaultMintValue_StringAttributeValue{
							StringAttributeValue: &nftmngrtypes.StringAttributeValue{
								Value: "Something per token You see if nothing set",
							},
						},
					},
					DataType:          "string",
					Required:          true,
					DisplayValueField: "Somthing per token You See",
					DisplayOption: &nftmngrtypes.DisplayOption{
						Opensea: &nftmngrtypes.OpenseaDisplayOption{
							DisplayType: "string",
							TraitType:   "Name of Something per token You See",
							MaxValue:    0,
						},
					},
					HiddenOveride:       false,
					HiddenToMarketplace: false,
				},
			},
			Actions: []*nftmngrtypes.Action{},
		},
		IsVerified:        false,
		MintAuthorization: "system",
	}

	return newSchema
}
