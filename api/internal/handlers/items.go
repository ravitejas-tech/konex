package handlers

import (
	"net/http"

	"github.com/pocketbase/pocketbase/core"

	"github.com/ravitejas/konex/api/internal/models"
	"github.com/ravitejas/konex/api/internal/services"
)

// handleListItems returns all items from the "items" collection.
func handleListItems(e *core.RequestEvent) error {
	items, err := services.ListItems(e.App)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return e.JSON(http.StatusOK, items)
}

// handleCreateItem creates a new item record.
func handleCreateItem(e *core.RequestEvent) error {
	var req models.CreateItemRequest
	if err := e.BindBody(&req); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	if req.Name == "" {
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": "name is required",
		})
	}

	item, err := services.CreateItem(e.App, req)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return e.JSON(http.StatusCreated, item)
}
