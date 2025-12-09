package clamurai

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/brunog-costa/clamurai/internal/inspector"
	"github.com/brunog-costa/clamurai/pkg/hash"
	"github.com/brunog-costa/clamurai/pkg/logger"
)

// Validate if config can be declared as is without breaking the middleware, else create a logic for that in the code
const (
	clamavAddress           = "localhost:3310"
	clamavReadTimeout       = 3600
	clamavConnectionTimeout = 90
	alertMode               = true
)

// Accepted parameters for plugin
type Config struct {
	ClamdAddress            string
	ClamavReadTimeout       uint64
	ClamavConnectionTimeout uint64
	AlertMode               bool
}

// Creates config with default value when none are provided by dynamic config
func CreateConfig() *Config {
	return &Config{
		ClamdAddress:            clamavAddress,
		ClamavReadTimeout:       clamavReadTimeout,
		ClamavConnectionTimeout: clamavConnectionTimeout,
		AlertMode:               alertMode,
	}
}

// Receives values from the plugin instantiation.
type Clamurai struct {
	next      http.Handler
	name      string
	inspector *inspector.Inspector
	alertMode bool
	logging   *logger.JSONLogger
}

// Sets up the middleware plugin for further usage
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	// This call is potentially fucking up the script
	hostname, _ := os.Hostname()

	// Initialize logger
	logging, err := logger.New(logger.Config{
		Middleware: "clamurai",
		Output:     os.Stdout,
		Hostname:   hostname,
	})
	if err != nil {
		return nil, err
	}

	logging.Info("Initializing clam client", logger.Fields("clamAddress", config.ClamdAddress))

	// Initialize inspector
	inspector, err := inspector.New(config.ClamdAddress, config.ClamavConnectionTimeout, config.ClamavReadTimeout)
	if err != nil {
		logging.Error("Failed to initialize inspector", logger.Fields("Error", err))
		return nil, err
	}

	// Return configured parameters
	return &Clamurai{
		next:      next,
		name:      name,
		inspector: inspector, // Change this for an inspector
		logging:   logging,
		alertMode: config.AlertMode,
	}, nil
}

// ProcessingResult holds results from parallel processing
type ProcessingResult struct {
	Clean     bool
	Signature string
	Hash      string
	Error     error
}

// processSequentially handles AV scanning and hashing sequentially
func (c *Clamurai) processSequentially(body []byte) ProcessingResult {
	// AV scanning
	c.logging.Info("Starting inspection", nil)
	clean, signature, err := c.inspector.InspectBody(body)
	if err != nil {
		c.logging.Error("Failed inspection", nil)
		return ProcessingResult{Error: err}
	}
	c.logging.Info("Finished inspection without errors", nil)

	// Hashing
	c.logging.Info("Starting hashsum", nil)
	hashSum := hash.HashSum("sha256", body)
	c.logging.Info("Finished hashsum", nil)

	return ProcessingResult{
		Clean:     clean,
		Signature: signature,
		Hash:      hashSum,
		Error:     nil,
	}
}

// Intercepting the requests and processing them before sending to next middleware or upstream
func (c *Clamurai) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	c.logging.Info("Building byte array from request body", nil)

	body, err := io.ReadAll(req.Body)
	if err != nil {
		c.logging.Error("Pre-scan failed", logger.Fields("Error", err))
		http.Error(rw, "Pre-scan failed, could not process body", http.StatusInternalServerError)
		return
	}
	req.Body.Close()

	// Restore request body for upstream
	req.Body = io.NopCloser(bytes.NewReader(body))
	req.ContentLength = int64(len(body))

	scanStart := time.Now()
	result := c.processSequentially(body)
	scanDuration := time.Since(scanStart)

	if result.Error != nil {
		c.logging.Error("Processing failed", logger.Fields("Error", result.Error, "Duration", scanDuration.Seconds()))
		http.Error(rw, "Scan failed", http.StatusInternalServerError)
		return
	}

	c.logging.Info("Processing completed", logger.Fields(
		"Clean", result.Clean,
		"Hash", result.Hash,
		"Duration", scanDuration.Seconds(),
	))

	if !result.Clean {
		c.logging.Warn("Malware detected", logger.Fields("Signature", result.Signature))
		if !c.alertMode {
			http.Error(rw, "Malware detected", http.StatusForbidden)
			return
		}
	}

	c.next.ServeHTTP(rw, req)
}
