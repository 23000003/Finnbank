package types

import (
	"time"
)

// ==================== Product Types ====================

type GetAllProductGraphQLResponse struct {
	Data struct {
		List []struct {
			ID    int64   `json:"id"`
			Name  string  `json:"name"`
			Info  string  `json:"info"`
			Price float64 `json:"price"`
		} `json:"list"`
	} `json:"data"`
	Errors any `json:"errors"`
}

type CreateProductGraphQLResponse struct {
	Data struct {
		Create struct {
			ID    int64   `json:"id"`
			Name  string  `json:"name"`
			Info  string  `json:"info"`
			Price float64 `json:"price"`
		} `json:"create"`
	} `json:"data"`
	Errors any `json:"errors"`
}

// =============== Opened Account Types ====================

type GetAllOpenedAccountsGraphQLResponse struct {
	Data struct {
		GetAll []struct {
			OpenedAccountID     int       `json:"openedaccount_id"`
			BankCardID          int       `json:"bankcard_id"`
			Balance             float64   `json:"balance"`
			AccountType         string    `json:"account_type"`
			DateCreated         time.Time `json:"date_created"`
			OpenedAccountStatus string    `json:"openedaccount_status"`
			AccountNumber       string    `json:"account_number"`
		} `json:"get_all"`
	} `json:"data"`
	Errors any `json:"errors"`
}

type GetOpenedAccountsGraphQLResponse struct {
	Data struct {
		GetById struct {
			OpenedAccountID     int       `json:"openedaccount_id"`
			BankCardID          int       `json:"bankcard_id"`
			Balance             float64   `json:"balance"`
			AccountType         string    `json:"account_type"`
			DateCreated         time.Time `json:"date_created"`
			OpenedAccountStatus string    `json:"openedaccount_status"`
			AccountNumber       string    `json:"account_number"`
		} `json:"get_by_id"`
	} `json:"data"`
	Errors any `json:"errors"`
}

type GetUserIdByOpenedAccountIdGraphQLResponse struct {
	Data struct {
		GetById struct {
			AccountID string `json:"account_id"`
		} `json:"get_by_id"`
	} `json:"data"`
	Errors any `json:"errors"`
}

type GetBothAccountNumberGraphQLResponse struct {
	Data struct {
		FindBothAccountNumber []struct {
			OpenedAccountID int    `json:"openedaccount_id"`
			AccountNumber   string `json:"account_number"`
		} `json:"find_both_account_num"`
	} `json:"data"`
	Errors any `json:"errors"`
}

type GetOpenedAccountIdGraphQLResponse struct {
	Data struct {
		FindByAccountNum int `json:"find_by_account_num"`
	} `json:"data"`
	Errors any `json:"errors"`
}

type CreateOpenedAccountsGraphQLResponse struct {
	Data struct {
		CreateAccount []struct {
			OpenedAccountID     int       `json:"openedaccount_id"`
			BankCardID          int       `json:"bankcard_id"`
			Balance             float64   `json:"balance"`
			AccountType         string    `json:"account_type"`
			DateCreated         time.Time `json:"date_created"`
			OpenedAccountStatus string    `json:"openedaccount_status"`
		} `json:"create_account"`
	} `json:"data"`
	Errors any `json:"errors"`
}

type UpdateOpenedAccountsGraphQLResponse struct {
	Data struct {
		UpdateAccountStatus string `json:"update_account_status"`
	} `json:"data"`
	Errors any `json:"errors"`
}

// ===================== Account Types ====================

type AccountLoginGraphQLResponse struct {
	Data struct {
		Login struct {
			AccessToken   string `json:"access_token"`
			FullName      string `json:"full_name"`
			AccountId     string `json:"account_id"`
			AccountStatus string `json:"account_status"`
		} `json:"login"`
	} `json:"data"`
	Errors any `json:"errors"`
}

type AccountSignUpGraphQLResponse struct {
	Data struct {
		CreateAccount struct {
			Email     string `json:"email"`
			AuthID    string `json:"auth_id"`
			AccountID string `json:"account_id"`
		} `json:"create_account"`
	} `json:"data"`
	Errors any `json:"errors"`
}

type GetAccountDetailsGraphQLResponse struct {
	Data struct {
		AccountById struct {
			Email         string    `json:"email"`
			First_Name    string    `json:"first_name"`
			Middle_Name   string    `json:"middle_name"`
			Last_Name     string    `json:"last_name"`
			PhoneNumber   string    `json:"phone_number"`
			DateCreated   time.Time `json:"date_created"`
			AccountNumber string    `json:"account_number"`
			NationalId    string    `json:"national_id"`
			AccountStatus string    `json:"account_status"`
			Address       string    `json:"address"`
			Nationality   string    `json:"nationality"`
			Birthdate     string    `json:"birthdate"`
			AccountType   string    `json:"account_type"`
		} `json:"account_by_id"`
	} `json:"data"`
	Errors any `json:"errors"`
}

type AccountNumberFetchResponse struct {
	Data struct {
		FetchedAccount struct {
			AccountID     string  `json:"account_id"`
			AccountStatus string  `json:"account_status"`
			AccountType   string  `json:"account_type"`
			Address       string  `json:"address"`
			AuthID        string  `json:"auth_id"`
			Balance       float64 `json:"balance"`
			Birthdate     string  `json:"birthdate"`
			DateCreated   string  `json:"date_created"`
			DateUpdated   string  `json:"date_updated"`
			Email         string  `json:"email"`
			FirstName     string  `json:"first_name"`
			HasCard       bool    `json:"has_card"`
			LastName      string  `json:"last_name"`
			MiddleName    string  `json:"middle_name"`
			NationalID    string  `json:"national_id"`
			Nationality   string  `json:"nationality"`
			PhoneNumber   string  `json:"phone_number"`
		} `json:"account_by_account_num"`
	} `json:"data"`
	Errors any `json:"errors"`
}
type EmailFetchResponse struct {
	Data struct {
		FetchedAccount struct {
			AccountID     string  `json:"account_id"`
			AccountStatus string  `json:"account_status"`
			AccountType   string  `json:"account_type"`
			Address       string  `json:"address"`
			AuthID        string  `json:"auth_id"`
			Balance       float64 `json:"balance"`
			Birthdate     string  `json:"birthdate"`
			DateCreated   string  `json:"date_created"`
			DateUpdated   string  `json:"date_updated"`
			Email         string  `json:"email"`
			FirstName     string  `json:"first_name"`
			HasCard       bool    `json:"has_card"`
			LastName      string  `json:"last_name"`
			MiddleName    string  `json:"middle_name"`
			NationalID    string  `json:"national_id"`
			Nationality   string  `json:"nationality"`
			PhoneNumber   string  `json:"phone_number"`
		} `json:"account_by_email"`
	} `json:"data"`
	Errors any `json:"errors"`
}

// ================= Transaction Types ====================

type GetAllTransactionsGraphQLResponse struct {
	Data struct {
		GetTransactionsByUserId []struct {
			TransactionID     int       `json:"transaction_id"`
			RefNo             string    `json:"ref_no"`
			SenderID          int       `json:"sender_id"`
			ReceiverID        int       `json:"receiver_id"`
			TransactionType   string    `json:"transaction_type"`
			Amount            float64   `json:"amount"`
			TransactionStatus string    `json:"transaction_status"`
			DateTransaction   time.Time `json:"date_transaction"`
			TransactionFee    float64   `json:"transaction_fee"`
			Notes             string    `json:"notes"`
		} `json:"getTransactionsByUserId"`
	} `json:"data"`
	Errors any `json:"errors"`
}

type GetAllTransactionsByTimeStampGraphQLResponse struct {
	Data struct {
		GetTransactionsByTimeStampByUserId []struct {
			TransactionID     int       `json:"transaction_id"`
			RefNo             string    `json:"ref_no"`
			SenderID          int       `json:"sender_id"`
			ReceiverID        int       `json:"receiver_id"`
			TransactionType   string    `json:"transaction_type"`
			Amount            float64   `json:"amount"`
			TransactionStatus string    `json:"transaction_status"`
			DateTransaction   time.Time `json:"date_transaction"`
			TransactionFee    float64   `json:"transaction_fee"`
			Notes             string    `json:"notes"`
		} `json:"getTransactionsByTimeStampByUserId"`
	} `json:"data"`
	Errors any `json:"errors"`
}

type CreateTransactionsGraphQLResponse struct {
	Data struct {
		CreateTransaction struct {
			TransactionID     int       `json:"transaction_id"`
			RefNo             string    `json:"ref_no"`
			SenderID          int       `json:"sender_id"`
			ReceiverID        int       `json:"receiver_id"`
			TransactionType   string    `json:"transaction_type"`
			Amount            float64   `json:"amount"`
			TransactionStatus string    `json:"transaction_status"`
			DateTransaction   time.Time `json:"date_transaction"`
			TransactionFee    float64   `json:"transaction_fee"`
			Notes             string    `json:"notes"`
		} `json:"createTransaction"`
	} `json:"data"`
	Errors any `json:"errors"`
}

// ===================== Notification Types ====================

type GetAllNotificationsGraphQLResponse struct {
	Data struct {
		GetAllNotificationByUserId []struct {
			NotifID       string    `json:"notif_id"`
			NotifType     string    `json:"notif_type"`
			NotifFromName string    `json:"notif_from_name"`
			Content       string    `json:"content"`
			IsRead        bool      `json:"is_read"`
			DateNotified  time.Time `json:"date_notified"`
		} `json:"getAllNotificationByUserId"`
	} `json:"data"`
	Errors any `json:"errors"`
}

type GetAllUnreadNotificationGraphQLResponse struct {
	Data struct {
		GetAllUnreadNotificationByUserId struct {
			TotalNotification  int `json:"total_notification"`
			UnreadNotification int `json:"unread_notification"`
		} `json:"getAllUnreadNotificationByUserId"`
	} `json:"data"`
	Errors any `json:"errors"`
}

type GetNotificationGraphQLResponse struct {
	Data struct {
		GetNotificationById struct {
			NotifID       string    `json:"notif_id"`
			NotifType     string    `json:"notif_type"`
			NotifFromName string    `json:"notif_from_name"`
			Content       string    `json:"content"`
			IsRead        bool      `json:"is_read"`
			DateNotified  time.Time `json:"date_notified"`
		} `json:"getNotificationById"`
	} `json:"data"`
	Errors any `json:"errors"`
}

type GetCreatedNotificationsGraphQLResponse struct {
	Data struct {
		GenerateNotificaton struct {
			NotifID string `json:"notif_id"`
		} `json:"generateNotification"`
	} `json:"data"`
	Errors any `json:"errors"`
}

type ReadNotificationGraphQLResponse struct {
	Data struct {
		ReadNotificationByUserId bool `json:"readNotificationByUserId"`
	} `json:"data"`
	Errors any `json:"errors"`
}

// ===================== Bankcard Types ====================

type GetAllBankCardsGraphQLResponse struct {
	Data struct {
		GetAllBankCard []struct {
			BankCardID  int       `json:"bankcard_id"`
			CardType    string    `json:"card_type"`
			CardNumber  string    `json:"card_number"`
			ExpiryDate  time.Time `json:"expiry_date"`
			DateCreated time.Time `json:"date_created"`
			CVV         string    `json:"cvv"`
		} `json:"get_all_bankcard"`
	} `json:"data"`
	Errors any `json:"errors"`
}

type UpdateBankCardGraphQLResponse struct {
	Data struct {
		UpdateBankcardExpiry []struct {
			BankCardID     int       `json:"bankcard_id"`
			BankCardType   string    `json:"bankcard_type"`
			BankCardNumber string    `json:"bankcard_number"`
			ExpiryDate     time.Time `json:"expiry_date"`
			DateCreated    time.Time `json:"date_created"`
			CVV            string    `json:"cvv"`
		} `json:"update_bankcard_expiry"`
	} `json:"data"`
	Errors any `json:"errors"`
}

type UpdateBankCardPinNumberGraphQLResponse struct {
	Data struct {
		UpdatePinNumber []struct {
			BankCardID     int       `json:"bankcard_id"`
			BankCardType   string    `json:"bankcard_type"`
			BankCardNumber string    `json:"bankcard_number"`
			ExpiryDate     time.Time `json:"expiry_date"`
			DateCreated    time.Time `json:"date_created"`
			CVV            string    `json:"cvv"`
		} `json:"update_pin_number"`
	} `json:"data"`
	Errors any `json:"errors"`
}

// ===================== Statement Types ====================

type GetStatementGraphQLResponse struct {
	Data struct {
		GenerateStatement struct {
			PdfBuffer string `json:"pdf_buffer"`
		} `json:"generate_statement"`
	} `json:"data"`
	Errors any `json:"errors"`
}

// ===================== RealTime Types ====================

type GetRealTimeNotification struct {
	NotifID       string    `json:"notif_id"`
	NotifType     string    `json:"notif_type"`
	NotifFromName string    `json:"notif_from_name"`
	NotifToID     string    `json:"notif_to_id"`
	Content       string    `json:"content"`
	IsRead        bool      `json:"is_read"`
	DateNotified  time.Time `json:"date_notified"`
}

type GetRealTimeTransaction struct {
	TransactionID     int       `json:"transaction_id"`
	RefNo             string    `json:"ref_no"`
	SenderID          int       `json:"sender_id"`
	ReceiverID        int       `json:"receiver_id"`
	TransactionType   string    `json:"transaction_type"`
	Amount            float64   `json:"amount"`
	TransactionStatus string    `json:"transaction_status"`
	DateTransaction   time.Time `json:"date_transaction"`
	TransactionFee    float64   `json:"transaction_fee"`
	Notes             string    `json:"notes"`
}
