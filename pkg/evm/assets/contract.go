package assets

import (
	"embed"
	"fmt"
)

//go:embed contract.abi
var contractABI embed.FS

func GetContractABIBytes() ([]byte, error) {
	var contractABIByte []byte

	contractABIByte, err := contractABI.ReadFile("contract.abi")
	if err != nil {
		return contractABIByte, fmt.Errorf("error on reading contract.abi file: %+v", err)
	}

	return contractABIByte, nil
}

func GetContractABIString() (abi string, err error) {
	var stringABI string
	abiBytes, err := GetContractABIBytes()
	if err != nil {
		return stringABI, err
	}

	stringABI = string(abiBytes)

	return stringABI, err
}

//go:embed contract.bin
var contractBIN embed.FS

func GetContractBINBytes() ([]byte, error) {
	var contractBINByte []byte

	contractBINByte, err := contractBIN.ReadFile("contract.bin")
	if err != nil {
		return contractBINByte, fmt.Errorf("error on reading contract.abi file: %+v", err)
	}

	return contractBINByte, nil
}

func GetContractBINString() (bin string, err error) {
	var stringBIN string
	binBytes, err := GetContractBINBytes()
	if err != nil {
		return stringBIN, err
	}

	stringBIN = string(binBytes)

	return stringBIN, err
}
