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

func (a *TransactionService) GetTransactionByOpenAccountId(ctx *gin.Context) {
	creditId, err := strconv.Atoi(ctx.Query("credit"))
	debitId, err1 := strconv.Atoi(ctx.Query("debit"))
	savingsId, err2 := strconv.Atoi(ctx.Query("savings"))
	limit, err3 := strconv.Atoi(ctx.Query("limit")) 

	if err != nil || err1 != nil || err2 != nil || err3 != nil {
		a.log.Info("Error converting ID to int: %v, %v, %v, %v", err, err1, err2, err3)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	query := fmt.Sprintf(`{
		getTransactionsByUserId(creditId: %d, debitId: %d, savingsId: %d, limit: %d) {
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
	}`, creditId, debitId, savingsId, limit)

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


func (a *TransactionService) GetRecentlySentByOpenAccountId(ctx *gin.Context) {
	creditId, err := strconv.Atoi(ctx.Query("credit"))
	debitId, err1 := strconv.Atoi(ctx.Query("debit"))
	savingsId, err2 := strconv.Atoi(ctx.Query("savings"))

	if err != nil || err1 != nil || err2 != nil {
		a.log.Info("Error converting ID to int: %v, %v, %v", err, err1, err2)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	query := fmt.Sprintf(`{
		getRecentlySent(creditId: %d, debitId: %d, savingsId: %d) {
			transaction_id
			sender_id
			receiver_id
		}
	}`, creditId, debitId, savingsId)

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

	var data t.GetRecentlySentGraphQLResponse
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

	ctx.JSON(http.StatusOK, gin.H{"data": data.Data.GetRecentlySent})
}

func (a *TransactionService) GetIsAccountAtLimit(ctx *gin.Context) {
	creditId, err := strconv.Atoi(ctx.Query("credit"))
	debitId, err1 := strconv.Atoi(ctx.Query("debit"))
	savingsId, err2 := strconv.Atoi(ctx.Query("savings"))
	accountType := ctx.Query("accountType")

	if err != nil || err1 != nil || err2 != nil {
		a.log.Info("Error converting ID to int: %v, %v, %v", err, err1, err2)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	query := fmt.Sprintf(`{
		getIsAccountAtLimit(creditId: %d, debitId: %d, savingsId: %d, account_type: "%s") 
	}`, creditId, debitId, savingsId, accountType)

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

	var data t.GetIsAccountAtLimitGraphQLResponse
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

	ctx.JSON(http.StatusOK, gin.H{"data": data.Data.GetIsAccountAtLimit})
}

func (a *TransactionService) GetTransactionByTimestamp(ctx *gin.Context) {
	creditId, err := strconv.Atoi(ctx.Query("credit"))
	debitId, err1 := strconv.Atoi(ctx.Query("debit"))
	savingsId, err2 := strconv.Atoi(ctx.Query("savings"))
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate") 

	if err != nil || err1 != nil || err2 != nil {
		a.log.Info("Error converting ID to int: %v, %v, %v", err, err1, err2 )
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	query := fmt.Sprintf(`{
		getTransactionsByTimeStampByUserId(creditId: %d, debitId: %d, savingsId: %d, startTime: "%s", endTime: "%s") {
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
	}`, creditId, debitId, savingsId, startDate, endDate)

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

	var data t.GetAllTransactionsByTimeStampGraphQLResponse
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

	ctx.JSON(http.StatusOK, gin.H{"data": data.Data.GetTransactionsByTimeStampByUserId})
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
		createTransaction( transaction: { sender_id: %d, receiver_id: %d, transaction_type: %s, amount: %f, transaction_fee: %f, notes: "%s" }) {
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

