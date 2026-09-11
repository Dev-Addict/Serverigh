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
	writeHandlers := s.handlers.Write()
	app.Get("/partials/folders", writeHandlers.Folders)
	app.Get("/preview", s.handlers.Preview)
	app.Get("/download", s.handlers.Download)
	app.Get("/raw", s.handlers.Raw)
	app.Post("/download/bulk", s.handlers.BulkDownload)

	app.Post("/settings", s.handlers.Settings)

	app.Post("/actions/mkdir", writeHandlers.Mkdir)
	app.Post("/actions/file", writeHandlers.CreateFile)
	app.Post("/actions/rename", writeHandlers.Rename)
	app.Post("/actions/move", writeHandlers.Move)
	app.Post("/actions/copy", writeHandlers.Copy)
	app.Post("/actions/duplicate", writeHandlers.Duplicate)
	app.Post("/actions/delete", writeHandlers.Delete)
	app.Post("/actions/trash/restore", writeHandlers.RestoreTrash)
	app.Post("/actions/trash/delete", writeHandlers.DeleteTrash)
	app.Post("/actions/bulk/copy", s.handlers.BulkCopy)
	app.Post("/actions/bulk/delete", s.handlers.BulkDelete)
	app.Post("/actions/bulk/duplicate", s.handlers.BulkDuplicate)
	app.Post("/actions/bulk/move", s.handlers.BulkMove)
	app.Post("/actions/upload", writeHandlers.Upload)

	app.Use(s.handlers.NotFound)
}
