package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/nolandseigler/wordser/wordserweb/internal/auth"
	"github.com/nolandseigler/wordser/wordserweb/internal/handlers"
	"github.com/nolandseigler/wordser/wordserweb/internal/storage/postgres"
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

	ctx := context.Background()
	// init and drop in storage.
	db, err := postgres.New(ctx, postgres.Config{})
	if err != nil {
		e.Logger.Fatal(err)
	}
	auth.New(auth.Config{}, NewStore(), db)
	// this going to trash prom??
	e.Use()

	e.Renderer = template.New()

	e.GET("/signup", handlers.GetSignupHandler)
	e.GET("/dashboard", handlers.GetDashboardHandler)

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

type TempKVStore struct {
	store *sync.Map
}

func (t *TempKVStore) Insert(key string, value string) error {
	t.store.Store(key, value)
	return nil
}
func (t *TempKVStore) Delete(key string) error {
	t.store.Delete(key)
	return nil
}
func (t *TempKVStore) Get(key string) (string, bool) {
	if val, ok := t.store.Load(key); ok {
		val, ok := val.(string)
		return val, ok
	}
	return "", false
}

func NewStore() *TempKVStore {
	return &TempKVStore{
		&sync.Map{},
	}
}
