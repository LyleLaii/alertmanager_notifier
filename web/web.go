package web

import (
	"alertmanager_notifier/log"
	"alertmanager_notifier/pkg/static"
	v1 "alertmanager_notifier/web/api/v1"
	"embed"
	"github.com/gin-gonic/gin"
)

//go:embed dist
var dist embed.FS

func RegisterAPI(r gin.IRouter, logger log.Logger) {
	r.POST("/login", v1.Login(logger))
	user := r.Group("/user")
	{
		user.GET("/listpage", v1.GetUserListPage(logger))
		user.POST("/remove", v1.DeleteUser(logger))
		user.POST("/batchremove", v1.DeleteUserBatch(logger))
		user.POST("/edit", v1.EditUser(logger))
		user.POST("/add", v1.AddUser(logger))
		user.GET("/alluid", v1.GetAllUID(logger))
		user.GET("/info", v1.GetUserInfo(logger))
	}
	receiverInfo := r.Group("/receiverInfo")
	{
		receiverInfo.GET("/listpage", v1.GetReceiverInfoListPage(logger))
		receiverInfo.POST("/remove", v1.DeleteReceiverInfo(logger))
		receiverInfo.POST("/batchremove", v1.DeleteReceiverInfoBatch(logger))
		receiverInfo.POST("/edit", v1.EditReceiverInfo(logger))
		receiverInfo.POST("/add", v1.AddReceiverInfo(logger))
	}
	rota := r.Group("/rota")
	{
		rota.GET("/listpage", v1.GetRotaListPage(logger))
		rota.POST("/remove", v1.DeleteRota(logger))
		rota.POST("/batchremove", v1.DeleteRotaBatch(logger))
		rota.POST("/edit", v1.EditRota(logger))
		rota.POST("/add", v1.AddRota(logger))
	}
}


func RegisterUI(r *gin.Engine, logger log.Logger) {
	r.Use(static.Serve("/", static.EmbedFolder(dist, "dist")))

	//r.StaticFS("/file",http.FS(dist))

	r.GET("/login", func(c *gin.Context) {
		c.Request.URL.Path = "/"
		r.HandleContext(c)
	})

	r.GET("/userinfo", func(c *gin.Context) {
		c.Request.URL.Path = "/"
		r.HandleContext(c)
	})

	r.GET("/userrota", func(c *gin.Context) {
		c.Request.URL.Path = "/"
		r.HandleContext(c)
	})

	r.GET("/receiverinfo", func(c *gin.Context) {
		c.Request.URL.Path = "/"
		r.HandleContext(c)
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}