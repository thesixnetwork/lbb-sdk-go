package main

import (
	"context"
	"fmt"

	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/client"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/metadata"
)

func main() {
	// Initialize Client
	client, err := client.NewCustomClient(
		context.Background(),
		"http://localhost:26657",
		"http://localhost:1317",
		"http://localhost:8545",
		"testnet",
	)
	if err != nil {
		panic(fmt.Sprintf("ERROR CREATE CLEINT %v", err))
	}

	a := account.NewAccount(client, "alice", account.TestMnemonic, "")

	// Create MetadataClient
	metadataClient := metadata.NewMetadataClient(*a)

	// Call GetNFTSchema with a sample nftSchemaCode
	nftSchemaCode := "TechSauceVV12.GlobalSummit2025"
	result, err := metadataClient.GetNFTSchema(nftSchemaCode)
	if err != nil {
		fmt.Printf("Error fetching NFT Schema: %v\n", err)
		return
	}

	_ = result
	// Print the result
	// fmt.Printf("NFT Schema Result: %+v \n", result)

	nftdata, err := metadataClient.GetNFTMetadata(nftSchemaCode, "1")
	if err != nil {
		fmt.Printf("Error fetching NFT Schema: %v\n", err)
		return
	}
	fmt.Printf("NFT Schema Result: %+v \n", nftdata)
}
