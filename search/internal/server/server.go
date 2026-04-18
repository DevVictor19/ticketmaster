package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DevVictor19/search/internal/brokers/cdc"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func Start() {
	cfg := loadConfig()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	eventDbConsumer := cdc.NewEventDbConsumer(cfg.Kafka.Listeners)
	eventDbConsumer.Start()

	e := echo.New()
	e.Use(middleware.RequestLogger())

	api := e.Group("/api/v1/search")
	api.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	sc := echo.StartConfig{
		Address: cfg.Server.Addr,
		BeforeServeFunc: func(s *http.Server) error {
			s.ReadTimeout = cfg.Server.ReadTimeout
			s.WriteTimeout = cfg.Server.WriteTimeout
			s.MaxHeaderBytes = cfg.Server.MaxHeaderBytes
			return nil
		},
	}

	if err := sc.Start(ctx, e); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	eventDbConsumer.Stop(shutdownCtx)
}
