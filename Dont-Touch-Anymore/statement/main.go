package main

import (
	"finnbank/common/grpc/statement"
	"finnbank/common/utils"
	"finnbank/internal-services/statement/server"
	"finnbank/internal-services/statement/service"
)

func main() {
	logger, err := utils.NewLogger()
	if err != nil {
		panic(err)
	}
	logger.Info("Starting the application...")

	statementService := service.StatementService{
		Logger:                              logger,
		UnimplementedStatementServiceServer: statement.UnimplementedStatementServiceServer{},
	}

	logger.Info("Starting the server...")
	logger.Info("Server running on localhost:8084")
	err = server.StartGrpcServer(statementService, logger)
	if err != nil {
		logger.Fatal("Failed to start gRPC server")
	}
}
