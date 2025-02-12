package service

// Comments with "@" are for Swagger, if samokan mo tanawn e erase lng and just use <<postman>> lng
// run ==swag init== if namoy bag o service ge add for the routers
// https://github.com/swaggo/swag for api operiations ( " @ " )

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"fmt"
)

// types => for data validation
type album struct {
    ID     string  `json:"id"`
    Title  string  `json:"title"`
    Artist string  `json:"artist"`
    Price  float64 `json:"price"`
}

type AccountService struct {}

type HelloRequest struct {
    Name    string `json:"name" binding:"required"`
    Message string `json:"message"`
}

type CreateAlbumRequest struct {
    Title  string  `json:"title" binding:"required"`
    Artist string  `json:"artist" binding:"required"`
    Price  float64 `json:"price" binding:"required"`
}

var albums = []album{
    {ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
    {ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
    {ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}


// @Summary Get root endpoint
// @Description Welcome message for root endpoint
// @Tags root
// @Produce plain
// @Success 200 {string} string
// @Router / [get]
func (c *AccountService) GetRoot(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Welcome to the Home Page!")
}

// @Summary Get all albums
// @Description Get a list of all albums
// @Tags albums
// @Produce json
// @Success 200 {array} album
// @Router /albums [get]
func (c *AccountService) GetAlbums(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, albums)
}

// @Summary Get album by ID
// @Description Get a single album by its ID
// @Tags albums
// @Produce json
// @Param id path string true "Album ID"
// @Success 200 {object} album
// @Failure 404 {object} map[string]string
// @Router /albums/{id} [get]
func (c *AccountService) GetAlbumsById(ctx *gin.Context) { // Dynamic through Params
    id := ctx.Param("id")

    // Loop through albums looking for matching id
    for _, a := range albums {
        if a.ID == id {
            ctx.JSON(http.StatusOK, a)
            return
        }
    }

    ctx.JSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

// @Summary Post hello endpoint
// @Description Post hello message with custom payload
// @Tags hello
// @Accept json
// @Produce json
// @Param request body HelloRequest true "Hello request"
// @Success 200 {object} map[string]string
// @Router /hello [post]
func (c *AccountService) PostHello(ctx *gin.Context) {

    // Get headers
    // userAgent := ctx.GetHeader("authorization")

    // parse request body and handle validation
    var request HelloRequest
    if err := ctx.ShouldBindJSON(&request); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // return response via body
    ctx.JSON(http.StatusOK, gin.H{
        "message": fmt.Sprintf("Hello, %s!", request.Name),
        "yourMessage": request.Message,
    })
}

// @Summary Create new album
// @Description Add a new album to the collection
// @Tags albums
// @Accept json
// @Produce json
// @Param album body CreateAlbumRequest true "Album to create"
// @Success 201 {object} album
// @Failure 400 {object} map[string]string
// @Router /albums [post]
func (c *AccountService) CreateAlbum(ctx *gin.Context) {
    var newAlbum CreateAlbumRequest

    // parse request body and handle validation
    if err := ctx.ShouldBindJSON(&newAlbum); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request format",
            "details": err.Error(),
        })
        return
    }

    // create new album with generated ID
    album := album{
        ID:     fmt.Sprintf("%d", len(albums) + 1),
        Title:  newAlbum.Title,
        Artist: newAlbum.Artist,
        Price:  newAlbum.Price,
    }

    // add to albums slice
    albums = append(albums, album)

    // return created album
    ctx.JSON(http.StatusCreated, album)
}