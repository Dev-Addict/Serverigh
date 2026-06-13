package httpserver

import "github.com/gofiber/fiber/v2"

func (s Server) registerRoutes(app *fiber.App) {
	app.Static("/static", "./web/static")

	app.Get("/", s.handlers.Root)

	app.Get("/browse", s.handlers.Browse)
	app.Get("/healthz", s.handlers.Health)

	app.Get("/partials/files", s.handlers.NotImplemented)
	app.Get("/partials/breadcrumbs", s.handlers.NotImplemented)
	app.Get("/preview", s.handlers.NotImplemented)
	app.Get("/download", s.handlers.NotImplemented)
	app.Get("/raw", s.handlers.NotImplemented)

	app.Post("/actions/mkdir", s.handlers.NotImplemented)
	app.Post("/actions/rename", s.handlers.NotImplemented)
	app.Post("/actions/move", s.handlers.NotImplemented)
	app.Post("/actions/copy", s.handlers.NotImplemented)
	app.Post("/actions/delete", s.handlers.NotImplemented)
	app.Post("/actions/upload", s.handlers.NotImplemented)

	app.Use(s.handlers.NotFound)
}
