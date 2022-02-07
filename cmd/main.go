package main

import (
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

	// MAP VARIABLES INTO MAP
	variables := map[string]string{
		"EXEC_MODE": execMode,
	}

	// SAVE VARIABLES INSIDE GIN CONTEXT
	r.Use(enviromentMiddleware(variables))

	// HTTP SERVER ROUTES
	r.GET("/heartbeat", pkg.HeartbeatEndpoint)

	// HEARTBEAT POSTING
	cron := gocron.NewScheduler(time.Local)
	pkg.HeartbeatCron(cron, appCron, app, execMode)
	cron.StartAsync()

	// HTTP SERVER
	r.Run(":" + listenPort)
}

func enviromentMiddleware(variables map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for key, value := range variables {
			c.Set(key, value)
			c.Next()
		}
	}
}
