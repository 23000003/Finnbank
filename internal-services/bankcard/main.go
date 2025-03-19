package main

import (
	"fmt"
	"internal-services/bankcard/service"
	"net/http"
)

func main() {
	bankCardService := service.BankCardService{}

	http.HandleFunc("/create-card", func(w http.ResponseWriter, r *http.Request) {
		accountHolder := r.URL.Query().Get("accountHolder")
		accountNumber := r.URL.Query().Get("accountNumber")
		birthdate := r.URL.Query().Get("birthdate")
		uuid := r.URL.Query().Get("uuid")

		if accountHolder == "" || accountNumber == "" || birthdate == "" || uuid == "" {
			http.Error(w, "Missing parameters", http.StatusBadRequest)
			return
		}

		card := bankCardService.CreateCard(accountHolder, accountNumber, birthdate, uuid)
		fmt.Fprintf(w, "Card created: %+v", card)
	})

	fmt.Println("BankCard service is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
