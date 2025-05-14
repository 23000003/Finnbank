package services

import (
	"bytes"
	"encoding/json"
	t "finnbank/api-gateway/types"
	"finnbank/common/utils"
	"fmt"
	"net/http"
	"strconv"

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

func (a *BankcardService) GetAllBankCardOfUserById(ctx *gin.Context) {
	id := ctx.Param("id")

	query := fmt.Sprintf(`{
		get_all_bankcard(user_id: "%s") {
			bankcard_id
			card_type
			card_number
			expiry_date
			date_created
			cvv
		}
	}`, id)

	qlrequestBody := map[string]interface{}{
		"query": query,
	}

	qlrequestBodyJSON, _ := json.Marshal(qlrequestBody)

	resp, err := http.Post(a.url, "application/json", bytes.NewBuffer(qlrequestBodyJSON))
	if err != nil {
		a.log.Info("Request Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	var data t.GetAllBankCardsGraphQLResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		a.log.Info("Error decoding response: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if data.Errors != nil {
		a.log.Info("GraphQL Errors: %v", data.Errors)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": data.Errors})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data.Data.GetAllBankCard})

}

func (a *BankcardService) UpdateBankcardExpiryDateByUserId(ctx *gin.Context) {
	id := ctx.Param("id")

	query := fmt.Sprintf(`{
		update_bankcard_expiry(bankcard_id: "%s") {
			bankcard_id
			bankcard_type
			bankcard_number
			expiry_date
			date_created
			cvv
		}
	}`, id)

	qlrequestBody := map[string]interface{}{
		"query": query,
	}

	qlrequestBodyJSON, _ := json.Marshal(qlrequestBody)

	resp, err := http.Post(a.url, "application/json", bytes.NewBuffer(qlrequestBodyJSON))
	if err != nil {
		a.log.Info("Request Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	var data t.UpdateBankCardGraphQLResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		a.log.Info("Error decoding response: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if data.Errors != nil {
		a.log.Info("GraphQL Errors: %v", data.Errors)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": data.Errors})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "Expiry date updated successfully"})
}

func (a *BankcardService) UpdateBankcardPinNumberById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		a.log.Info("Invalid ID format: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	newPin := ctx.Param("new-pin")

	query := fmt.Sprintf(`{
		update_pin_number(bankcard_id: %d, pin_number: "%s") {
			bankcard_id
			bankcard_type
			bankcard_number
			expiry_date
			date_created
			cvv
		}
	}`, id, newPin)

	qlrequestBody := map[string]interface{}{
		"query": query,
	}

	qlrequestBodyJSON, _ := json.Marshal(qlrequestBody)

	resp, err := http.Post(a.url, "application/json", bytes.NewBuffer(qlrequestBodyJSON))
	if err != nil {
		a.log.Info("Request Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	var data t.UpdateBankCardGraphQLResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		a.log.Info("Error decoding response: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if data.Errors != nil {
		a.log.Info("GraphQL Errors: %v", data.Errors)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": data.Errors})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "Pin number updated successfully"})
}