package middleware

import (
	"net/http"
	"runtime/debug"
	"strings"

	"alertmanager_notifier/log"

	"github.com/gin-gonic/gin"
)

// GinRecovery catch panic and logs
func GinRecovery(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				DebugStack := ""
				for _, v := range strings.Split(string(debug.Stack()), "\n") {
					DebugStack += v + "<br>"
				}
				logger.Error("gin", DebugStack)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
