package store

import "time"

type CreateStoreRequest struct {
	Name string `json:"name" validate:"required"`
}

type StoreResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	UserID    uint      `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type GetAllStoreQuery struct {
	Page    int    `json:"page" validate:"min=1"`
	Limit   int    `json:"limit" validate:"oneof=5 10 20 50"`
	OrderBy string `json:"orderBy" validate:"oneof=created_at updated_at name"`
	Sort    string `json:"sort" validate:"oneof=asc desc"`
	Search  string `json:"search"`
}

type Meta struct {
	Page         int `json:"page"`
	Limit        int `json:"limit"`
	TotalPages   int `json:"totalPages"`
	TotalRecords int `json:"totalRecords"`
}
