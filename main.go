package main

import (
	"discord-22-api/config"
	"discord-22-api/ctx"
	"discord-22-api/router"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{DisableTimestamp: true})

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Format: `{"id":"${id}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
		`"status":${status},"latency":"${latency}","bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n"}))
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ex := &ctx.Extended{Context: c}
			return next(ex)
		}
	})
	e.Use(middleware.RequestID())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := c.Request().Header.Get("X-22-KEY")
			if key == config.Instance().RootApiKey {
				return next(c)
			}
			return c.NoContent(http.StatusUnauthorized)
		}
	})
	router.Assign(e)
	e.Logger.Fatal(e.Start(":8888"))
}
