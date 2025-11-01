package logger

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

const (
	hostname   string = "test-hostname"
	middleware string = "test-middleware"
)

// Test helper
type testOutput struct {
	buffer *bytes.Buffer
}

// Checks the lenght of the written buffer?
func (t *testOutput) Write(p []byte) (n int, err error) {
	return t.buffer.Write(p)
}

func (t *testOutput) String() string {
	return t.buffer.String()
}

func (t *testOutput) Lines() []string {
	return strings.Split(strings.TrimSpace(t.buffer.String()), "\n")
}

func newTestOutput() *testOutput {
	return &testOutput{buffer: &bytes.Buffer{}}
}

func TestNewJSONLogger(t *testing.T) {
	// Define test variables here
	tests := []struct {
		name       string
		config     Config
		wantFields map[string]interface{}
	}{
		// Create test cases here
		{
			name: "basic logger creation",
			config: Config{
				Middleware: "clamurai-test",
				Hostname:   "clamurai-vm",
			},
			wantFields: map[string]interface{}{
				"middleware": "clamurai-test",
				"hostname":   "clamurai-vm",
			},
		},
		{
			name: "custom output logger",
			config: Config{
				Middleware: "output-test",
				Output:     &bytes.Buffer{},
			},
			wantFields: map[string]interface{}{
				"middleware": "output-test",
			},
		},
	}

	// Runs each test case against defined parameters
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewJSONLogger(tt.config)

			if logger.middleware != tt.wantFields["middleware"] {
				t.Errorf("Expected value %s for middleware, got %s", tt.wantFields, logger.middleware)
			}
		})
	}

}

func TestFieldsHelper(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		expected map[string]interface{}
	}{
		{
			name:     "valid key-value pairs",
			input:    []interface{}{"key1", "value1", "key2", 42},
			expected: map[string]interface{}{"key1": "value1", "key2": 42},
		},
		{
			name:     "non-string key",
			input:    []interface{}{42, "value", "key", "valid"},
			expected: map[string]interface{}{"key": "valid"},
		},
		{
			name:     "odd number of arguments",
			input:    []interface{}{"key1", "value1", "key2"},
			expected: map[string]interface{}{"key1": "value1"},
		},
		{
			name:     "empty input",
			input:    []interface{}{},
			expected: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Fields(tt.input...)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d fields, got %d", len(tt.expected), len(result))
			}

			for key, expectedValue := range tt.expected {
				if result[key] != expectedValue {
					t.Errorf("For key %s, expected %v, got %v", key, expectedValue, result[key])
				}
			}
		})
	}
}

func TestLogLevels(t *testing.T) {
	timeFunction := time.Now().Format(time.RFC3339)
	output := newTestOutput()

	logger := NewJSONLogger(Config{
		Middleware: middleware,
		Hostname:   hostname,
		Output:     output,
	})

	tests := []struct {
		name     string
		logFunc  func()
		expected LogEntry
	}{
		{
			name:    "info level",
			logFunc: func() { logger.Info("info message", nil) },
			expected: LogEntry{
				Timestamp:  timeFunction,
				Level:      "INFO",
				Message:    "info message",
				Middleware: middleware,
				Hostname:   hostname,
			},
		},
		{
			name:    "error level",
			logFunc: func() { logger.Error("error message", nil) },
			expected: LogEntry{
				Timestamp:  timeFunction,
				Level:      "ERROR",
				Message:    "error message",
				Middleware: middleware,
				Hostname:   hostname,
			},
		},
		{
			name:    "debug level",
			logFunc: func() { logger.Debug("debug message", nil) },
			expected: LogEntry{
				Timestamp:  timeFunction,
				Level:      "DEBUG",
				Message:    "debug message",
				Middleware: middleware,
				Hostname:   hostname,
			},
		},
		{
			name:    "warn level",
			logFunc: func() { logger.Warn("warning message", nil) },
			expected: LogEntry{
				Timestamp:  timeFunction,
				Level:      "WARNING",
				Message:    "warning message",
				Middleware: middleware,
				Hostname:   hostname,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output.buffer.Reset()
			tt.logFunc()

			lines := output.Lines()
			if len(lines) != 1 {
				t.Fatalf("Expected 1 log line, got %d", len(lines))
			}

			var entry LogEntry
			if err := json.Unmarshal([]byte(lines[0]), &entry); err != nil {
				t.Fatalf("Failed to unmarshal log entry: %v", err)
			}

			if entry.Level != tt.expected.Level {
				t.Errorf("Expected level %s, got %s", tt.expected.Level, entry.Level)
			}
			if entry.Message != tt.expected.Message {
				t.Errorf("Expected message %s, got %s", tt.expected.Message, entry.Message)
			}
			if entry.Timestamp != tt.expected.Timestamp {
				t.Errorf("Expected timestamp %s, got %s", tt.expected.Timestamp, entry.Timestamp)
			}
		})
	}
}
