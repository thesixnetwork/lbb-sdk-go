package config

import (
	"cosmossdk.io/x/tx/signing"
	amino "github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/gogoproto/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	enccodec "github.com/evmos/evmos/v20/encoding/codec"
	"github.com/evmos/evmos/v20/ethereum/eip712"
	evmtypes "github.com/evmos/evmos/v20/x/evm/types"

	// cosmos modules
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	// SixProtocol modules
	nftadminmoduletypes "github.com/thesixnetwork/six-protocol/v4/x/nftadmin/types"
	nftmngrmoduletypes "github.com/thesixnetwork/six-protocol/v4/x/nftmngr/types"
	nftoraclemoduletypes "github.com/thesixnetwork/six-protocol/v4/x/nftoracle/types"
	protocoladminmoduletypes "github.com/thesixnetwork/six-protocol/v4/x/protocoladmin/types"
	tokenmngrmoduletypes "github.com/thesixnetwork/six-protocol/v4/x/tokenmngr/types"
)

// encodingConfig creates a new EncodingConfig and returns it
func MakeConfig() sdktestutil.TestEncodingConfig {
	cdc := amino.NewLegacyAmino()
	signingOptions := signing.Options{
		AddressCodec: address.Bech32Codec{
			Bech32Prefix: sdk.GetConfig().GetBech32AccountAddrPrefix(),
		},
		ValidatorAddressCodec: address.Bech32Codec{
			Bech32Prefix: sdk.GetConfig().GetBech32ValidatorAddrPrefix(),
		},
		CustomGetSigners: map[protoreflect.FullName]signing.GetSignersFunc{
			evmtypes.MsgEthereumTxCustomGetSigner.MsgType: evmtypes.MsgEthereumTxCustomGetSigner.Fn,
		},
	}

	interfaceRegistry, _ := types.NewInterfaceRegistryWithOptions(types.InterfaceRegistryOptions{
		ProtoFiles:     proto.HybridResolver,
		SigningOptions: signingOptions,
	})
	codec := amino.NewProtoCodec(interfaceRegistry)
	enccodec.RegisterLegacyAminoCodec(cdc)
	enccodec.RegisterInterfaces(interfaceRegistry)

	// Register standard Cosmos modules
	authtypes.RegisterInterfaces(interfaceRegistry)
	banktypes.RegisterInterfaces(interfaceRegistry)

	// Register SixProtocol modules
	nftmngrmoduletypes.RegisterInterfaces(interfaceRegistry)
	tokenmngrmoduletypes.RegisterInterfaces(interfaceRegistry)
	nftadminmoduletypes.RegisterInterfaces(interfaceRegistry)
	nftoraclemoduletypes.RegisterInterfaces(interfaceRegistry)
	protocoladminmoduletypes.RegisterInterfaces(interfaceRegistry)

	// This is needed for the EIP712 txs because currently is using
	// the deprecated method legacytx.StdSignBytes
	legacytx.RegressionTestingAminoCodec = cdc
	eip712.SetEncodingConfig(cdc, interfaceRegistry)

	return sdktestutil.TestEncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             codec,
		TxConfig:          tx.NewTxConfig(codec, tx.DefaultSignModes),
		Amino:             cdc,
	}
}
