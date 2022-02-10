package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"

	"github.com/dmonteroh/distributed-resource-collector/internal"
)

func HeartbeatEndpoint(c *gin.Context) {
	defer internal.RecoverEndpoint(c)
	execMode := c.MustGet("EXEC_MODE").(string)
	if execMode == "DEBUG" {
		c.JSON(200, internal.GetServerStats().String())
	} else {
		c.JSON(200, internal.GetServerStats())
	}
}

func recoverHeartbeat() {
	if err := recover(); err != nil {
		fmt.Println("RECOVER HEARTBEAT")
		msg := "Error: [Recovered] "
		switch errType := err.(type) {
		case string:
			msg += err.(string)
		case *json.SyntaxError:
			msg += errType.Error()
		default:
		}
		fmt.Println(msg)
	}
}

func sendHeartbeat(url string, execMode string) {
	defer recoverHeartbeat()
	body := internal.GetServerStats()
	if execMode == "DEBUG" {
		fmt.Println("DEUBG MODE - POST")
		fmt.Println(body.String())
	}
	res, err := http.Post(url, "application/json", bytes.NewBuffer(body.Marshal()))
	if err != nil {
		fmt.Println(res)
		fmt.Println(err)
		panic(err)
	}

	defer res.Body.Close()
	fmt.Println(res.Body)
}

func recoverCron() {
	if err := recover(); err != nil {
		fmt.Println("RECOVER HEARTBEAT")
		msg := "Error: [Recovered] "
		switch errType := err.(type) {
		case string:
			msg += err.(string)
		case *json.SyntaxError:
			msg += errType.Error()
		default:
		}
		fmt.Println(msg)
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
