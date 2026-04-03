package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/DevVictor19/event/internal/handlers"
	"github.com/DevVictor19/event/internal/repositories"
	"github.com/DevVictor19/event/internal/usecases"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func Start() {
	cfg := loadConfig()
	db := connectDB(cfg.Postgres)

	e := echo.New()
	e.Use(middleware.RequestLogger())

	eventRepository := repositories.NewEventRepository(db)
	findEventByUuidUC := usecases.NewFindEventByUuidUC(eventRepository)
	eventsHandler := handlers.NewEventsHandler(findEventByUuidUC)

	api := e.Group("/api/v1/management/events")
	api.GET("/:uuid", eventsHandler.FindByUUID)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

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
}
