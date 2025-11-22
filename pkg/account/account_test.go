package account

import (
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	client "github.com/thesixnetwork/lbb-sdk-go/client"
)

const (
	// NOTE: (@ddeedev) this mnemonic is for testing purposes only. Do NOT use it in production.
	TestMnemonic         = "history perfect across group seek acoustic delay captain sauce audit carpet tattoo exhaust green there giant cluster want pond bulk close screen scissors remind"
	TestPassword         = "testpassword"
	InvalidMnemonic      = "invalid mnemonic phrase"
	TestPrivateKey       = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	TestPrivateKeyWith0x = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
)

// Helper function to create a test account service
// This creates a minimal account service that doesn't require cosmos SDK context
func createTestAccountService() AccountI {
	// Create a simple context without full SDK initialization to avoid codec issues
	ctx := client.Context{}
	return &Account{
		ctx:                  ctx,
		accountName:          "test-account",
		accountAddressPrefix: "6x",
	}
}

func TestGenerateMnemonic(t *testing.T) {
	t.Run("Generate valid mnemonic", func(t *testing.T) {
		mnemonic, err := GenerateMnemonic()

		require.NoError(t, err, "Should generate mnemonic without error")
		require.NotEmpty(t, mnemonic, "Generated mnemonic should not be empty")

		words := strings.Split(mnemonic, " ")
		assert.Equal(t, 24, len(words), "Generated mnemonic should have 24 words")

		// Validate the generated mnemonic
		account := createTestAccountService()
		assert.True(t, account.ValidateMnemonic(mnemonic), "Generated mnemonic should be valid")

		t.Logf("Generated mnemonic: %s...", getFirstWords(mnemonic, 3))
	})

	t.Run("Generate unique mnemonics", func(t *testing.T) {
		const numTests = 5
		mnemonics := make(map[string]bool)

		for range numTests {
			mnemonic, err := GenerateMnemonic()
			require.NoError(t, err, "Should generate mnemonic without error")

			// Check uniqueness
			assert.False(t, mnemonics[mnemonic], "Generated mnemonic should be unique")
			mnemonics[mnemonic] = true
		}

		assert.Equal(t, numTests, len(mnemonics), "All generated mnemonics should be unique")
		t.Logf("Generated %d unique mnemonics", numTests)
	})
}

func TestValidateMnemonic(t *testing.T) {
	account := createTestAccountService()

	testCases := []struct {
		name     string
		mnemonic string
		expected bool
	}{
		{"Valid 12-word test mnemonic", TestMnemonic, true},
		{"Invalid mnemonic phrase", InvalidMnemonic, false},
		{"Empty mnemonic", "", false},
		{"Single word", "abandon", false},
		{"Valid 24-word generated mnemonic", "", true}, // Will be set in test
		{"Actually valid repeated words", "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon", true},
	}

	// Generate a valid 24-word mnemonic for testing
	validMnemonic, err := GenerateMnemonic()
	require.NoError(t, err)
	testCases[4].mnemonic = validMnemonic

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := account.ValidateMnemonic(tc.mnemonic)
			assert.Equal(t, tc.expected, result, "Validation result should match expected")

			if tc.expected {
				t.Logf("Correctly validated mnemonic: %s...", getFirstWords(tc.mnemonic, 3))
			} else {
				t.Log("Correctly rejected invalid mnemonic")
			}
		})
	}
}

func TestCreateEVMAccountFromMnemonic(t *testing.T) {
	account := createTestAccountService()

	t.Run("Create EVM account from valid mnemonic", func(t *testing.T) {
		address, err := account.CreateEVMAccountFromMnemonic(TestMnemonic, TestPassword)

		require.NoError(t, err, "Should create EVM account without error")
		assert.NotEqual(t, common.Address{}, address, "Generated address should not be empty")
		assert.True(t, common.IsHexAddress(address.Hex()), "Generated address should be valid hex")

		t.Logf("Generated EVM address: %s", address.Hex())
	})

	t.Run("Different passwords generate different addresses", func(t *testing.T) {
		address1, err := account.CreateEVMAccountFromMnemonic(TestMnemonic, "password1")
		require.NoError(t, err)

		address2, err := account.CreateEVMAccountFromMnemonic(TestMnemonic, "password2")
		require.NoError(t, err)

		assert.NotEqual(t, address1, address2, "Different passwords should generate different addresses")
		t.Logf("Address1: %s, Address2: %s", address1.Hex(), address2.Hex())
	})

	t.Run("Invalid mnemonic should return error", func(t *testing.T) {
		_, err := account.CreateEVMAccountFromMnemonic(InvalidMnemonic, TestPassword)
		assert.Error(t, err, "Should return error for invalid mnemonic")
		assert.Contains(t, err.Error(), "invalid mnemonic", "Error should mention invalid mnemonic")
	})
}

func TestGetPrivateKeyFromMnemonic(t *testing.T) {
	account := createTestAccountService()

	t.Run("Extract private key from valid mnemonic", func(t *testing.T) {
		privateKey, err := account.GetPrivateKeyFromMnemonic(TestMnemonic, TestPassword)

		require.NoError(t, err, "Should extract private key without error")
		assert.NotEmpty(t, privateKey, "Private key should not be empty")
		assert.Equal(t, 64, len(privateKey), "Private key should be 64 characters (32 bytes hex)")

		t.Logf("Generated private key: %s...", privateKey[:16])
	})

	t.Run("Invalid mnemonic should return error", func(t *testing.T) {
		_, err := account.GetPrivateKeyFromMnemonic(InvalidMnemonic, TestPassword)
		assert.Error(t, err, "Should return error for invalid mnemonic")
		assert.Contains(t, err.Error(), "invalid mnemonic", "Error should mention invalid mnemonic")
	})
}

func TestCreateEVMAccountFromPrivateKey(t *testing.T) {
	account := createTestAccountService()

	t.Run("Create account from private key without 0x prefix", func(t *testing.T) {
		address, err := account.CreateEVMAccountFromPrivateKey(TestPrivateKey, TestPassword)

		require.NoError(t, err, "Should create account without error")
		assert.NotEmpty(t, address, "Address should not be empty")
		assert.True(t, common.IsHexAddress(address), "Should be valid hex address")
		assert.True(t, strings.HasPrefix(address, "0x"), "Address should have 0x prefix")

		// This should generate the known address for this test private key
		expectedAddress := "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
		assert.Equal(t, expectedAddress, address, "Should generate expected address")

		t.Logf("Generated address: %s", address)
	})

	t.Run("Create account from private key with 0x prefix", func(t *testing.T) {
		address1, err1 := account.CreateEVMAccountFromPrivateKey(TestPrivateKey, TestPassword)
		address2, err2 := account.CreateEVMAccountFromPrivateKey(TestPrivateKeyWith0x, TestPassword)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.Equal(t, address1, address2, "Should generate same address regardless of 0x prefix")

		t.Log("0x prefix handling works correctly")
	})

	t.Run("Invalid private key formats should return error", func(t *testing.T) {
		invalidKeys := []string{
			"invalid_hex",
			"123",    // Too short
			"xyz123", // Invalid hex characters
			"not_a_hex_string",
		}

		for _, invalidKey := range invalidKeys {
			_, err := account.CreateEVMAccountFromPrivateKey(invalidKey, TestPassword)
			assert.Error(t, err, "Should return error for invalid private key: %s", invalidKey)
		}

		t.Log("Invalid private key formats correctly rejected")
	})
}

func TestNewHDPathIterator(t *testing.T) {
	t.Run("Test BIP44 path format", func(t *testing.T) {
		expectedPath := GetFullBIP44Path()
		// Don't assume specific values, just check it's a valid path
		assert.NotEmpty(t, expectedPath, "BIP44 path should not be empty")
		assert.Contains(t, expectedPath, "m/", "Should be a derivation path")
		assert.Equal(t, uint32(60), CoinType, "Coin type should be 60 for Ethereum")

		t.Logf("BIP44 path: %s", expectedPath)
	})

	t.Run("Test HD path iterator functionality", func(t *testing.T) {
		basePath := "m/44'/60'/0'/0"

		iterator, err := NewHDPathIterator(basePath)
		require.NoError(t, err, "Should create iterator without error")
		require.NotNil(t, iterator, "Iterator should not be nil")

		// Test path generation by collecting multiple paths
		var paths []string
		const interationNum int = 5
		for range interationNum {
			path := iterator()
			paths = append(paths, path.String())
		}

		// Verify we got 5 paths
		require.Len(t, paths, 5, "Should generate 5 paths")

		// Verify they follow the expected pattern
		assert.Equal(t, "m/44'/60'/0'/0", paths[0], "First path should be the base path")
		assert.Equal(t, "m/44'/60'/0'/1", paths[1], "Second path should increment account index")
		assert.Equal(t, "m/44'/60'/0'/2", paths[2], "Third path should continue incrementing")
		assert.Equal(t, "m/44'/60'/0'/3", paths[3], "Fourth path should continue incrementing")
		assert.Equal(t, "m/44'/60'/0'/4", paths[4], "Fifth path should continue incrementing")

		// Verify all paths are unique
		pathSet := make(map[string]bool)
		for _, path := range paths {
			assert.False(t, pathSet[path], "Path %s should be unique", path)
			pathSet[path] = true
		}

		t.Logf("HD path iterator works correctly:")
		t.Logf("  Base: %s", basePath)
		for i, path := range paths {
			t.Logf("  Path %d: %s", i, path)
		}
	})

	t.Run("Test invalid HD path handling", func(t *testing.T) {
		invalidPaths := []string{
			"invalid/path",
			"not-a-derivation-path",
			"",
			"m/invalid",
		}

		for _, invalidPath := range invalidPaths {
			_, err := NewHDPathIterator(invalidPath)
			assert.Error(t, err, "Should return error for invalid path: %s", invalidPath)
		}

		t.Log("Invalid HD paths correctly rejected")
	})

	t.Run("Test coin type constant", func(t *testing.T) {
		assert.Equal(t, uint32(60), CoinType, "Coin type should be 60 for Ethereum")
		t.Log("Ethereum coin type (60) validated")
	})
}

func TestAccountConsistency(t *testing.T) {
	account := createTestAccountService()

	t.Run("Private key from mnemonic matches address generation", func(t *testing.T) {
		// Generate address from mnemonic
		address1, err := account.CreateEVMAccountFromMnemonic(TestMnemonic, TestPassword)
		require.NoError(t, err)

		// Get private key from same mnemonic
		privateKey, err := account.GetPrivateKeyFromMnemonic(TestMnemonic, TestPassword)
		require.NoError(t, err)

		// Generate address from private key
		address2Hex, err := account.CreateEVMAccountFromPrivateKey(privateKey, TestPassword)
		require.NoError(t, err)

		address2 := common.HexToAddress(address2Hex)
		assert.Equal(t, address1, address2, "Addresses should match when derived from same mnemonic")

		t.Logf("Consistency validated: %s", address1.Hex())
	})
}

func TestWithEnvironmentVariables(t *testing.T) {
	// Set up environment variables
	os.Setenv("EVM_FIVENET_RPC", "https://rpc-evm.fivenet.sixprotocol.net:443")
	os.Setenv("FIVENET_RPC", "https://rpc1.fivenet.sixprotocol.net:443")
	os.Setenv("FIVENET_API", "https://api1.fivenet.sixprotocol.net:443")
	defer func() {
		os.Unsetenv("EVM_FIVENET_RPC")
		os.Unsetenv("FIVENET_RPC")
		os.Unsetenv("FIVENET_API")
	}()

	t.Run("Test environment variables are set", func(t *testing.T) {
		t.Logf("Environment context:")
		t.Logf("  EVM_FIVENET_RPC: %s", os.Getenv("EVM_FIVENET_RPC"))
		t.Logf("  FIVENET_RPC: %s", os.Getenv("FIVENET_RPC"))
		t.Logf("  FIVENET_API: %s", os.Getenv("FIVENET_API"))

		// Test basic account functionality without full SDK context
		mnemonic, err := GenerateMnemonic()
		require.NoError(t, err)

		account := createTestAccountService()
		isValid := account.ValidateMnemonic(mnemonic)
		assert.True(t, isValid, "Generated mnemonic should be valid")

		t.Log("Environment test completed successfully")
	})
}

// Benchmark tests
func BenchmarkGenerateMnemonic(b *testing.B) {
	for range b.N {
		_, err := GenerateMnemonic()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCreateEVMAccountFromMnemonic(b *testing.B) {
	account := createTestAccountService()
	for range b.N {
		_, err := account.CreateEVMAccountFromMnemonic(TestMnemonic, TestPassword)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCreateEVMAccountFromPrivateKey(b *testing.B) {
	account := createTestAccountService()
	for range b.N {
		_, err := account.CreateEVMAccountFromPrivateKey(TestPrivateKey, TestPassword)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNewHDPathIterator(b *testing.B) {
	basePath := "m/44'/60'/0'/0"
	for range b.N {
		iterator, err := NewHDPathIterator(basePath)
		if err != nil {
			b.Fatal(err)
		}
		_ = iterator()
	}
}

// Helper function to get first N words from mnemonic for logging
func getFirstWords(mnemonic string, n int) string {
	words := strings.Split(mnemonic, " ")
	if len(words) <= n {
		return mnemonic
	}
	return strings.Join(words[:n], " ")
}
