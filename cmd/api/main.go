package main

import (
	"dashboard-service/internal/config"
	cpuHandlers "dashboard-service/internal/domain/cpu/handlers"
	cpuRoutines "dashboard-service/internal/domain/cpu/routines"
	kubernetesHandlers "dashboard-service/internal/domain/kubernetes/handlers"
	kubernetesRoutines "dashboard-service/internal/domain/kubernetes/routines"
	memoryHandlers "dashboard-service/internal/domain/memory/handlers"
	memoryRoutines "dashboard-service/internal/domain/memory/routines"
	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()

	app.GET("/cpu", cpuHandlers.FetchCPUMetrics)
	app.GET("/memory", memoryHandlers.FetchMemoryMetrics)
	app.GET("/kubernetes", kubernetesHandlers.FetchKubernetesMetrics)

	app.Use(middleware.CORS())

	// Connect to mongo database
	config.ConnectDB()

	var intervalSeconds int64 = 500

	// Start a Goroutine to make API calls every delay seconds
	go memoryRoutines.MakeMemoryRoutine(intervalSeconds)
	go cpuRoutines.MakeCPURoutine(intervalSeconds)
	go kubernetesRoutines.MakeKubernetesRoutine(intervalSeconds)

	app.Logger.Fatal(app.Start(":8081"))
}
