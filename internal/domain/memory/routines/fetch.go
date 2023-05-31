package routine

import (
	"context"
	"dashboard-service/internal/config"
	"dashboard-service/internal/domain/memory/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Memory struct {
	Used      float64 `json:"used"`
	Committed float64 `json:"committed"`
	Total     float64 `json:"total"`
}

var memoryCollection *mongo.Collection = config.GetCollection(config.DB, "memory")

func MakeMemoryRoutine(seconds int64) {
	for {
		fetchMemoryMetrics()

		// Wait for 5 seconds before making the next API call
		time.Sleep(time.Duration(seconds) * time.Second)
	}
}

func fetchMemoryMetrics() {
	// Make the API call
	response, err := http.Get(config.EnvGatewayServiceURI() + "/jmx/memory")
	if err != nil {
		fmt.Println("Error making API call:", err)
	} else {
		fmt.Println("API response:", response.Status)

		defer response.Body.Close()

		// Read the response body
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error making API call:", err)
		}

		// Parse the response JSON
		memoryResponse := Memory{}
		err = json.Unmarshal(body, &memoryResponse)
		if err != nil {
			fmt.Println("Error unmarshall response:", err)
		}

		fmt.Println("API response body:", memoryResponse)

		memory := models.MemoryModel{
			ID:        primitive.NewObjectID(),
			Used:      memoryResponse.Used,
			Committed: memoryResponse.Committed,
			Total:     memoryResponse.Total,
		}

		result, err := memoryCollection.InsertOne(context.TODO(), memory)
		if err != nil {
			fmt.Println("Error inserting memory document to mongoDb:", err)
		}

		fmt.Println("Successfully inserted memory document:", result.InsertedID)
	}
}
