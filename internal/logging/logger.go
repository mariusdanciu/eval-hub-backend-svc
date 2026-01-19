package logging

import (
	"log/slog"
	"net/http"

	"github.ibm.com/julpayne/eval-hub-backend-svc/internal/constants"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates and returns a new structured logger using zap as the underlying
// logging implementation, wrapped with slog's interface. The logger is configured
// with production settings and ISO8601 time encoding for consistent log formatting.
//
// Returns:
//   - *slog.Logger: A structured logger instance that can be used throughout the application
//   - error: An error if the logger could not be initialized
func NewLogger() (*slog.Logger, error) {
	var logConfig zap.Config
	logConfig = zap.NewProductionConfig()
	logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapLog, err := logConfig.Build()
	if err != nil {
		return nil, err
	}
	return slog.New(zapslog.NewHandler(zapLog.Core())), nil
}

// LoggerWithRequest enhances a logger with request-specific fields
func LoggerWithRequest(logger *slog.Logger, r *http.Request) *slog.Logger {
	// Extract RequestID from X-Global-Transaction-Id header, or generate a UUID if not present
	requestID := r.Header.Get("X-Global-Transaction-Id")
	if requestID == "" {
		requestID = uuid.New().String()
	}

	// Add request_id to logger using With
	enhancedLogger := logger.With(constants.LOG_REQUEST_ID, requestID)

	// Extract and add HTTP method and URI if they exist
	method := r.Method
	if method != "" {
		enhancedLogger = enhancedLogger.With(constants.LOG_METHOD, method)
	}

	uri := ""
	if r.URL != nil {
		uri = r.URL.Path
	}
	if uri == "" {
		uri = r.RequestURI
	}
	if uri != "" {
		enhancedLogger = enhancedLogger.With(constants.LOG_URI, uri)
	}

	// Extract and add HTTP request fields to logger if they exist
	userAgent := r.Header.Get("User-Agent")
	if userAgent != "" {
		enhancedLogger = enhancedLogger.With(constants.LOG_USER_AGENT, userAgent)
	}

	remoteAddr := r.RemoteAddr
	if remoteAddr != "" {
		enhancedLogger = enhancedLogger.With(constants.LOG_REMOTE_ADR, remoteAddr)
	}

	// Extract remote_user from URL user info or header
	remoteUser := ""
	if r.URL != nil && r.URL.User != nil {
		remoteUser = r.URL.User.Username()
	}
	if remoteUser == "" {
		remoteUser = r.Header.Get("Remote-User")
	}
	if remoteUser != "" {
		enhancedLogger = enhancedLogger.With(constants.LOG_USER, remoteUser)
	}

	referer := r.Header.Get("Referer")
	if referer != "" {
		enhancedLogger = enhancedLogger.With(constants.LOG_REFERER, referer)
	}

	return enhancedLogger
}
