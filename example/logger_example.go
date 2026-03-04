package main

import (
	"errors"
	"os"
	"time"

	"github.com/thesixnetwork/lbb-sdk-go/pkg/logger"
)

func main() {
	println("=== LBB SDK Logger Examples (Powered by Zerolog) ===\n")

	// Example 1: Basic Logging
	basicLogging()

	// Example 2: Log Levels
	logLevels()

	// Example 3: Structured Logging
	structuredLogging()

	// Example 4: Contextual Logging
	contextualLogging()

	// Example 5: Error Logging
	errorLogging()

	// Example 6: Custom Logger Instances
	customLoggers()

	// Example 7: Configuration Options
	configurationOptions()

	// Example 8: Component-Specific Logging
	componentLogging()

	// Example 9: Performance Example
	performanceExample()

	println("\n=== All Examples Complete ===")
}

func basicLogging() {
	println("\n--- Example 1: Basic Logging ---")

	logger.Info("This is an info message")
	logger.Debug("This debug message won't show (default level is INFO)")
	logger.Warn("This is a warning")
	logger.Error("This is an error (non-fatal)")

	// Compatibility with fmt package
	logger.Printf("Formatted message: %s = %d", "answer", 42)
	logger.Println("Simple println message")
}

func logLevels() {
	println("\n--- Example 2: Log Levels ---")

	// Save original level
	original := logger.New(logger.WithLevel(logger.INFO))
	defer logger.SetLevel(logger.INFO)

	logger.Info("Current level: INFO (debug won't show)")
	logger.Debug("You won't see this debug message")

	// Set to DEBUG level
	logger.SetLevel(logger.DEBUG)
	logger.Info("Changed level to DEBUG")
	logger.Debug("Now you can see debug messages!")

	// Set to WARN level
	logger.SetLevel(logger.WARN)
	logger.Info("This info message is hidden (level is WARN)")
	logger.Warn("But warnings still show")
	logger.Error("And errors too")

	// Restore
	logger.SetLevel(logger.INFO)
	_ = original
}

func structuredLogging() {
	println("\n--- Example 3: Structured Logging ---")

	// Add single field
	log := logger.GetDefault().WithField("user_id", "12345")
	log.Info("User logged in")

	// Add multiple fields
	log = logger.GetDefault().WithFields(map[string]interface{}{
		"user_id":   "12345",
		"username":  "alice",
		"ip":        "192.168.1.100",
		"method":    "POST",
		"endpoint":  "/api/v1/login",
		"timestamp": time.Now().Unix(),
	})
	log.Info("API request received")

	// Chain fields
	logger.GetDefault().
		WithField("module", "evm").
		WithField("action", "deploy").
		WithField("gas_used", 21000).
		Info("Contract deployed successfully")
}

func contextualLogging() {
	println("\n--- Example 4: Contextual Logging ---")

	// Create a logger with context for a request
	requestID := "req-abc-123"
	requestLog := logger.GetDefault().WithField("request_id", requestID)

	requestLog.Info("Processing request")
	requestLog.WithField("step", "validation").Info("Validating input")
	requestLog.WithField("step", "processing").Info("Processing data")
	requestLog.WithField("step", "complete").Info("Request completed")

	// Simulate transaction logging
	txHash := "0x1234567890abcdef"
	txLog := logger.GetDefault().
		WithField("tx_hash", txHash).
		WithField("chain", "six-protocol")

	txLog.Info("Broadcasting transaction")
	txLog.WithField("gas_price", "20gwei").Info("Transaction sent")
	txLog.WithField("block", 12345).Info("Transaction confirmed")
}

func errorLogging() {
	println("\n--- Example 5: Error Logging ---")

	// Log error with context
	err := errors.New("connection timeout")

	logger.GetDefault().
		WithError(err).
		WithField("host", "rpc.example.com").
		WithField("port", 8545).
		Error("Failed to connect to RPC")

	// Simulate contract deployment error
	deployErr := errors.New("insufficient funds for gas")
	logger.GetDefault().
		WithError(deployErr).
		WithFields(map[string]interface{}{
			"contract":     "NFT",
			"required_gas": 500000,
			"balance":      100000,
		}).
		Error("Contract deployment failed")

	// Warning without error
	logger.GetDefault().
		WithField("retry_count", 3).
		WithField("max_retries", 5).
		Warn("Retrying failed operation")
}

func customLoggers() {
	println("\n--- Example 6: Custom Logger Instances ---")

	// Create a logger with custom configuration
	debugLog := logger.New(
		logger.WithLevel(logger.DEBUG),
		logger.WithPrefix("DEBUG"),
		logger.WithColors(true),
		logger.WithTime(true),
	)

	debugLog.Debug("This is a debug logger")
	debugLog.Info("It shows all log levels")

	// Create a logger for a specific component
	evmLog := logger.New(
		logger.WithPrefix("EVM"),
		logger.WithLevel(logger.INFO),
	)

	evmLog.Info("EVM module initialized")
	evmLog.WithField("address", "0xABC...").Info("Account created")

	// Create a logger without colors (for file output)
	fileLog := logger.New(
		logger.WithColors(false),
		logger.WithPrefix("FILE"),
	)

	fileLog.Info("This log has no colors (suitable for files)")
}

func configurationOptions() {
	println("\n--- Example 7: Configuration Options ---")

	// Disable colors
	logger.SetColors(false)
	logger.Info("Log without colors")

	// Re-enable colors
	logger.SetColors(true)
	logger.Info("Log with colors restored")

	// Set prefix
	logger.SetPrefix("APP")
	logger.Info("Log with prefix")

	// Clear prefix
	logger.SetPrefix("")
	logger.Info("Log without prefix")

	// Disable timestamps (note: this might not work as expected with current implementation)
	logger.Info("Timestamp control via zerolog configuration")
}

func componentLogging() {
	println("\n--- Example 8: Component-Specific Logging ---")

	// Simulate different SDK components
	accountLog := logger.New(logger.WithPrefix("Account"))
	clientLog := logger.New(logger.WithPrefix("Client"))
	evmLog := logger.New(logger.WithPrefix("EVM"))
	nftLog := logger.New(logger.WithPrefix("NFT"))

	accountLog.Info("Creating new account from mnemonic")
	accountLog.WithField("address", "six1abc...").Info("Account created")

	clientLog.Info("Connecting to RPC endpoint")
	clientLog.WithField("endpoint", "http://localhost:26657").Info("Connected")

	evmLog.Info("Deploying smart contract")
	evmLog.WithFields(map[string]interface{}{
		"contract": "MyToken",
		"gas":      3000000,
	}).Info("Contract deployed")

	nftLog.Info("Minting NFT")
	nftLog.WithFields(map[string]interface{}{
		"token_id":  "1",
		"schema":    "my_schema",
		"recipient": "six1xyz...",
	}).Info("NFT minted successfully")
}

func performanceExample() {
	println("\n--- Example 9: Performance Example ---")

	start := time.Now()
	iterations := 10000

	// Zerolog is designed for high-performance logging
	log := logger.New(
		logger.WithPrefix("PERF"),
		logger.WithLevel(logger.INFO),
	)

	for i := 0; i < iterations; i++ {
		log.WithFields(map[string]interface{}{
			"iteration": i,
			"batch":     i / 1000,
			"timestamp": time.Now().Unix(),
		}).Debug("Performance test log") // Debug logs filtered out by level
	}

	elapsed := time.Since(start)
	logger.Info("Logged %d messages in %v (avg: %v per message)",
		iterations,
		elapsed,
		elapsed/time.Duration(iterations),
	)
}

// Example: Simulating a service with structured logging
type UserService struct {
	log *logger.Logger
}

func NewUserService() *UserService {
	return &UserService{
		log: logger.New(
			logger.WithPrefix("UserService"),
			logger.WithLevel(logger.INFO),
		),
	}
}

func (s *UserService) CreateUser(userID string, email string) error {
	log := s.log.WithFields(map[string]interface{}{
		"user_id": userID,
		"email":   email,
	})

	log.Info("Creating user")

	// Simulate validation
	if email == "" {
		err := errors.New("email is required")
		log.WithError(err).Error("User creation failed")
		return err
	}

	// Simulate database operation
	time.Sleep(10 * time.Millisecond)

	log.Info("User created successfully")
	return nil
}

// Example: Using logger with file output
func fileLoggingExample() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Error("Failed to open log file: %v", err)
		return
	}
	defer file.Close()

	fileLogger := logger.New(
		logger.WithOutput(file),
		logger.WithColors(false), // No colors for file output
		logger.WithPrefix("APP"),
	)

	fileLogger.Info("This log goes to a file")
	fileLogger.WithField("event", "startup").Info("Application started")
}

// Example: Environment-based configuration
func initLoggerForEnvironment() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	switch env {
	case "production":
		logger.SetLevel(logger.INFO)
		logger.SetColors(false)
		logger.SetPrefix("PROD")
		logger.Info("Logger configured for production")

	case "development":
		logger.SetLevel(logger.DEBUG)
		logger.SetColors(true)
		logger.SetPrefix("DEV")
		logger.Info("Logger configured for development")

	case "test":
		logger.SetLevel(logger.WARN)
		logger.SetColors(false)
		logger.SetPrefix("TEST")
		logger.Info("Logger configured for testing")

	default:
		logger.SetLevel(logger.INFO)
		logger.Info("Logger configured with defaults")
	}
}

// Example: Advanced zerolog usage
func advancedZerologExample() {
	log := logger.New(logger.WithPrefix("Advanced"))

	// Get the underlying zerolog.Logger for advanced features
	zlog := log.GetZerolog()

	// Use zerolog's advanced features
	zlog.Info().
		Str("module", "evm").
		Int("gas_used", 21000).
		Dur("duration", 100*time.Millisecond).
		Bool("success", true).
		Hex("tx_hash", []byte{0x12, 0x34, 0x56, 0x78}).
		Msg("Transaction completed")
}
