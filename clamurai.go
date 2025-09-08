package clamurai

import (
	"context"
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
	DevMode                 bool
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

// Validates midleware configuration before starts

func validateConfig(config *Config) error {
	return nil
}

// Sets up the middleware plugin for further usage
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	// Checks input at midleware initialization
	err := validateConfig(config)
	if err != nil && config.DevMode {
		config = CreateConfig()
	} else {
		panic("Please validate ")
	}

	// Initialize logger
	logging := logger.NewJSONLogger(name)
	logging.InfoWithFields("Initializing clam client", logger.Fields("clamAddress", config.ClamdAddress))

	// Return configured parameters
	return &Clamurai{
		next:         next,
		name:         name,
		clamavClient: client, // Change this for an inspector
		logging:      logging,
	}, nil
}

// Intercepting the requests and processing them before sending to next middleware or upstream
func (c *Clamurai) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	c.logging.Info("Building byte array from request body")

	body, err := io.ReadAll(req.Body)
	if err != nil {
		// This will exit if there's an error
		c.logging.ErrorWithFields("Pre-scan failed", logger.Fields("Error", err))
		http.Error(rw, "Pre-scan failed, could not process body", http.StatusInternalServerError)
	}

	scanStart := time.Now()
	scanFinish := time.Since(scanStart).Seconds()

	// Split processing in two stages:
	// inspection handles clam and hashing if needed
	// url re-write receives http request and makes adjustments for the destination

	// If its all cool, log the result and fwd the request
	c.next.ServeHTTP(rw, req)
}
