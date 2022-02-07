package pkg

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dmonteroh/distributed-resource-collector/internal"
)

func HeartbeatEndpoint(c *gin.Context) {
	c.JSON(200, internal.GetServerStats())
}

func HeartbeatDebugEndpoint(c *gin.Context) {
	c.JSON(200, internal.GetServerStats().String())
}

func RecoverHeartbeat() {
	if r := recover(); r != nil {
		fmt.Println("Recovered from ", r)
	}
}

func SendHeartbeat(url string, execMode string) {
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
