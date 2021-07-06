package main

import (
	"alertmanager_notifier/config"
	"alertmanager_notifier/log"
	"alertmanager_notifier/metrics"
	"alertmanager_notifier/middleware"
	"alertmanager_notifier/modules/db"
	"alertmanager_notifier/notifiers"
	alertApi "alertmanager_notifier/notifiers/api"
	"alertmanager_notifier/notifiers/kafka"
	"alertmanager_notifier/notifiers/shell"
	"alertmanager_notifier/notifiers/webhook"
	"alertmanager_notifier/pkg/utils"
	"alertmanager_notifier/pkg/version"
	"alertmanager_notifier/router"
	"alertmanager_notifier/template"
	"alertmanager_notifier/web"
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const SERVERNAME = "AlertManager-Notifier"

func buildReceiverIntegrations(rs []config.Receiver, tmpl *template.Template, logger log.Logger) ([]notifiers.NotifyIntegration, error) {
	var (
		errs                 []notifiers.ErrorInfo
		notifierIntegrations []notifiers.NotifyIntegration
		add                  = func(name string, receiverType string, f func() (notifiers.Notifier, error)) {
			n, e := f()
			if e != nil {
				errs = append(errs, notifiers.ErrorInfo{
					Name: name,
					Err:  fmt.Sprintf("%v", e),
				})
				return
			}
			notifierIntegrations = append(notifierIntegrations, notifiers.NotifyIntegration{
				Name:         name,
				ReceiverType: receiverType,
				Notifier:     n,
			})
		}
	)
	for _, cs := range rs {
		for _, c := range cs.KafkaConfigs {
			add(cs.Name, cs.ReceiverType, func() (notifiers.Notifier, error) { return kafka.New(cs.Name, &c, tmpl, logger) })
		}
		for _, c := range cs.WebhookConfigs {
			add(cs.Name, cs.ReceiverType, func() (notifiers.Notifier, error) { return webhook.New(cs.Name, &c, tmpl, logger) })
		}
		for _, c := range cs.ShellConfigs {
			add(cs.Name, cs.ReceiverType, func() (notifiers.Notifier, error) { return shell.New(cs.Name, &c, tmpl, logger) })
		}
	}
	// TODO fix it
	if len(errs) > 0 {
		logger.Panic("Notify", fmt.Sprintf("Init notify intergration get errors: %+v", errs))
	}
	return notifierIntegrations, nil
}

func main(){
	os.Exit(run())
}

func run() int {
	var (
		cfg = kingpin.Flag("config.file", "AlertManager-Notifier configuration file path. Default is ./conf/settings.yaml").Default("./conf/settings.yaml").String()
		port = kingpin.Flag("web.port", "Port to listen on for the web interface and API. Default is :8080").Default(":8080").String()
		maxBackups = kingpin.Flag(config.LogMaxBackupsFlagName, config.LogMaxBackupsFlagHelp).Default("5").Int()
		maxDays = kingpin.Flag(config.LogMaxDaysFlagName, config.LogMaxDaysFlagHelp).Default("30").Int()
	)
	runConfig := config.RunningConfig{MaxBackups: *maxBackups, MaxDays: *maxDays}
	config.AddFlags(kingpin.CommandLine, &runConfig)
	kingpin.Version(version.Print())
	kingpin.CommandLine.GetFlag("help").Short('h')
	kingpin.Parse()

	if err := utils.InitConfig(*cfg); err != nil {
		panic(err)
	}

	// TODO: Should Gin use another logger? use same file"
	loggerConfig := log.ConfigZap(SERVERNAME, &runConfig)
	logger := log.NewZapSugarLogger(loggerConfig)
	//ginLogger := log.NewZapSugarLoggerGin(loggerConfig)

	logger.Info(SERVERNAME, fmt.Sprintf("Version InfoContext: %s", version.InfoContext()))
	logger.Info(SERVERNAME, fmt.Sprintf("Build InfoContext: %s", version.BuildContext()))
	logger.Info(SERVERNAME, fmt.Sprintf("ListenPort: %s, RunMode: %v, UserRota: %v", *port, runConfig.RunMode.String(), utils.GetUseRote()))

	db.InitDBHandler(logger, &runConfig)

	tmpl, _ := template.GenTempInstance("data/custom.tpl")

	var receivers []config.Receiver
	if err := viper.UnmarshalKey("receivers", &receivers); err != nil {
		logger.Error("GetReceivers", fmt.Sprintf("Error load receivers config: %s", err))
	}

	nis, _ := buildReceiverIntegrations(receivers, tmpl, logger)


	routers := router.InitRouter(logger, &runConfig)
	//pprof.Register(routers)
	metrics.Register(routers)
	router.AddMiddleware(routers, middleware.MwPrometheusHTTP)
	alertApi.Register(routers.Group("/api/v1"), nis, logger)
	web.RegisterAPI(routers.Group("/api/v1"), logger)
	web.RegisterUI(routers, logger)

	srv := &http.Server{
		Addr:    *port,
		Handler: routers,
	}
	srvc := make(chan struct{})

	go func() {
		logger.Info(SERVERNAME, fmt.Sprintf("Listening port: %s", *port))
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Error(SERVERNAME, fmt.Sprintf("Listen error: %+v", err))
			close(srvc)
		}
		defer func() {
			if err := srv.Close(); err != nil {
				logger.Error(SERVERNAME, fmt.Sprintf("Error on closing the server: %+v", err))
			}
		}()
	}()

	var (
		hup      = make(chan os.Signal, 1)
		hupReady = make(chan bool)
		term     = make(chan os.Signal, 1)
	)
	signal.Notify(hup, syscall.SIGHUP)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-hupReady
		for {
			select {
			case <-hup:
				// TODO
				// ignore error, already logged in `reload()`
				logger.Info(SERVERNAME,"receive hup signal")
			}
		}
	}()

	// Wait for reload or termination signals.
	close(hupReady) // Unblock SIGHUP handler.

	for {
		select {
		case <-term:
			logger.Info(SERVERNAME, "Received SIGTERM, exiting gracefully...")
			return 0
		case <-srvc:
			return 1
		}
	}
}