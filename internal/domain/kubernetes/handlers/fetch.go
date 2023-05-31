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

	// Group the pods by pod_name
	//groupedPods := make(map[string][]KubernetesResponse)
	//for _, pod := range items {
	//	result := KubernetesResponse{
	//		PodName:   pod.PodName,
	//		PodStatus: pod.PodStatus,
	//		Metrics:   pod.Containers,
	//	}
	//	groupedPods[pod.PodName] = append(groupedPods[pod.PodName], result)
	//}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	// Send the response
	return c.JSON(http.StatusOK, items)
}
