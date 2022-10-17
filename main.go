package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

type LogSupportHandler = func(zerolog.Logger) echo.HandlerFunc

// log output counter(≒access counter)
var logCount = 1

func main() {
	e := echo.New()
	logger := zerolog.New(os.Stdout)

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI: true,
		LogLatency: true,
		LogRemoteIP: true,
		LogUserAgent: true,
		LogError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logCount++
			// Access logs are output no error level to distinguish them from info logs.
			logger.Log().
				Str("logtime", v.StartTime.Format("2006-01-02 15:04:05")).
				Str("URI", v.URI).
				Int("Status", v.Status).
				Str("latency", v.Latency.String()).
				Str("remoteIP", v.RemoteIP).
				Str("User-Agent", v.UserAgent).
				Send()

			return nil
		},
	}))

	e.GET("/", helloWorld(logger))
	e.GET("/search", noQueryRows(logger))
	e.GET("/auth", authenticationFailed(logger))
	e.GET("/fatal", intentionalFatal(logger))

	e.Logger.Fatal(e.Start(":1323"))
}

// helloWorld responsed with the message 'Hello, World!'
// Outputs info level logs using zerolog.
func helloWorld (l zerolog.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		l.Info().Int("LogNo", logCount).Send()
		return c.String(http.StatusOK, "Hello, World!")
	}
}

// noQueryRows behaves as if record not found with fetch to database.
// Outputs warn level logs using zerolog.
func noQueryRows (l zerolog.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		l.Warn().Int("LogNo", logCount).Send()
		return c.String(http.StatusSeeOther, "records not found.")
	}
}

// authenticationFailed behaves as if the authentication failed.
func authenticationFailed (l zerolog.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		l.Error().Int("LogNo", logCount).Send()
		return c.String(http.StatusUnauthorized, "authentication failed.")
	}
}

// intentionalFatal as the name implies, intentionally generates an error.
// Used for verification of recovery and life/death monitoring.
func intentionalFatal (l zerolog.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Fatal("Intentional error")
		// Doesn't pass bacause os.Exit(1) is called just before.
		return c.String(http.StatusInternalServerError, "intentional error.")
	}
}
