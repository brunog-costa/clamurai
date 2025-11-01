package inspector

import (
	"bytes"
	"testing"
)

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

func TestInspectBody(t *testing.T) {
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
