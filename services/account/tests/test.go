package main

import (
	"context"

	"finnbank/services/common/grpc/account"
	"finnbank/services/common/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// DEPRECATED, DIDN'T KNOW I CAN USE POSTMAN

func main() {
	// Connect to the gRPC server
	logger, _ := utils.NewLogger()
	conn, err := grpc.NewClient("localhost:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("Could not connect to gRPC server: %s", err)
		return
	}
	defer conn.Close()

	client := account.NewAccountServiceClient(conn)

	logger.Info("Test 1: AddAccount")
	res, err := TestAddAccount(client, logger)
	if err != nil {
		logger.Error("Test 1 not passed: %v", err)
		return
	}
	logger.Debug("Test 1 Passed: %v", res)
}

// Test function for AddAccount
func TestAddAccount(c account.AccountServiceClient, logger *utils.Logger) (*account.AddAccountResponse, error) {
	request := &account.AddAccountRequest{
		Email:       "testuse7@example.com",
		FullName:    "Test User",
		PhoneNumber: "09399883951",
		Password:    "SecurePassword123",
		Address:     "123 Test Street, Test City",
		AccountType: "Personal",
	}

	res, err := c.AddAccount(context.Background(), request)
	if err != nil {
		logger.Error("Error when calling AddAccount: %s", err)
	}
	return res, nil
}
