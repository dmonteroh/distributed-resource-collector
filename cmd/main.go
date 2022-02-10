package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"

	"github.com/dmonteroh/distributed-resource-collector/internal"
	"github.com/dmonteroh/distributed-resource-collector/pkg"
)

func main() {
	r := gin.Default()

	// ENVIROMENTAL VARIABLES
	execMode := internal.GetEnv("EXEC_MODE", "DEBUG")
	listenPort := internal.GetEnv("INTERNAL_PORT", "8081")
	appProtocol := internal.GetEnv("APP_PROTOCOL", "http")
	appIP := internal.GetEnv("APP_IP", "localhost:8080")
	appUrl := internal.GetEnv("APP_URL", "collector")
	app := strings.Join([]string{appProtocol, strings.Join([]string{appIP, appUrl}, "/")}, "://")
	appCron, _ := strconv.Atoi(internal.GetEnv("APP_CRON", "30"))
	heartbeat, _ := strconv.ParseBool(internal.GetEnv("HEARTBEAT", "true"))

	// MAP VARIABLES INTO MAP
	variables := map[string]string{
		"EXEC_MODE": execMode,
	}

	// SAVE VARIABLES INSIDE GIN CONTEXT
	r.Use(internal.EnviromentMiddleware(variables))

	// HTTP SERVER ROUTES
	r.GET("/heartbeat", pkg.HeartbeatEndpoint)
	r.POST("/latency", pkg.LatencyEndpoint)

	// HEARTBEAT POSTING
	if heartbeat {
		fmt.Println("INITIATE HEARTBEAT SCHEDULER")
		cron := gocron.NewScheduler(time.Local)
		pkg.HeartbeatCron(cron, appCron, app, execMode)
		cron.StartAsync()
	}

	// HTTP SERVER
	r.Run(":" + listenPort)
}
