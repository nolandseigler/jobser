package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/nolandseigler/wordser/wordserweb/internal/template"
)

func main() {
	// Setup
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(echoprometheus.NewMiddleware("wordserweb"))
	e.GET("/metrics", echoprometheus.NewHandler())

	e.Renderer = template.New()

	e.GET("/signup", func(c echo.Context) error {
		return c.Render(http.StatusOK, "signup", []string{"fucker"})
	})

	e.GET("/dashboard", func(c echo.Context) error {

		requestURL := "http://wordser:8080/echo"
		req, err := http.NewRequest(
			http.MethodPost,
			requestURL,
			bytes.NewReader([]byte(`{"text": "A message from CS361"}`)),
		)
		if err != nil {
			return c.JSON(500, fmt.Sprintf(`{"err": "%s"}`, err))
		}
		
		req.Header.Set("Content-Type", "application/json")

		client := http.Client{
			Timeout: 30 * time.Second,
		}

		e.Logger.Info("sending message to webser service!")
		res, err := client.Do(req)
		if err != nil {
			return c.JSON(500, fmt.Sprintf(`{"err": "%s"}`, err))
		}

		e.Logger.Info(res)
		
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
