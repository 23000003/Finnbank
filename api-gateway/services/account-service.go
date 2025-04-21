package services

import (
	"finnbank/common/utils"
	"net/http"
	"bytes"
	"fmt"
	"encoding/json"
	t "finnbank/api-gateway/types"
	"github.com/gin-gonic/gin"
)

// TO REFACTOR SCHEMA BASE ON README

// Account (
// 	Account_ID PK,
// 	Email VARCHAR(255) NOT NULL UNIQUE,
// 	First_Name VARCHAR(100),
// 	Last_Name VARCHAR(100),
// 	Surname VARCHAR(100),
// 	Phone_Number VARCHAR(20),
// 	Password VARCHAR(255) NOT NULL,   -- Store hashed passwords, not plain text
// 	Address TEXT,
// 	nationalIdNumber VARCHAR(20) UNIQUE,  -- Unique national ID number for each user
// 	Account_Number VARCHAR(20) UNIQUE,  -- Unique account number for each user
// 	Account_Type ENUM('Business', 'Personal'),
// 	Account_Status ENUM('Active', 'Suspended', 'Closed'),
//       Birthdate DATE (NEWW) add ??
// 	Date_Created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
// ) 


type AccountService struct {
	log *utils.Logger
	url string
}

func newAccountService(log *utils.Logger) *AccountService {
	return &AccountService{
		log: log,
		url: "http://localhost:8083/graphql/account",
	}
}

func (a *AccountService) LoginUser(ctx *gin.Context) {

	var req t.LoginAccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a.log.Info("LoginUser: %v", req)

	// http://localhost:8083/graphql/account?query=mutation+_{login(account: { email: "", password: "" })
	// {access_token, email, auth_id}}
	query := fmt.Sprintf(`mutation {
		login(account: { email: "%s", password: "%s" } ) {
			access_token
			full_name
			account_id
		}
	}`, req.Email, req.Password)

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

	var data t.AccountLoginGraphQLResponse
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

	ctx.JSON(http.StatusOK, gin.H{"data": data.Data.Login})

}

func (a *AccountService) SignupUser(ctx *gin.Context) {

	var req t.SignupAccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// http://localhost:8083/graphql/account?query=mutation+_{create_account(account:
	// { email: "", password: "", first_name: "", last_name: """, phone_number: "", address: "", account_type: "", nationality: "})
	// {email, auth_id, full_name}}
	query := fmt.Sprintf(`mutation {
		create_account( account : { email: "%s", password: "%s", first_name: "%s", last_name: "%s", phone_number: "%s", address: "%s", account_type: "%s", nationality: "%s" } ) {
			email
			auth_id
		}
	}`, req.Email, req.Password, req.FirstName, req.LastName, req.PhoneNumber, req.Address, req.AccountType, req.Nationality)

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

	var data t.AccountSignUpGraphQLResponse
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

	ctx.JSON(http.StatusOK, gin.H{"data": "Registered successfully"})
}


func (a *AccountService) GetAccountNumberById(ctx *gin.Context) {

	id := ctx.Param("id");

	a.log.Info("GetOpenedAccount: %v", id)

	// http://localhost:8083/graphql/account?query={account_by_id(id:"38fba771-37f4-49c8-b5b2-634dfc3871da"){account_id, email, auth_id}}
	// To add : account status, nationalIdNumber
	query := fmt.Sprintf(`{
		account_by_id(id: "%s") {
			email
			full_name
			phone_number
			date_created
			account_number
			address
			nationality,
			account_type
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

	var data t.GetAccountDetailsGraphQLResponse
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

	// Mock since this not in account schema
	data.Data.AccountById.AccountStatus = "Active"
	data.Data.AccountById.NationalId = "6342123456789"

	ctx.JSON(http.StatusOK, gin.H{"data": data.Data.AccountById})
}
