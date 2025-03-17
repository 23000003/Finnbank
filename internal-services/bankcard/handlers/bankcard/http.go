package bankcard

import (
	"encoding/json"
	"net/http"
)

/***********
	Use this for TESTS
	Http Handlers for Account Service
	This is optional since we are using grpc for communication not http
***********/

func CreateCardHandler(w http.ResponseWriter, r *http.Request) {
	var Request struct {
		Account_holder string `json:"Account_holder"`
		Account_number string `json:"Account_number"`
		birthdate      string `json:"birthdate"`
		uuid           string `json:"uuid"`
	}

	if err := json.NewDecoder(r.Body).Decode(&Request); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	bankservice := services.bankcard{}
	newCard := bankservice.CreateCard(Request.Account_holder, Request.Account_number, Request.birthdate, Request.uuid)

	if err := db.SaveCard(newCard); err != nil {
		http.Error(w, "Failed to save card", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newCard)
}
