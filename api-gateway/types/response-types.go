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
	Errors interface{} `json:"errors"`
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
	Errors interface{} `json:"errors"`
}

// =============== Opened Account Types ====================

type GetAllOpenedAccountsGraphQLResponse struct {
	Data struct {
		GetAll []struct {
			OpenedAccountID     int       `json:"openedaccount_id"`
			BankCardID          *int      `json:"bankcard_id"`
			Balance             float64   `json:"balance"`
			AccountType         string    `json:"account_type"`
			DateCreated         time.Time `json:"date_created"`
			OpenedAccountStatus string    `json:"openedaccount_status"`
		} `json:"get_all"`
	} `json:"data"`
	Errors interface{} `json:"errors"`
}

type GetOpenedAccountsGraphQLResponse struct {
	Data struct {
		 GetById struct {
			OpenedAccountID     int       `json:"openedaccount_id"`
			BankCardID          *int      `json:"bankcard_id"`
			Balance             float64   `json:"balance"`
			AccountType         string    `json:"account_type"`
			DateCreated         time.Time `json:"date_created"`
			OpenedAccountStatus string    `json:"openedaccount_status"`
		} `json:"get_by_id"`
	} `json:"data"`
	Errors interface{} `json:"errors"`
}

type CreateOpenedAccountsGraphQLResponse struct {
	Data struct {
			CreateAccount struct {
			OpenedAccountID     int       `json:"openedaccount_id"`
			BankCardID          *int      `json:"bankcard_id"`
			Balance             float64   `json:"balance"`
			AccountType         string    `json:"account_type"`
			DateCreated         time.Time `json:"date_created"`
			OpenedAccountStatus string    `json:"openedaccount_status"`
		} `json:"create_account"`
	} `json:"data"`
	Errors interface{} `json:"errors"`
}

type UpdateOpenedAccountsGraphQLResponse struct {
	Data struct {
			UpdateAccountStatus string `json:"update_account_status"`
	} `json:"data"`
	Errors interface{} `json:"errors"`
}


// ===================== Account Types ====================

type AccountLoginGraphQLResponse struct {
	Data struct {
		Login struct {
			AccessToken   string    `json:"access_token"`
			FullName     	string    `json:"full_name"`
			AccountId     string    `json:"account_id"`
		} `json:"login"`
	} `json:"data"`
	Errors interface{} `json:"errors"`
}

type AccountSignUpGraphQLResponse struct {
	Data struct {
		CreateAccount struct {
			AccessToken   string      `json:"access_token"`
			Email     		string    `json:"email"`
			AuthID     		string    `json:"auth_id"`
		} `json:"create_account"`
	} `json:"data"`
	Errors interface{} `json:"errors"`
}


