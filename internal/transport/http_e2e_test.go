package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/zhenisduissekov/another-dummy-service/internal/common/server"
	"github.com/zhenisduissekov/another-dummy-service/internal/repository/inmem"
	"github.com/zhenisduissekov/another-dummy-service/internal/services"
)

type HttpTestSuite struct {
	suite.Suite
	portService PortService
	httpServer  HttpServer
}

func NewTestSuite() *HttpTestSuite {
	s := &HttpTestSuite{}

	portStore := inmem.NewPortStore()

	s.portService = services.NewPortService(portStore)
	s.httpServer = NewHttpServer(s.portService)

	return s
}

func TestHttpTestSuite(t *testing.T) {
	suite.Run(t, NewTestSuite())
}

func (suite *HttpTestSuite) TestUploadPorts_badJSON() {
	// Load test fixtures
	portsRequest := []byte(`blabla`)
	portsResponse := []byte(`{"slug":"invalid json","message":"Bad request","httpStatus":400,"details":null,"timestamp":"2025-01-03T09:01:56Z"}`)

	// Execute the request
	responseData, res := suite.executeUploadPortsRequest(portsRequest)

	// Validate response
	suite.validateResponse(res, portsResponse, responseData, http.StatusBadRequest)
}

func (suite *HttpTestSuite) TestUploadPorts() {
	// Load test fixtures
	portsRequest := suite.loadFixture("testfixtures/ports_request.json")
	portsResponse := suite.loadFixture("testfixtures/ports_response.json")

	// Execute the request
	responseData, res := suite.executeUploadPortsRequest(portsRequest)

	// Validate response
	suite.validateResponse(res, portsResponse, responseData, http.StatusOK)

	// Validate stored ports
	suite.validateStoredPorts(portsRequest)
}

func (suite *HttpTestSuite) loadFixture(path string) []byte {
	data, err := os.ReadFile(path)
	require.NoError(suite.T(), err, "Failed to load fixture: %s", path)
	return data
}

func (suite *HttpTestSuite) executeUploadPortsRequest(portsRequest []byte) ([]byte, *http.Response) {
	// Create POST /ports request
	req := httptest.NewRequest(http.MethodPost, "/ports", bytes.NewBuffer(portsRequest))
	w := httptest.NewRecorder()

	// Run the request
	suite.httpServer.UploadPorts(w, req)

	res := w.Result()
	defer func() {
		err := res.Body.Close()
		require.NoError(suite.T(), err, "Failed to close response body")
	}()

	data, err := io.ReadAll(res.Body)
	require.NoError(suite.T(), err, "Failed to read response body")

	return data, res
}

func (suite *HttpTestSuite) validateResponse(res *http.Response, portsResponse, actualResponse []byte, statusResponse int) {
	require.Equal(suite.T(), statusResponse, res.StatusCode, "Unexpected status code")

	var expected, actual server.ResponseOK
	require.NoError(suite.T(), json.Unmarshal(portsResponse, &expected), "Failed to unmarshal expected response")
	require.NoError(suite.T(), json.Unmarshal(actualResponse, &actual), "Failed to unmarshal actual response")

	// Normalize timestamps
	now := time.Now().String()
	expected.Timestamp = now
	actual.Timestamp = now

	require.Equal(suite.T(), expected, actual, "Response does not match expected")
}

func (suite *HttpTestSuite) validateStoredPorts(portsRequest []byte) {
	storedPortsTotal, err := suite.portService.CountPorts(context.Background())
	require.NoError(suite.T(), err, "Failed to count stored ports")

	expectedPortsCount := countJSONPorts(suite.T(), portsRequest)
	require.Equal(suite.T(), expectedPortsCount, storedPortsTotal, "Stored ports count does not match")
}

func countJSONPorts(t *testing.T, data []byte) int {
	t.Helper()
	var ports map[string]struct{}
	err := json.Unmarshal(data, &ports)
	require.NoError(t, err)

	return len(ports)
}
