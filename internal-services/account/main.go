package main

import (
	"context"
	pb "finnbank/common/grpc/auth"
	"finnbank/common/utils"
	"finnbank/internal-services/account/auth"
	"finnbank/internal-services/account/db"
	"finnbank/internal-services/account/server"
	"finnbank/internal-services/account/service"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	logger, err := utils.NewLogger()
	if err != nil {
		panic(err)
	}
	db, err := db.ConnectToDb(context.Background())
	accountService := service.AuthService{
		Logger: logger,
		DB:     db,
		Helper: &auth.AuthHelpers{
			DB: db,
		},
		UnimplementedAuthServiceServer: pb.UnimplementedAuthServiceServer{},
	}
	if err != nil {
		logger.Fatal("Failed to connect to database")
		return
	}
	go func() {
		if err := server.GrpcServer(accountService, logger); err != nil {
			logger.Fatal("Failed to start gRPC server")
			return
		}
	}()
	logger.Info("Starting the server...")
	logger.Info("Server running on localhost:8082")
	wg.Wait()
}
