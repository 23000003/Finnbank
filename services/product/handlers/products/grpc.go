package products

import (
	"finnbank/services/common/grpc/products"
	"finnbank/services/product/types"
)

type ProductGrpcHandler struct {
	productService types.ProductService
	products.UnimplementedProductServiceServer
}

func ProductGrpcService() {
	// gRPCHandler := &ProductGrpcHandler{}
}