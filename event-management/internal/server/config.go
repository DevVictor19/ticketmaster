package server

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Postgres *Postgres
	Server   *Server
}

type Postgres struct {
	Host            string
	User            string
	Password        string
	DBName          string
	Port            string
	SSLMode         string
	TimeZone        string
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	MaxIdleConns    int
	ConnMaxIdleTime time.Duration
}

func (p *Postgres) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		p.Host, p.User, p.Password, p.DBName, p.Port, p.SSLMode, p.TimeZone)
}

type Server struct {
	Mode           string
	Addr           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
	CtxTimeout     time.Duration
}

func loadConfig() *Config {
	return &Config{
		Postgres: &Postgres{
			Host:            getString("POSTGRES_HOST"),
			User:            getString("POSTGRES_USER"),
			Password:        getString("POSTGRES_PASSWORD"),
			DBName:          getString("POSTGRES_DBNAME"),
			Port:            getString("POSTGRES_PORT"),
			SSLMode:         getString("POSTGRES_SSLMODE"),
			TimeZone:        getString("POSTGRES_TIMEZONE"),
			MaxOpenConns:    getInt("POSTGRES_MAX_OPEN_CONNS"),
			ConnMaxLifetime: time.Duration(getInt("POSTGRES_CONN_MAX_LIFETIME")) * time.Second,
			MaxIdleConns:    getInt("POSTGRES_MAX_IDLE_CONNS"),
			ConnMaxIdleTime: time.Duration(getInt("POSTGRES_CONN_MAX_IDLE_TIME")) * time.Second,
		},
		Server: &Server{
			Mode:           getString("SERVER_MODE"),
			Addr:           getString("SERVER_ADDR"),
			ReadTimeout:    time.Duration(getInt("SERVER_READ_TIMEOUT")) * time.Second,
			WriteTimeout:   time.Duration(getInt("SERVER_WRITE_TIMEOUT")) * time.Second,
			MaxHeaderBytes: getInt("SERVER_MAX_HEADER_BYTES"),
			CtxTimeout:     time.Duration(getInt("SERVER_CTX_TIMEOUT")) * time.Second,
		},
	}
}

func getString(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("missing %s on .env file\n", key)
	}

	return val
}

func getInt(key string) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("missing %s on .env file\n", key)
	}

	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		log.Fatal(err)
	}

	return valAsInt
}
