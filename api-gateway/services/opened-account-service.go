package services

import (
	"bytes"
	"encoding/json"
	"finnbank/common/utils"
	"fmt"
	"net/http"
	"strconv"
	t "finnbank/api-gateway/types"
	"github.com/gin-gonic/gin"
)

type OpenedAccountService struct {
	log *utils.Logger
	url string
}

func newOpenedAccountService(log *utils.Logger) *OpenedAccountService {
	return &OpenedAccountService{
		log: log,
		url: "http://localhost:8083/graphql/opened-account",
	}
}

func (a *OpenedAccountService) GetAllOpenedAccountsByUserId(ctx *gin.Context) {
	
	id := ctx.Param("id");
	a.log.Info("GetAllOpenedAccountsByUserId: %s", id)

	query := fmt.Sprintf(`{
		get_all(account_id: "%s") {
			openedaccount_id
			bankcard_id
			balance
			account_type
			date_created
			openedaccount_status
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

	var data t.GetAllOpenedAccountsGraphQLResponse
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

	a.log.Info("Response data: %+v", data.Data.GetAll)
	ctx.JSON(http.StatusOK, gin.H{"data": data.Data.GetAll})
}


func (a *OpenedAccountService) GetOpenedAccountOfUserById(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account id"})
		return
	}

	a.log.Info("GetOpenedAccount: %v", id)

	query := fmt.Sprintf(`{
		get_by_id(openedaccount_id: %d) {
			openedaccount_id
			bankcard_id
			balance
			account_type
			date_created
			openedaccount_status
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

	var data t.GetOpenedAccountsGraphQLResponse
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

	ctx.JSON(http.StatusOK, gin.H{"data": data.Data.GetById})
}


func (a *OpenedAccountService) OpenAnAccountByAccountType(ctx *gin.Context) {

	var req t.CreateOpenAccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a.log.Info("OpenAnAccountByAccountType: %s %s %f", req.AccountId, req.AccountType, req.Balance)

	// http://localhost:8083/graphql/opened-account?query=mutation+_{create_account(account_id:1,account_type:"string",balance:1.99){<entities>}}
	query := fmt.Sprintf(`mutation {
		create_account(account_id: "%s", account_type: "%s", balance: %f) {
			openedaccount_id
			bankcard_id
			balance
			account_type
			date_created
			openedaccount_status
		}
	}`, req.AccountId, req.AccountType, req.Balance)

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

	var data t.CreateOpenedAccountsGraphQLResponse
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

	ctx.JSON(http.StatusOK, gin.H{"data": data.Data.CreateAccount})

}

func (a *OpenedAccountService) UpdateOpenedAccountStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	status := ctx.Param("status")

	openedAccountId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	a.log.Info("Updating Opened Account ID: %d to status: %s", openedAccountId, status)

	query := fmt.Sprintf(
		`mutation {
			update_account_status(openedaccount_id: %d, openedaccount_status: "%s")
	}`, openedAccountId, status)

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

	var data t.UpdateOpenedAccountsGraphQLResponse
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

	ctx.JSON(http.StatusOK, gin.H{"data": data.Data.UpdateAccountStatus})
}
