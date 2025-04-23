package services

import (
	"finnbank/common/utils"
	"fmt"
	"net/http"
	"encoding/json"
	"bytes"
	"github.com/gin-gonic/gin"
	t "finnbank/api-gateway/types"
)

type TransactionService struct {
	log *utils.Logger
	url string
}

func newTransactionService(log *utils.Logger) *TransactionService {
	return &TransactionService{
		log: log,
		url: "http://localhost:8083/graphql/transaction",
	}
}

func (a *TransactionService) GetTransactionByUserId(ctx *gin.Context) {
	id := ctx.Param("id");

	query := fmt.Sprintf(`{
		getTransactionsByUserId(userId: "%s") {
			transaction_id
			ref_no
			sender_id
			receiver_id
			transaction_type
			amount
			transaction_status
			date_transaction
			transaction_fee
			notes
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

	var data t.GetAllTransactionsGraphQLResponse
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

	ctx.JSON(http.StatusOK, gin.H{"data": data.Data.GetTransactionsByUserId})
}

func (a *TransactionService) CreateTransaction(ctx *gin.Context) {
	var req t.CreateTransactionRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// http://localhost:8083/graphql/transaction?query=mutation+_{createTransaction
	// (sender_id:"1",receiver_id:"1", transaction_type:"string", amount:1.99, transaction_fee:1.99, notes: "Thanis")
	// {transaction_id, ref_no, sender_id, receiver_id, transaction_type, amount, transaction_status, date_transaction, transaction_fee, notes}}
	query := fmt.Sprintf(`mutation {
		createTransaction( transaction: { sender_id: "%s",receiver_id: "%s", transaction_type: %s, amount: %f, transaction_fee: %f, notes: "%s" }) {
			transaction_id
			ref_no
			sender_id
			receiver_id
			transaction_type
			amount
			transaction_status
			date_transaction
			transaction_fee
			notes
		}
	}`, req.SenderId, req.ReceiverId, req.TransactionType, req.Amount, req.TransactionFee, req.Notes)

	qlrequestBody := map[string]interface{}{
		"query": query,
	}
	qlrequestBodyJSON, _ := json.Marshal(qlrequestBody)

	resp, err := http.Post(a.url, "application/json", bytes.NewBuffer(qlrequestBodyJSON))
	if err != nil {
		a.log.Info("Request Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	defer resp.Body.Close()

	a.log.Info("Response: %v", resp.Body)

	var data t.CreateTransactionsGraphQLResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		a.log.Info("Error decoding response: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	a.log.Info("%v ======= DATA", data)

	if data.Errors != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": data.Errors})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data.Data.CreateTransaction})
}

