package middleware

import (
	pkgLogger "Fix-Go-Fiber-Backend/pkg/logger"

	"github.com/gofiber/fiber/v2/middleware/logger"
)

func NewLoggerMiddleware() logger.Config {
	return pkgLogger.NewFiberLogger()
}