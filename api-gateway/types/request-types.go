package types

import "time"

// ==================== Product Types ====================

type CreateProductRequest struct {
	Name  string  `json:"name" binding:"required"`
	Info  string  `json:"info" binding:"required"`
	Price float64 `json:"price" binding:"required"`
}

type UpdateProductRequest struct {
	Name  string  `json:"name"`
	Info  string  `json:"info"`
	Price float64 `json:"price"`
}

// ==================== Opened Account Types ====================

type GetAllOpenedAccountRequest struct {
	AccountId string `json:"id"`
}

type CreateOpenAccountRequest struct {
	AccountId   string  `json:"account_id"`
}

type UpdateOpenAccountRequest struct {
	OpenedAccountId     int    `json:"openedaccount_id"`
	OpenedAccountStatus string `json:"openedaccount_status"`
}

// ===================== Account Types ====================

type LoginAccountRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupAccountRequest struct {
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	FirstName   string    `json:"first_name"`
	MiddleName  string    `json:"middle_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	Address     string    `json:"address"`
	AccountType string    `json:"account_type"`
	NationalID  string    `json:"national_id"`
	Nationality string    `json:"nationality"`
	Birthdate   time.Time `json:"birthdate"`
}

// ====================== Transaction Types ======================

type CreateTransactionRequest struct {
	SenderId        int   `json:"sender_id"`
	ReceiverId      int   `json:"receiver_id"`
	TransactionType string   `json:"transaction_type"`
	Amount          float64  `json:"amount"`
	TransactionFee  float64  `json:"transaction_fee"`
	Notes           string   `json:"notes"`
}

// ====================== Notifcation Types ======================

type CreateNotificationRequest struct {
	NotifType     string `json:"notif_type"`
	NotifToID     string `json:"notif_to_id"`
	NotifFromName string `json:"notif_from_name"`
	Content       string `json:"content"`
}