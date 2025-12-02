package assets

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type NFT struct {
	NFTCaller     // Read-only binding to the contract
	NFTTransactor // Write-only binding to the contract
	NFTFilterer   // Log filterer for contract events
}

type NFTCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

type NFTTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

type NFTFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

type NFTSession struct {
	Contract     *NFT           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

type NFTCallerSession struct {
	Contract *NFTCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

type NFTTransactorSession struct {
	Contract     *NFTTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

type NFTRaw struct {
	Contract *NFT // Generic contract binding to access the raw methods on
}

type NFTCallerRaw struct {
	Contract *NFTCaller // Generic read-only contract binding to access the raw methods on
}

type NFTTransactorRaw struct {
	Contract *NFTTransactor // Generic write-only contract binding to access the raw methods on
}

func NewNFT(address common.Address, backend bind.ContractBackend) (*NFT, error) {
	contract, err := bindNFT(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NFT{NFTCaller: NFTCaller{contract: contract}, NFTTransactor: NFTTransactor{contract: contract}, NFTFilterer: NFTFilterer{contract: contract}}, nil
}

func NewNFTCaller(address common.Address, caller bind.ContractCaller) (*NFTCaller, error) {
	contract, err := bindNFT(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NFTCaller{contract: contract}, nil
}

func NewNFTTransactor(address common.Address, transactor bind.ContractTransactor) (*NFTTransactor, error) {
	contract, err := bindNFT(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NFTTransactor{contract: contract}, nil
}

func NewNFTFilterer(address common.Address, filterer bind.ContractFilterer) (*NFTFilterer, error) {
	contract, err := bindNFT(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NFTFilterer{contract: contract}, nil
}

func bindNFT(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	stringABI, err := GetContractABIString()
	if err != nil {
		return nil, err
	}
	parsed, err := abi.JSON(strings.NewReader(stringABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NFT *NFTRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NFT.Contract.NFTCaller.contract.Call(opts, result, method, params...)
}


// TODO:: Create Contract Method according to ABI
