package inspector

import (
	"context"
	"errors"
	"fmt"
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
}

// Builds new inspector from config
func NewInspector(clamAddress string, readTimeOut uint64, connectionTimeout uint64) *Inspector {
	// Provisions custom logger for inspector
	log := logger.NewJSONLogger("Inspector")

	// Creates a client and perform a healthcheck in order to validate if server has capacity
	client := clamav.NewClamClient("tcp", clamAddress, time.Duration(connectionTimeout)*time.Second, time.Duration(readTimeOut)*time.Second)

	err := clamHealthCheck(client)
	if err != nil {
		log.ErrorWithFields("Found problems with inspector configuration", logger.Fields("Error", err))
		panic("Failed to start inspector, exiting")
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

	// Check if inner loop is exiting without redundant inspection
	for attempt := 0; attempt < maxRetries; attempt++ {
		// Wait for either the result or context cancellation
		select {
		case <-ctx.Done():
			return false, fmt.Errorf("inspection timeout: %w", ctx.Err())
		case <-ticker.C:
			results, err = i.client.Instream(body)
			if err == nil {
				i.logger.Info("Body inspection completed")
			} else {
				i.logger.ErrorWithFields("Failed to scan ", logger.Fields("Error", err, "Attempt", attempt))
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
	}

	return clean, nil
}
