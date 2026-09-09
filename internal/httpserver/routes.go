package httpserver

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	fiberfs "github.com/gofiber/fiber/v2/middleware/filesystem"

	"serverigh/web"
)

func (s Server) registerRoutes(app *fiber.App) {
	app.Use("/static", fiberfs.New(fiberfs.Config{
		Root:       http.FS(web.Static),
		PathPrefix: "static",
	}))

	app.Get("/", s.handlers.Root)

	app.Get("/browse", s.handlers.Browse)
	app.Get("/healthz", s.handlers.Health)

	app.Get("/partials/files", s.handlers.Files)
	app.Get("/partials/breadcrumbs", s.handlers.Breadcrumbs)
	app.Get("/partials/search", s.handlers.Search)
	app.Get("/partials/folders", s.handlers.Folders)
	app.Get("/preview", s.handlers.Preview)
	app.Get("/download", s.handlers.Download)
	app.Get("/raw", s.handlers.Raw)

	app.Post("/actions/mkdir", s.handlers.Mkdir)
	app.Post("/actions/file", s.handlers.CreateFile)
	app.Post("/actions/rename", s.handlers.Rename)
	app.Post("/actions/move", s.handlers.Move)
	app.Post("/actions/copy", s.handlers.Copy)
	app.Post("/actions/duplicate", s.handlers.Duplicate)
	app.Post("/actions/delete", s.handlers.Delete)
	app.Post("/actions/upload", s.handlers.Upload)

	app.Use(s.handlers.NotFound)
}
