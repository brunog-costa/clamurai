package logger

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"time"
)

// Abstracts functions for enabling mock
type Logger interface {
	Info(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Debug(msg string, fields map[string]interface{})
}

// New JSON Logger struct
type JSONLogger struct {
	middleware string
	hostname   string
	logger     *log.Logger
}

// Configuration exposed for new JSON Logger
type Config struct {
	Middleware string
	Output     io.Writer
	Hostname   string
}

// Logger staple fields
type LogEntry struct {
	Timestamp  string                 `json:"timestamp"`
	Level      string                 `json:"level"`
	Message    string                 `json:"message"`
	Middleware string                 `json:"middleware"`
	Hostname   string                 `json:"hostname"`
	Fields     map[string]interface{} `json:"fields,omitempty"`
}

// Initialize the logger
func New(config Config) (*JSONLogger, error) {
	if config.Output == nil {
		config.Output = os.Stdout
	}

	hostname := config.Hostname
	if hostname == "" {
		var err error
		hostname, err = os.Hostname()
		if err != nil {
			return nil, err
		}
	}

	return &JSONLogger{
		middleware: config.Middleware,
		hostname:   hostname,
		logger:     log.New(config.Output, "", 0),
	}, nil
}

func Fields(kv ...interface{}) map[string]interface{} {
	fields := make(map[string]interface{})
	for i := 0; i < len(kv)-1; i += 2 {
		key, ok := kv[i].(string)
		if !ok {
			continue // skip non-string keys
		}
		fields[key] = kv[i+1]
	}
	return fields
}

// Receives message and extra fields forwards it to stdout
func (l *JSONLogger) logWithFields(level, msg string, fields map[string]interface{}) {
	entry := LogEntry{
		Timestamp:  time.Now().Format(time.RFC3339),
		Level:      level,
		Middleware: l.middleware,
		Hostname:   l.hostname,
		Message:    msg,
		Fields:     fields,
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		// Fallback to simple text logging if JSON marshaling fails
		l.logger.Printf("[%s] %s: %s (JSON marshal error: %v)", level, l.middleware, msg, err)
		return
	}

	l.logger.Println(string(jsonData))
}

// Helper methods for each log level
func (l *JSONLogger) Info(msg string, fields map[string]interface{}) {
	l.logWithFields("INFO", msg, fields)
}

func (l *JSONLogger) Error(msg string, fields map[string]interface{}) {
	l.logWithFields("ERROR", msg, fields)
}

func (l *JSONLogger) Debug(msg string, fields map[string]interface{}) {
	l.logWithFields("DEBUG", msg, fields)
}

func (l *JSONLogger) Warn(msg string, fields map[string]interface{}) {
	l.logWithFields("WARNING", msg, fields)
}
