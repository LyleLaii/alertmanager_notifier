package template

import (
	"alertmanager_notifier/notifiers"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)





func TestShell(t *testing.T) {

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

	var alertInfo notifiers.AlertWebhookInfo
	json.Unmarshal([]byte(alertDataString), &alertInfo)

	testData := notifiers.AlertMessage{
		Receiver:     "test",
		ReceiverDate: time.Time{},
		AlertInfo: alertInfo,

	}

	tmpl, err := GenTempInstance("data/custom.tpl")

	t.Run("CreateTemplate", func(t *testing.T) {

		if err != nil {
			t.Errorf("create template get err: %+v", err)
		}
	})

	t.Run("Parse Data", func(t *testing.T) {
		r, err := tmpl.ParseTmplString(1, "{{ .Receiver }}", testData)
		if err != nil {
			t.Fatalf("parse data get err: %+v", err)
		}
		if r != testData.Receiver {
			t.Errorf("parse data failed, except %s get %s",testData.Receiver, r)
		}
	})

	t.Run("Parse Data twice", func(t *testing.T) {
		r, err := tmpl.ParseTmplString(0, "{{ $labels.instance }}", testData)
		if err != nil {
			t.Fatalf("parse data get err: %+v", err)
		}
		if r != "1.1.1.1" {
			t.Errorf("parse data failed, except %s get %s","1.1.1.1", r)
		}

		r, err = tmpl.ParseTmplString(1, "{{ $labels.instance }}", testData)
		if err != nil {
			t.Fatalf("parse data get err: %+v", err)
		}
		if r != "1.1.1.2" {
			t.Errorf("parse data failed, except %s get %s","1.1.1.2", r)
		}
	})

	t.Run("Parse functiong", func(t *testing.T) {
		r, err := tmpl.ParseTmplString(0, "{{ uuid }}", testData)
		if err != nil {
			t.Fatalf("parse data get err: %+v", err)
		}
		fmt.Println(r)
	})

}
