package clamurai

import (
	"context"
	"errors"
	"io"
	"net/http"
	"regexp"

	"github.com/brunog-costa/clamurai/internal/inspector"
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

// Validates midleware configuration before starts
func validateConfig(config *Config) error {
	// Check if clamAddress is valid dns or ipv4 - tks leetcode for this one
	octectsPattern := `([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])`
	domainPattern := `^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`
	v4 := regexp.MustCompile("^" + octectsPattern + `(\.` + octectsPattern + `){3}$`)
	dns := regexp.MustCompile(domainPattern)
	// Gotta find a way of negatting those
	if v4.Match([]byte(config.ClamdAddress)) || dns.Match([]byte(config.ClamdAddress)) {
		return errors.New("invalid clam address")
	}

	// Check if connection timeout and read timeout are valid values
	return nil
}

// Sets up the middleware plugin for further usage
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	// Initialize logger
	logging := logger.NewJSONLogger(name)
	logging.InfoWithFields("Initializing clam client", logger.Fields("clamAddress", config.ClamdAddress))

	// Checks input at midleware initialization
	err := validateConfig(config)
	if err != nil {
		panic("Please validate dynamic configurations file and try again")
	}

	// Initialize inspector
	inspector := inspector.NewInspector(config.ClamdAddress, config.ClamavConnectionTimeout, config.ClamavReadTimeout)

	// Return configured parameters
	return &Clamurai{
		next:      next,
		name:      name,
		inspector: inspector, // Change this for an inspector
		logging:   logging,
		alertMode: config.AlertMode,
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

	// Split processing into two parallel  stages:
	// 1 - thread inspects body
	// 2 - if no shasum is found on headers, perform hashsum
	clean, err := c.inspector.InspectBody(body)
	if err != nil {
		c.logging.ErrorWithFields("Inspector failed to verify body with", logger.Fields("Error", err))
	}

	// If its all cool, log the result and fwd the request
	if clean || c.alertMode {
		c.next.ServeHTTP(rw, req)
	}

	// By default block requests
	return
}
