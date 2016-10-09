package jobs

import (
	"github.com/bamzi/jobrunner"
	"github.com/gin-gonic/gin"
)

func StartWebInterface() {
	routes := gin.Default()
	routes.GET("/jobs/json", jobJSON)
	routes.Run(":5565")
}

func jobJSON(c *gin.Context) {
	c.JSON(200, jobrunner.StatusJson())
}
