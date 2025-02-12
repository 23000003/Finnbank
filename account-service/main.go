package main

import (
	"account-service/utils"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

type Account struct {
	Email         string    `json:"email"`
	Full_Name     string    `json:"full_name"`
	Phone         string    `json:"phone_number"`
	Password      string    `json:"password"`
	HasCard       bool      `json:"has_card"`
	AccountNumber string    `json:"account_number"`
	Address       string    `json:"address"`
	Balance       float64   `json:"balance"`
	AccountType   string    `json:"account_type"`
	DateCreated   time.Time `json:"date_created"`
	DateUpdated   time.Time `json:"date_updated"`
}

func supabaseInit() *supabase.Client {
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

func getAccounts(client *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, count, err := client.From("account").Select("*", "exact", false).Execute()
		if err != nil {
			c.JSON(int(count), gin.H{"error": "Failed to fetch accounts"})
			return
		}
		c.IndentedJSON(http.StatusOK, data)
	}
}

func dataCheck(account Account) bool {
	return account.Email != "" &&
		account.Full_Name != "" &&
		account.Phone != "" &&
		account.Password != "" &&
		account.AccountNumber != "" &&
		account.AccountType != ""
}

func addAccount(client *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newAcc Account
		if err := c.ShouldBindJSON(&newAcc); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON request", "details": err.Error()})
			return
		}
		if !dataCheck(newAcc) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing data"})
		}
		response, count, err := client.From("account").Insert(newAcc, false, "", "", "exact").Execute()
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
}

func getRoot(c *gin.Context) {
	c.String(http.StatusOK, "Welcome to the Home Page!")
}

func getHello(c *gin.Context) {
	c.String(http.StatusOK, "Hello, World!")
}

func main() {

	client := supabaseInit()
	router := gin.New()
	logger, err := utils.NewLogger()
	if err != nil {
		panic(err)
	}
	logger.Info("Starting the application...")

	serviceAPI := router.Group("/api/account-service") // base path
	{
		serviceAPI.GET("/", getRoot)
		serviceAPI.GET("/accounts", getAccounts(client))
		serviceAPI.POST("/register-acc", addAccount(client))
		serviceAPI.GET("/hello", getHello)
	}

	if err := router.Run("localhost:8081"); err != nil {
		logger.Fatal("Failed to start server: %v", err)
	}
}
