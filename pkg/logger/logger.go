package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// Logger struct
type JSONLogger struct {
	middleware string
	hostname   string
	logger     *log.Logger
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
func NewJSONLogger(middlewareName string) *JSONLogger {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = middlewareName
	}

	return &JSONLogger{
		middleware: middlewareName,
		hostname:   hostname,
		logger:     log.New(os.Stdout, "", 0), // no prefix, we’ll format ourselves
	}
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
		fmt.Fprintf(os.Stderr, "failed to marshal log: %v\n", err)
		return
	}

	l.logger.Println(string(jsonData))
}

// Helper methods for each log level
func (l *JSONLogger) Info(msg string) {
	l.logWithFields("INFO", msg, map[string]interface{}{})

}

func (l *JSONLogger) InfoWithFields(msg string, fields map[string]interface{}) {
	l.logWithFields("INFO", msg, fields)

}

func (l *JSONLogger) Error(msg string) {
	l.logWithFields("ERROR", msg, map[string]interface{}{})
}

func (l *JSONLogger) ErrorWithFields(msg string, fields map[string]interface{}) {
	l.logWithFields("ERROR", msg, fields)
}

func (l *JSONLogger) Debug(msg string) {
	l.logWithFields("DEBUG", msg, map[string]interface{}{})
}

func (l *JSONLogger) DebugWithFields(msg string, fields map[string]interface{}) {
	l.logWithFields("DEBUG", msg, fields)
}

func (l *JSONLogger) Warn(msg string) {
	l.logWithFields("WARNING", msg, map[string]interface{}{})
}

func (l *JSONLogger) WarnWithFields(msg string, fields map[string]interface{}) {
	l.logWithFields("WARNING", msg, fields)
}
