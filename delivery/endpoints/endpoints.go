package endpoints

import (
	"github.com/gin-gonic/gin"
	"planner/delivery/rest"
)

func EndPoints(r *gin.Engine, c *rest.Service) {
	r.POST("/planner/add", c.CreateTasks)
	r.GET("/planner/get", c.GetTasks)
	r.GET("/planner/get/:id", rest.GetTaskByID)
	r.PUT("/planner/:id", c.UpdateTasks)
	r.DELETE("/planner/:id", c.DelTask)
}
