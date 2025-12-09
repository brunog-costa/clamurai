package inspector

import (
	"errors"
	"io"
	"testing"

	"github.com/brunog-costa/clamurai/pkg/logger"
	"github.com/hq0101/go-clamav/pkg/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Creating clam client mock and implementing instream function
type MockClamClient struct {
	mock.Mock
}

func (m *MockClamClient) Instream(body []byte) ([]cli.ScanResult, error) {
	args := m.Called(body)
	return args.Get(0).([]cli.ScanResult), args.Error(1)
}

func (m *MockClamClient) Stats() (*cli.ClamdStats, error) {
	return cli.ParseStatStr(""), nil
}

func (m *MockClamClient) Ping() (string, error) {
	return "", nil
}

// Helper function for returning Instream results
func createScanResult(status string) []cli.ScanResult {
	// Implement other values
	switch status {
	case "clean":
		return []cli.ScanResult{
			{
				Status: "OK",
				Virus:  "",
			},
		}
	case "malware":
		return []cli.ScanResult{
			{
				Status: "MALWARE",
				Virus:  "test.malware.signature",
			},
		}

	}

	return []cli.ScanResult{
		{
			Status: "",
			Virus:  "",
		},
	}
}

func createTestInspector(client ClamClient) *Inspector {

	log, err := logger.New(logger.Config{
		Middleware: "test-inspector",
		Output:     io.Discard,
		Hostname:   "test-host",
	})
	if err != nil {
		panic(err) // This should not happen in tests
	}

	return &Inspector{
		client: client,
		logger: log,
	}
}

func TestNew(t *testing.T) {
	t.Run("should create inspector with valid parameters", func(t *testing.T) {
		// Arrange
		clamAddress := "localhost:3310"
		readTimeout := uint64(3600)
		connectionTimeout := uint64(90)

		// Act
		inspector, err := New(clamAddress, readTimeout, connectionTimeout)

		// Assert
		assert.NoError(t, err, "Should not return error")
		assert.NotNil(t, inspector, "Inspector should not be nil")
		assert.NotNil(t, inspector.client, "Client should not be nil")
		assert.NotNil(t, inspector.logger, "Logger should not be nil")
	})
}

func TestInspectBody(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		expectedOutput bool
		expectedError  error
	}{
		{
			name:           "clean file upload",
			body:           "clean",
			expectedOutput: true,
			expectedError:  nil,
		},
		{
			name:           "malware file upload",
			body:           "malware",
			expectedOutput: false,
			expectedError:  nil,
		},
		{
			name:           "invalid file upload",
			body:           "invalid",
			expectedOutput: false,
			expectedError:  errors.New("Connection Timed Out"),
		},
	}

	mockClient := &MockClamClient{}
	inspector := createTestInspector(mockClient)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock expectations
			if !tc.expectedOutput {
				mockClient.On("Instream", []byte(tc.body)).Return(createScanResult(tc.body), tc.expectedError).Times(3)
			} else {
				mockClient.On("Instream", []byte(tc.body)).Return(createScanResult(tc.body), nil).Once()
			}

			// Act
			output, _, err := inspector.InspectBody([]byte(tc.body))

			// Assert
			if err != nil {
				assert.Error(t, err)
				assert.False(t, output, "When scan is incomplete, inspector should return an error")
			} else {
				switch tc.body {
				case "clean":
					assert.NoError(t, err)
					assert.True(t, output, "Should return true whenever file is clean")

				case "malware":
					assert.NoError(t, err)
					assert.False(t, output, "Should return false whenever file is not clean")

				}
			}
		})
	}
}
