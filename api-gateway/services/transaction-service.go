package services

import (
	"finnbank/common/utils"

	"github.com/gin-gonic/gin"
)

type TransactionService struct {
	log *utils.Logger
	url string
}

func NewTransactionService(log *utils.Logger) *TransactionService {
	return &TransactionService{
		log: log,
		url: "http://localhost:8083/graphql/transaction",
	}
}

func (a *TransactionService) GetAllTransaction(*gin.Context) {
}

func (a *TransactionService) GetTransaction(*gin.Context) {
}

func (a *TransactionService) GenerateTransaction(*gin.Context) {
}
