package pkg

import (
	"github.com/gin-gonic/gin"

	"github.com/dmonteroh/distributed-resource-collector/internal"
)

func HeartbeatEndpoint(c *gin.Context) {
	c.JSON(200, internal.GetServerStats())
}
