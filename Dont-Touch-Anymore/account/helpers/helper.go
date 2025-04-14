package helpers

import (
	"math/rand"
	"time"
)

// genAccNum generates a random 16-digit account number as a string.
func GenAccNum() (string, error) {
	rand.Seed(time.Now().UnixNano())
	var accNum string
	for i := 0; i < 16; i++ {
		accNum += string(rune('0' + rand.Intn(10))) // Generates a random digit from 0-9
	}
	return accNum, nil
}
