package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

// CPUMetric struct from gateway-service metric proxy
type CPUMetric struct {
	Percentage float64 `json:"percentage"`
}

// FetchCPUMetrics fetches CPU metrics from gateway-service
func FetchCPUMetrics(c echo.Context) error {
	// Make a GET request to another service running on localhost:8080
	response, err := http.Get("http://localhost:8080/jmx/cpu")
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// Parse the response JSON
	cpuMetric := CPUMetric{}
	err = json.Unmarshal(body, &cpuMetric)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, cpuMetric)
}
