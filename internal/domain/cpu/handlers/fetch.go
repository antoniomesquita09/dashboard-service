package handlers

import (
	"dashboard-service/internal/config"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"math"
	"net/http"

	"dashboard-service/internal/domain/cpu/models"
)

var cpuCollection *mongo.Collection = config.GetCollection(config.DB, "cpu")

// FetchCPUMetrics fetches CPU metrics from gateway-service
func FetchCPUMetrics(c echo.Context) error {
	ctx := c.Request().Context()
	// Retrieve all items from the collection
	cursor, err := cpuCollection.Find(ctx, bson.D{})
	if err != nil {
		fmt.Println(err)
		c.NoContent(http.StatusInternalServerError)
	}
	defer cursor.Close(ctx)

	var items []models.CPUModel

	for cursor.Next(ctx) {
		var result models.CPUModel
		err := cursor.Decode(&result)
		if err != nil {
			fmt.Println(err)
			c.NoContent(http.StatusInternalServerError)
		}
		multiplied := result.Percentage * 100
		roundedPercentage := math.Round(multiplied*100) / 100

		parsedResult := models.CPUModel{
			ID:         result.ID,
			DateTime:   result.DateTime,
			Percentage: roundedPercentage,
		}
		items = append(items, parsedResult)
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	// Send the response
	return c.JSON(http.StatusOK, items)
}
