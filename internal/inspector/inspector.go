package inspector

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/brunog-costa/clamurai/pkg/logger"
	"github.com/hq0101/go-clamav/pkg/clamav"
	"github.com/hq0101/go-clamav/pkg/cli"
)

type Inspector struct {
	client *clamav.ClamClient
	logger *logger.JSONLogger
}

type InspectorWorker interface {
	clamHealthCheck(client *clamav.ClamClient) error
	InspectBody(body []byte) (bool, error)
}

// Builds new inspector from config - might return an nil pointer  need to check this out
func New(clamAddress string, readTimeOut uint64, connectionTimeout uint64) *Inspector {
	// Provisions custom logger for inspector
	hostname, _ := os.Hostname()

	// Initialize logger
	log := logger.NewJSONLogger(logger.Config{
		Middleware: "inspector",
		Output:     io.MultiWriter(),
		Hostname:   hostname,
	})

	// Creates a client and perform a healthcheck in order to validate if server has capacity
	client := clamav.NewClamClient("tcp", clamAddress, time.Duration(connectionTimeout)*time.Second, time.Duration(readTimeOut)*time.Second)

	err := clamHealthCheck(client)
	if err != nil {
		log.Error("Found problems with inspector configuration", logger.Fields("Error", err))
		panic("Failed to start inspector, exiting midleware")
	} else {
		return &Inspector{
			client: client,
			logger: log,
		}
	}
}

// Performs health check to validate the client
func clamHealthCheck(client *clamav.ClamClient) error {
	// Change for ping?
	serverStatus, err := client.Stats()
	if err != nil {
		return err
	} else if serverStatus.ThreadsMax == 0 {
		// validate more items from stats
		return errors.New("no available threads on remote anti-virus server")
	}

	return nil
}

// Inspects request body byte buffer with retries and timeout
func (i *Inspector) InspectBody(body []byte) (bool, error) {
	const maxRetries = 3
	const baseDelay = 100 * time.Millisecond

	// Runs retry with tickers in order to not block threads
	ticker := time.NewTicker(baseDelay)

	// Create a ratio of time x body size in order to improve timeout (considering retries)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var results []cli.ScanResult
	var err error
	var clean bool

	scanStart := time.Now()

	// Check if inner loop is exiting without redundant inspection
	for attempt := 0; attempt < maxRetries; attempt++ {
		// Wait for either the result or context cancellation
		select {
		case <-ctx.Done():
			return false, fmt.Errorf("inspection timeout: %w", ctx.Err())
		case <-ticker.C:
			results, err = i.client.Instream(body)
			if err == nil {
				i.logger.Info("Body inspection completed", nil)
			} else {
				i.logger.Error("Failed to scan ", logger.Fields("Error", err, "Attempt", attempt))
				attempt++
				continue
			}
		}
	}

	// Check results in order to validate status (one inspection can hold multiple statuses)
	for _, result := range results {
		if result.Status == "CLEAN" {
			clean = true
		} else {
			clean = false
		}
		// log results in an uniform manner
		i.logger.Info("inspection completed", logger.Fields(
			"scanResult", result.Status,
			"virusSignature", result.Virus,
			"scanDuration", time.Since(scanStart).Seconds(),
		))
	}

	return clean, nil
}
