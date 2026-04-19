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
	"github.com/DevVictor19/search/internal/handlers"
	"github.com/DevVictor19/search/internal/repositories"
	"github.com/DevVictor19/search/internal/services"
	"github.com/DevVictor19/search/internal/usecases"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func Start() {
	cfg := loadConfig()
	esClient := getEsClient(cfg.Elasticsearch)
	db := connectDB(cfg.Postgres)
	rdb := getRedisClient(cfg.Redis)
	defer rdb.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	venueRpo := repositories.NewVenueRepository(db)
	searchEngineSvc := services.NewSearchEngineService(esClient, venueRpo)
	searchCacheSvc := services.NewSearchCacheService(rdb)
	eventDbConsumer := cdc.NewEventDbConsumer(cfg.Kafka.Listeners, searchEngineSvc)
	eventDbConsumer.Start()

	searchEventsUC := usecases.NewSearchEventsUC(searchEngineSvc, searchCacheSvc)
	eventsHandler := handlers.NewEventsHandler(searchEventsUC)

	e := echo.New()
	e.Use(middleware.RequestLogger())

	api := e.Group("/api/v1/search")
	api.GET("/events", eventsHandler.SearchEvents)

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
