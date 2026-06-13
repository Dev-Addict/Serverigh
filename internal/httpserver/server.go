package httpserver

import (
	"github.com/gofiber/fiber/v2"

	"serverigh/internal/config"
	"serverigh/internal/httpserver/handlers"
)

type Server struct {
	handlers handlers.Handlers
}

func New(cfg config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:               "Serverigh",
		DisableStartupMessage: true,
	})

	server := Server{
		handlers: handlers.New(cfg),
	}

	server.registerRoutes(app)

	return app
}
