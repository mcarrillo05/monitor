package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mcarrillo05/monitor/model"
)

type Agent struct {
	model.Agent
	model.SNMPObject
}

//AddAgent registers an agent recieved in JSON request.
func AddAgent(c *gin.Context) {
	var agent model.Agent
	if err := c.Bind(&agent); err != nil {
		JSONError(http.StatusBadRequest, err, c)
		return
	}
	err := agent.Add()
	if err != nil {
		status := http.StatusInternalServerError
		if err == model.ErrDuplicate {
			status = http.StatusConflict
		}
		JSONError(status, err, c)
		return
	}
	c.JSON(http.StatusOK, agent)
}

//DeleteAgent deletes an agent using param ip from request.
func DeleteAgent(c *gin.Context) {
	agent := model.Agent{
		IP: c.Param("ip"),
	}
	err := agent.Delete()
	if err != nil {
		status := http.StatusInternalServerError
		if err == model.ErrNotFound {
			status = http.StatusNotFound
		}
		JSONError(status, err, c)
		return
	}
	c.JSON(http.StatusOK, agent)
}

//GetAgent returns information about agent using param ip from request.
func GetAgent(c *gin.Context) {
	agent := model.Agent{
		IP: c.Param("ip"),
	}
	err := agent.Get()
	if err != nil {
		status := http.StatusInternalServerError
		if err == model.ErrNotFound {
			status = http.StatusNotFound
		}
		JSONError(status, err, c)
		return
	}
	var resp interface{}
	if c.Query("query") != "" {
		object, err := agent.GetOID(c.Query("query"))
		if err != nil {
			status := http.StatusInternalServerError
			switch err {
			case model.ErrNotFound:
				status = http.StatusNotFound
			case model.ErrInvalidResource:
				status = http.StatusBadRequest
			}
			JSONError(status, err, c)
			return
		}
		var a Agent
		a.Agent = agent
		a.SNMPObject = object
		resp = a
	} else {
		resp = agent
	}
	c.JSON(http.StatusOK, resp)
}

//GetAgents returns all agents registered.
func GetAgents(c *gin.Context) {
	var agents []model.Agent
	agents, err := model.GetAllAgents()
	if err != nil {
		JSONError(http.StatusInternalServerError, err, c)
		return
	}
	c.JSON(http.StatusOK, agents)
}
