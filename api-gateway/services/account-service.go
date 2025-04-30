package services

import (
	"bytes"
	"encoding/json"
	t "finnbank/api-gateway/types"
	"finnbank/common/utils"
	"fmt"
	"io"
	"net/http"

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
			account_id
		}
	}`, req.Email, req.Password, req.FirstName, req.LastName, req.PhoneNumber, req.Address, req.AccountType, req.Nationality)

	qlrequestBody := map[string]any{
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

	ctx.JSON(http.StatusOK, gin.H{"message": "Registered successfully", "data": data.Data.CreateAccount})
}

func (a *AccountService) GetUserAccountById(ctx *gin.Context) {

	id := ctx.Param("id")

	a.log.Info("GetOpenedAccount: %v", id)

	// http://localhost:8083/graphql/account?query={account_by_id(id:"38fba771-37f4-49c8-b5b2-634dfc3871da"){account_id, email, auth_id}}
	// To add : account status, nationalIdNumber
	query := fmt.Sprintf(`{
		account_by_id(id: "%s") {
			account_id
			account_status
			account_type
			address
			auth_id
			birthdate
			date_created
			date_updated
			email
			first_name
			has_card
			last_name
			middle_name
			national_id
			nationality
			phone_number
		}
	}`, id)

	qlrequestBody := map[string]any{
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

func (a *AccountService) GetUserAccountByAccountNumber(c *gin.Context) {
	acc_num := c.Param("account_number")

	query := fmt.Sprintf(`{
  account_by_account_num(
    account_number: "%s"
  )
  {
    account_id
    account_status
    account_type
    address
    auth_id
    balance
    birthdate
    date_created
    date_updated
    email
    first_name
    has_card
    last_name
    middle_name
    national_id
    nationality
    phone_number
  }
}`, acc_num)

	qlrequestBody := map[string]any{
		"query": query,
	}
	qlrequestBodyJSON, _ := json.Marshal(qlrequestBody)

	resp, err := http.Post(a.url, "application/json", bytes.NewBuffer(qlrequestBodyJSON))
	if err != nil {
		a.log.Info("Request Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		a.log.Info("Read Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	var res t.AccountNumberFetchResponse

	if err := json.Unmarshal(body, &res); err != nil {
		a.log.Info("Unmarshal Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}
	if res.Errors != nil {
		a.log.Info("GraphQL Error: %v", res.Errors)
		c.JSON(http.StatusBadRequest, gin.H{"error": res.Errors})
		return
	}
	fmt.Println("RAW GQL Response:", string(body))

	c.JSON(http.StatusOK, gin.H{"account": res.Data.FetchedAccount})
}
func (a *AccountService) GetUserAccountByEmail(c *gin.Context) {
	email := c.Param("email")

	query := fmt.Sprintf(`{
		account_by_email(
		  email: "%s"
		)
		{
		  account_id
		  account_status
		  account_type
		  address
		  auth_id
		  balance
		  birthdate
		  date_created
		  date_updated
		  email
		  first_name
		  has_card
		  last_name
		  middle_name
		  national_id
		  nationality
		  phone_number
		}
	  }`, email)

	qlrequestBody := map[string]any{
		"query": query,
	}
	qlrequestBodyJSON, _ := json.Marshal(qlrequestBody)

	resp, err := http.Post(a.url, "application/json", bytes.NewBuffer(qlrequestBodyJSON))
	if err != nil {
		a.log.Info("Request Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		a.log.Info("Read Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	var res t.EmailFetchResponse

	if err := json.Unmarshal(body, &res); err != nil {
		a.log.Info("Unmarshal Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}
	if res.Errors != nil {
		a.log.Info("GraphQL Error: %v", res.Errors)
		c.JSON(http.StatusBadRequest, gin.H{"error": res.Errors})
		return
	}
	c.JSON(http.StatusOK, gin.H{"account": res.Data.FetchedAccount})

}
func (a *AccountService) UpdateUserPassword(c *gin.Context) {
	var req struct {
		AuthId      string `json:"auth_id"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		a.log.Info("Request Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	query := fmt.Sprintf(`mutation {
	update_password(UpdatePasswordInput:{
		auth_id: "%s"
		old_password: "%s"
		new_password: "%s"
		}) {
		account_id
		}
	}
	`, req.AuthId, req.OldPassword, req.NewPassword)

	qlRequestBody := map[string]any{
		"query": query,
	}
	qlRequestJSON, _ := json.Marshal(qlRequestBody)

	resp, err := http.Post(a.url, "application/json", bytes.NewBuffer(qlRequestJSON))
	if err != nil {
		a.log.Info("GraphQL Request Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to connect to auth service"})
		return
	}
	defer resp.Body.Close()

	// Decode response
	var res struct {
		Data struct {
			UpdatedPassword struct {
				AccountID string `json:"account_id"`
			} `json:"update_password"`
		} `json:"data"`
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		a.log.Info("Read Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}
	if err := json.Unmarshal(body, &res); err != nil {
		a.log.Info("Unmarshal Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	fmt.Println("RAW GQL Response:", string(body))
	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully", "id": res.Data.UpdatedPassword.AccountID})
}
