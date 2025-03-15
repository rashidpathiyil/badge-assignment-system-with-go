package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// LogLevel defines the severity level for logging
type LogLevel int

const (
	// LogLevelOff disables all logging
	LogLevelOff LogLevel = iota
	// LogLevelError only logs errors
	LogLevelError
	// LogLevelWarning logs warnings and errors
	LogLevelWarning
	// LogLevelInfo logs info, warnings, and errors
	LogLevelInfo
	// LogLevelDebug logs debug info and all above
	LogLevelDebug
	// LogLevelTrace logs detailed trace info and all above
	LogLevelTrace
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case LogLevelOff:
		return "OFF"
	case LogLevelError:
		return "ERROR"
	case LogLevelWarning:
		return "WARN"
	case LogLevelInfo:
		return "INFO"
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelTrace:
		return "TRACE"
	default:
		return fmt.Sprintf("LEVEL(%d)", l)
	}
}

// Logger provides functions for logging at different levels
type Logger struct {
	mu       sync.Mutex
	level    LogLevel
	prefix   string
	output   io.Writer
	logger   *log.Logger
	enabled  bool
	maskData bool
}

var (
	// Default global logger instance
	defaultLogger *Logger
	// Global logger mutex
	globalMu sync.Mutex
)

func init() {
	// Initialize the default logger
	defaultLogger = NewLogger("BADGE-SYSTEM", LogLevelInfo)
}

// NewLogger creates a new logger with the specified prefix and level
func NewLogger(prefix string, level LogLevel) *Logger {
	return &Logger{
		level:    level,
		prefix:   prefix,
		output:   os.Stdout,
		logger:   log.New(os.Stdout, "", log.LstdFlags),
		enabled:  level > LogLevelOff,
		maskData: false,
	}
}

// SetLevel changes the log level of the logger
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
	l.enabled = level > LogLevelOff
}

// SetOutput changes the output destination of the logger
func (l *Logger) SetOutput(output io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = output
	l.logger = log.New(output, "", log.LstdFlags)
}

// SetMaskData controls whether sensitive data should be masked in logs
func (l *Logger) SetMaskData(mask bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.maskData = mask
}

// formatMessage formats a log message with timestamp, prefix, and level
func (l *Logger) formatMessage(level LogLevel, format string, args ...interface{}) string {
	message := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	return fmt.Sprintf("%s [%s] [%s] %s", timestamp, l.prefix, level.String(), message)
}

// log logs a message at the specified level
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if !l.enabled || l.level < level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	message := l.formatMessage(level, format, args...)

	// If sensitive data should be masked, do it here
	if l.maskData {
		// Replace patterns like "password": "secret" with "password": "[MASKED]"
		message = maskSensitiveData(message)
	}

	l.logger.Println(message)
}

// maskSensitiveData replaces sensitive data in log messages
func maskSensitiveData(message string) string {
	sensitiveFields := []string{"password", "token", "secret", "key", "credential"}

	for _, field := range sensitiveFields {
		keyPattern := fmt.Sprintf(`"%s":`, field)
		if strings.Contains(message, keyPattern) {
			// Match anything between quotes after the field
			parts := strings.SplitN(message, keyPattern, 2)
			if len(parts) != 2 {
				continue
			}

			// Try to find the value and mask it
			if strings.Contains(parts[1], `"`) {
				valueParts := strings.SplitN(parts[1], `"`, 3)
				if len(valueParts) >= 3 {
					message = parts[0] + keyPattern + `"[MASKED]"` + valueParts[2]
				}
			}
		}
	}

	return message
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(LogLevelError, format, args...)
}

// Warning logs a warning message
func (l *Logger) Warning(format string, args ...interface{}) {
	l.log(LogLevelWarning, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(LogLevelInfo, format, args...)
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(LogLevelDebug, format, args...)
}

// Trace logs a trace message with extremely detailed information
func (l *Logger) Trace(format string, args ...interface{}) {
	l.log(LogLevelTrace, format, args...)
}

// GetLogger returns the default logger
func GetLogger() *Logger {
	globalMu.Lock()
	defer globalMu.Unlock()
	return defaultLogger
}

// SetDefaultLevel sets the log level for the default logger
func SetDefaultLevel(level LogLevel) {
	globalMu.Lock()
	defer globalMu.Unlock()
	defaultLogger.SetLevel(level)
}

// SetDefaultOutput sets the output destination for the default logger
func SetDefaultOutput(output io.Writer) {
	globalMu.Lock()
	defer globalMu.Unlock()
	defaultLogger.SetOutput(output)
}

// Error logs an error message using the default logger
func Error(format string, args ...interface{}) {
	GetLogger().Error(format, args...)
}

// Warning logs a warning message using the default logger
func Warning(format string, args ...interface{}) {
	GetLogger().Warning(format, args...)
}

// Info logs an info message using the default logger
func Info(format string, args ...interface{}) {
	GetLogger().Info(format, args...)
}

// Debug logs a debug message using the default logger
func Debug(format string, args ...interface{}) {
	GetLogger().Debug(format, args...)
}

// Trace logs a trace message using the default logger
func Trace(format string, args ...interface{}) {
	GetLogger().Trace(format, args...)
}
