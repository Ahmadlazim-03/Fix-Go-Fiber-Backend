package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetupMiddleware(app *fiber.App) {
	// Panic recovery middleware
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization,X-Requested-With",
	}))

	// Logger middleware
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	// Custom logging middleware for requests
	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		
		// Process request
		err := c.Next()
		
		// Log request details
		duration := time.Since(start)
		log.Printf("[%s] %s %s - %d - %v\n", 
			c.Method(), 
			c.Path(), 
			c.IP(), 
			c.Response().StatusCode(), 
			duration,
		)
		
		return err
	})
}