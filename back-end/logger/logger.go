package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type SimpleLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
	logToFile   bool
	logFile     *os.File
}

type Config struct {
	Level    string
	FilePath string
}

// NewSimpleLogger creates a new simple logger
func NewSimpleLogger(cfg Config) (*SimpleLogger, error) {
	logger := &SimpleLogger{}

	var logFile *os.File
	var err error

	if cfg.FilePath != "" {
		dir := filepath.Dir(cfg.FilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}
		logFile, err = os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		logger.logFile = logFile
		logger.logToFile = true
	}

	if logger.logToFile {
		logger.infoLogger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.errorLogger = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.debugLogger = log.New(logFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		logger.infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.debugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	return logger, nil
}

// Info logs info level messages
func (l *SimpleLogger) Info(message string, fields ...interface{}) {
	if len(fields) > 0 {
		message = fmt.Sprintf("%s | %v", message, fields)
	}
	l.infoLogger.Println(message)
}

// Error logs error level messages
func (l *SimpleLogger) Error(message string, err error, fields ...interface{}) {
	errorMsg := message
	if err != nil {
		errorMsg = fmt.Sprintf("%s: %v", message, err)
	}
	if len(fields) > 0 {
		errorMsg = fmt.Sprintf("%s | %v", errorMsg, fields)
	}
	l.errorLogger.Println(errorMsg)
}

// Debug logs debug level messages
func (l *SimpleLogger) Debug(message string, fields ...interface{}) {
	if len(fields) > 0 {
		message = fmt.Sprintf("%s | %v", message, fields)
	}
	l.debugLogger.Println(message)
}

// LogRequest logs HTTP requests
func (l *SimpleLogger) LogRequest(method, path string, statusCode int, duration time.Duration) {
	l.Info(fmt.Sprintf("Request: %s %s | Status: %d | Duration: %v",
		method, path, statusCode, duration))
}

// Close closes the log file if it exists
func (l *SimpleLogger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}
