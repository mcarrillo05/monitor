package main

import (
	"flag"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/mcarrillo05/monitor/controller"
)

var view = ""

func init() {
	flagView := flag.String("view", "", "path to frontend view")
	flag.Parse()
	view = *flagView
}

func main() {
	router := gin.Default()
	router.StaticFile("/", view+"/index.html")
	router.Use(static.Serve("/", static.LocalFile(view, false))) //Serving static files
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
