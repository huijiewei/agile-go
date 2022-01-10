package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/huijiewei/agile-go/agile-framework/auth"
	"github.com/huijiewei/agile-go/agile-framework/problem"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Disposition"},
		AllowCredentials: true,
	}))

	e.Use(auth.AgileAuthWithConfig(auth.AgileAuthConfig{
		Skipper: func(c echo.Context) bool {
			if strings.Contains(c.Request().URL.Path, "swagger") {
				return true
			}

			return false
		},
		Validator: func(clientId string, accessToken string) any {
			return nil
		},
	}))

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if prb, ok := err.(*problem.Problem); ok {
			if !c.Response().Committed {
				if _, err := prb.WriteTo(c.Response()); err != nil {
					c.Logger().Error(err)
				}
			}
		} else {
			e.DefaultHTTPErrorHandler(err, c)
		}
	}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	go func() {
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
