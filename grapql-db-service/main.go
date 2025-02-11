package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"graphql-db-service/utils"
)

type album struct {
    ID     string  `json:"id"`
    Title  string  `json:"title"`
    Artist string  `json:"artist"`
    Price  float64 `json:"price"`
}

var albums = []album{
    {ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
    {ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
    {ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func getAlbums(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, albums)
}

func getRoot(c *gin.Context) {
    c.String(http.StatusOK, "Welcome to the Home Page!")
}

func getHello(c *gin.Context) {
    c.String(http.StatusOK, "Hello, World!")
}

func main() {

    router := gin.New()
    logger, err := utils.NewLogger()
    if err != nil {
        panic(err)
    }
    logger.Info("Starting the application...")

	serviceAPI := router.Group("/api/db-service") // base path
    {
        serviceAPI.GET("/", getRoot)
        serviceAPI.GET("/albums", getAlbums)
        serviceAPI.GET("/hello", getHello)
    }

    if err := router.Run("localhost:8082"); err != nil {
        logger.Fatal("Failed to start server: %v", err)
    }
}