package services

import (
	"bytes"
	"encoding/json"
	"finnbank/common/utils"
	t "finnbank/api-gateway/types"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductService struct {
	log *utils.Logger
	url string
}

// ***********************************************
// Test Communication
// See graphql resolvers for each query and mutation call
// ***********************************************

func NewProductService(log *utils.Logger) *ProductService {
	return &ProductService{
		log: log,
		url: "http://localhost:8083/graphql/product",
	}
}

func (p *ProductService) GetAllProduct(ctx *gin.Context) {

	// Define the GraphQL query
	query := `{
		list {
			id
			name
			info
			price
		}
	}`

	qlrequestBody := map[string]interface{}{
		"query": query,
	}

	qlrequestBodyJSON, _ := json.Marshal(qlrequestBody)

	resp, err := http.Post(p.url, "application/json", bytes.NewBuffer(qlrequestBodyJSON))
	if err != nil {
		p.log.Info("Request Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	defer resp.Body.Close()

	// decode the response from ql api
	var data t.GetAllProductGraphQLResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		p.log.Info("Error decoding response: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	p.log.Info("%v ======= DATA", data)

	ctx.JSON(http.StatusOK, gin.H{"data": data.Data.List})
}

// same as getall just a little bit of change
func (p *ProductService) GetByIdProduct(ctx *gin.Context) {
}

func (p *ProductService) CreateProduct(ctx *gin.Context) {
	var req t.CreateProductRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//http://localhost:8083/graphql/product?query=mutation+_{create(name:"a Kaoala",info:"Inca on verbena (wiki)",price:1.99){id,name,info,price}}
	query := fmt.Sprintf(`mutation {
		create(name: "%s", info: "%s", price: %f) {
			id
			name
			info
			price
		}
	}`, req.Name, req.Info, req.Price)

	qlrequestBody := map[string]interface{}{
		"query": query,
	}
	qlrequestBodyJSON, _ := json.Marshal(qlrequestBody)

	resp, err := http.Post(p.url, "application/json", bytes.NewBuffer(qlrequestBodyJSON))
	if err != nil {
		p.log.Info("Request Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	defer resp.Body.Close()

	var data t.CreateProductGraphQLResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		p.log.Info("Error decoding response: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	p.log.Info("%v ======= DATA", data)

	if data.Errors != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": data.Errors})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data.Data.Create})
}

// same as post just a little bit of change
func (p *ProductService) UpdateProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	p.log.Info("ID: %v", id)

	var req t.UpdateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name == "" && req.Info == "" && req.Price == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "At least one field must be provided for update"})
		return
	}
}

// same as post just a little bit of change
func (p *ProductService) DeleteProduct(ctx *gin.Context) {
}
