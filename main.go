package main

import (
	"api-test/common"
	module "api-test/modules"
	"api-test/v1/students"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"
)

func main() {
	r := gin.Default()
	r.Use(common.LoggerInitializer())
	r.Use(cors.Default())
	module.DbConnection()
	apiRoutes := r.Group("/v1")
	studentRoute := apiRoutes.Group("/student/")
	{
		studentRoute.POST("/:id", students.GetStudent(module.DB))
	}
	r.NoRoute(func(c *gin.Context) {
		common.WriteToLog(c, c.Request.RequestURI+" doesnt exist on the serve", zapcore.InfoLevel, "INFO")
		c.SecureJSON(http.StatusNotImplemented, gin.H{"Error": 501, "msg": "Resource couldn't be found!"})
	})
	gin.SetMode(gin.ReleaseMode)
	r.Run(":6060")
}
