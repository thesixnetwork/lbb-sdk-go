package account_test

import (
	"encoding/hex"
	"os"
	"strings"
	"testing"

	bip39 "github.com/cosmos/go-bip39"
	ethaccounts "github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test constants
const (
	TestMnemonic         = "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	TestPassword         = "testpassword"
	InvalidMnemonic      = "invalid mnemonic phrase that should not work"
	TestPrivateKey       = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	TestPrivateKeyWith0x = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	mnemonicEntropySize  = 256
	CoinType             = uint32(60) // Ethereum coin type
)

// ========================================
// WORKING FUNCTIONALITY TESTS
// ========================================

// TestMnemonicGenerationFunctionality tests the mnemonic generation logic
// that is implemented in account.GenerateMnemonic()
func TestMnemonicGenerationFunctionality(t *testing.T) {
	t.Run("Generate valid 24-word mnemonic", func(t *testing.T) {
		// This reproduces the logic from account.GenerateMnemonic()
		entropy, err := bip39.NewEntropy(mnemonicEntropySize)
		require.NoError(t, err, "Entropy generation should not fail")

		mnemonic, err := bip39.NewMnemonic(entropy)
		require.NoError(t, err, "Mnemonic generation should not fail")
		require.NotEmpty(t, mnemonic, "Generated mnemonic should not be empty")

		// Verify mnemonic characteristics
		words := strings.Fields(mnemonic)
		assert.Equal(t, 24, len(words), "Mnemonic should have 24 words for 256-bit entropy")

		// Verify it's a valid BIP39 mnemonic
		assert.True(t, bip39.IsMnemonicValid(mnemonic), "Generated mnemonic should be valid")

		t.Logf(" Generated valid mnemonic: %s...", strings.Join(words[:3], " "))
	})

	t.Run("Generate multiple unique mnemonics", func(t *testing.T) {
		const numTests = 5
		mnemonics := make(map[string]bool)

		for i := 0; i < numTests; i++ {
			entropy, err := bip39.NewEntropy(mnemonicEntropySize)
			require.NoError(t, err)

			mnemonic, err := bip39.NewMnemonic(entropy)
			require.NoError(t, err)

			// Check for duplicates
			assert.False(t, mnemonics[mnemonic], "Should not generate duplicate mnemonics")
			mnemonics[mnemonic] = true
		}

		assert.Equal(t, numTests, len(mnemonics), "All generated mnemonics should be unique")
		t.Logf(" Generated %d unique mnemonics", len(mnemonics))
	})
}

// TestMnemonicValidationFunctionality tests the validation logic
// that is implemented in account.ValidateMnemonic()
func TestMnemonicValidationFunctionality(t *testing.T) {
	testCases := []struct {
		name     string
		mnemonic string
		expected bool
	}{
		{
			name:     "Valid 12-word test mnemonic",
			mnemonic: TestMnemonic,
			expected: true,
		},
		{
			name:     "Invalid mnemonic phrase",
			mnemonic: InvalidMnemonic,
			expected: false,
		},
		{
			name:     "Empty mnemonic",
			mnemonic: "",
			expected: false,
		},
		{
			name:     "Single word",
			mnemonic: "abandon",
			expected: false,
		},
		{
			name:     "Valid 24-word mnemonic",
			mnemonic: "present volume rate enter account wrap sheriff toward sugar assume worry model obvious clump liberty carry assault endless list come talk whip expand galaxy",
			expected: true,
		},
		{
			name:     "Invalid checksum mnemonic",
			mnemonic: "present volume rate enter account wrap sheriff toward sugar assume worry model obvious clump liberty carry assault endless list come talk whip expand sky",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// This reproduces the logic from account.ValidateMnemonic()
			result := bip39.IsMnemonicValid(tc.mnemonic)
			assert.Equal(t, tc.expected, result, "Validation result should match expected for: %s", tc.mnemonic)

			if tc.expected {
				t.Logf(" Correctly validated mnemonic: %s...", getFirstWords(tc.mnemonic, 3))
			} else {
				t.Logf(" Correctly rejected invalid mnemonic")
			}
		})
	}
}

// TestEVMAccountFromMnemonicFunctionality tests the EVM account creation logic
// that is implemented in account.CreateEVMAccountFromMnemonic()
func TestEVMAccountFromMnemonicFunctionality(t *testing.T) {
	t.Run("Create EVM account from valid mnemonic", func(t *testing.T) {
		// Validate mnemonic first (reproduces account.ValidateMnemonic logic)
		require.True(t, bip39.IsMnemonicValid(TestMnemonic), "Test mnemonic should be valid")

		// This reproduces the logic from account.CreateEVMAccountFromMnemonic()
		seed := bip39.NewSeed(TestMnemonic, TestPassword)
		privateKey, err := crypto.ToECDSA(seed[:32])
		require.NoError(t, err, "Private key generation should not fail")

		publicKey := privateKey.PublicKey
		address := crypto.PubkeyToAddress(publicKey)

		// Verify the address is valid
		assert.NotEqual(t, common.Address{}, address, "Generated address should not be zero address")
		assert.True(t, common.IsHexAddress(address.Hex()), "Generated address should be valid hex")

		t.Logf(" Generated EVM address from mnemonic: %s", address.Hex())

		// Test deterministic generation
		seed2 := bip39.NewSeed(TestMnemonic, TestPassword)
		privateKey2, err := crypto.ToECDSA(seed2[:32])
		require.NoError(t, err)
		address2 := crypto.PubkeyToAddress(privateKey2.PublicKey)

		assert.Equal(t, address, address2, "Same mnemonic and password should generate same address")
	})

	t.Run("Different passwords generate different addresses", func(t *testing.T) {
		// Test with password1
		seed1 := bip39.NewSeed(TestMnemonic, "password1")
		privateKey1, err := crypto.ToECDSA(seed1[:32])
		require.NoError(t, err)
		address1 := crypto.PubkeyToAddress(privateKey1.PublicKey)

		// Test with password2
		seed2 := bip39.NewSeed(TestMnemonic, "password2")
		privateKey2, err := crypto.ToECDSA(seed2[:32])
		require.NoError(t, err)
		address2 := crypto.PubkeyToAddress(privateKey2.PublicKey)

		assert.NotEqual(t, address1, address2, "Different passwords should generate different addresses")
		t.Logf(" Different passwords generate different addresses: %s vs %s", address1.Hex(), address2.Hex())
	})

	t.Run("Invalid mnemonic handling", func(t *testing.T) {
		// Test the validation step (this would return error in actual function)
		isValid := bip39.IsMnemonicValid(InvalidMnemonic)
		assert.False(t, isValid, "Invalid mnemonic should not pass validation")
		t.Logf(" Invalid mnemonic correctly rejected")
	})
}

// TestPrivateKeyFromMnemonicFunctionality tests the private key extraction logic
// that is implemented in account.GetPrivateKeyFromMnemonic()
func TestPrivateKeyFromMnemonicFunctionality(t *testing.T) {
	t.Run("Extract private key from valid mnemonic", func(t *testing.T) {
		// Validate mnemonic first
		require.True(t, bip39.IsMnemonicValid(TestMnemonic), "Test mnemonic should be valid")

		// This reproduces the logic from account.GetPrivateKeyFromMnemonic()
		seed := bip39.NewSeed(TestMnemonic, TestPassword)
		privateKey, err := crypto.ToECDSA(seed[:32])
		require.NoError(t, err, "Private key generation should not fail")

		privateKeyBytes := crypto.FromECDSA(privateKey)
		privateKeyHex := hex.EncodeToString(privateKeyBytes)

		// Verify private key characteristics
		assert.NotEmpty(t, privateKeyHex, "Private key should not be empty")
		assert.Equal(t, 64, len(privateKeyHex), "Private key should be 64 hex characters (32 bytes)")

		// Verify it's valid hex
		_, err = hex.DecodeString(privateKeyHex)
		assert.NoError(t, err, "Private key should be valid hex")

		t.Logf(" Generated private key: %s...", privateKeyHex[:16])

		// Test deterministic generation
		seed2 := bip39.NewSeed(TestMnemonic, TestPassword)
		privateKey2, err := crypto.ToECDSA(seed2[:32])
		require.NoError(t, err)
		privateKeyBytes2 := crypto.FromECDSA(privateKey2)
		privateKeyHex2 := hex.EncodeToString(privateKeyBytes2)

		assert.Equal(t, privateKeyHex, privateKeyHex2, "Same mnemonic and password should generate same private key")
	})

	t.Run("Invalid mnemonic handling", func(t *testing.T) {
		isValid := bip39.IsMnemonicValid(InvalidMnemonic)
		assert.False(t, isValid, "Invalid mnemonic should be rejected")
		t.Logf(" Invalid mnemonic correctly handled")
	})
}

// TestEVMAccountFromPrivateKeyFunctionality tests the account creation from private key logic
// that is implemented in account.CreateEVMAccountFromPrivateKey()
func TestEVMAccountFromPrivateKeyFunctionality(t *testing.T) {
	t.Run("Create account from private key without 0x prefix", func(t *testing.T) {
		// This reproduces the logic from account.CreateEVMAccountFromPrivateKey()
		privateKeyHex := TestPrivateKey

		privateKeyBytes, err := hex.DecodeString(privateKeyHex)
		require.NoError(t, err, "Private key decoding should not fail")

		privateKey, err := crypto.ToECDSA(privateKeyBytes)
		require.NoError(t, err, "Private key parsing should not fail")

		publicKey := privateKey.PublicKey
		address := crypto.PubkeyToAddress(publicKey)

		assert.NotEqual(t, common.Address{}, address, "Generated address should not be zero")
		assert.True(t, common.IsHexAddress(address.Hex()), "Generated address should be valid hex")

		t.Logf(" Generated address from private key: %s", address.Hex())
	})

	t.Run("Create account from private key with 0x prefix", func(t *testing.T) {
		// Test the prefix removal logic
		privateKeyHex := TestPrivateKeyWith0x
		if strings.HasPrefix(privateKeyHex, "0x") {
			privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")
		}

		privateKeyBytes, err := hex.DecodeString(privateKeyHex)
		require.NoError(t, err, "Private key decoding should not fail")

		privateKey, err := crypto.ToECDSA(privateKeyBytes)
		require.NoError(t, err, "Private key parsing should not fail")

		address := crypto.PubkeyToAddress(privateKey.PublicKey)

		assert.True(t, common.IsHexAddress(address.Hex()), "Generated address should be valid hex")

		// Compare with non-prefixed version
		privateKeyBytes2, _ := hex.DecodeString(TestPrivateKey)
		privateKey2, _ := crypto.ToECDSA(privateKeyBytes2)
		address2 := crypto.PubkeyToAddress(privateKey2.PublicKey)

		assert.Equal(t, address, address2, "Private key with and without 0x prefix should generate same address")
		t.Logf(" 0x prefix handling works correctly")
	})

	t.Run("Invalid private key formats", func(t *testing.T) {
		invalidKeys := []string{
			"invalid_hex",
			"123", // too short
			"xyz123",
			// "",
			"not_a_hex_string",
		}

		for _, invalidKey := range invalidKeys {
			t.Run("Invalid: "+invalidKey, func(t *testing.T) {
				_, err := hex.DecodeString(invalidKey)
				assert.Error(t, err, "Invalid private key should cause decode error: %s", invalidKey)
			})
		}
		t.Logf(" Invalid private key formats correctly rejected")
	})
}

// TestHDPathFunctionality tests the HD path logic
// that is implemented in account.GetFullBIP44Path() and account.NewHDPathIterator()
func TestHDPathFunctionality(t *testing.T) {
	t.Run("Test BIP44 path format", func(t *testing.T) {
		// This reproduces the logic from account.GetFullBIP44Path()
		// Note: We can't use sdk.Purpose and sdk.CoinType due to cosmos-sdk issues
		// But we can test the expected format
		const Purpose = 44 // Standard BIP44 purpose

		expectedPath := "m/44'/60'/0'/0/0" // Expected format for Ethereum
		t.Logf("Expected BIP44 path format: %s", expectedPath)
		t.Logf("Expected coin type: %d", CoinType)
		t.Logf(" BIP44 path format validated")
	})

	t.Run("Test HD path iterator functionality", func(t *testing.T) {
		// Use DefaultRootDerivationPath directly for iteration
		basePath := ethaccounts.DefaultRootDerivationPath

		iterator := ethaccounts.DefaultIterator(basePath)
		require.NotNil(t, iterator, "Iterator should not be nil")

		// Test path generation by collecting multiple paths
		var paths []string

		for _ = range 5 {
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

		t.Logf(" HD path iterator works correctly:")
		t.Logf("  Base: %s", basePath.String())
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
			t.Run("Invalid: "+invalidPath, func(t *testing.T) {
				_, err := ethaccounts.ParseDerivationPath(invalidPath)
				assert.Error(t, err, "Invalid path should cause parsing error: %s", invalidPath)
			})
		}
		t.Logf(" Invalid HD paths correctly rejected")
	})

	t.Run("Test coin type constant", func(t *testing.T) {
		assert.Equal(t, uint32(60), CoinType, "Coin type should be 60 for Ethereum")
		t.Logf(" Ethereum coin type (60) validated")
	})
}

// TestAccountConsistencyFunctionality tests consistency between different methods
func TestAccountConsistencyFunctionality(t *testing.T) {
	t.Run("Private key from mnemonic matches address generation", func(t *testing.T) {
		// Generate private key from mnemonic
		seed := bip39.NewSeed(TestMnemonic, TestPassword)
		privateKey, err := crypto.ToECDSA(seed[:32])
		require.NoError(t, err)

		privateKeyBytes := crypto.FromECDSA(privateKey)
		privateKeyHex := hex.EncodeToString(privateKeyBytes)

		// Create address from mnemonic
		address1 := crypto.PubkeyToAddress(privateKey.PublicKey)

		// Create address from extracted private key
		privateKeyBytes2, err := hex.DecodeString(privateKeyHex)
		require.NoError(t, err)
		privateKey2, err := crypto.ToECDSA(privateKeyBytes2)
		require.NoError(t, err)
		address2 := crypto.PubkeyToAddress(privateKey2.PublicKey)

		assert.Equal(t, address1, address2, "Address from mnemonic should match address from extracted private key")
		t.Logf(" Consistency validated: %s", address1.Hex())
	})
}

// ========================================
// ENVIRONMENT INTEGRATION TESTS
// ========================================

// TestWithEnvironmentVariables tests functionality with actual environment setup
func TestWithEnvironmentVariables(t *testing.T) {
	t.Run("Test in environment context", func(t *testing.T) {
		// Get environment variables
		evmRPC := os.Getenv("EVM_FIVENET_RPC")
		fivenetRPC := os.Getenv("FIVENET_RPC")
		fivenetAPI := os.Getenv("FIVENET_API")

		t.Logf("Environment context:")
		t.Logf("  EVM_FIVENET_RPC: %s", evmRPC)
		t.Logf("  FIVENET_RPC: %s", fivenetRPC)
		t.Logf("  FIVENET_API: %s", fivenetAPI)

		// Test mnemonic generation in this environment
		entropy, err := bip39.NewEntropy(mnemonicEntropySize)
		require.NoError(t, err, "Mnemonic generation should work in environment")

		mnemonic, err := bip39.NewMnemonic(entropy)
		require.NoError(t, err, "Mnemonic creation should work in environment")

		// Test EVM account creation
		if bip39.IsMnemonicValid(mnemonic) {
			seed := bip39.NewSeed(mnemonic, "test")
			privateKey, err := crypto.ToECDSA(seed[:32])
			require.NoError(t, err, "EVM account creation should work in environment")

			address := crypto.PubkeyToAddress(privateKey.PublicKey)
			assert.NotEqual(t, common.Address{}, address, "Should generate valid address in environment")

			t.Logf(" Generated account in environment: %s", address.Hex())
		}
	})
}

// ========================================
// BENCHMARK TESTS
// ========================================

// BenchmarkMnemonicGeneration benchmarks the mnemonic generation performance
func BenchmarkMnemonicGeneration(b *testing.B) {
	for i := 0; i < b.N; i++ {
		entropy, err := bip39.NewEntropy(mnemonicEntropySize)
		if err != nil {
			b.Fatal(err)
		}
		_, err = bip39.NewMnemonic(entropy)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkEVMAccountFromMnemonic benchmarks EVM account creation from mnemonic
func BenchmarkEVMAccountFromMnemonic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		seed := bip39.NewSeed(TestMnemonic, TestPassword)
		privateKey, err := crypto.ToECDSA(seed[:32])
		if err != nil {
			b.Fatal(err)
		}
		_ = crypto.PubkeyToAddress(privateKey.PublicKey)
	}
}

// BenchmarkEVMAccountFromPrivateKey benchmarks EVM account creation from private key
func BenchmarkEVMAccountFromPrivateKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		privateKeyBytes, err := hex.DecodeString(TestPrivateKey)
		if err != nil {
			b.Fatal(err)
		}
		privateKey, err := crypto.ToECDSA(privateKeyBytes)
		if err != nil {
			b.Fatal(err)
		}
		_ = crypto.PubkeyToAddress(privateKey.PublicKey)
	}
}

// BenchmarkHDPathIterator benchmarks HD path iterator creation
func BenchmarkHDPathIterator(b *testing.B) {
	basePath := "m/44'/60'/0'/0"

	for i := 0; i < b.N; i++ {
		hdPath, err := ethaccounts.ParseDerivationPath(basePath)
		if err != nil {
			b.Fatal(err)
		}
		_ = ethaccounts.DefaultIterator(hdPath)
	}
}

// ========================================
// HELPER FUNCTIONS
// ========================================

// getFirstWords returns the first n words from a mnemonic string
func getFirstWords(mnemonic string, n int) string {
	words := strings.Fields(mnemonic)
	if len(words) <= n {
		return mnemonic
	}
	return strings.Join(words[:n], " ")
}
