package main

import (
	"dashboard-service/internal/config"
	cpuHandlers "dashboard-service/internal/domain/cpu/handlers"
	kubernetesHandlers "dashboard-service/internal/domain/kubernetes/handlers"
	memoryRoutine "dashboard-service/internal/domain/memory/routines"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/memory", cpuHandlers.FetchCPUMetrics)
	e.GET("/kubernetes", kubernetesHandlers.FetchKubernetesMetrics)

	// Start a Goroutine to make API calls every 5 seconds
	go memoryRoutine.MakeMemoryRoutine()

	config.ConnectDB()

	e.Logger.Fatal(e.Start(":8081"))
}
