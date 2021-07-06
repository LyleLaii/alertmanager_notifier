package config

import (
	"time"
)

// ReceiversMap receiver map
var ReceiversMap = make(map[string]interface{})

// Receiver receiver config
type Receiver struct {
	Name           string          `mapstructure:"name"`
	ReceiverType   string          `mapstructure:"receiver_type"`
	KafkaConfigs   []KafkaConfig   `mapstructure:"kafka_config,omitempty"`
	WebhookConfigs []WebhookConfig `mapstructure:"webhook_config,omitempty"`
	WeChatConfigs  []WeChatConfig  `mapstructure:"wechat_config,omitempty"`
	ShellConfigs   []ShellConfig   `mapstructure:"shell_config,omitempty"`
}

// Notifier notifier instance
type Notifier interface {
	Notify(interface{}) error
	GenerateMessage(interface{}) interface{}
}

// KafkaConfig kafka config
type KafkaConfig struct {
	Topic        string            `mapstructure:"topic"`
	Broker       []string          `mapstructure:"Broker"`
	BrokerConfig BrokerConfig      `mapstructure:"broker_config,omitempty"`
	MsgContent   map[string]string `mapstructure:"msg_content,omitempty"`
}

// BrokerConfig kafka broker client config
type BrokerConfig struct {
	DialTimeout time.Duration `mapstructure:"dial_timeout,omitempty"`
}

// WebhookConfig webhook config simple
type WebhookConfig struct {
	URL string `mapstructure:"url"`
}

// WeChatConfig WorkWechat config
type WeChatConfig struct {
	CorpID         string `mapstructure:"corpid"`
	AgentID        string `mapstructure:"agentid"`
	CorpSecret     string `mapstructure:"corpsecret"`
	DefaultToUser  string `mapstructure:"default_touser,omitempty"`
	DefaultToParty string `mapstructure:"default_toparty,omitempty"`
	DefaultToTag   string `mapstructure:"default_toparty,omitempty"`
}

type ShellConfig struct {
	Command string `mapstructure:"command"`
	Args []string `mapstructure:"args"`
	//StaticArgs map[string]string `mapstructure:"static_args"`
}