package product

import "time"

type ProductRequest struct {
	Name  string `json:"name" validate:"required"`
	Price int    `json:"price" validate:"required"`
}

type ProductResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	StoreID   uint      `json:"storeId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
