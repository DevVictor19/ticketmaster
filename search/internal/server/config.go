package server

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server *Server
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
		log.Fatalf("missing %s env\n", key)
	}

	return val
}

func getInt(key string) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("missing %s env\n", key)
	}

	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		log.Fatal(err)
	}

	return valAsInt
}
