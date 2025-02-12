package service

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

type AccountService struct {
	Client *supabase.Client
}

type Account struct {
	Email         string  `json:"email"`
	Full_Name     string  `json:"full_name"`
	Phone         string  `json:"phone_number"`
	Password      string  `json:"password"`
	HasCard       bool    `json:"has_card"`
	AccountNumber string  `json:"account_number"`
	Address       string  `json:"address"`
	Balance       float64 `json:"balance"`
	AccountType   string  `json:"account_type"`
}

func SupabaseInit() *supabase.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Missing env files")
	}
	url := os.Getenv("DB_URL")
	key := os.Getenv("DB_KEY")

	if url == "" || key == "" {
		log.Fatalf("Supabase URL and Keys missing")
	}
	client, err := supabase.NewClient(url, key, &supabase.ClientOptions{})
	if err != nil {
		log.Fatalf("Failed to initialize Supabase client: %v", err)
	}

	return client
}

// @Summary Get accounts
// @Description Fetch all accounts
// @Tags accounts
// @Produce json
// @Success 200 {array} Account
// @Router /accounts [get]
func (s *AccountService) GetAccounts(c *gin.Context) {
	data, count, err := s.Client.From("account").Select("*", "exact", false).Execute()
	if err != nil  || count == 0{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch accounts"})
		return
	}
	c.IndentedJSON(http.StatusOK, data)
}

// @Summary Register an account
// @Description Create a new account
// @Tags accounts
// @Accept json
// @Produce json
// @Param request body Account true "Account data"
// @Success 201 {object} Account
// @Failure 400 {object} map[string]string
// @Router /register-acc [post]
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
	response, count, err := s.Client.From("account").Insert(newAcc, false, "", "", "exact").Execute()
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

	c.JSON(http.StatusOK, gin.H{"message": "Account created successfully", "account": insertedAcc})
}

func dataCheck(account Account) bool {
	return account.Email != "" &&
		account.Full_Name != "" &&
		account.Phone != "" &&
		account.Password != "" &&
		account.AccountNumber != "" &&
		account.AccountType != ""
}
