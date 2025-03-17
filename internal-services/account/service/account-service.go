package service

import (
	"crypto/rand"
	"encoding/json"
	"finnbank/common/utils"
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supabase-community/auth-go"
	"github.com/supabase-community/auth-go/types"
	"github.com/supabase-community/supabase-go"
)

type AccountService struct {
	Client *supabase.Client
	Auth   auth.Client
	Logger *utils.Logger
}

type Account struct {
	Email       string `json:"email"`
	Full_Name   string `json:"full_name"`
	Phone       string `json:"phone_number"`
	Password    string `json:"password"`
	Address     string `json:"address"`
	AccountType string `json:"account_type"`
}

// @Summary Get accounts
// @Description Fetch all accounts
// @Tags accounts
// @Produce json
// @Success 200 {array} Account
// @Router /accounts [get]
func (s *AccountService) GetAccounts(c *gin.Context) {
	var accGot []Account
	data, count, err := s.Client.From("account").Select("*", "exact", false).Execute()
	s.Logger.Debug("%s", string(data))
	if err != nil || count == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": data})
		return
	}
	if err := json.Unmarshal(data, &accGot); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data", "details": err.Error()})
	}
	c.IndentedJSON(http.StatusOK, accGot)
}

// @Summary Get an account by ID
// @Description Fetch an account using the account number
// @Tags accounts
// @Produce json
// @Param acc_num path string true "Account Number"
// @Success 200 {object} Account
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "Account not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /fetch-acc/{acc_num} [get]
func (s *AccountService) GetAccoutById(c *gin.Context) {
	accNum := c.Param("acc_num")
	var accGot []Account
	s.Logger.Debug("%s", accNum)
	response, count, err := s.Client.From("account").
		Select("*", "exact", false).
		Eq("account_number", accNum).
		Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch account"})
		return
	}
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account does not exist"})
		return
	}
	if err := json.Unmarshal(response, &accGot); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data", "details": err.Error()})
	}
	c.IndentedJSON(http.StatusOK, accGot[0])
}

// @Summary Update HasCard field
// @Description Updates the has_card attribute of an account using account_number
// @Tags accounts
// @Accept json
// @Produce json
// @Param acc_num path string true "Account Number"
// @Param request body map[string]bool true "Updated has_card value"
// @Success 200 {object} map[string]interface{} "Account updated OK"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /update_acc/{acc_num} [patch]
func (s *AccountService) UpdateHasCard(c *gin.Context) {
	accountNum := c.Param("acc_num")
	var updateData map[string]interface{}
	s.Logger.Debug("%s", accountNum)

	if s.Client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database client not initialized"})
		return
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Server Request " + err.Error()})
		return
	}

	response, count, err := s.Client.From("account").
		Update(updateData, "", "exact").Eq("account_number", accountNum).
		Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed updating account " + err.Error()})
		return
	}
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find account "})
		return
	}
	var updatedAcc []Account
	err = json.Unmarshal(response, &updatedAcc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data", "details": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Account Updated Succesfully", "response": updatedAcc})

}

// @Summary Delete a user
// @Description Deletes a user from the "account" table and removes them from Supabase Auth.
// @Tags Users
// @Accept json
// @Produce json
// @Param email path string true "User Email"
// @Success 200 {object} map[string]string "Successfully deleted user"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /delete-user/{email} [delete]
func (s *AccountService) DeleteUser(c *gin.Context) {

	email, exists := c.Get("email")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UUID not found in context"})
		c.Abort()
		return
	}

	res, _, err := s.Client.From("account").
		Delete("", "exact").
		Eq("email", email.(string)).
		Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed Deleting Account: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted User successfully: ", "response": string(res)})

}

// @Summary Register an account
// @Description Create a new account
// @Tags accounts
// @Accept json
// @Produce json
// @Param request body Account true "Account data"
// @Success 201 {object} Account
// @Failure 400 {object} map[string]string
// @Router /register [post]
func (s *AccountService) AddAccount(c *gin.Context) {
	var newAcc Account
	if err := c.ShouldBindJSON(&newAcc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON request", "details": err.Error()})
		return
	}
	if !dataCheck(newAcc) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing data"})
		return
	}
	accNum, err := generateAccountNumber(s.Client)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate account number"})
	}

	authRes, err := s.Auth.Signup(types.SignupRequest{
		Email:    newAcc.Email,
		Password: newAcc.Password,
		Data: map[string]interface{}{
			"account_number": accNum,
			"phone":          newAcc.Phone,
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user in Auth: " + err.Error()})
		return
	}

	accountData := map[string]interface{}{
		"email":          newAcc.Email,
		"full_name":      newAcc.Full_Name,
		"phone_number":   newAcc.Phone,
		"has_card":       false,
		"account_number": accNum,
		"address":        newAcc.Address,
		"balance":        0.00,
		"account_type":   newAcc.AccountType,
	}
	response, count, err := s.Client.From("account").
		Insert(accountData, false, "", "", "exact").
		Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No account was inserted"})
		return
	}
	var insertedAcc []Account
	err = json.Unmarshal(response, &insertedAcc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account created successfully", "account": insertedAcc, "user": authRes})
}

func generateAccountNumber(client *supabase.Client) (string, error) {
	for {
		num, err := rand.Int(rand.Reader, big.NewInt(9999999999999999))
		if err != nil {
			return "", err
		}

		accountNum := num.String()

		for len(accountNum) < 16 {
			accountNum = "0" + accountNum
		}

		_, count, err := client.From("account").Select("account_number", "exact", false).Eq("account_number", accountNum).Execute()
		if err != nil {
			return "", err
		}

		if count == 0 {
			return accountNum, nil
		}
	}
}

func dataCheck(account Account) bool {
	return account.Email != "" &&
		account.Full_Name != "" &&
		account.Phone != "" &&
		account.Password != "" &&
		account.AccountType != ""
}
