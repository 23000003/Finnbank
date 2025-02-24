package main


/**
	TEST SERVICE
**/

import (
	"finnbank/services/common/utils"
	
)


func main() {
	logger, err := utils.NewLogger()
	if err != nil {
		panic(err)
	}
	go RunHttpServer(logger, ":9080")
	RunGrpcServer(logger, ":9000")
}