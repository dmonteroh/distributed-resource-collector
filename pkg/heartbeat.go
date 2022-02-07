package pkg

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"

	"github.com/dmonteroh/distributed-resource-collector/internal"
)

func HeartbeatEndpoint(c *gin.Context) {
	execMode, ok := c.MustGet("EXEC_MODE").(string)
	if !ok {
		println("Error: EXEC_MODE not available in Middleware")
	} else {
		if execMode == "DEBUG" {
			c.JSON(200, internal.GetServerStats().String())
		} else {
			c.JSON(200, internal.GetServerStats())
		}
	}
}

func RecoverHeartbeat() {
	if r := recover(); r != nil {
		fmt.Println("Recovered from ", r)
	}
}

func sendHeartbeat(url string, execMode string) {
	defer RecoverHeartbeat()
	body := internal.GetServerStats()
	if execMode == "DEBUG" {
		fmt.Println("DEUBG MODE - POST")
		fmt.Println(body.String())
	}
	res, err := http.Post(url, "application/json", bytes.NewBuffer(body.Marshal()))
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	fmt.Println(res.Body)
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

func HeartbeatCron(cron *gocron.Scheduler, seconds int, app string, execMode string) {

	defer recoverCron()
	cronRes, cronErr := cron.Every(seconds).Seconds().Do(sendHeartbeat, app, execMode)
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
