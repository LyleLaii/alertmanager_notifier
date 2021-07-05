# A simple build script, just for specify version info
CGO_ENABLED=1 GOOS=linux  go build -o out/alertmanager_notifier \
        -ldflags "-X alertmanager_notifier/pkg/version.Version=`cat VERSION` -X alertmanager_notifier/pkg/version.Revision=`git rev-parse HEAD` -X alertmanager_notifier/pkg/version.Branch=`git rev-parse --abbrev-ref HEAD` -X alertmanager_notifier/pkg/version.BuildUser=`whoami` -X alertmanager_notifier/pkg/version.BuildDate=`date +%Y%m%d-%H:%M:%S`"  \
        -v -a -trimpath cmd/alertmanager_notifier/main.go
