package auth

import (
	"github.com/huijiewei/agile-go/agile-framework/problem"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	HeaderClientId    = "X-Client-Id"
	HeaderAccessToken = "X-Access-Token"
)

type (
	AgileAuthConfig struct {
		Skipper middleware.Skipper

		Validator AgileAuthValidator

		ContextKey string
	}

	AgileAuthValidator func(string, string) any
)

var (
	DefaultAgileAuthConfig = AgileAuthConfig{
		Skipper:    middleware.DefaultSkipper,
		ContextKey: "user",
	}
)

func AgileAuth(fn AgileAuthValidator) echo.MiddlewareFunc {
	c := DefaultAgileAuthConfig
	c.Validator = fn
	return AgileAuthWithConfig(c)
}

func AgileAuthWithConfig(config AgileAuthConfig) echo.MiddlewareFunc {

	if config.Validator == nil {
		panic("echo: agile-auth middleware requires a validator function")
	}

	if config.Skipper == nil {
		config.Skipper = DefaultAgileAuthConfig.Skipper
	}

	if config.ContextKey == "" {
		config.ContextKey = DefaultAgileAuthConfig.ContextKey
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			clientId := c.Request().Header.Get(HeaderClientId)
			accessToken := c.Request().Header.Get(HeaderAccessToken)

			if len(clientId) > 0 && len(accessToken) > 0 {
				user := config.Validator(clientId, accessToken)

				if user != nil {
					c.Set(config.ContextKey, user)

					return next(c)
				}
			}

			return problem.ErrUnauthorized
		}
	}
}
