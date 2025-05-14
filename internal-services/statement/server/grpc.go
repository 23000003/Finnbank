package server

import (
	"finnbank/common/grpc/statement"
	"finnbank/common/utils"
	"finnbank/internal-services/statement/service"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

// START THE GRPC SERVER HERE PLEAESE AND RETURN THE CLIENT
func StartGrpcServer(s service.StatementService, logger *utils.Logger) error {
	// START THE GRPC SERVER HERE PLEAESE AND RETURN THE CLIENT
	lis, err := net.Listen("tcp", "localhost:9004")
	if err != nil {
		logger.Fatal("Could not start gRPC server on port 9004: %s", err)
		return err
	}
	logger.Info("Port 9004 listening success")
	grpcServer := grpc.NewServer()
	statement.RegisterStatementServiceServer(grpcServer, &s)
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
