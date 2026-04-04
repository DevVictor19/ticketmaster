package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/DevVictor19/booking/internal/handlers"
	"github.com/DevVictor19/booking/internal/repositories"
	"github.com/DevVictor19/booking/internal/services"
	"github.com/DevVictor19/booking/internal/usecases"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func Start() {
	cfg := loadConfig()
	db := connectDB(cfg.Postgres)
	rdb := getRedisClient(cfg.Redis)
	defer rdb.Close()

	e := echo.New()
	e.Use(middleware.RequestLogger())

	ticketRepo := repositories.NewTicketRepository(db)
	ticketLockSvc := services.NewTicketLockService(rdb)
	reserveTicketUC := usecases.NewReserveTicketUseCase(ticketRepo, ticketLockSvc)
	confirmTicketUC := usecases.NewConfirmTicketUseCase(ticketRepo, ticketLockSvc)
	ticketsHandler := handlers.NewTicketsHandler(reserveTicketUC, confirmTicketUC)

	api := e.Group("/api/v1/booking")
	ticketsApi := api.Group("/tickets")
	ticketsApi.POST("/reserve", ticketsHandler.ReserveTicket)
	ticketsApi.POST("/confirm", ticketsHandler.ConfirmTicket)

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
