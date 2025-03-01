package types

import (
	"context"
	"finnbank/services/common/grpc/products"
)

type Product struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Info  string  `json:"info,omitempty"`
	Price float32 `json:"price"`
}

type IProductService interface {
	GetAllProducts(context.Context) []*products.Product
	GetByIdProduct(context.Context, int64) *products.Product
	CreateProduct(context.Context, *products.Product) (*products.Product, error)
	UpdateProduct(context.Context, *products.Product) (*products.Product, error)
	DeleteProduct(context.Context, *products.Product) (string, error)
}