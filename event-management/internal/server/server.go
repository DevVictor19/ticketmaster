package server

import (
	"fmt"
	"log"

	"github.com/DevVictor19/event/internal/entities"
	"github.com/DevVictor19/event/internal/handlers"
	"github.com/DevVictor19/event/internal/repositories"
	"github.com/DevVictor19/event/internal/usecases"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
	SSLMode  string
	TimeZone string
}

func (p *Postgres) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		p.Host, p.User, p.Password, p.DBName, p.Port, p.SSLMode, p.TimeZone)
}

type Env struct{}

func Start() {
	config := &Postgres{
		Host:     "localhost",
		User:     "main",
		Password: "password",
		DBName:   "event-management",
		Port:     "5432",
		SSLMode:  "disable",
		TimeZone: "UTC",
	}
	dsn := config.DSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(
		entities.Venue{},
		entities.Performer{},
		entities.Event{},
		entities.Ticket{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	e := echo.New()
	e.Use(middleware.RequestLogger())

	eventRepository := repositories.NewEventRepository(db)
	findEventByUuidUC := usecases.NewFindEventByUuidUC(eventRepository)
	eventsHandler := handlers.NewEventsHandler(findEventByUuidUC)

	api := e.Group("/api/v1/management/events")
	api.GET("/:uuid", eventsHandler.FindByUUID)

	if err := e.Start(":8000"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
