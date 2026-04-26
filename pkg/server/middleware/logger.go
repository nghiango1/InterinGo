package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// LogFormatterParams is the structure any formatter will be handed when time to log comes
type LogFormatterParams struct {
	Request *http.Request

	// TimeStamp shows the time after the server returns a response.
	TimeStamp time.Time
	// StatusCode is HTTP response code.
	StatusCode int
	// Latency is how much time the server cost to process a certain request.
	Latency time.Duration
	// ClientIP equals Context's ClientIP method.
	ClientIP string
	// Method is the HTTP method given to the request.
	Method string
	// Path is a path the client requests.
	Path string
	// ErrorMessage is set if error has occurred in processing the request.
	ErrorMessage string
	// BodySize is the size of the Response Body
	BodySize int
	// Keys are the keys set on the request's context.
	Keys map[any]any
}

// defaultLogFormatter is the default log format function Logger middleware uses.
func slogFormater(param LogFormatterParams) {
	switch {
	case param.Latency > time.Minute:
		param.Latency = param.Latency.Truncate(time.Second * 10)
	case param.Latency > time.Second:
		param.Latency = param.Latency.Truncate(time.Millisecond * 10)
	case param.Latency > time.Millisecond:
		param.Latency = param.Latency.Truncate(time.Microsecond * 10)
	}

	var level slog.Level
	code := param.StatusCode
	switch {
	case code >= http.StatusContinue && code < http.StatusOK:
		fallthrough
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		fallthrough
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		level = slog.LevelInfo
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		level = slog.LevelWarn
	default:
		level = slog.LevelError
	}

	slog.Log(context.Background(), level, "GIN logger",
		"statusCode", param.StatusCode,
		"latency", param.Latency,
		"clientIp", param.ClientIP,
		"method", param.Method,
		"path", param.Path,
		"error", param.ErrorMessage,
	)
}

func Logger(c *gin.Context) {
	// Start timer
	start := time.Now()
	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery

	// Process request
	c.Next()

	param := LogFormatterParams{
		Request: c.Request,
		Keys:    c.Keys,
	}

	// Stop timer
	param.TimeStamp = time.Now()
	param.Latency = param.TimeStamp.Sub(start)

	param.ClientIP = c.ClientIP()
	param.Method = c.Request.Method
	param.StatusCode = c.Writer.Status()
	param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

	param.BodySize = c.Writer.Size()

	if raw != "" {
		path = path + "?" + raw
	}

	param.Path = path

	slogFormater(param)
}
