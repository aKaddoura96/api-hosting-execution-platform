package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

// LogLevel represents the severity of a log message
type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
	FATAL LogLevel = "FATAL"
)

// Logger handles structured logging
type Logger struct {
	serviceName string
	level       LogLevel
}

// NewLogger creates a new logger instance
func NewLogger(serviceName string) *Logger {
	return &Logger{
		serviceName: serviceName,
		level:       INFO, // Default to INFO
	}
}

// SetLevel sets the minimum log level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// log formats and writes a log message
func (l *Logger) log(level LogLevel, message string, fields map[string]interface{}) {
	// Get caller information
	_, file, line, ok := runtime.Caller(2)
	caller := "unknown"
	if ok {
		caller = fmt.Sprintf("%s:%d", file, line)
	}

	// Format timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	// Build log entry
	entry := fmt.Sprintf("[%s] [%s] [%s] %s", 
		timestamp,
		level,
		l.serviceName,
		message,
	)

	// Add fields if provided
	if len(fields) > 0 {
		entry += " |"
		for k, v := range fields {
			entry += fmt.Sprintf(" %s=%v", k, v)
		}
	}

	// Add caller for WARN, ERROR, FATAL
	if level == WARN || level == ERROR || level == FATAL {
		entry += fmt.Sprintf(" | caller=%s", caller)
	}

	// Write to stdout
	fmt.Println(entry)

	// Exit on FATAL
	if level == FATAL {
		os.Exit(1)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(message string, fields ...map[string]interface{}) {
	if l.level == DEBUG {
		var f map[string]interface{}
		if len(fields) > 0 {
			f = fields[0]
		}
		l.log(DEBUG, message, f)
	}
}

// Info logs an info message
func (l *Logger) Info(message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(INFO, message, f)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(WARN, message, f)
}

// Error logs an error message
func (l *Logger) Error(message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(ERROR, message, f)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(message string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(FATAL, message, f)
}

// LogHTTPRequest logs an HTTP request
func (l *Logger) LogHTTPRequest(method, path, remoteAddr string, statusCode int, duration time.Duration) {
	l.Info("HTTP Request", map[string]interface{}{
		"method":     method,
		"path":       path,
		"remote":     remoteAddr,
		"status":     statusCode,
		"duration_ms": duration.Milliseconds(),
	})
}

// LogError logs an error with context
func (l *Logger) LogError(err error, context string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	} else {
		f = make(map[string]interface{})
	}
	f["error"] = err.Error()
	f["context"] = context
	l.Error("Error occurred", f)
}

// Global logger instance (can be replaced with service-specific loggers)
var defaultLogger *Logger

func init() {
	defaultLogger = NewLogger("default")
}

// Package-level functions that use the default logger
func Debug(message string, fields ...map[string]interface{}) {
	defaultLogger.Debug(message, fields...)
}

func Info(message string, fields ...map[string]interface{}) {
	defaultLogger.Info(message, fields...)
}

func Warn(message string, fields ...map[string]interface{}) {
	defaultLogger.Warn(message, fields...)
}

func Error(message string, fields ...map[string]interface{}) {
	defaultLogger.Error(message, fields...)
}

func Fatal(message string, fields ...map[string]interface{}) {
	defaultLogger.Fatal(message, fields...)
}

// SetDefaultLogger sets the default logger
func SetDefaultLogger(logger *Logger) {
	defaultLogger = logger
	log.SetOutput(os.Stdout)
	log.SetFlags(0) // Disable default log prefix as we handle it
}
