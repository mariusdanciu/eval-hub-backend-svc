package features

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.ibm.com/julpayne/eval-hub-backend-svc/internal/config"
	"github.ibm.com/julpayne/eval-hub-backend-svc/internal/logging"
	"github.ibm.com/julpayne/eval-hub-backend-svc/internal/server"

	"github.com/cucumber/godog"
)

type apiFeature struct {
	server     *server.Server
	httpServer *http.Server
	client     *http.Client
	response   *http.Response
	body       []byte
	baseURL    string
}

func createServer(port int) (*server.Server, error) {
	logger, err := logging.NewLogger()
	if err != nil {
		return nil, err
	}
	return server.NewServer(logger, &config.Config{Service: &config.ServiceConfig{Port: port}})
}

func (a *apiFeature) theServiceIsRunning(ctx context.Context) error {
	origPort := 8080
	newServer, err := createServer(origPort)
	if err != nil {
		return err
	}
	a.server = newServer

	// Create a test server
	handler, err := a.server.SetupRoutes()
	a.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", origPort),
		Handler: handler,
	}

	// Start server in background
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", origPort))
	if err != nil {
		return err
	}

	port := listener.Addr().(*net.TCPAddr).Port
	a.baseURL = fmt.Sprintf("http://localhost:%d", port)

	go func() {
		a.httpServer.Serve(listener)
	}()

	// Wait for server to start
	time.Sleep(200 * time.Millisecond)

	a.client = &http.Client{
		Timeout: 5 * time.Second,
	}

	return nil
}

func (a *apiFeature) iSendARequestTo(method, path string) error {
	url := fmt.Sprintf("%s%s", a.baseURL, path)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	a.response, err = a.client.Do(req)
	if err != nil {
		return err
	}

	a.body, err = io.ReadAll(a.response.Body)
	if err != nil {
		return err
	}
	a.response.Body.Close()

	return nil
}

func (a *apiFeature) theResponseStatusShouldBe(status int) error {
	if a.response.StatusCode != status {
		return fmt.Errorf("expected status %d, got %d", status, a.response.StatusCode)
	}
	return nil
}

func (a *apiFeature) theResponseShouldBeJSON() error {
	contentType := a.response.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return fmt.Errorf("expected JSON content type, got %s", contentType)
	}

	var js interface{}
	if err := json.Unmarshal(a.body, &js); err != nil {
		return fmt.Errorf("response is not valid JSON: %v", err)
	}

	return nil
}

func (a *apiFeature) theResponseShouldContainWithValue(key, value string) error {
	var data map[string]interface{}
	if err := json.Unmarshal(a.body, &data); err != nil {
		return err
	}

	if data[key] != value {
		return fmt.Errorf("expected %s to be %s, got %v", key, value, data[key])
	}

	return nil
}

func (a *apiFeature) theResponseShouldContain(key string) error {
	var data map[string]interface{}
	if err := json.Unmarshal(a.body, &data); err != nil {
		return err
	}

	if _, ok := data[key]; !ok {
		return fmt.Errorf("response does not contain key: %s", key)
	}

	return nil
}

func (a *apiFeature) theResponseShouldContainPrometheusMetrics() error {
	bodyStr := string(a.body)
	if !strings.Contains(bodyStr, "# HELP") || !strings.Contains(bodyStr, "# TYPE") {
		return fmt.Errorf("response does not appear to be Prometheus metrics format")
	}
	return nil
}

func (a *apiFeature) theMetricsShouldInclude(metricName string) error {
	bodyStr := string(a.body)
	if !strings.Contains(bodyStr, metricName) {
		return fmt.Errorf("metrics do not include %s", metricName)
	}
	return nil
}

func (a *apiFeature) theMetricsShouldShowRequestCountFor(path string) error {
	bodyStr := string(a.body)
	// Check if metrics contain the path
	if !strings.Contains(bodyStr, path) {
		return fmt.Errorf("metrics do not show requests for path %s", path)
	}
	return nil
}

func (a *apiFeature) resetResponse(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
	a.response = nil
	a.body = nil
	return ctx, nil
}

func (a *apiFeature) cleanup(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
	if a.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		a.httpServer.Shutdown(ctx)
	}
	return ctx, nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	api := &apiFeature{}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		return api.resetResponse(ctx, sc)
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		return api.cleanup(ctx, sc, err)
	})

	ctx.Step(`^the service is running$`, api.theServiceIsRunning)
	ctx.Step(`^I send a (GET|POST|PUT|DELETE) request to "([^"]*)"$`, api.iSendARequestTo)
	ctx.Step(`^the response status should be (\d+)$`, api.theResponseStatusShouldBe)
	ctx.Step(`^the response should be JSON$`, api.theResponseShouldBeJSON)
	ctx.Step(`^the response should contain "([^"]*)" with value "([^"]*)"$`, api.theResponseShouldContainWithValue)
	ctx.Step(`^the response should contain "([^"]*)"$`, api.theResponseShouldContain)
	ctx.Step(`^the response should contain Prometheus metrics$`, api.theResponseShouldContainPrometheusMetrics)
	ctx.Step(`^the metrics should include "([^"]*)"$`, api.theMetricsShouldInclude)
	ctx.Step(`^the metrics should show request count for "([^"]*)"$`, api.theMetricsShouldShowRequestCountFor)
}
