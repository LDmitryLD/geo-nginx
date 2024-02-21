package config

import (
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
)

const (
	AppName = "APP_NAME"

	serverPort         = "SERVER_PORT"
	envShutdownTimeout = "SHUTDOWN_TIMEOUT"

	parseShutdownTimeoutError    = "config: parse server shutdown timeout error"
	parseRpcShutdownTimeoutError = "config: parse rpc server shutdown timeout error"
)

type AppConf struct {
	AppName string
	Server  Server
	DB      DB
	Cache   Cache
	Logger  Logger
}

type DB struct {
	Driver   string
	Name     string
	User     string
	Password string
	Host     string
	Port     string
	MaxConn  int
	Timeout  int
}

type Server struct {
	Port            string
	ShutdownTimeout time.Duration
}

type Cache struct {
	Host string
	Port string
}

type Logger struct {
	Level string
}

func NewAppConf() AppConf {
	port := os.Getenv(serverPort)

	return AppConf{
		AppName: os.Getenv("APP_NAME"),
		Server: Server{
			Port: port,
		},
		DB: DB{
			Driver:   os.Getenv("DB_DRIVER"),
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
		},
		Cache: Cache{
			Host: os.Getenv("REDIS_HOST"),
			Port: os.Getenv("REDIS_PORT"),
		},
	}
}

func (a *AppConf) Init(logger *zap.Logger) {
	shutDownTimeOut, err := strconv.Atoi(os.Getenv(envShutdownTimeout))
	if err != nil {
		logger.Fatal(parseShutdownTimeoutError)
	}
	shutDownTimeout := time.Duration(shutDownTimeOut) * time.Second

	dbTimeout, err := strconv.Atoi(os.Getenv("DB_TIMEOUT"))
	if err != nil {
		logger.Fatal("config: parse db timeout err", zap.Error(err))
	}
	dbMaxConn, err := strconv.Atoi(os.Getenv("MAX_CONN"))
	if err != nil {
		logger.Fatal("config: parse db max connection err", zap.Error(err))
	}
	a.DB.Timeout = dbTimeout
	a.DB.MaxConn = dbMaxConn

	a.Server.ShutdownTimeout = shutDownTimeout

	a.Logger.Level = os.Getenv("LOGGER_LEVEL")
}
