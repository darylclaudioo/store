package entity

type Products struct {
	ID    string
	Name  string
	Stock int
}

type GetProductsRequest struct {
	Page  int
	Limit int
	Name  string
}

type UpdateProductsRequest struct {
	Name  string
	Stock *int
}