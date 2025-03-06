package handlers

import (
	"ascenda-loyalty-assignment/internal/services/hotel_service"
	"ascenda-loyalty-assignment/pkg/logging"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	hotelsDataFileName      = "hotels.json"
	suppliersDataFileName   = "suppliers.json"
	defaultTimeOutInSeconds = 60
)

type HotelQueryParams struct {
	HotelIDs       []string `form:"hotelIds"`
	DestinationIDs []int    `form:"destinationIds"`
}

func GetAllHotels(logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var queryParams HotelQueryParams
		if err := c.ShouldBindQuery(&queryParams); err != nil {
			logger.Error(fmt.Sprintf("Invalid request params %v", c.Params), err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request params"})
			return
		}
		wd, err := os.Getwd()
		if err != nil {
			logger.Error("Error getting working directory", err)
			c.Status(http.StatusInternalServerError)
			return
		}
		hotelService := hotel_service.NewHotelService(logger, nil, c)
		hotelDataFilePath := filepath.Join(wd, "internal", "data", hotelsDataFileName)
		hotels, err := hotelService.GetHotels(hotelDataFilePath, queryParams.HotelIDs, queryParams.DestinationIDs)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, hotels)
	}
}

func UpdateHotelData(logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		wd, err := os.Getwd()
		if err != nil {
			logger.Error("Error getting working directory", err)
			c.Status(http.StatusInternalServerError)
			return
		}
		client := &http.Client{
			Timeout: time.Second * defaultTimeOutInSeconds,
		}
		hotelService := hotel_service.NewHotelService(logger, client, c)
		suppliersDataFilePath := filepath.Join(wd, "internal", "data", suppliersDataFileName)
		hotelDataFilePath := filepath.Join(wd, "internal", "data", hotelsDataFileName)
		fetchedSources, err := hotelService.UpdateHotelsFromSuppliers(suppliersDataFilePath, hotelDataFilePath)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Hotel data updated successfully", "sources": fetchedSources})
	}
}
