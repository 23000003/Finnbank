package main

import (
	pb "finnbank/common/grpc/auth"
	"finnbank/common/utils"
	"finnbank/internal-services/account/server"
	"finnbank/internal-services/account/service"
	"sync"

	"github.com/joho/godotenv"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	godotenv.Load()
	logger, err := utils.NewLogger()
	if err != nil {
		panic(err)
	}
	accountService := service.AuthService{
		Logger:                         logger,
		UnimplementedAuthServiceServer: pb.UnimplementedAuthServiceServer{},
	}
	go func() {
		if err := server.GrpcServer(accountService, logger); err != nil {
			logger.Fatal("Failed to start gRPC server")
			return
		}
	}()
	logger.Info("Starting the server...")
	logger.Info("Server running on localhost:9002")
	wg.Wait()
}
