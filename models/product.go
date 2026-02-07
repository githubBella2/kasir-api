package models
type Product struct {
	ID    int    `json:"id" validate:"min:10"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}