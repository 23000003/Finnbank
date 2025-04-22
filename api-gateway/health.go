package main

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Server is OK")
}

func GraphQLAPIHealthCheck(ctx *gin.Context) {
	resp, err := http.Get("http://localhost:8083/graphql/health")
	if err != nil {
		ctx.String(http.StatusInternalServerError, "GraphQL server error")
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK {
		ctx.String(http.StatusOK, string(body))
	} else {
		ctx.String(http.StatusServiceUnavailable, "GraphQL server is down")
	}
}
