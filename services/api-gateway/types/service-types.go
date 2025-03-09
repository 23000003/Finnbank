package types

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