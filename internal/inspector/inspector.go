// inspector.go
package inspector

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/brunog-costa/clamurai/pkg/logger"
	"github.com/hq0101/go-clamav/pkg/clamav"
	"github.com/hq0101/go-clamav/pkg/cli"
)

// Define interfaces for the external dependencies
type ClamClient interface {
	Stats() (*cli.ClamdStats, error)
	Instream(body []byte) ([]cli.ScanResult, error)
}

type Inspector struct {
	client ClamClient
	logger *logger.JSONLogger
}

type InspectorWorker interface {
	clamHealthCheck(client ClamClient) error
	InspectBody(body []byte) (bool, error)
}

// Builds new inspector from config
func New(clamAddress string, readTimeOut uint64, connectionTimeout uint64) *Inspector {
	hostname, _ := os.Hostname()

	log := logger.NewJSONLogger(logger.Config{
		Middleware: "inspector",
		Output:     io.MultiWriter(),
		Hostname:   hostname,
	})

	// Use provided client or create real one
	client := clamav.NewClamClient("tcp", clamAddress,
		time.Duration(connectionTimeout)*time.Second,
		time.Duration(readTimeOut)*time.Second)

	return &Inspector{
		client: client,
		logger: log,
	}
}

// Inspects request body byte buffer with retries and timeout
func (i *Inspector) InspectBody(body []byte) (bool, error) {
	const maxRetries = 3

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var results []cli.ScanResult
	var err error
	scanStart := time.Now()

	for attempt := 0; attempt < maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return false, fmt.Errorf("inspection timeout: %w", ctx.Err())
		default:
			results, err = i.client.Instream(body)
			if err == nil {
				i.logger.Info("Body inspection completed", nil)
				break // Exit retry loop on success
			}

			i.logger.Error("Failed to scan", logger.Fields("Error", err, "Attempt", attempt))

			// Wait before retry, but respect context timeout
			if attempt < maxRetries-1 {
				select {
				case <-time.After(time.Duration(attempt+1) * 100 * time.Millisecond):
				case <-ctx.Done():
					return false, fmt.Errorf("inspection timeout during retry: %w", ctx.Err())
				}
			}
		}
	}

	if err != nil {
		return false, err
	}

	// Process results
	clean := true
	for _, result := range results {
		if result.Status != "CLEAN" {
			clean = false
		}

		i.logger.Info("inspection completed", logger.Fields(
			"scanResult", result.Status,
			"virusSignature", result.Virus,
			"scanDuration", time.Since(scanStart).Seconds(),
		))
	}

	return clean, nil
}
