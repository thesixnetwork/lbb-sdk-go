package metadata

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
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

func (m MetadataClient) WaitForTransaction(txhash string) error {
	fmt.Printf("Waiting for transaction %s to be mined...\n", txhash)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timeout := time.After(20 * time.Second)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for transaction to be mined")
		case <-ticker.C:
			output, err := authtx.QueryTx(m.CosmosClientCTX, txhash)
			if err != nil {
				return err
			}
			if output.Empty() {
				return fmt.Errorf("no transaction found with hash %s", txhash)
			}
			return nil
		}
		// continue
	}
}
