package main

/**
	TEST SERVICE
**/

import (
	"finnbank/common/utils"
	handler "finnbank/internal-services/product/handlers/products"
	prodService "finnbank/internal-services/product/services"
	"net"

	"google.golang.org/grpc"
)

func RunGrpcServer(l *utils.Logger, address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		l.Error("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Register services/handlers here
	productService := prodService.ProductServiceInstance()
	handler.ConfigureProductGrpcServices(grpcServer, productService, l)

	l.Info("gRPC server running on port: %s", address)

	return grpcServer.Serve(lis)
}
