package services

import (
	"finnbank/common/utils"

	"github.com/gin-gonic/gin"
)

type BankcardService struct {
	log *utils.Logger
	url string
}

func newBankcardService(log *utils.Logger) *BankcardService {
	return &BankcardService{
		log: log,
		url: "http://localhost:8083/graphql/bankcard",
	}
}

func (a *BankcardService) GetUserBankcard(*gin.Context) {
}

func (a *BankcardService) GenerateBankcardForUser(*gin.Context) {
}

func (a *BankcardService) RenewBankcardForUser(*gin.Context) {
}
