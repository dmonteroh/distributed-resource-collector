package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/dmonteroh/distributed-resource-collector/pkg"
)

func main() {
	r := gin.Default()
	println(os.Getenv("EXEC_MODE"))
	if os.Getenv("EXEC_MODE") == "DEBUG" {
		println("DEUBG MODE")
		r.GET("/heartbeat", pkg.HeartbeatDebugEndpoint)
	} else {
		r.GET("/heartbeat", pkg.HeartbeatEndpoint)
	}

	r.Run(":" + os.Getenv("INTERNAL_PORT"))
}
