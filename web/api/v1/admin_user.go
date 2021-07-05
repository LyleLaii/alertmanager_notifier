package v1

import (
	"alertmanager_notifier/log"
	"alertmanager_notifier/metrics"
	"alertmanager_notifier/modules/db"
	e "alertmanager_notifier/pkg/err"
	"alertmanager_notifier/pkg/utils"
	"fmt"
	"github.com/spf13/viper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userinfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login check user info for login
func Login(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var u userinfo
		var status bool
		if err := c.BindJSON(&u); err != nil {
			metrics.GaugeVecAPIError.WithLabelValues("login").Inc()
			logger.Warn("ADMINUSER", fmt.Sprintf("login bind json find err: %v", err))
			c.JSON(http.StatusBadRequest, gin.H{
				"code": e.INVALID_PARAMS,
				"msg":  e.GetMsg(e.INVALID_PARAMS),
			})
			return
		}
		if viper.IsSet("static_login") {
			passwordEncodeStr := viper.GetString("static_login")
			status = utils.Check(u.Password, passwordEncodeStr)
		} else {
			status = db.CheckAdmin(u.Username, u.Password)
		}

		if status {
			c.JSON(http.StatusOK, gin.H{
				"code": e.SUCCESS,
				"msg":  e.GetMsg(e.SUCCESS),
				"data": gin.H{
					"username": u.Username,
				},
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": e.ERROR_AUTH,
				"msg":  e.GetMsg(e.ERROR_AUTH),
			})
		}
	}
}
