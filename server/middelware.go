package server

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NoError(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_ = next(c)
		return nil
	}
}

func Logger() echo.MiddlewareFunc {
	return logger
}
func logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		res := c.Response()

		start := time.Now()
		err := next(c)
		elapsed := time.Since(start)

		var logger *zerolog.Event
		if res.Status >= 500 {
			logger = log.Error()
		} else {
			logger = log.Info()
		}

		logger.Str("protocol", req.Proto)
		logger.Str("method", req.Method)
		logger.Str("url", req.URL.String())
		logger.Str("from", req.RemoteAddr)
		logger.Str("user-agent", req.UserAgent())

		logger.Int("status-code", res.Status)
		logger.Dur("elapsed", elapsed)

		if err != nil {
			logger.Err(err)
		}

		logger.Send()

		return err
	}
}
