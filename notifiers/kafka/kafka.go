package kafka

import (
	"alertmanager_notifier/config"
	"alertmanager_notifier/log"
	"alertmanager_notifier/metrics"
	"alertmanager_notifier/notifiers"
	"alertmanager_notifier/pkg/utils"
	"alertmanager_notifier/template"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

// Kafka kafka notify instance
type Kafka struct {
	name   string
	config *config.KafkaConfig
	client sarama.SyncProducer
	tmpl *template.Template
	logger log.Logger
}

// New return a new kafka instance
func New(name string, c *config.KafkaConfig, tmpl *template.Template, logger log.Logger) (*Kafka, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          //赋值为-1：这意味着producer在follower副本确认接收到数据后才算一次发送完成。
	config.Producer.Partitioner = sarama.NewRandomPartitioner //写到随机分区中，默认设置8个分区
	config.Producer.Return.Successes = true
	if c.BrokerConfig.DialTimeout == 0 {
		config.Net.DialTimeout = 30 * time.Second
	} else {
		config.Net.DialTimeout = c.BrokerConfig.DialTimeout
	}
	client, err := sarama.NewSyncProducer(c.Broker, config)
	if err != nil {
		return nil, err
	}
	return &Kafka{
		name:   name,
		config: c,
		client: client,
		tmpl: tmpl,
		logger: logger,
	}, nil
}

// Notify Send notify
func (k *Kafka) Notify(am *notifiers.AlertMessage) {
	ms := k.GenerateMessage(am)
	for _, m := range ms {
		jsonstr, _ := json.Marshal(m)
		k.logger.Debug("kafkaserver",
			fmt.Sprintf("Send msg to %s: %v", k.name, string(jsonstr)))
		msg := &sarama.ProducerMessage{
			Key:   sarama.StringEncoder("key_" + m["messageid"]),
			Topic: k.config.Topic,
			Value: sarama.ByteEncoder(jsonstr),
		}
		metrics.CountVecNotifier.WithLabelValues("kafka").Inc()
		s := time.Now()
		_, _, err := k.client.SendMessage(msg)
		elapsed := time.Since(s)
		if err != nil {
			k.logger.Error("kafkaserver", err)
			metrics.CountVecErrorNotifier.WithLabelValues("kafka").Inc()
		}
		k.logger.Debug("kafkaserver", fmt.Sprintf("It took %v to send msg", elapsed))
		metrics.HistogramVecNofityDuration.WithLabelValues("kafka").Observe(float64(elapsed.Milliseconds()))
	}
}

// GenerateMessage generate alert message
// 占位符替换，待重构
func (k *Kafka) GenerateMessage(am *notifiers.AlertMessage) (messages []map[string]string) {
	var m map[string]string
	for _, alert := range am.AlertInfo.Alerts {
		m = make(map[string]string)
		var alertinfo string
		if am.AlertInfo.Status == "firing" {
			alertinfo = fmt.Sprintf("%s [Firing since %s] %s",
				k.config.Info,
				utils.TransTimeZoneAuto(alert.StartsAt),
				alert.Annotations["message"])
		}
		if am.AlertInfo.Status == "resolved" {
			alertinfo = fmt.Sprintf("%s [Resolved at %s] %s",
				k.config.Info,
				utils.TransTimeZoneAuto(alert.EndsAt),
				alert.Annotations["message"])
		}
		m["receiver"] = alertinfo
		messages = append(messages, m)
	}
	return
}
