package main

import (
	"github.com/gin-gonic/gin"

	"github.com/dmonteroh/distributed-resource-collector/pkg"
)

func main() {
	r := gin.Default()

	r.GET("/heartbeat", pkg.HeartbeatEndpoint)

	r.Run(":8080")
}
