package middleware

import (
	"fmt"
	"time"

	"alertmanager_notifier/log"

	"github.com/gin-gonic/gin"
)

func GinLogger(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		latencyTime := time.Since(startTime)
		reqMethod := c.Request.Method
		reqURL := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		logger.Info("Gin", fmt.Sprintf("%s|%s|%s|%v|%s|%s",clientIP,reqMethod,reqURL,statusCode,latencyTime,userAgent))
		//logger.WithFields(logrus.Fields{
		//	"status_code":  statusCode,
		//	"latency_time": latencyTime,
		//	"client_ip":    clientIP,
		//	"req_method":   reqMethod,
		//	"req_uri":      reqURL,
		//	"module":       "Gin",
		//}).Info()
	}
}
