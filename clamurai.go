package clamurai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/brunog-costa/clamurai/pkg/logger"
	"github.com/hq0101/go-clamav/pkg/clamav"
)

// Validate if config can be declared as is without breaking the middleware, else create a logic for that in the code
const (
	clamavAddress           = "localhost:3310"
	clamavReadTimeout       = 3600
	clamavConnectionTimeout = 90
)

// Accepted parameters for plugin
type Config struct {
	ClamdAddress            string
	ClamavReadTimeout       uint64
	ClamavConnectionTimeout uint64
}

// Creates config with default value when none are provided by dynamic config
func CreateConfig() *Config {
	return &Config{
		ClamdAddress:            clamavAddress,
		ClamavReadTimeout:       clamavReadTimeout,
		ClamavConnectionTimeout: clamavConnectionTimeout,
	}
}

// Receives values from the plugin instantiation.
type Clamurai struct {
	next         http.Handler
	name         string
	clamavClient *clamav.ClamClient
	logging      *logger.JSONLogger
}

// Sets up the middleware plugin for further usage
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config == nil {
		config = CreateConfig()
	}

	// Initialize logger
	logging := logger.NewJSONLogger(name)
	logging.InfoWithFields("Initializing clam client", logger.Fields("clamAddress", config.ClamdAddress))

	// Instantiates the clam client
	client := clamav.NewClamClient(
		"tcp",
		config.ClamdAddress,
		time.Duration(config.ClamavConnectionTimeout)*time.Second,
		time.Duration(config.ClamavReadTimeout)*time.Second,
	)

	// Exits traefik if can't reach clam tcp port
	serverStats, err := client.Stats()
	if err != nil {
		logging.ErrorWithFields("Could not get server stats, paniking out", logger.Fields("Error", err))
		panic(err)
	}

	logging.InfoWithFields("Server Stats", logger.Fields(
		"ThreadsIdle", serverStats.ThreadsIdle,
		"ThreadsLive", serverStats.ThreadsLive,
		"MemFree", serverStats.MemFree,
		"MemUsed", serverStats.MemUsed))

	// Return configured parameters
	return &Clamurai{
		next:         next,
		name:         name,
		clamavClient: client,
		logging:      logging,
	}, nil
}

// Intercepting the requests and processing them before sending to next middleware or upstream
func (c *Clamurai) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	fullURL := fmt.Sprintf("https://%s%s", req.Host, req.URL)

	// Log incoming AWS headers - perhaps migrate to a for loop and print each header
	c.logging.InfoWithFields("Initializing scan of incoming artifact for url", logger.Fields("url", fullURL))

	c.logging.Info("Building byte array from request body")
	body, err := io.ReadAll(req.Body)
	if err != nil {
		// This will exit if there's an error
		c.logging.ErrorWithFields("Pre-scan failed", logger.Fields("Error", err))
		http.Error(rw, "Pre-scan failed, could not process body", http.StatusInternalServerError)
	}

	scanStart := time.Now()

	// Sends ByteStrem to clam scan
	results, err := c.clamavClient.Instream(body)
	if err != nil {
		// If file can't be inspected, rejects the request with 500 (Fail Close) -> gotta improve response writter handle
		c.logging.ErrorWithFields("In stream scan failed", logger.Fields("Error", err))
		http.Error(rw, "Could not process body in time", http.StatusBadGateway)
		return
	}

	req.Body = io.NopCloser(bytes.NewReader(body))

	// // Inspects results for malware
	for _, result := range results {
		if result.Status == "FOUND" {
			c.logging.WarnWithFields("Virus found on request body, refusing request", logger.Fields(
				"Signature", result.Virus,
				"Object", req.URL.String(),
				"SHA256", req.Header.Get("X-Amz-Content-Sha256"),
				"ScanDuration", time.Since(scanStart).Seconds(),
			))

			http.Error(rw, "Forbidden", http.StatusForbidden)
			return
		}
	}
	// Sends request downstream if the file is clean
	c.logging.InfoWithFields("Clean file, sending to downstream", logger.Fields(
		"Object", req.URL.String(),
		"SHA256", req.Header.Get("X-Amz-Content-Sha256"),
		"ScanDuration", time.Since(scanStart).Seconds(),
	))

	c.next.ServeHTTP(rw, req)
}
