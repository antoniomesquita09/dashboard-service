package main

import (
	"dashboard-service/internal/config"
	cpuHandlers "dashboard-service/internal/domain/cpu/handlers"
	kubernetesHandlers "dashboard-service/internal/domain/kubernetes/handlers"
	memoryRoutine "dashboard-service/internal/domain/memory/routines"
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()
	app.Use(echoprometheus.NewMiddleware("myapp"))

	go func() {
		metrics := echo.New()                                // this Echo will run on separate port 8081
		metrics.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics
		if err := metrics.Start(":8082"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	app.GET("/memory", cpuHandlers.FetchCPUMetrics)
	app.GET("/kubernetes", kubernetesHandlers.FetchKubernetesMetrics)

	// Connect to mongo database
	config.ConnectDB()

	// Start a Goroutine to make API calls every 5 seconds
	go memoryRoutine.MakeMemoryRoutine()

	app.Logger.Fatal(app.Start(":8081"))
}
