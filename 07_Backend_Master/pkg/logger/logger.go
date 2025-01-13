package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Logger wraps zerolog.Logger
type Logger struct {
	logger zerolog.Logger
}

// Config holds logger configuration
type Config struct {
	Level      string
	File       string
	TimeFormat string
	Console    bool
}

// New creates a new logger instance
func New(cfg Config) (*Logger, error) {
	// Set default time format
	if cfg.TimeFormat == "" {
		cfg.TimeFormat = time.RFC3339
	}

	// Parse log level
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}

	// Configure writers
	var writers []io.Writer

	// Add console writer if enabled
	if cfg.Console {
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: cfg.TimeFormat,
		}
		writers = append(writers, consoleWriter)
	}

	// Add file writer if specified
	if cfg.File != "" {
		file, err := os.OpenFile(cfg.File, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			return nil, err
		}
		writers = append(writers, file)
	}

	// Create multi-writer if multiple writers exist
	var writer io.Writer
	if len(writers) > 1 {
		writer = io.MultiWriter(writers...)
	} else if len(writers) == 1 {
		writer = writers[0]
	} else {
		writer = os.Stdout // Default to stdout if no writers specified
	}

	// Create logger
	logger := zerolog.New(writer).
		Level(level).
		With().
		Timestamp().
		Logger()

	return &Logger{logger: logger}, nil
}

// Debug logs debug level message
func (l *Logger) Debug(msg string, fields map[string]interface{}) {
	event := l.logger.Debug()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// Info logs info level message
func (l *Logger) Info(msg string, fields map[string]interface{}) {
	event := l.logger.Info()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// Warn logs warn level message
func (l *Logger) Warn(msg string, fields map[string]interface{}) {
	event := l.logger.Warn()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// Error logs error level message
func (l *Logger) Error(msg string, err error, fields map[string]interface{}) {
	event := l.logger.Error()
	if err != nil {
		event = event.Err(err)
	}
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// Fatal logs fatal level message and exits
func (l *Logger) Fatal(msg string, err error, fields map[string]interface{}) {
	event := l.logger.Fatal()
	if err != nil {
		event = event.Err(err)
	}
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// WithFields adds fields to the logger
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	contextLogger := l.logger.With()
	for k, v := range fields {
		contextLogger = contextLogger.Interface(k, v)
	}
	return &Logger{logger: contextLogger.Logger()}
}

// WithRequestID adds request ID to the logger
func (l *Logger) WithRequestID(requestID string) *Logger {
	return &Logger{logger: l.logger.With().Str("request_id", requestID).Logger()}
}

// WithUser adds user information to the logger
func (l *Logger) WithUser(userID string) *Logger {
	return &Logger{logger: l.logger.With().Str("user_id", userID).Logger()}
}
