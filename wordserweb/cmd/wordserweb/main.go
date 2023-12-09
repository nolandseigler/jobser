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
	authpkg "github.com/nolandseigler/wordser/wordserweb/internal/auth"
	"github.com/nolandseigler/wordser/wordserweb/internal/handlers"
	"github.com/nolandseigler/wordser/wordserweb/internal/static"
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

	e.Renderer = template.New()
	static.RegisterStaticFS(e)
	

	ctx := context.Background()

	dbCfg, err := postgres.ConfigFromEnv()
	if err != nil {
		e.Logger.Fatal(err)
	}
	db, err := postgres.New(ctx, dbCfg)
	if err != nil {
		e.Logger.Fatal(err)
	}

	authCfg, err := authpkg.ConfigFromEnv()
	if err != nil {
		e.Logger.Fatal(err)
	}

	auth, err := authpkg.New(ctx, authCfg, NewStore(), db)
	if err != nil {
		e.Logger.Fatal(err)
	}
	// this going to trash prom??
	e.Use(auth.ValidateJWTMiddleWare)


	e.GET("/signup", handlers.GetSignupHandler)
	e.POST("/signup", handlers.PostSignupHandler(auth, db))
	e.GET("/login", handlers.GetLoginHandler)
	e.POST("/login", handlers.PostLoginHandler(auth, db))
	e.GET("/dashboard", handlers.GetDashboardHandler)
	e.GET("/translate", handlers.GetTranslateHandler)

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
