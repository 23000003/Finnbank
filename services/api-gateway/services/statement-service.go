package services

import (
	"finnbank/services/common/utils"
	"github.com/gin-gonic/gin"
)

type StatementService struct {
	log *utils.Logger
	url string
}

func NewStatementService(log *utils.Logger) *StatementService {
	return &StatementService{
		log: log,
		url: "http://localhost:8083/graphql/statement",
	}
}

func (a *StatementService) GenerateStatement(*gin.Context){
}
