package handlers

import (
	"net/http"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// RegisterRoutes binds all custom API endpoints to the PocketBase router.
// Routes are namespaced under /api/custom/ to avoid conflicts with
// PocketBase's built-in /api/ endpoints.
func RegisterRoutes(se *core.ServeEvent) {
	// --- Public routes (no auth required) ---
	se.Router.GET("/api/custom/health", handleHealth)

	// --- Auth-protected routes ---
	protected := se.Router.Group("/api/custom")
	protected.Bind(apis.RequireAuth())

	protected.GET("/items", handleListItems)
	protected.POST("/items", handleCreateItem)
}

// handleHealth returns a simple health check response.
func handleHealth(e *core.RequestEvent) error {
	return e.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}
