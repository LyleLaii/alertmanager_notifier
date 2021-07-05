package wechat

import (
	"alertmanager_notifier/config"
	"alertmanager_notifier/log"
	"alertmanager_notifier/metrics"
	"alertmanager_notifier/notifiers"
	"alertmanager_notifier/pkg/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// 不可使用，待重构

// WeChat send alert by wechat
type WeChat struct {
	name          string
	config        *config.WeChatConfig
	client        *http.Client
	accessToken   string
	accessTokenAt time.Time
	URL           *url.URL
	logger        log.Logger
}

// New return a wechat instance
func New(name string, c *config.WeChatConfig, logger log.Logger) (*WeChat, error) {
	// TODO: 初始化时是否获取一次Token？
	u, err := url.Parse("https://qyapi.weixin.qq.com/cgi-bin")
	if err != nil {
		logger.Warn("WeChat", fmt.Sprintf("Parse url err: %s", err))
		return nil, err
	}
	return &WeChat{
		name:   name,
		config: c,
		client: &http.Client{},
		URL:    u,
		logger: logger,
	}, nil
}

type token struct {
	AccessToken string `json:"access_token"`
}

func (w *WeChat) getAccessToken() {
	u := *w.URL
	u.Path += "gettoken"
	q := u.Query()
	q.Set("corpid", w.config.CorpID)
	q.Set("corpsecret", w.config.CorpSecret)
	u.RawQuery = q.Encode()
	requestGet, _ := http.NewRequest("Get", u.String(), nil)
	resp, err := w.client.Do(requestGet)
	if err != nil {
		w.logger.Warn("WeChat", fmt.Sprintf("get accesstoken error: %s!", err.Error()))
		return
	}
	defer resp.Body.Close()
	var wechatToken token
	if err := json.NewDecoder(resp.Body).Decode(&wechatToken); err != nil {
		w.logger.Warn("WeChat", fmt.Sprintf("get accesstoken error: %s!", err.Error()))
	}
	w.accessToken = wechatToken.AccessToken
	w.accessTokenAt = time.Now()
}

type respMessage struct {
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
	InvalidUser  string `json:"invaliduser"`
	InvalidParty string `json:"invalidparty"`
	InvalidTag   string `json:"invalidtag"`
}

// Notify send alert message
func (w *WeChat) Notify(am *notifiers.AlertMessage) {
	if w.accessToken == "" || time.Since(w.accessTokenAt) > 2*time.Hour {
		w.getAccessToken()
	}
	u := *w.URL
	u.Path += "message/send"
	q := u.Query()
	q.Set("access_token", w.accessToken)
	u.RawQuery = q.Encode()
	ms := w.GenerateMessage(am)
	for _, m := range ms {
		requestPost, _ := http.NewRequest("POST", u.User.String(), bytes.NewReader(m))
		metrics.CountVecNotifier.WithLabelValues("webhook").Inc()
		resp, err := w.client.Do(requestPost)

		if err != nil {
			metrics.CountVecErrorNotifier.WithLabelValues("webhook").Inc()
			w.logger.Warn("WeChat", fmt.Sprintf("wechat send message error: %s!", err.Error()))
			continue
		}
		defer resp.Body.Close()
		var respMsg respMessage
		// TODO: 重试判断？ webhook
		json.NewDecoder(resp.Body).Decode(&respMsg)
	}
}

// GenerateMessage generate alert message
func (w *WeChat) GenerateMessage(am *notifiers.AlertMessage) (messages [][]byte) {
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
		data["touser"] = am.Receiver
		data["content"] = alertinfo
		data["msgtype"] = "text"
		data["agentid"] = w.config.AgentID

		jsonData, _ := json.Marshal(data)
		messages = append(messages, jsonData)
	}
	return
}
