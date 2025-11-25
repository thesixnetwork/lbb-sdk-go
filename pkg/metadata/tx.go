package metadata

import (
	clienttx "github.com/cosmos/cosmos-sdk/client/tx"
)

type MetadataTxFactory struct {
	TxFactory clienttx.Factory
}

func NewFactory
