package shell

import (
	"alertmanager_notifier/config"
	"alertmanager_notifier/log"
	"alertmanager_notifier/notifiers"
	"alertmanager_notifier/template"
	"encoding/json"
	"testing"
	"time"
)

func TestShell(t *testing.T) {

	shellConfig := config.ShellConfig{
		Command: "echo",
		//Args:    []string{"test", "{{ .Receiver }}", "{{ $labels.instance }}", "{{ uuid }}", "{{ uuid32 }}", "{{ template \"test\" }}", "{{ template \"test1\" $alert.Status }}"},
		Args: []string{"test", "{{ .Receiver }}", "{{ $labels.instance }}", "{{ uuid }}"},
		//StaticArgs: map[string]string{"test": "a"},
	}

	tmpl, _ := template.GenTempInstance("")
	shellInstance := Shell{name: "test", config: &shellConfig, tmpl: tmpl, logger: log.NewNopLogger()}

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

	t.Run("simple_test", func(t *testing.T) {
		shellInstance.Notify(&testAm)
	})

	//t.Run("template_test", func(t *testing.T) {
	//  r := shellInstance.parseArgs(&testAm)
	//  fmt.Printf("%+v", r)
	//})

}
