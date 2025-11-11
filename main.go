package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"urlfetcher/fetcher"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.POST("/execute", fetcher.ExecuteHandler)

	const port = ":8080"
	fmt.Printf("server is running on %s port", port)

	err := router.Run(port)
	if err != nil {
		fmt.Errorf("failed to start server %v", err)
		return
	}
}
