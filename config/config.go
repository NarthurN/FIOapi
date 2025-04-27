package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBconfig
	ServerConfig
	ClientConfig
}

type DBconfig struct {
	DBpath     string
	DBhost     string
	DBport     string
	DBuser     string
	DBpassword string
	DBname     string
	DBsslmode  string
}

type ServerConfig struct {
	ServerHost         string
	ServerPort         string
	ServerMode         string
	ServerReadTimeout  time.Duration
	ServerWriteTimeout time.Duration
	ServerIdleTimeout  time.Duration
}

type ClientConfig struct {
	AgePath       string
	GenderPath    string
	NatioPath     string
	ClientTimeout time.Duration
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Невозможно загрузить config: %s", err)
	}

	config := &Config{}

	config.DBhost = os.Getenv("DB_HOST")
	config.DBport = os.Getenv("DB_PORT")
	config.DBuser = os.Getenv("DB_USER")
	config.DBpassword = os.Getenv("DB_PASSWORD")
	config.DBname = os.Getenv("DB_NAME")
	config.DBsslmode = os.Getenv("DB_SSLMODE")
	config.DBpath = fmt.Sprintf(
		"user=%s password=%s host=%s port=%s database=%s sslmode=%s",
		config.DBuser, config.DBpassword, config.DBhost, config.DBport, config.DBname, config.DBsslmode,
	)

	config.ServerHost = os.Getenv("SERVER_HOST")
	config.ServerPort = os.Getenv("SERVER_PORT")
	config.ServerMode = os.Getenv("SERVER_MODE")
	var err error
	config.ServerReadTimeout, err = time.ParseDuration(os.Getenv("SERVER_READ_TIMEOUT"))
	if err != nil {
		log.Fatal("Не смог прочитать SERVER_READ_TIMEOUT:", err)
	}
	config.ServerWriteTimeout, err = time.ParseDuration(os.Getenv("SERVER_WRITE_TIMEOUT"))
	if err != nil {
		log.Fatal("Не смог прочитать SERVER_WRITE_TIMEOUT:", err)
	}
	config.ServerIdleTimeout, err = time.ParseDuration(os.Getenv("SERVER_IDLE_TIMEOUT"))
	if err != nil {
		log.Fatal("Не смог прочитать SERVER_IDLE_TIMEOUT:", err)
	}

	config.AgePath = os.Getenv("AGE_PATH")
	config.GenderPath = os.Getenv("SEX_PATH")
	config.NatioPath = os.Getenv("NATIO_PATH")
	config.ClientTimeout, err = time.ParseDuration(os.Getenv("CLIENT_TIMEOUT"))
	if err != nil {
		log.Fatal("Не смог прочитать CLIENT_TIMEOUT:", err)
	}

	return config
}
