package main

/**
	TEST SERVICE
**/

import (
	"net"
	"finnbank/services/product/utils"
	"google.golang.org/grpc"
)

func RunGrpcServer(l *utils.Logger, address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		l.Error("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	l.Info("gRPC server running on %s", address)

	return grpcServer.Serve(lis)
}
