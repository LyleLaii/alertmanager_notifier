package kafka

import (
	"alertmanager_notifier/config"
	"alertmanager_notifier/log"
	"alertmanager_notifier/notifiers"
	"alertmanager_notifier/template"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)


func TestKafka(t *testing.T) {


	kafkaConfig := config.KafkaConfig{
		Topic:        "test",
		Broker:       []string{"127.0.0.1:9092"},
		BrokerConfig: config.BrokerConfig{},
		MsgContent: map[string]string{"receiver": "{{ .Receiver }}", "instance": "{{ $labels.instance }}"},
	}

	tmpl, _ := template.GenTempInstance("")
	kafkaInstance, _ := New("test", &kafkaConfig, tmpl, log.NewNopLogger())

	alertDataString := `{
        "status":"firing",
        "receiver":"test",
        "alerts":[
            {
                "status":"firing",
                "labels":{
                    "instance":"1.1.1.1",
                    "namespace":"ns1",
                    "severity": "error"
                },
                "annotations":{
                    "message":"this is a test"
                }
            },
            {
                "status":"firing",
                "labels":{
                    "instance":"1.1.1.2",
                    "namespace":"ns2"
                },
                "annotations":{
                    "message":"this is a test"
                }
            }
        ]
    }`

	var alert notifiers.AlertWebhookInfo
	json.Unmarshal([]byte(alertDataString), &alert)

	testAm := notifiers.AlertMessage{
		Receiver:     "test",
		ReceiverDate: time.Time{},
		AlertInfo:    alert,
	}

	//fmt.Printf("%+v", testAm)

	t.Run("GenerateMessage_test", func(t *testing.T) {
		messages := kafkaInstance.GenerateMessage(&testAm, kafkaInstance.tmpl)
		fmt.Printf("%+v", messages)
		jsonstr, _ := json.Marshal(messages)
		fmt.Printf("%s", jsonstr)
	})

	t.Run("Notify_test", func(t *testing.T) {
		kafkaInstance.Notify(&testAm)
	})

	//t.Run("template_test", func(t *testing.T) {
	//  r := shellInstance.parseArgs(&testAm)
	//  fmt.Printf("%+v", r)
	//})

}
