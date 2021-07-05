package router

import (
	"alertmanager_notifier/config"
	"alertmanager_notifier/log"
	"alertmanager_notifier/middleware"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

// InitRouter init router config
func InitRouter(logger log.Logger, rConf *config.RunningConfig) *gin.Engine {
	var mode = ""

	switch rConf.RunMode.String() {
	case "dev":
		mode = "debug"
	case "debug":
		mode = "debug"
		gin.DisableConsoleColor()
	case "release":
		mode = "release"
		gin.DisableConsoleColor()
		gin.DefaultWriter = ioutil.Discard
	default:
		gin.DisableConsoleColor()
	}

	gin.SetMode(mode)

	r := gin.Default()
	r.Use(middleware.Cors())
	r.Use(middleware.GinLogger(logger))
	r.Use(middleware.GinRecovery(logger))

	r.GET("/ping", func(c *gin.Context) {
		//输出json结果给调用方
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	return r
}

func AddMiddleware(r gin.IRouter, hfs ...gin.HandlerFunc) {
	r.Use(hfs...)
}