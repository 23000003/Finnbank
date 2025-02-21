package types

import (
	"context"
	"finnbank/services/common/grpc/products"
)

type Product struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Info  string  `json:"info,omitempty"`
	Price float64 `json:"price"`
}

type ProductService interface {
	GetAllProducts(context.Context, *products.Product) error
	GetByIdProduct(context.Context, *products.Product) error
	CreateProduct(context.Context, *products.Product) error
	UpdateProduct(context.Context, *products.Product) error
	DeleteProduct(context.Context, *products.Product) error
}