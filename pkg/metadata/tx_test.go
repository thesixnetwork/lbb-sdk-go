package metadata_test

import (
	"testing"

	nftmngrtypes "github.com/thesixnetwork/six-protocol/v4/x/nftmngr/types"

	"github.com/thesixnetwork/lbb-sdk-go/account"
	"github.com/thesixnetwork/lbb-sdk-go/pkg/metadata"
)

func TestMetadataMsg_BuildMintMetadataWithInfoMsg(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		a             account.Account
		nftSchemaCode string
		// Named input parameters for target function.
		tokenID string
		info    metadata.CertificateInfo
		want    *nftmngrtypes.MsgCreateMetadata
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := metadata.NewMetadataMsg(tt.a, tt.nftSchemaCode)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NewMetadataMsg failed: %v", err)
				}
				return
			}
			got, gotErr := m.BuildMintMetadataWithInfoMsg(tt.tokenID, tt.info)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("BuildMintMetadataWithInfoMsg() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("BuildMintMetadataWithInfoMsg() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("BuildMintMetadataWithInfoMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}
