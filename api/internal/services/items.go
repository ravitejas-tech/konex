package services

import (
	"fmt"

	"github.com/pocketbase/pocketbase/core"

	"github.com/ravitejas/konex/api/internal/models"
)

const ItemsCollection = "items"

// ListItems retrieves all items from the "items" collection, ordered by creation date descending.
func ListItems(app core.App) ([]models.ItemResponse, error) {
	records, err := app.FindAllRecords(ItemsCollection, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list items: %w", err)
	}

	items := make([]models.ItemResponse, 0, len(records))
	for _, r := range records {
		items = append(items, models.ItemResponse{
			ID:        r.Id,
			Name:      r.GetString("name"),
			CreatedAt: r.GetString("created"),
			UpdatedAt: r.GetString("updated"),
		})
	}

	return items, nil
}

// CreateItem creates a new record in the "items" collection.
func CreateItem(app core.App, req models.CreateItemRequest) (*models.ItemResponse, error) {
	collection, err := app.FindCollectionByNameOrId(ItemsCollection)
	if err != nil {
		return nil, fmt.Errorf("collection %q not found: %w", ItemsCollection, err)
	}

	record := core.NewRecord(collection)
	record.Set("name", req.Name)

	if err := app.Save(record); err != nil {
		return nil, fmt.Errorf("failed to save item: %w", err)
	}

	return &models.ItemResponse{
		ID:        record.Id,
		Name:      record.GetString("name"),
		CreatedAt: record.GetString("created"),
		UpdatedAt: record.GetString("updated"),
	}, nil
}
