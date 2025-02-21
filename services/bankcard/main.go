package main

import (
	"finnbank/services/common/utils"
)

func main() {
	logger, err := utils.NewLogger()
	if err != nil {
		panic(err)
	}
	logger.Info("Starting the application...")

	RunHttpServer(logger)
}
