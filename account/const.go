package account

import "math/big"

// Chain ID constants
const (
	// COSMOS FORMAT CHAINID
	ChainNameMainnet = "sixnet"
	ChainNameTestnet = "fivenet"

	// DEVELOPMENT COSMOS CHAIN ID
	ChainNameLocalnet = "testnet"

	// REGISTERED EVM CHAIN ID
	ChainIDMainnet = 98
	ChainIDTestnet = 150

	// DEVELOPMENT EVM CHIAN ID
	ChainIDLocalnet = 666

	// CHAID ID EPOCH TO PREVENT DUPLICATE IN/IF MIGRATION FROCESS
	ChainIDEpoch = 1
)

const (
	// NOTE: (@ddeedev) this mnemonic is for testing purposes only. Do NOT use it in production.
	TestMnemonic         = "history perfect across group seek acoustic delay captain sauce audit carpet tattoo exhaust green there giant cluster want pond bulk close screen scissors remind"
	TestPassword         = "testpassword"
	InvalidMnemonic      = "invalid mnemonic phrase"
	TestPrivateKey       = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	TestPrivateKeyWith0x = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
)

type ChainIDTable map[string]*big.Int

var ChainIDMapping ChainIDTable

func init() {
	ChainIDMapping = ChainIDTable{
		ChainNameMainnet:  big.NewInt(ChainIDMainnet),
		ChainNameTestnet:  big.NewInt(ChainIDTestnet),
		ChainNameLocalnet: big.NewInt(ChainIDLocalnet),
	}
}
