package notifiers

import (
	"alertmanager_notifier/modules/db"
	"github.com/prometheus/alertmanager/template"
	"strings"
	"time"
)

// Notifier Notifier instance
type Notifier interface {
	Notify(*AlertMessage)
}

// NotifyIntegration notify integration
type NotifyIntegration struct {
	Name         string
	ReceiverType string
	Notifier     Notifier
}

type ErrorInfo struct {
	Name string
	Err  string
}


// Notify notify instance send message
func (n NotifyIntegration) Notify(am *AlertMessage) {
	n.Notifier.Notify(am)
}

// AlertMessage alertmessage format
type AlertMessage struct {
	Receiver     string
	ReceiverDate time.Time
	AlertInfo    AlertWebhookInfo
}

// AlertWebhookInfo alertmanager webhook format
type AlertWebhookInfo template.Data


// SearchReceiver transe webhook receiver to user
func SearchReceiver(receiver string, strtime string, receiverType string) (receiverName string) {
	var receivers []string
	r := strings.Split(receiver, ":")
	department := r[0]
	group := r[1]
	u := db.GetStandbyUser(department, group, strtime)
	switch receiverType {
	case "uid":
		receivers = u
	case "comment":
		receivers = db.GetComment(u)
	default:
		receivers = db.GetReceiverNameByType(u, receiverType)
	}
	receiverName = strings.Join(receivers, ",")
	return
}

//// InitNotifyIntergration generate notify integration
//func InitNotifyIntergration(rs []config.Receiver) {
//	var (
//		errs                 []ErrorInfo
//		notifierintergration []NotifyIntegration
//		add                  = func(name string, receivertype string, f func() (Notifier, error)) {
//			n, e := f()
//			if e != nil {
//				errs = append(errs, ErrorInfo{
//					Name: name,
//					Err:  fmt.Sprintf("%v", e),
//				})
//				return
//			}
//			notifierintergration = append(notifierintergration, NotifyIntegration{
//				Name:         name,
//				ReceiverType: receivertype,
//				Notifier:     n,
//			})
//		}
//	)
//	for _, cs := range rs {
//		for _, c := range cs.KafkaConfigs {
//			add(cs.Name, cs.ReceiverType, func() (Notifier, error) { return kafka.New(cs.Name, &c) })
//		}
//		for _, c := range cs.WebhookConfigs {
//			add(cs.Name, cs.ReceiverType, func() (Notifier, error) { return webhook.New(cs.Name, &c) })
//		}
//		for _, c := range cs.ShellConfigs {
//			add(cs.Name, cs.ReceiverType, func() (Notifier, error) { return shell.New(cs.Name, &c) })
//		}
//	}
//	if len(errs) > 0 {
//		log.Panic("Notify", fmt.Sprintf("Init notify intergration get errors: %+v", errs))
//	}
//	NotifyIntergrations = notifierintergration
//}
//
//// NotifyIntergrationMap notify instance by map
//type NotifyIntergrationMap map[string][]NotifyIntegration
//
//// NotifyMap notify integration map
//var NotifyMap NotifyIntergrationMap
//
//// InitNotifyIntergrationMap generate notify integration by map format
//func InitNotifyIntergrationMap(rs []config.Receiver) {
//	var (
//		errs        []ErrorInfo
//		notifiermap NotifyIntergrationMap
//		add         = func(name string, receivertype string, f func() (Notifier, error)) {
//			n, e := f()
//			if e != nil {
//				errs = append(errs, ErrorInfo{
//					Name: name,
//					Err:  fmt.Sprintf("%v", e),
//				})
//				return
//			}
//			notifiermap[name] = append(NotifyMap[name], NotifyIntegration{
//				Name:         name,
//				ReceiverType: receivertype,
//				Notifier:     n,
//			})
//		}
//	)
//
//	for _, cs := range rs {
//		for _, c := range cs.KafkaConfigs {
//			add(cs.Name, cs.ReceiverType, func() (Notifier, error) { return kafka.New(cs.Name, &c) })
//		}
//		for _, c := range cs.WebhookConfigs {
//			add(cs.Name, cs.ReceiverType, func() (Notifier, error) { return webhook.New(cs.Name, &c) })
//		}
//	}
//
//	if len(errs) > 0 {
//		log.Panic("Notify", fmt.Sprintf("Init notify intergration get errors: %+v", errs))
//	}
//	NotifyMap = notifiermap
//}
