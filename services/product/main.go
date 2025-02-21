package main


/**
	TEST SERVICE
**/

import (
	"finnbank/services/product/utils"
	
)


func main() {
	logger, err := utils.NewLogger()
	if err != nil {
		panic(err)
	}
	RunHttpServer(logger)
	RunGrpcServer(logger, ":9000")
}