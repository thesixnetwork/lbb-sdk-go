# Logger Package

A structured logging package for the LBB SDK, built on top of [zerolog](https://github.com/rs/zerolog) - one of the fastest and most efficient logging libraries for Go.

> **Note:** This logger wraps zerolog and delegates level management directly to it, avoiding redundant level tracking for optimal performance.

## Features

- **High Performance**: Built on zerolog, offering zero-allocation structured logging
- **Structured Logging**: Add context fields to logs for better observability
- **Multiple Log Levels**: DEBUG, INFO, WARN, ERROR, FATAL
- **Colored Console Output**: Beautiful, human-readable logs for development
- **JSON Output Support**: Machine-readable logs for production via zerolog
- **Contextual Logging**: Create child loggers with persistent fields
- **Flexible Configuration**: Control log level, output, colors, timestamps
- **Compatibility**: Drop-in replacement for `fmt.Print*` functions

## Installation

The logger is already included in the SDK and uses zerolog as a dependency.

## Quick Start

```go
import "github.com/thesixnetwork/lbb-sdk-go/pkg/logger"

// Use package-level functions with default logger
logger.Info("Application started")
logger.Debug("Debug information: %s", debugInfo)
logger.Warn("Warning: %s", warningMsg)
logger.Error("Error occurred: %v", err)

// Fatal logs and exits the application
logger.Fatal("Critical error: %v", criticalErr)
```

## Creating Custom Loggers

```go
// Create a logger with custom configuration
log := logger.New(
    logger.WithLevel(logger.DEBUG),
    logger.WithPrefix("MyService"),
    logger.WithColors(true),
    logger.WithTime(true),
)

log.Info("Service initialized")
```

## Structured Logging

Add context to your logs with fields:

```go
// Add a single field
log := logger.GetDefault().WithField("user_id", "12345")
log.Info("User logged in")

// Add multiple fields
log := logger.GetDefault().WithFields(map[string]interface{}{
    "user_id": "12345",
    "ip": "192.168.1.1",
    "method": "POST",
})
log.Info("API request received")

// Add error context
if err != nil {
    log.WithError(err).Error("Failed to process request")
}
```

## Log Levels

```go
// Set global log level
logger.SetLevel(logger.DEBUG)  // Show all logs
logger.SetLevel(logger.INFO)   // Hide debug logs (default)
logger.SetLevel(logger.WARN)   // Show only warnings and errors
logger.SetLevel(logger.ERROR)  // Show only errors and fatal

// Log at different levels
logger.Debug("Detailed debugging info")
logger.Info("Informational message")
logger.Warn("Warning message")
logger.Error("Error message")
logger.Fatal("Fatal error - exits program")
```

## Configuration Options

### Output Destination

```go
import "os"

// Log to a file
file, _ := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
logger.SetOutput(file)

// Log to stderr
logger.SetOutput(os.Stderr)
```

### Colors

```go
// Enable/disable colored output
logger.SetColors(true)   // Colored console output (default)
logger.SetColors(false)  // Plain text output
```

### Timestamps

```go
// Enable/disable timestamps
logger.SetTime(true)   // Show timestamps (default)
logger.SetTime(false)  // Hide timestamps
```

### Component/Prefix

```go
// Set a component name for all logs
logger.SetPrefix("EVM")
logger.Info("Transaction sent")  // Will show [EVM] component tag
```

### Get Current Level

```go
// Get the current log level (returns zerolog.Level)
log := logger.GetDefault()
currentLevel := log.GetLevel()

if currentLevel <= zerolog.DebugLevel {
    // Debug is enabled
}
```

## Advanced Usage with Zerolog

Access the underlying zerolog logger for advanced features:

```go
log := logger.GetDefault()
zlog := log.GetZerolog()

// Use zerolog's advanced features
zlog.Info().
    Str("module", "evm").
    Int("gas_used", 21000).
    Dur("duration", elapsed).
    Msg("Transaction complete")
```

## JSON Output for Production

For production environments, you can configure pure JSON output:

```go
import (
    "os"
    "github.com/rs/zerolog"
    "github.com/thesixnetwork/lbb-sdk-go/pkg/logger"
)

// Configure for production: JSON output, no colors
log := logger.New(
    logger.WithLevel(logger.INFO),
    logger.WithOutput(os.Stdout),
)

// Override with pure zerolog for JSON
jsonLogger := zerolog.New(os.Stdout).With().Timestamp().Logger()
// ... use jsonLogger directly for JSON output
```

## Migration from fmt Package

The logger provides compatibility helpers:

```go
// Before:
fmt.Println("Hello world")
fmt.Printf("User: %s", username)

// After:
logger.Println("Hello world")
logger.Printf("User: %s", username)
```

## Best Practices

### 1. Use Structured Fields Instead of String Formatting

**Good:**
```go
logger.GetDefault().
    WithField("user_id", userID).
    WithField("action", "login").
    Info("User action")
```

**Bad:**
```go
logger.Info("User %s performed action %s", userID, action)
```

### 2. Create Component-Specific Loggers

```go
// In account package
var log = logger.New(logger.WithPrefix("Account"))

func CreateAccount() {
    log.Info("Creating new account")
}

// In evm package
var log = logger.New(logger.WithPrefix("EVM"))

func DeployContract() {
    log.Info("Deploying contract")
}
```

### 3. Use Contextual Loggers in Request Handlers

```go
func HandleRequest(requestID string) {
    log := logger.GetDefault().WithField("request_id", requestID)
    
    log.Info("Processing request")
    // ... do work ...
    log.Info("Request completed")
}
```

### 4. Log Errors with Context

```go
if err := processTransaction(tx); err != nil {
    logger.GetDefault().
        WithError(err).
        WithField("tx_hash", tx.Hash()).
        Error("Failed to process transaction")
}
```

### 5. Set Appropriate Log Levels

- **DEBUG**: Verbose information for debugging (disabled in production)
- **INFO**: General informational messages (API calls, successful operations)
- **WARN**: Warning messages that don't stop execution
- **ERROR**: Errors that are handled but should be investigated
- **FATAL**: Critical errors that require immediate termination

## Environment-Based Configuration

```go
import "os"

func InitLogger() {
    if os.Getenv("ENV") == "production" {
        logger.SetLevel(logger.INFO)
        logger.SetColors(false)
        // Consider using pure JSON output
    } else {
        logger.SetLevel(logger.DEBUG)
        logger.SetColors(true)
    }
}
```

## Performance Considerations

- Zerolog is designed for zero-allocation logging in hot paths
- Use structured fields instead of `fmt.Sprintf` in log messages
- Consider setting `logger.INFO` or higher in production to reduce log volume
- Zerolog can handle millions of log entries per second

## Examples

### Basic CLI Application

```go
package main

import "github.com/thesixnetwork/lbb-sdk-go/pkg/logger"

func main() {
    logger.SetPrefix("CLI")
    logger.Info("Application starting...")
    
    if err := run(); err != nil {
        logger.Fatal("Application failed: %v", err)
    }
    
    logger.Info("Application finished successfully")
}
```

### Service with Structured Logging

```go
package service

import "github.com/thesixnetwork/lbb-sdk-go/pkg/logger"

type Service struct {
    log *logger.Logger
}

func NewService(name string) *Service {
    return &Service{
        log: logger.New(
            logger.WithPrefix(name),
            logger.WithLevel(logger.INFO),
        ),
    }
}

func (s *Service) Process(userID string, data []byte) error {
    log := s.log.WithField("user_id", userID)
    log.Info("Processing data")
    
    if err := validate(data); err != nil {
        log.WithError(err).Error("Validation failed")
        return err
    }
    
    log.WithField("size", len(data)).Info("Data processed successfully")
    return nil
}
```

## Design Decisions

### Using Zerolog's Level Directly

The logger delegates level management entirely to zerolog rather than maintaining a duplicate level field. This design choice:

1. **Eliminates Redundancy**: No need to keep logger's level in sync with zerolog's level
2. **Better Performance**: One less field to maintain and copy in contextual loggers
3. **Simplicity**: Zerolog already implements level filtering perfectly
4. **Direct Access**: Use `GetLevel()` to query the current level from zerolog

### Why Zerolog?

We chose zerolog for the LBB SDK because:

1. **Performance**: Zero-allocation, fastest JSON logging in Go
2. **Simplicity**: Clean, chainable API
3. **Flexibility**: Easy to switch between console and JSON output
4. **Battle-tested**: Used by many production systems
5. **Active maintenance**: Regular updates and community support
6. **No Overhead**: Direct delegation to zerolog means minimal abstraction cost

## References

- [Zerolog Documentation](https://github.com/rs/zerolog)
- [Zerolog Performance Benchmarks](https://github.com/rs/zerolog#benchmarks)
- [Structured Logging Best Practices](https://www.honeycomb.io/blog/structured-logging-and-your-team)

## License

Same as the LBB SDK license.