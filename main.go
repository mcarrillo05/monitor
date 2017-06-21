package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mcarrillo05/monitor/controller"
)

func main() {
	router := gin.Default()
	apiRoutes := router.Group("api")
	{
		agent := apiRoutes.Group("agent")
		{
			agent.POST("", controller.AddAgent)
			agent.GET("", controller.GetAgents)
			agent.GET(":ip", controller.GetAgent)
			agent.DELETE(":ip", controller.DeleteAgent)
		}
	}
	router.Run(":8000")
}
