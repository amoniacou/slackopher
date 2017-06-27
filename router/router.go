package router

import (
	"github.com/amoniacou/slackopher/controllers"
	"github.com/gin-gonic/gin"
)

func WebServer() *gin.Engine {
	r := gin.Default()
	r.GET("/stats", controllers.StatsHandler)
	return r
}
