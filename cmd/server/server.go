package main

import (
	"ascenda-loyalty-assignment/internal/handlers"
	"ascenda-loyalty-assignment/pkg/logging"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const port = ":8000"

func main() {
	logger := logging.LogrusLogger()

	router := gin.Default()

	router.GET("/hotels", handlers.GetAllHotels(logger))
	router.POST("/update_data", handlers.UpdateHotelData(logger))

	err := http.ListenAndServe(port, router)

	if err != nil {
		logger.Critical("HTTP server error: %v", err)
	} else {
		logger.Info(fmt.Sprintf("Server is running on port %s", port))
	}

}
