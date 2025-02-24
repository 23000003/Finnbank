package products

import (
	"context"
	"finnbank/services/common/grpc/products"
	"finnbank/services/common/utils"
	"finnbank/services/product/types"
	"math/rand"
	"google.golang.org/grpc"
)

type StructProductGrpcHandler struct {
	productService types.IProductService
	products.UnimplementedProductServiceServer
	logger *utils.Logger
}

func ConfigureProductGrpcServices(grpc *grpc.Server, productService types.IProductService, logger *utils.Logger) {
	gRPCHandler := &StructProductGrpcHandler{
		productService: productService,
		logger: logger,
	}
	products.RegisterProductServiceServer(grpc, gRPCHandler)
}

func (h *StructProductGrpcHandler) GetAllProducts(ctx context.Context, req *products.GetAllProductsRequest) (*products.GetProductResponse, error) {
	
	h.logger.Info("GetAllProducts gRPC request %v", req)
	h.logger.Info("GetAllProducts gRPC Context %v", ctx)
	data := h.productService.GetAllProducts(ctx)
	res := &products.GetProductResponse{Product: data}
	return res, nil
}

func (h *StructProductGrpcHandler) GetByIdProduct(ctx context.Context, req *products.ProductByIdRequest) (*products.GetSingleProductResponse, error) {
	id := req.ID;
	data := h.productService.GetByIdProduct(ctx, id)
	res := &products.GetSingleProductResponse{Product: data}
	return res, nil
}

func (h *StructProductGrpcHandler) CreateProduct(ctx context.Context, req *products.CreateProductRequest) (*products.GetSingleProductResponse, error) {
	product := &products.Product{
		ID:    int64(rand.Intn(100000)), // change to get latest ID
		Name:  req.Name,
		Info:  req.Info,
		Price: req.Price,
	}
	data, err := h.productService.CreateProduct(ctx, product)
	res := &products.GetSingleProductResponse{Product: data}
	return res, err
}

func (h *StructProductGrpcHandler) UpdateProduct(ctx context.Context, req *products.ProductByIdRequest) (*products.GetSingleProductResponse, error) {
	// well just call its service and update :D
	// just the same concept as create prod 
	return nil, nil
}

func (h *StructProductGrpcHandler) DeleteProduct(ctx context.Context, req *products.ProductByIdRequest) (*products.GetSingleProductResponse, error) {
	// well just call its service and delete :D 
	// just the same concept as create prod 
	return nil, nil
}