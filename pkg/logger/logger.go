package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

// LogLevel represents the severity level of a log message
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// Logger wraps zerolog for SDK-wide logging
type Logger struct {
	zlog      zerolog.Logger
	prefix    string
	useColors bool
	output    io.Writer
}

// Option is a functional option for configuring the logger
type Option func(*Logger)

// WithLevel sets the minimum log level
func WithLevel(level LogLevel) Option {
	return func(l *Logger) {
		l.zlog = l.zlog.Level(toZerologLevel(level))
	}
}

// WithOutput sets the output writer
func WithOutput(w io.Writer) Option {
	return func(l *Logger) {
		l.zlog = l.zlog.Output(w)
	}
}

// WithColors enables or disables colored output
func WithColors(enabled bool) Option {
	return func(l *Logger) {
		l.useColors = enabled
		if enabled {
			consoleWriter := zerolog.ConsoleWriter{
				Out:        l.output,
				TimeFormat: "2006-01-02 15:04:05",
				NoColor:    false,
			}
			l.zlog = zerolog.New(consoleWriter).With().Timestamp().Logger()
		} else {
			l.zlog = zerolog.New(l.output).With().Timestamp().Logger()
		}
	}
}

// WithTime enables or disables timestamp in logs
func WithTime(enabled bool) Option {
	return func(l *Logger) {
		if enabled {
			l.zlog = l.zlog.With().Timestamp().Logger()
		}
	}
}

// WithPrefix sets a prefix/component name for all log messages
func WithPrefix(prefix string) Option {
	return func(l *Logger) {
		l.prefix = prefix
		if prefix != "" {
			l.zlog = l.zlog.With().Str("component", prefix).Logger()
		}
	}
}

// New creates a new Logger instance
func New(opts ...Option) *Logger {
	// Default: console output with colors
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
		NoColor:    false,
	}

	l := &Logger{
		zlog:      zerolog.New(consoleWriter).With().Timestamp().Logger().Level(zerolog.InfoLevel),
		useColors: true,
		output:    os.Stdout,
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

// Default logger instance
var defaultLogger = New()

// SetDefault sets the default logger instance
func SetDefault(l *Logger) {
	defaultLogger = l
}

// GetDefault returns the default logger instance
func GetDefault() *Logger {
	return defaultLogger
}

// toZerologLevel converts our LogLevel to zerolog.Level
func toZerologLevel(level LogLevel) zerolog.Level {
	switch level {
	case DEBUG:
		return zerolog.DebugLevel
	case INFO:
		return zerolog.InfoLevel
	case WARN:
		return zerolog.WarnLevel
	case ERROR:
		return zerolog.ErrorLevel
	case FATAL:
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.zlog.Debug().Msgf(format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.zlog.Info().Msgf(format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.zlog.Warn().Msgf(format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.zlog.Error().Msgf(format, args...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.zlog.Fatal().Msgf(format, args...)
}

// Debug logs a debug message using the default logger
func Debug(format string, args ...interface{}) {
	defaultLogger.Debug(format, args...)
}

// Info logs an info message using the default logger
func Info(format string, args ...interface{}) {
	defaultLogger.Info(format, args...)
}

// Warn logs a warning message using the default logger
func Warn(format string, args ...interface{}) {
	defaultLogger.Warn(format, args...)
}

// Error logs an error message using the default logger
func Error(format string, args ...interface{}) {
	defaultLogger.Error(format, args...)
}

// Fatal logs a fatal message and exits using the default logger
func Fatal(format string, args ...interface{}) {
	defaultLogger.Fatal(format, args...)
}

// WithField returns a new logger with the given field added
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		zlog:      l.zlog.With().Interface(key, value).Logger(),
		prefix:    l.prefix,
		useColors: l.useColors,
		output:    l.output,
	}
}

// WithFields returns a new logger with multiple fields added
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	ctx := l.zlog.With()
	for key, value := range fields {
		ctx = ctx.Interface(key, value)
	}
	return &Logger{
		zlog:      ctx.Logger(),
		prefix:    l.prefix,
		useColors: l.useColors,
		output:    l.output,
	}
}

// WithError adds an error field to the logger
func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		zlog:      l.zlog.With().Err(err).Logger(),
		prefix:    l.prefix,
		useColors: l.useColors,
		output:    l.output,
	}
}

// SetLevel sets the minimum log level for the default logger
func SetLevel(level LogLevel) {
	defaultLogger.zlog = defaultLogger.zlog.Level(toZerologLevel(level))
}

// GetLevel returns the current log level
func (l *Logger) GetLevel() zerolog.Level {
	return l.zlog.GetLevel()
}

// SetOutput sets the output writer for the default logger
func SetOutput(w io.Writer) {
	// Recreate logger with new output
	consoleWriter := zerolog.ConsoleWriter{
		Out:        w,
		TimeFormat: "2006-01-02 15:04:05",
		NoColor:    false,
	}
	defaultLogger.zlog = zerolog.New(consoleWriter).With().Timestamp().Logger()
}

// SetColors enables or disables colored output for the default logger
func SetColors(enabled bool) {
	writer := os.Stdout
	if enabled {
		consoleWriter := zerolog.ConsoleWriter{
			Out:        writer,
			TimeFormat: "2006-01-02 15:04:05",
			NoColor:    false,
		}
		defaultLogger.zlog = zerolog.New(consoleWriter).With().Timestamp().Logger()
	} else {
		// Plain output without console writer
		defaultLogger.zlog = zerolog.New(writer).With().Timestamp().Logger()
	}
}

// SetTime enables or disables timestamp for the default logger
func SetTime(enabled bool) {
	// Recreate logger with or without timestamp
	writer := os.Stdout
	consoleWriter := zerolog.ConsoleWriter{
		Out:        writer,
		TimeFormat: "2006-01-02 15:04:05",
		NoColor:    false,
	}
	if enabled {
		defaultLogger.zlog = zerolog.New(consoleWriter).With().Timestamp().Logger()
	} else {
		defaultLogger.zlog = zerolog.New(consoleWriter)
	}
}

// SetPrefix sets the prefix/component for the default logger
func SetPrefix(prefix string) {
	defaultLogger.prefix = prefix
	if prefix != "" {
		defaultLogger.zlog = defaultLogger.zlog.With().Str("component", prefix).Logger()
	}
}

// GetZerolog returns the underlying zerolog.Logger for advanced usage
func (l *Logger) GetZerolog() *zerolog.Logger {
	return &l.zlog
}
