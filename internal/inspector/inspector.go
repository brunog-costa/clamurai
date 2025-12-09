// inspector.go
package inspector

import (
	"context"
	"fmt"
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
func New(clamAddress string, readTimeOut uint64, connectionTimeout uint64) (*Inspector, error) {
	hostname, _ := os.Hostname()

	log, err := logger.New(logger.Config{
		Middleware: "inspector",
		Output:     os.Stdout,
		Hostname:   hostname,
	})
	if err != nil {
		return nil, err
	}

	// Use provided client or create real one
	client := clamav.NewClamClient("tcp", clamAddress,
		time.Duration(connectionTimeout)*time.Second,
		time.Duration(readTimeOut)*time.Second)

	return &Inspector{
		client: client,
		logger: log,
	}, nil
}

// Inspects request body byte buffer with retries and timeout
func (i *Inspector) InspectBody(body []byte) (bool, string, error) {
	const maxRetries = 3

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var signature string
	var clean bool
	var err error

	// Custom inspection with retry loop
	for attempt := 0; attempt < maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return false, "", fmt.Errorf("inspection timeout: %w", ctx.Err())
		default:
			results, err := i.client.Instream(body)
			if err == nil {
				i.logger.Info("Body inspection completed", nil)
				for _, result := range results {
					i.logger.Info("Body inspection completed", logger.Fields("status", result.Status))
					if result.Status != "OK" {
						return false, result.Virus, nil
					} else {
						return true, result.Virus, nil
					}
				}
				// Exit retry loop on success
				i.logger.Info("Inspector finished checking body", nil)
				break
			}

			i.logger.Error("Failed to scan", logger.Fields("Error", err, "Attempt", attempt))

			// Wait before retry, but respect context timeout
			if attempt < maxRetries-1 {
				select {
				case <-time.After(time.Duration(attempt+1) * 1000 * time.Millisecond):
				}
			}
		}
	}

	// Give found data back to main
	return clean, signature, err
}
