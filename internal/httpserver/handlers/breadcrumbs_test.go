package handlers

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"

	"serverigh/internal/httpserver/view"
)

func TestBreadcrumbsRendersActivePath(t *testing.T) {
	root := t.TempDir()
	makeTestDir(t, root, "docs/guides")
	h := testHandlersWithRoot(t, root)

	app := fiber.New()
	app.Get("/partials/breadcrumbs", h.Breadcrumbs)

	resp, body := testRequest(
		t,
		app,
		http.MethodGet,
		"/partials/breadcrumbs?path=/docs/guides",
	)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	for _, part := range []string{"Root", "docs", "guides"} {
		if !strings.Contains(body, part) {
			t.Fatalf("expected breadcrumb %q in response, got %q", part, body)
		}
	}

	if !strings.Contains(body, `aria-current="page"`) {
		t.Fatalf("expected current breadcrumb marker, got %q", body)
	}

	if !strings.Contains(body, `hx-get="/partials/files?path=%2Fdocs"`) {
		t.Fatalf("expected htmx breadcrumb link, got %q", body)
	}

	if !strings.Contains(body, `id="breadcrumb-back"`) {
		t.Fatalf("expected breadcrumb back control, got %q", body)
	}

	if !strings.Contains(body, `href="/browse?path=%2Fdocs"`) {
		t.Fatalf("expected parent fallback link, got %q", body)
	}

	if !strings.Contains(body, `data-history-back`) {
		t.Fatalf("expected browser history hook, got %q", body)
	}
}

func TestBreadcrumbsRejectsTraversal(t *testing.T) {
	h := testHandlers(t)
	app := fiber.New()
	app.Get("/partials/breadcrumbs", h.Breadcrumbs)

	resp, _ := testRequest(
		t,
		app,
		http.MethodGet,
		"/partials/breadcrumbs?path=../outside",
	)

	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", resp.StatusCode)
	}
}

func TestParentPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{name: "root", path: "/", want: "/"},
		{name: "nested", path: "/docs/guides", want: "/docs"},
		{name: "single segment", path: "/docs", want: "/"},
		{name: "unclean", path: "docs//guides", want: "/docs"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := view.ParentPath(tt.path); got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}
