package services

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"finnbank/common/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type BankcardService struct {
	log *utils.Logger
	url string
}

func NewBankcardService(log *utils.Logger) *BankcardService {
	return &BankcardService{
		log: log,
		url: "http://localhost:8083/graphql/bankcard",
	}
}

func generateBankcardNumber(first_name, last_name, birth_date string) string {
	input := strings.ToLower(first_name + last_name + birth_date)

	hash := sha256.Sum256([]byte(input))

	hexString := hex.EncodeToString(hash[:])

	cardNumber := ""

	for _, char := range hexString {
		if len(cardNumber) == 12 {
			break
		}

		if char >= '0' && char <= '9' {
			cardNumber += string(char)
		}
	}

	for len(cardNumber) < 12 {
		cardNumber += "0"
	}

	return cardNumber
}

func (a *BankcardService) GetUserBankcard(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GetUserBankcard not implemented yet"})
}

func (a *BankcardService) GenerateBankcardForUser(c *gin.Context) {
	var request struct {
		Fname string `json:"First_Name" binding:"required"`
		Lname string `json:"Last_Name" binding:"required"`
		Bdate string `json:"Date_Created" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		a.log.Error("Invalid request data:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})

		return
	}

	cardNumber := generateBankcardNumber(request.Fname, request.Lname, request.Bdate)

	// Respond with the generated card number
	c.JSON(http.StatusOK, gin.H{
		"card_number": cardNumber,
	})
}

func (a *BankcardService) RenewBankcardForUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "RenewBankcardForUser not implemented"})
}

func GetBankCardOfUserById(db interface{}, accountID int) (map[string]interface{}, error) {
	conn := db.(*sql.DB) // adjust if you're using a custom DB wrapper
	query := "SELECT BankCard_ID, Card_Number, Expiry, Card_Type FROM BankCard WHERE Account_ID = ?"
	row := conn.QueryRow(query, accountID)

	var cardID int
	var number, cardType string
	var expiry sql.NullTime

	err := row.Scan(&cardID, &number, &expiry, &cardType)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"BankCard_ID": cardID,
		"Card_Number": number,
		"Expiry":      expiry.Time,
		"Card_Type":   cardType,
	}, nil
}

func CreateBankCardForUser(db interface{}, args map[string]interface{}) (string, error) {
	conn := db.(*sql.DB)
	query := `INSERT INTO BankCard (Card_Number, Expiry, Account_ID, Card_Type) VALUES (?, ?, ?, ?)`
	// You should generate card number securely
	cardNumber := "1234567890123456" // mock
	expiry := args["expiry"].(string)
	accountID := args["account_id"].(int)
	cardType := args["card_type"].(string)

	_, err := conn.Exec(query, cardNumber, expiry, accountID, cardType)
	if err != nil {
		return "", err
	}
	return "Bank card created successfully", nil
}

func UpdateBankcardExpiryDateByUserId(db interface{}, args map[string]interface{}) (string, error) {
	conn := db.(*sql.DB)
	query := `UPDATE BankCard SET Expiry = ? WHERE Account_ID = ?`
	expiry := args["new_expiry"].(string)
	accountID := args["account_id"].(int)

	_, err := conn.Exec(query, expiry, accountID)
	if err != nil {
		return "", err
	}
	return "Expiry date updated", nil
}
