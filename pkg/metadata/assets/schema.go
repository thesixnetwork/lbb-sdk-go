package assets

import (
	"embed"
	"fmt"
)

//go:embed schema.json
var schema embed.FS

func GetJSONSchema() ([]byte, error) {
	var schemaByte []byte

	schemaByte, err := schema.ReadFile("schema.json")
	if err != nil {
		return schemaByte, fmt.Errorf("error on reading schema.json file: %+v", err)
	}

	return schemaByte, nil
}
