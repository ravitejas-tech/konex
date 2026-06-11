package models

// ItemResponse is the DTO returned by custom item endpoints.
type ItemResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created"`
	UpdatedAt string `json:"updated"`
}

// CreateItemRequest is the expected payload for creating a new item.
type CreateItemRequest struct {
	Name string `json:"name"`
}
