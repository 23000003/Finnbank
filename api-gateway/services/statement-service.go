package services

import (
	"finnbank/common/utils"

	"github.com/gin-gonic/gin"
)

type StatementService struct {
	log *utils.Logger
	url string
}

func newStatementService(log *utils.Logger) *StatementService {
	return &StatementService{
		log: log,
		url: "http://localhost:8083/graphql/statement",
	}
}

func (a *StatementService) GenerateStatement(*gin.Context) {
}
