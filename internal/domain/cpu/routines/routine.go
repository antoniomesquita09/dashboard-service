package routines

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"dashboard-service/internal/config"
	"dashboard-service/internal/domain/cpu/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CPU struct {
	Percentage float64 `json:"percentage"`
}

var cpuCollection *mongo.Collection = config.GetCollection(config.DB, "cpu")

func MakeCPURoutine(seconds int64) {
	for {
		fetchCPUMetrics()

		// Wait for 5 seconds before making the next API call
		time.Sleep(time.Duration(seconds) * time.Second)
	}
}

func fetchCPUMetrics() {
	response, err := http.Get(config.EnvGatewayServiceURI() + "/jmx/cpu")
	if err != nil {
		fmt.Println("Error while trying to get jmx CPU metrics")
	}
	defer response.Body.Close()

	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error while trying to read CPU metrics")
	}

	// Parse the response JSON
	cpuResponse := CPU{}
	err = json.Unmarshal(body, &cpuResponse)
	if err != nil {
		fmt.Println("Error while trying to unmarshal CPU metrics")
	}

	cpu := models.CPUModel{
		ID:         primitive.NewObjectID(),
		Percentage: cpuResponse.Percentage,
		DateTime:   formattedTime,
	}

	result, err := cpuCollection.InsertOne(context.TODO(), cpu)
	if err != nil {
		fmt.Println("Error inserting cpu document to mongoDb:", err)
	}

	fmt.Println("Successfully inserted cpu document:", result.InsertedID)

}
