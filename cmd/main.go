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

	// HTTP SERVER ROUTES
	if execMode == "DEBUG" {
		println("DEUBG MODE")
		r.GET("/heartbeat", pkg.HeartbeatDebugEndpoint)
	} else {
		r.GET("/heartbeat", pkg.HeartbeatEndpoint)
	}

	// HEARTBEAT POSTING
	//go heartbeatCron(appCron, app, execMode)
	// go func() {
	// 	cron := gocron.NewScheduler(time.UTC)
	// 	cron.Every(appCron).Seconds().Do(pkg.SendHeartbeat(app, execMode))
	// 	cron.StartBlocking()
	// }()
	cron := gocron.NewScheduler(time.Local)
	heartbeatCron(cron, appCron, app, execMode)
	cron.StartAsync()

	// HTTP SERVER
	r.Run(":" + listenPort)
}

func recoverCron() {
	if r := recover(); r != nil {
		fmt.Println("Recovered from ", r)
	}
}

func debugCron(cronRes *gocron.Job) {
	fmt.Println("RUN COUNT:", cronRes.RunCount())
	fmt.Println("NEXT RUN: ", cronRes.NextRun())
}

func heartbeatCron(cron *gocron.Scheduler, seconds int, app string, execMode string) {

	defer recoverCron()
	cronRes, cronErr := cron.Every(seconds).Seconds().Do(pkg.SendHeartbeat, app, execMode)
	if cronErr != nil {
		panic(cronErr)
	}
	go func() {
		if execMode == "DEBUG" {
			_, cronErrDebug := cron.Every(seconds).Seconds().Do(debugCron, cronRes)
			if cronErrDebug != nil {
				panic(cronErr)
			}
		}
	}()
}
