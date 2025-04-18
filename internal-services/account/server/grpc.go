package server

import (
	"finnbank/internal-services/account/service"
	"finnbank/common/grpc/auth"
	"finnbank/common/utils"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

// Set up grpc connection

func GrpcServer(s service.AuthService, logger *utils.Logger) error {
	lis, err := net.Listen("tcp", "localhost:9002")
	if err != nil {
		logger.Fatal("Could not start gRPC server on port 9002: %s", err)
		return err
	}
	logger.Info("Port 9002 listening success")
	grpcServer := grpc.NewServer()
	auth.RegisterAuthServiceServer(grpcServer, &s)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		logger.Warn("Shutting down gRPC server...")
		grpcServer.GracefulStop()
		lis.Close()
		os.Exit(0)
	}()
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("Failed to start gRPC server: %s", err)
		return err
	}
	return nil
}
