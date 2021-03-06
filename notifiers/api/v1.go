package api

import (
	"alertmanager_notifier/config"
	"alertmanager_notifier/log"
	"alertmanager_notifier/metrics"
	"alertmanager_notifier/notifiers"
	e "alertmanager_notifier/pkg/err"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AlertWebhook(ni []notifiers.NotifyIntegration, logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var awi notifiers.AlertWebhookInfo
		if err := c.BindJSON(&awi); err != nil {
			metrics.GaugeVecAPIError.WithLabelValues("alert").Inc()
			logger.Warn("ALERT", fmt.Sprintf("alertwebhook bind json find err: %v", err))
			c.JSON(http.StatusBadRequest, gin.H{
				"code": e.INVALID_PARAMS,
				"msg":  e.GetMsg(e.INVALID_PARAMS),
			})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{
			"code": e.ACCEPTED,
			"msg":  e.GetMsg(e.ACCEPTED),
		})
		am := &notifiers.AlertMessage{}
		am.ReceiverDate = time.Now()
		am.AlertInfo = awi
		for _, n := range ni {
			// TODO： 先搜索用户，再遍历对应接受通道？
			am.Receiver = notifiers.SearchReceiver(awi.Receiver, am.ReceiverDate.Format(config.TimeLayout), n.ReceiverType)
			if am.Receiver == "" {
				logger.Warn("Notify",
					fmt.Sprintf("Do not send msg cause notify: %s and Receiver: %v and ReceiverType: %v with Date: %s do not find receiver",
						n.Name,
						awi.Receiver,
						n.ReceiverType,
						am.ReceiverDate.Format(config.TimeLayoutTZ)))
				continue
			}
			// TODO: TBD how to send multi alert
			go n.Notify(am)
		}
	}
}

func Register(r gin.IRouter, ni []notifiers.NotifyIntegration, logger log.Logger) {
	r.POST("/alert", AlertWebhook(ni, logger))
}
