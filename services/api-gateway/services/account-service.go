package services

import (
	"finnbank/services/common/utils"
	"github.com/gin-gonic/gin"
)

type AccountService struct {
	log *utils.Logger
	url string
}

func NewAccountService(log *utils.Logger) *AccountService {
	return &AccountService{
		log: log,
		url: "http://localhost:8083/graphql/account",
	}
}

func (a *AccountService) LoginUser(*gin.Context) {
}

func (a *AccountService) SignupUser(*gin.Context) {
}