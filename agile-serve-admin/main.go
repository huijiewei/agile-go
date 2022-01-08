package main

import (
	"net/http"
	"strings"

	"github.com/huijiewei/agile-go/agile-framework/auth"
	"github.com/huijiewei/agile-go/agile-framework/problem"
	_ "github.com/huijiewei/agile-go/agile-serve-admin/docs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"
)

func main() {
	e := echo.New()

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

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
