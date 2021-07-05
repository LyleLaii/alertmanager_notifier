package webhook

import (
	"alertmanager_notifier/config"
	"alertmanager_notifier/log"
	"alertmanager_notifier/metrics"
	"alertmanager_notifier/notifiers"
	"alertmanager_notifier/pkg/utils"
	"alertmanager_notifier/template"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Webhook send alert by webhook
type Webhook struct {
	name   string
	config *config.WebhookConfig
	client *http.Client
	tmpl *template.Template
	logger log.Logger
}

// New return a http instance
// just return a simple client for test
func New(name string, c *config.WebhookConfig, tmpl *template.Template, logger log.Logger) (*Webhook, error) {
	return &Webhook{
		name:   name,
		config: c,
		client: &http.Client{},
		tmpl: tmpl,
		logger: logger,
	}, nil
}

// Notify send alert message
func (w Webhook) Notify(am *notifiers.AlertMessage) {
	ms := w.GenerateMessage(am)
	for _, m := range ms {
		requestPost, err := http.NewRequest("POST", w.config.URL, bytes.NewReader(m))
		metrics.CountVecNotifier.WithLabelValues("webhook").Inc()
		resp, err := w.client.Do(requestPost)
		if err != nil {
			metrics.CountVecErrorNotifier.WithLabelValues("webhook").Inc()
			w.logger.Warn("WebHook", fmt.Sprintf("server post to %s error: %s!", w.config.URL, err.Error()))
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			bodyContent, _ := ioutil.ReadAll(resp.Body)
			metrics.CountVecErrorNotifier.WithLabelValues("webhook").Inc()
			w.logger.Warn("WebHook", fmt.Sprintf("server post to %s status: %d , resdata: %s", w.config.URL, resp.StatusCode, string(bodyContent)))
		} else {
			w.logger.Info("WebHook", fmt.Sprintf("server post to %s status: %d", w.config.URL, resp.StatusCode))
		}
	}
}

// GenerateMessage generate alert message
// 待重构
func (w *Webhook) GenerateMessage(am *notifiers.AlertMessage) (messages [][]byte) {
	var data map[string]interface{}
	for _, alert := range am.AlertInfo.Alerts {
		data = make(map[string]interface{})
		var alertinfo string
		if am.AlertInfo.Status == "firing" {
			alertinfo = fmt.Sprintf("[Firing since %s] %s",
				utils.TransTimeZoneAuto(alert.StartsAt),
				alert.Annotations["message"])
		}
		if am.AlertInfo.Status == "resolved" {
			alertinfo = fmt.Sprintf("[Resolved at %s] %s",
				utils.TransTimeZoneAuto(alert.EndsAt),
				alert.Annotations["message"])
		}
		data["owener"] = am.Receiver
		data["alertinfo"] = alertinfo
		data["time"] = time.Now()
		jsonData, _ := json.Marshal(data)
		messages = append(messages, jsonData)
	}
	return
}
