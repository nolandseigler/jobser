package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/nolandseigler/jobser/jobserweb/internal/template"
)

func main() {
	// Setup
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(echoprometheus.NewMiddleware("jobserweb"))
	e.GET("/metrics", echoprometheus.NewHandler())

	e.Renderer = template.New()

	e.GET("/signup", func(c echo.Context) error {
		return c.Render(http.StatusOK, "signup", []string{"fucker"})
	})

	e.GET("/dashboard", func(c echo.Context) error {
		return c.Render(http.StatusOK, "dashboard", "")
	})

	// Start server
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
