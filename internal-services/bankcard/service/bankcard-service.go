package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type BankCard struct {
	ID             int
	Account_holder string
	Account_number string
	Card_number    string
	Expiry_date    string
}

type BankCardService struct{}

func generateSomeHash(fullname, birthdate, uuid string) string {
	data := fullname + birthdate + uuid
	hash := sha256.Sum256([]byte(data))

	return hex.EncodeToString(hash[:])
}

func (b *BankCardService) GenerateCardNumber(fullname, birthdate, uuid string) string {
	hashedData := generateSomeHash(fullname, birthdate, uuid)

	cardNumber := hashedData[:15]

	var digits []int

	for _, char := range cardNumber {
		num, err := strconv.Atoi(string(char))

		if err != nil {
			num = int(char) % 10
		}

		digits = append(digits, num)
	}

	// Luhn Algo Checksum suggested by ChatGPT for generating the bank card number
	sum := 0
	for i := 0; i < 15; i++ {
		num := digits[i]

		if i%2 == 0 {
			num *= 2

			if num > 9 {
				num -= 9
			}
		}

		sum += num
	}

	checksum := (10 - (sum % 10)) % 10
	digits = append(digits, checksum)

	var result strings.Builder

	for _, digit := range digits {
		result.WriteString(strconv.Itoa(digit))
	}

	return result.String()
}

// CreateCard generates and returns a new bank card
func (b *BankCardService) CreateCard(accountHolder, accountNumber, birthdate, uuid string) BankCard {
	return BankCard{
		ID:             rand.Intn(10000),
		Account_holder: accountHolder,
		Account_number: accountNumber,
		Card_number:    b.GenerateCardNumber(accountHolder, birthdate, uuid),
		Expiry_date:    fmt.Sprintf("%02d/%d", rand.Intn(12)+1, time.Now().Year()+3),
	}
}
