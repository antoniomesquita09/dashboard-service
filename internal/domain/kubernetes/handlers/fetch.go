package handlers

import (
	"dashboard-service/internal/config"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"dashboard-service/internal/domain/kubernetes/models"
)

//type KubernetesResponse struct {
//	PodName   string                  `json:"pod_name"`
//	PodStatus string                  `json:"pod_status"`
//	Metrics   []models.ContainerModel `json:"metrics"`
//}

var kubernetesCollection *mongo.Collection = config.GetCollection(config.DB, "kubernetes")

func FetchKubernetesMetrics(c echo.Context) error {
	ctx := c.Request().Context()
	// Retrieve all items from the collection
	cursor, err := kubernetesCollection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var items []models.PodModel

	for cursor.Next(ctx) {
		var result models.PodModel
		err := cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, result)
	}

	// Create an array of arrays to store the grouped objects
	var groups [][]models.PodModel

	// Iterate over the objects and group them by pod_name
	for _, obj := range items {
		// Check if a group with the same pod_name already exists
		found := false
		for i, group := range groups {
			if group[0].PodName == obj.PodName {
				groups[i] = append(groups[i], obj)
				found = true
				break
			}
		}

		// If a group doesn't exist, create a new one
		if !found {
			groups = append(groups, []models.PodModel{obj})
		}
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	// Send the response
	return c.JSON(http.StatusOK, groups)
}
