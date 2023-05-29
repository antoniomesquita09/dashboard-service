package routine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Memory struct {
	Used      float64 `json:"used"`
	Committed float64 `json:"committed"`
	Total     float64 `json:"total"`
}

func MakeMemoryRoutine() {
	for {
		// Make the API call
		response, err := http.Get("http://localhost:8080/jmx/memory") // Replace with your API URL
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
			memory := Memory{}
			err = json.Unmarshal(body, &memory)
			if err != nil {
				fmt.Println("Error unmarshall response:", err)
			}
			fmt.Println("API response body:", memory)
			// Add your code to process the API response here
		}

		// Wait for 5 seconds before making the next API call
		time.Sleep(5 * time.Second)
	}
}
