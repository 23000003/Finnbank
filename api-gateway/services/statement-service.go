package services

import (
	"encoding/json"
	"finnbank/common/utils"
	"fmt"
	"net/http"
	"strconv"
	"bytes"
	"github.com/gin-gonic/gin"
	t "finnbank/api-gateway/types"
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

func (a *StatementService) GenerateStatement(ctx *gin.Context) {
	creditId, err := strconv.Atoi(ctx.Query("credit"))
	debitId, err1 := strconv.Atoi(ctx.Query("debit"))
	savingsId, err2 := strconv.Atoi(ctx.Query("savings"))
	start := ctx.Query("start_date")
	end := ctx.Query("end_date")

	if err != nil || err1 != nil || err2 != nil {
		a.log.Info("Error converting ID to int: %v, %v, %v", err, err1, err2)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// http://localhost:8083/graphql/statement?query={generate_statement
	// (credit_id:1, debit_id:2, savings_id:3, start_date:"HEY", end_date:"HEY"){pdf_buffer}}
	query := fmt.Sprintf(`{
		generate_statement(credit_id: %d, debit_id: %d, savings_id: %d, start_date: "%s", end_date: "%s") {
			pdf_buffer
		}
	}`, creditId, debitId, savingsId, start, end)

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

	var data t.GetStatementGraphQLResponse
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

	ctx.JSON(http.StatusOK, gin.H{"data": data.Data.GenerateStatement})
}
