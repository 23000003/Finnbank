package types

// ==================== Product Types ====================

type CreateProductRequest struct {
	Name string `json:"name" binding:"required"`
	Info string `json:"info" binding:"required"`
	Price float64 `json:"price" binding:"required"`
}

type UpdateProductRequest struct {
	Name string `json:"name"`
	Info string `json:"info"`
	Price float64 `json:"price"`
}

// ==================== Opened Account Types ====================

type GetAllOpenedAccountRequest struct {
	AccountId int `json:"id"`
}

type CreateOpenAccountRequest struct {
	AccountId     int       `json:"account_id"`
	Balance       float64   `json:"balance"`
	AccountType   string    `json:"account_type"`
}

type UpdateOpenAccountRequest struct {
	OpenedAccountId     	int       `json:"openedaccount_id"`
	OpenedAccountStatus   string    `json:"openedaccount_status"`

// ===================== Account Types ====================

type LoginAccountRequest struct {
	Email     	string       `json:"email"`
	Password   	string    `json:"password"`
}

type SignupAccountRequest struct {
	Email     		string    `json:"email"`
	Password   		string    `json:"password"`
	FirstName 		string    `json:"first_name"`
	LastName 			string    `json:"last_name"`
	PhoneNumber 	string    `json:"phone_number"`
	Address 			string    `json:"address"`
	AccountType 	string    `json:"account_type"`
	Nationality 	string    `json:"nationality"`
}