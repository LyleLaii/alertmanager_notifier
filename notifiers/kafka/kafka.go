package kafka

import (
	"alertmanager_notifier/config"
	"alertmanager_notifier/log"
	"alertmanager_notifier/metrics"
	"alertmanager_notifier/notifiers"
	"alertmanager_notifier/template"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

const executorName = "kafkaExecutor"

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
	ms := k.GenerateMessage(am, k.tmpl)
	for _, m := range ms {
		jsonstr, _ := json.Marshal(m)
		k.logger.Debug(executorName,
			fmt.Sprintf("Send msg to %s: %v", k.name, string(jsonstr)))
		msg := &sarama.ProducerMessage{
			// Key:   sarama.StringEncoder(),
			Topic: k.config.Topic,
			Value: sarama.ByteEncoder(jsonstr),
		}
		metrics.CountVecNotifier.WithLabelValues(executorName).Inc()
		s := time.Now()
		_, _, err := k.client.SendMessage(msg)
		elapsed := time.Since(s)
		if err != nil {
			k.logger.Error(executorName, err)
			metrics.CountVecErrorNotifier.WithLabelValues(executorName).Inc()
		}
		k.logger.Debug(executorName, fmt.Sprintf("It took %v to send msg", elapsed))
		metrics.HistogramVecNofityDuration.WithLabelValues(executorName).Observe(float64(elapsed.Milliseconds()))
	}
}

// GenerateMessage generate alert message
func (k *Kafka) GenerateMessage(am *notifiers.AlertMessage, tmpl *template.Template) (messages []map[string]string) {
	var m map[string]string
	for index, _ := range am.AlertInfo.Alerts {
		m = make(map[string]string)
		for key, v := range k.config.MsgContent {
			value, err := tmpl.ParseTmplString(index, v, am)
			if err != nil {
				k.logger.Error(executorName, err)
			}
			m[key] = value
		}

		messages = append(messages, m)
	}
	return
}
