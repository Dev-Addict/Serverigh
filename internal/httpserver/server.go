package httpserver

import (
	"github.com/gofiber/fiber/v2"

	"serverigh/internal/apperror"
	"serverigh/internal/config"
	"serverigh/internal/httpserver/handlers"
)

type Server struct {
	handlers handlers.Handlers
}

func New(cfg config.Config) (*fiber.App, error) {
	app := fiber.New(fiber.Config{
		AppName:               "Serverigh",
		DisableStartupMessage: true,
	})

	routeHandlers, err := handlers.New(cfg)
	if err != nil {
		return nil, apperror.WrapOperation(
			apperror.CodeServer,
			"server error",
			"initialize handlers",
			err,
		)
	}

	server := Server{
		handlers: routeHandlers,
	}

	registerMiddleware(app)
	server.registerRoutes(app)

	return app, nil
}
