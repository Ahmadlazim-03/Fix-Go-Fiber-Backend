package logger

import (
	"os"
	"time"

	"Fix-Go-Fiber-Backend/pkg/config"

	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sirupsen/logrus"
)

func NewLogrus(cfg *config.Config) *logrus.Logger {
	log := logrus.New()
	
	// Set log level
	if cfg.App.Debug {
		log.SetLevel(logrus.DebugLevel)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}
	
	// Set formatter
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	
	// Set output
	log.SetOutput(os.Stdout)
	
	return log
}

func NewFiberLogger() logger.Config {
	return logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Asia/Jakarta",
	}
}