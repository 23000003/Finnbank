package services

import (
	"context"
	"finnbank/services/common/grpc/products"
	"fmt"
)

type ProductService struct {}

func ProductServiceInstance() *ProductService {
	return &ProductService{}
}

var productsData = []*products.Product{
	{
		ID:    1,
		Name:  "Chicha Morada",
		Info:  "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
		Price: 7.99,
	},
	{
		ID:    2,
		Name:  "Chicha de jora",
		Info:  "Chicha de jora is a corn beer chicha prepared by germinating maize, extracting the malt sugars, boiling the wort, and fermenting it in large vessels (traditionally huge earthenware vats) for several days (wiki)",
		Price: 5.95,
	},
	{
		ID:    3,
		Name:  "Pisco",
		Info:  "Pisco is a colorless or yellowish-to-amber colored brandy produced in winemaking regions of Peru and Chile (wiki)",
		Price: 9.95,
	},
}

func (p *ProductService) GetAllProducts(ctx context.Context) []*products.Product {
	
	protoProducts := make([]*products.Product, len(productsData))
	for i, p := range productsData {
		protoProducts[i] = &products.Product{
			ID:    p.ID,
			Name:  p.Name,
			Info:  p.Info,
			Price: p.Price,
		}
	}

	return protoProducts
}

func (p *ProductService) GetByIdProduct(ctx context.Context, ID int64) *products.Product {
	for _, product := range productsData {
		if int64(product.ID) == ID {
			return &products.Product{
				ID:    product.ID,
				Name:  product.Name,
				Info:  product.Info,
				Price: product.Price,
			}
		}
	}
	return nil
}

func (p *ProductService) CreateProduct(ctx context.Context, prod *products.Product) (*products.Product, error) {
	// check if product with same name already exists
	for _, p := range productsData {
		if p.Name == prod.Name {
			return nil, fmt.Errorf("product with name '%s' already exists", prod.Name)
		}
	}
	productsData = append(productsData, prod)
	return prod, nil
}

func (p *ProductService) UpdateProduct(ctx context.Context, prod *products.Product) (*products.Product, error) {
	// same concent as create prod but with update and a little adjustment
	return nil, nil
}

func (p *ProductService) DeleteProduct(ctx context.Context, prod *products.Product) (string, error) {
	// just delete the product if it exists if not return an error same as products
	return "delete successsful", nil
}