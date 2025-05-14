package services

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"golang.org/x/crypto/bcrypt"
)

func generateRandomNumber(length int) (string, error) {
	const digits = "0123456789"
	ref := make([]byte, length)
	for i := range ref {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", fmt.Errorf("failed to generate ref_no: %w", err)
		}
		ref[i] = digits[n.Int64()]
	}
	return string(ref), nil
}

// These are just helper functions
func verifyPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

func hashPassword(plainPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
