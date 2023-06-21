package handlers

import (
	"dashboard-service/internal/config"
	"fmt"
	"math"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"dashboard-service/internal/domain/memory/models"
)

var memoryCollection *mongo.Collection = config.GetCollection(config.DB, "memory")

func FetchMemoryMetrics(c echo.Context) error {
	ctx := c.Request().Context()
	// Retrieve all items from the collection
	cursor, err := memoryCollection.Find(ctx, bson.D{})
	if err != nil {
		fmt.Println(err)
		c.NoContent(http.StatusInternalServerError)
	}
	defer cursor.Close(ctx)

	var items []models.MemoryModel

	for cursor.Next(ctx) {
		var result models.MemoryModel
		err := cursor.Decode(&result)
		if err != nil {
			fmt.Println(err)
			c.NoContent(http.StatusInternalServerError)
		}

		parsedResult := models.MemoryModel{
			ID:        result.ID,
			DateTime:  result.DateTime,
			Committed: bytesToMegabytes(result.Committed),
			Used:      bytesToMegabytes(result.Used),
			Total:     bytesToMegabytes(result.Total),
		}
		items = append(items, parsedResult)
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	// Send the response
	return c.JSON(http.StatusOK, items)
}

func bytesToMegabytes(bytes float64) float64 {
	megabytes := bytes / (1024 * 1024)
	rounded := math.Round(megabytes*100) / 100
	return rounded
}
